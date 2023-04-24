package parser

import (
	"fmt"
	"strconv"

	"github.com/odas0r/yail/ast"
	"github.com/odas0r/yail/lexer"
	"github.com/odas0r/yail/token"
)

// We use iota to give the following constants incrementing numbers as values.
//
// Which numbers we use doesn’t matter, but the order and the relation to each
// other do.
const (
	_ int = iota
	LOWEST
	EQUALS      // ==
	LESSGREATER // > or <
	SUM         // +
	PRODUCT     // *
	PREFIX      // -X or !X
	CALL        // myFunction(X)
)

var precedences = map[token.TokenType]int{
	token.EQ:       EQUALS,      // == and != have the same precedence
	token.NOT_EQ:   EQUALS,      // == and != have the same precedence
	token.LT:       LESSGREATER, // <
	token.GT:       LESSGREATER, // >
	token.PLUS:     SUM,         // +
	token.MINUS:    SUM,         // -
	token.SLASH:    PRODUCT,     // /
	token.ASTERISK: PRODUCT,     // *
	token.LPAREN:   CALL,
}

type Parser struct {
	l *lexer.Lexer

	errors    []string
	curToken  token.Token
	peekToken token.Token

	prefixParseFns map[token.TokenType]prefixParserFn
	infixParseFns  map[token.TokenType]infixParseFn
}

type (
	prefixParserFn func() ast.Expression
	infixParseFn   func(ast.Expression) ast.Expression
)

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	// Read two tokens, so curToken and peekToken are both set
	p.nextToken()
	p.nextToken()

	p.prefixParseFns = make(map[token.TokenType]prefixParserFn)
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.FLOAT, p.parseFloatLiteral)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.TRUE, p.parseBoolean)
	p.registerPrefix(token.FALSE, p.parseBoolean)
	p.registerPrefix(token.LPAREN, p.parseGroupExpression)
	p.registerPrefix(token.IF, p.parseIfStatement)
	p.registerPrefix(token.STRING, p.parseStringLiteral)

	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NOT_EQ, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)
	p.registerInfix(token.LPAREN, p.parseCallExpression)

	return p
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.IDENT:
		return p.parseVariableStatement()
	case token.STRUCTS:
		return p.parseStructsStatement()
	default:
		return p.parseExpressionStatement()
	}
}

// ---------------------- parsers ----------------------

func (p *Parser) parseVariableStatement() ast.Statement {
	varStmt := &ast.VarStatement{
		Token: p.curToken,
	}

	if p.peekTokenIs(token.ASSIGN) {
		varStmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		p.nextToken()

		// VECTOR
		if p.peekTokenIs(token.LBRACE) {
			p.nextToken()
			p.nextToken()

			vecValues := p.parseExpressionList()

			vecStmt := &ast.VectorStatement{
				Token:  varStmt.Token,
				Name:   varStmt.Name,
				Values: vecValues,
				Size: &ast.IntegerLiteral{
					Token: token.Token{Type: token.INT, Literal: strconv.Itoa(len(vecValues))},
					Value: int64(len(vecValues)),
				},
			}

			// Expect the '}' token
			if !p.expectPeek(token.RBRACE) {
				return nil
			}

			if p.peekTokenIs(token.SEMICOLON) {
				p.nextToken()
			}

			return vecStmt
		} else {
			// VARIABLE
			p.nextToken()

			varStmt.Value = p.parseExpression(LOWEST)

			// Expect the semicolon
			if p.peekTokenIs(token.SEMICOLON) {
				p.nextToken()
			}

			return varStmt
		}
	}

	// Expect the variable name
	if !p.expectPeek(token.IDENT) {
		return nil
	}
	varStmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if p.peekTokenIs(token.LBRACKET) {
		vecSize := p.parseVectorSize()

		// read the [<expr>] part
		vecStmt := &ast.VectorStatement{
			Token: varStmt.Token,
			Name:  varStmt.Name,
			Size:  vecSize,
		}

		if p.peekTokenIs(token.SEMICOLON) {
			sizeLiteral, ok := vecStmt.Size.(*ast.IntegerLiteral)

			if ok {
				vecStmt.Values = make([]ast.Expression, sizeLiteral.Value)
				for i := 0; i < int(sizeLiteral.Value); i++ {
					vecStmt.Values[i] = p.defaultValueForType(vecStmt.Token)
				}
			}

			return vecStmt
		} else if p.peekTokenIs(token.ASSIGN) {
			p.nextToken()

			// Expect the '{' token
			if !p.expectPeek(token.LBRACE) {
				return nil
			}
			p.nextToken()

			vecStmt.Values = p.parseExpressionList()

			// Expect the '}' token
			if !p.expectPeek(token.RBRACE) {
				return nil
			}

			// Expect the semicolon
			if p.peekTokenIs(token.SEMICOLON) {
				p.nextToken()
			}

			return vecStmt
		}
	}

	if p.peekTokenIs(token.ASSIGN) {
		p.nextToken()
		p.nextToken()

		varStmt.Value = p.parseExpression(LOWEST)
	} else {
		varStmt.Value = p.defaultValueForType(varStmt.Token)
	}

	// Expect the semicolon
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return varStmt
}

func (p *Parser) parseVectorSize() ast.Expression {
	if !p.expectPeek(token.LBRACKET) {
		return nil
	}

	var size ast.Expression

	if !p.peekTokenIs(token.RBRACKET) {
		p.nextToken()
		size = p.parseExpression(LOWEST)
	}

	p.nextToken()

	if p.peekTokenIs(token.SEMICOLON) {
		size = &ast.IntegerLiteral{Token: token.Token{Type: token.INT, Literal: "1"}, Value: 1}
	}

	return size
}

func (p *Parser) parseExpressionList() []ast.Expression {
	expressions := []ast.Expression{}

	if p.curTokenIs(token.RBRACE) {
		return expressions
	}

	expressions = append(expressions, p.parseExpression(LOWEST))

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		expressions = append(expressions, p.parseExpression(LOWEST))
	}

	return expressions
}

func (p *Parser) parseStructsStatement() *ast.StructsDefinition {
	sd := &ast.StructsDefinition{Token: p.curToken}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}
	p.nextToken()

	var structs []*ast.StructLiteral
	structNames := make(map[string]bool)

	for !p.peekTokenIs(token.RBRACE) {
		ps := p.parseStructLiteral()

		// check if struct already exists
		name := ps.TokenLiteral()
		if _, exists := structNames[name]; exists {
			p.errors = append(p.errors, fmt.Sprintf("Duplicate struct name '%s' in structs", name))
			return nil
		}
		structNames[name] = true

		structs = append(structs, ps)
	}

	// save parsed structs
	sd.Structs = structs

	if !p.expectPeek(token.RBRACE) {
		return nil
	}

	return sd
}

func (p *Parser) parseStructLiteral() *ast.StructLiteral {
	sl := &ast.StructLiteral{Token: p.curToken}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}
	// expect attribute type
	if !p.expectPeek(token.IDENT) {
		return nil
	}

	var attributes []*ast.Attribute
	attributeNames := make(map[string]bool)

	attr := &ast.Attribute{Token: p.curToken}

	// expect attribute name
	if !p.expectPeek(token.IDENT) {
		return nil
	}

	// check if attribute already exists
	name := p.curToken.Literal
	if _, exists := attributeNames[name]; exists {
		p.errors = append(p.errors, fmt.Sprintf("Duplicate attribute name '%s' in struct %s", name, sl.TokenLiteral()))
		return nil
	}
	attributeNames[name] = true

	attr.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if p.peekTokenIs(token.LBRACKET) {
		attr.IsVector = true
		p.nextToken()

		if !p.peekTokenIs(token.RBRACKET) {
			p.nextToken()
			attr.Size = p.parseExpression(LOWEST)
		}

		if !p.expectPeek(token.RBRACKET) {
			return nil
		}
	}

	attributes = append(attributes, attr)

	if p.peekTokenIs(token.COMMA) {
		p.nextToken() // consume the comma
		p.nextToken() // consume the comma

		// iterate through the attributes till you find a semicolon or a type
		for {
			var attrToken token.Token

			// if it's a type identifier, break
			if p.curTokenIs(token.IDENT) {
				attrToken = p.curToken

				if !p.expectPeek(token.IDENT) {
					return nil
				}
			} else {
				attrToken = attr.Token
			}

			attr := &ast.Attribute{
				Token: attrToken, // take the token from the first attribute (type)
				Name:  &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal},
			}

			// check if attribute already exists
			name := p.curToken.Literal
			if _, exists := attributeNames[name]; exists {
				p.errors = append(p.errors, fmt.Sprintf("Duplicate attribute name '%s' in struct %s", name, sl.TokenLiteral()))
				return nil
			}
			attributeNames[name] = true

			if p.peekTokenIs(token.LBRACKET) {
				// set attribute as vector
				attr.IsVector = true

				p.nextToken()

				if !p.peekTokenIs(token.RBRACKET) {
					p.nextToken()
					attr.Size = p.parseExpression(LOWEST)
				}

				// Expect the ']' token
				if !p.expectPeek(token.RBRACKET) {
					return nil
				}
			}

			attributes = append(attributes, attr)

			if p.peekTokenIs(token.COMMA) {
				p.nextToken()
				p.nextToken()
			} else {
				// Expect ';' token to end the attribute
				if !p.expectPeek(token.SEMICOLON) {
					return nil
				}
				break
			}
		}
	} else {
		// Expect ';' token to end the attribute
		if !p.expectPeek(token.SEMICOLON) {
			return nil
		}
	}

	sl.Attributes = attributes

	// Expect '} ;' tokens
	if !p.expectPeek(token.RBRACE) {
		return nil
	}
	if !p.expectPeek(token.SEMICOLON) {
		return nil
	}

	if !p.peekTokenIs(token.RBRACE) {
		p.nextToken()
	}

	return sl
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}

	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}

	leftExp := prefix()

	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}

		p.nextToken()

		leftExp = infix(leftExp)
	}

	return leftExp
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.curToken}

	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	lit.Value = value

	return lit
}

func (p *Parser) parseFloatLiteral() ast.Expression {
	fl := &ast.FloatLiteral{Token: p.curToken}

	value, err := strconv.ParseFloat(p.curToken.Literal, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as float", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	fl.Value = value

	return fl
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}

	p.nextToken()

	expression.Right = p.parseExpression(PREFIX)

	return expression
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}

	precedence := p.curPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)

	return expression
}

func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{
		Token: p.curToken,
		Value: p.curTokenIs(token.TRUE),
	}
}

func (p *Parser) parseGroupExpression() ast.Expression {
	p.nextToken()

	exp := p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return exp
}

func (p *Parser) parseIfStatement() ast.Expression {
	expression := &ast.IfExpression{Token: p.curToken}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	p.nextToken()
	expression.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	expression.Consequence = p.parseBlockStatement()

	if p.peekTokenIs(token.ELSE) {
		p.nextToken()

		if !p.expectPeek(token.LBRACE) {
			return nil
		}

		expression.Alternative = p.parseBlockStatement()
	}

	return expression
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.curToken}
	block.Statements = []ast.Statement{}

	p.nextToken()

	for !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}

	return block
}

// func (p *Parser) parseFunctionLiteral() ast.Expression {
// 	lit := &ast.FunctionLiteral{Token: p.curToken}
//
// 	if !p.expectPeek(token.LPAREN) {
// 		return nil
// 	}
//
// 	lit.Parameters = p.parseFunctionParameters()
//
// 	if !p.expectPeek(token.LBRACE) {
// 		return nil
// 	}
//
// 	lit.Body = p.parseBlockStatement()
//
// 	return lit
// }
//
// func (p *Parser) parseFunctionParameters() []*ast.Identifier {
// 	identifiers := []*ast.Identifier{}
//
// 	if p.peekTokenIs(token.RPAREN) {
// 		p.nextToken()
// 		return identifiers
// 	}
//
// 	p.nextToken()
//
// 	ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
// 	identifiers = append(identifiers, ident)
//
// 	for p.peekTokenIs(token.COMMA) {
// 		p.nextToken()
// 		p.nextToken()
//
// 		ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
// 		identifiers = append(identifiers, ident)
// 	}
//
// 	if !p.expectPeek(token.RPAREN) {
// 		return nil
// 	}
//
// 	return identifiers
// }

func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	exp := &ast.CallExpression{Token: p.curToken, Function: function}
	exp.Arguments = p.parseCallArguments()

	return exp
}

func (p *Parser) parseCallArguments() []ast.Expression {
	args := []ast.Expression{}

	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return args
	}

	p.nextToken()
	args = append(args, p.parseExpression(LOWEST))

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		args = append(args, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return args
}

func (p *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) defaultValueForType(t token.Token) ast.Expression {
	switch t.Literal {
	case "int":
		return &ast.IntegerLiteral{Token: token.Token{Type: token.INT, Literal: "0"}, Value: 0}
	case "float":
		return &ast.FloatLiteral{Token: token.Token{Type: token.FLOAT, Literal: "0"}, Value: 0}
	case "bool":
		return &ast.Boolean{Token: token.Token{Type: token.FALSE, Literal: "false"}, Value: false}
	default:
		return nil
	}
}

// ---------------------- errors ----------------------

func (p *Parser) peekError(t token.TokenType) {
	msg := "expected next token to be %s, got %s instead"
	p.errors = append(p.errors, fmt.Sprintf(msg, t, p.peekToken.Type))
}
func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}

// ---------------------- helpers ----------------------

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParserFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}

	return LOWEST
}

func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}

	return LOWEST
}
