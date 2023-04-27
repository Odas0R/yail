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
	INDEX       // array[index]
	ACCESSOR    // a.b
)

var precedences = map[token.TokenType]int{
	token.EQ:     EQUALS, // == and != have the same precedence
	token.NOT_EQ: EQUALS, // == and != have the same precedence

	token.LT:  LESSGREATER, // <
	token.GT:  LESSGREATER, // >
	token.LTE: LESSGREATER, // <=
	token.GTE: LESSGREATER, // >=
	token.AND: LESSGREATER, // and
	token.OR:  LESSGREATER, // or

	token.PLUS:     SUM,      // +
	token.MINUS:    SUM,      // -
	token.SLASH:    PRODUCT,  // /
	token.ASTERISK: PRODUCT,  // *
	token.LPAREN:   CALL,     // myFunction(X)
	token.LBRACKET: INDEX,    // array[index]
	token.ACCESSOR: ACCESSOR, // a.b
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

	p.prefixParseFns = make(map[token.TokenType]prefixParserFn)
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.FLOAT, p.parseFloatLiteral)
	p.registerPrefix(token.STRING, p.parseStringLiteral)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.TRUE, p.parseBoolean)
	p.registerPrefix(token.FALSE, p.parseBoolean)
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)
	p.registerPrefix(token.IF, p.parseIfExpression)

	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NOT_EQ, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)

	p.registerInfix(token.LTE, p.parseInfixExpression)
	p.registerInfix(token.GTE, p.parseInfixExpression)
	p.registerInfix(token.AND, p.parseInfixExpression)
	p.registerInfix(token.OR, p.parseInfixExpression)

	p.registerInfix(token.LPAREN, p.parseCallExpression)
	p.registerInfix(token.LBRACKET, p.parseIndexExpression)
	p.registerInfix(token.ACCESSOR, p.parseAccessorExpression)

	// Read two tokens, so curToken and peekToken are both set
	p.nextToken()
	p.nextToken()

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
		if p.peekTokenIs(token.LPAREN) { // FUNCTION
			if p.isNextTokenFunctionStatement() {
				return p.parseFunctionStatement()
			} else {
				return p.parseExpressionStatement()
			}
		} else if p.peekTokenIs(token.IDENT) || p.peekTokenIs(token.ASSIGN) { // VARIABLE
			return p.parseVariableStatement()
		}
	case token.GLOBAL:
		return p.parseGlobalStatement()
	case token.CONST:
		return p.parseConstStatement()
	case token.STRUCTS:
		return p.parseStructsStatement()
	case token.WHILE:
		return p.parseWhileStatement()
	case token.FOR:
		return p.parseForStatement()
	}

	targetExpression := p.parseExpression(LOWEST)

	if p.peekTokenIs(token.INCREMENT) {
		p.nextToken()
		return p.parseIncrementStatement(targetExpression)
	} else if p.peekTokenIs(token.DECREMENT) {
		p.nextToken()
		return p.parseDecrementStatement(targetExpression)
	} else if p.peekTokenIs(token.PLUS_EQ) {
		p.nextToken()
		return p.parsePlusEqualsStatement(targetExpression)
	} else if p.peekTokenIs(token.MULT_EQ) {
		p.nextToken()
		return p.parseMultEqualsStatement(targetExpression)
	} else if p.peekTokenIs(token.MINUS_EQ) {
		p.nextToken()
		return p.parseMinusEqualsStatement(targetExpression)
	} else if p.peekTokenIs(token.ASSIGN) {
		p.nextToken()
		return p.parseAssignmentStatement(targetExpression)
	}

	stmt := &ast.ExpressionStatement{Token: p.curToken, Expression: targetExpression}

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// ---------------------- parsers ----------------------

func (p *Parser) parseVariableStatement() ast.Statement {

	if p.peekTokenIs(token.ASSIGN) {
		return p.parseVariableAssignStatment()
	}

	var statements []ast.Statement

	curToken := p.curToken
	curType := p.parseType()

	for {
		// Expect the variable name
		if !p.expectPeek(token.IDENT) {
			return nil
		}

		name := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

		// ----------------
		// Vector
		// ----------------

		if p.peekTokenIs(token.LBRACKET) {
			statements = append(statements, p.parseVectorStatement(curToken, curType, name))
		} else {

			// ----------------
			// Variable
			// ----------------

			if p.peekTokenIs(token.ASSIGN) {
				p.nextToken()
				p.nextToken()

				expr := p.parseExpression(LOWEST)

				stmt := &ast.VariableStatement{
					Token: curToken,
					Name:  name,
					Type:  curType,
					Value: expr,
				}

				statements = append(statements, stmt)
			} else {
				stmt := &ast.VariableStatement{
					Token: curToken,
					Name:  name,
					Type:  curType,
					Value: p.defaultValueForType(curToken),
				}

				statements = append(statements, stmt)
			}
		}

		// If the next token is a comma, continue parsing additional variable declarations
		if p.peekTokenIs(token.COMMA) {
			p.nextToken() // Consume the ',' token
		} else {
			break
		}
	}

	// Expect the semicolon
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	if len(statements) == 1 {
		return statements[0]
	} else {
		return &ast.BlockStatement{Statements: statements}
	}
}

func (p *Parser) parseVectorStatement(curToken token.Token, curType *ast.Identifier, name *ast.Identifier) ast.Statement {
	vecSize := p.parseVectorSize()

	vecStmt := &ast.VectorStatement{
		Token: curToken,
		Name:  name,
		Size:  vecSize,
		Type:  curType,
	}

	if !p.peekTokenIs(token.ASSIGN) {
		sizeLiteral, ok := vecStmt.Size.(*ast.IntegerLiteral)
		if ok {
			values := make([]ast.Expression, sizeLiteral.Value)
			for i := 0; i < int(sizeLiteral.Value); i++ {
				values[i] = p.defaultValueForType(vecStmt.Token)
			}

			vecStmt.Values = &ast.VectorLiteral{
				Token:  token.Token{Type: token.LBRACE, Literal: "{"},
				Values: values,
			}
		}
	} else {
		p.nextToken()

		// Expect the '{' token
		if !p.expectPeek(token.LBRACE) {
			return nil
		}
		p.nextToken()

		vecStmt.Values = p.parseVectorLiteral()

		// Set the size if it wasn't set before
		if vl, ok := vecStmt.Values.(*ast.VectorLiteral); ok {
			vecStmt.Size = &ast.IntegerLiteral{
				Token: token.Token{Type: token.INT, Literal: strconv.Itoa(len(vl.Values))},
				Value: int64(len(vl.Values)),
			}
		}

		// Expect the '}' token
		if !p.expectPeek(token.RBRACE) {
			return nil
		}
	}

	return vecStmt
}

func (p *Parser) parseVariableAssignStatment() ast.Statement {
	curToken := p.curToken

	name := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	p.nextToken()

	if p.peekTokenIs(token.LBRACE) {
		p.nextToken()
		p.nextToken()

		values := p.parseVectorLiteral()

		if !p.expectPeek(token.RBRACE) {
			return nil
		}

		if p.peekTokenIs(token.SEMICOLON) {
			p.nextToken()
		}

		return &ast.AssignmentStatement{
			Token: curToken,
			Left:  name,
			Value: values,
		}
	} else {
		p.nextToken()

		value := p.parseExpression(LOWEST)

		// Expect the semicolon
		if p.peekTokenIs(token.SEMICOLON) {
			p.nextToken()
		}

		return &ast.AssignmentStatement{
			Token: curToken,
			Left:  name,
			Value: value,
		}
	}
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

func (p *Parser) parseVectorLiteral() ast.Expression {
	curToken := p.curToken
	expressions := []ast.Expression{}

	if p.curTokenIs(token.RBRACE) {
		return &ast.VectorLiteral{Token: curToken, Values: expressions}
	}

	expressions = append(expressions, p.parseExpression(LOWEST))

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		expressions = append(expressions, p.parseExpression(LOWEST))
	}

	return &ast.VectorLiteral{Token: curToken, Values: expressions}
}

func (p *Parser) parseStructsStatement() *ast.StructsStatement {
	sd := &ast.StructsStatement{Token: p.curToken}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}
	p.nextToken()

	var structs []*ast.Struct

	for !p.peekTokenIs(token.RBRACE) {
		sl := &ast.Struct{Token: p.curToken}

		if !p.expectPeek(token.LBRACE) {
			return nil
		}

		sl.Attributes = p.parseStructAttributes()

		if !p.peekTokenIs(token.RBRACE) {
			p.nextToken()
		}

		structs = append(structs, sl)
	}

	// save parsed structs
	sd.Structs = structs

	if !p.expectPeek(token.RBRACE) {
		return nil
	}

	return sd
}

func (p *Parser) parseStructAttributes() []*ast.Attribute {
	attributes := []*ast.Attribute{}

	if p.peekTokenIs(token.RBRACE) {
		p.nextToken()
		return attributes
	}
	p.nextToken()

	topAttr := &ast.Attribute{
		Token: p.curToken,
		Type:  p.parseType(),
	}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	topAttr.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if p.peekTokenIs(token.LBRACKET) {
		p.nextToken()
		p.nextToken()

		if p.curTokenIs(token.RBRACKET) {
			topAttr.Size = &ast.IntegerLiteral{Token: token.Token{Type: token.INT, Literal: "1"}, Value: 1}
		} else {
			topAttr.Size = p.parseExpression(LOWEST)

			if !p.expectPeek(token.RBRACKET) {
				return nil
			}
		}
		topAttr.IsVector = true
	}

	attributes = append(attributes, topAttr)

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()

		attr := &ast.Attribute{
			Token: p.curToken,
		}

		if !p.peekTokenIs(token.IDENT) {
			attr.Type = topAttr.Type // set type top attribute type
			attr.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		} else {
			attr.Type = p.parseType()

			if !p.expectPeek(token.IDENT) {
				return nil
			}

			attr.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		}

		// IS ARRAY
		if p.peekTokenIs(token.LBRACKET) {
			p.nextToken()
			p.nextToken()

			if p.curTokenIs(token.RBRACKET) {
				attr.Size = &ast.IntegerLiteral{Token: token.Token{Type: token.INT, Literal: "1"}, Value: 1}
			} else {
				attr.Size = p.parseExpression(LOWEST)

				if !p.expectPeek(token.RBRACKET) {
					return nil
				}
			}

			attr.IsVector = true
		}

		attributes = append(attributes, attr)
	}

	if !p.expectPeek(token.SEMICOLON) {
		return nil
	}

	if !p.expectPeek(token.RBRACE) {
		return nil
	}

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return attributes
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

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()

	exp := p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return exp
}

func (p *Parser) parseIfExpression() ast.Expression {
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
		var stmt ast.Statement

		if p.curTokenIs(token.LOCAL) {
			stmt = p.parseLocalStatement()
		} else {
			stmt = p.parseStatement()
		}

		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}

	return block
}

func (p *Parser) parseVariableBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.curToken}
	block.Statements = []ast.Statement{}

	p.nextToken()

	for !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()

		_, isBlock := stmt.(*ast.BlockStatement)

		if isBlock {
			for _, s := range stmt.(*ast.BlockStatement).Statements {
				_, isVariable := s.(*ast.VariableStatement)
				_, isVector := s.(*ast.VectorStatement)
				_, isAssign := stmt.(*ast.AssignmentStatement)

				if !isVariable && !isVector && !isAssign {
					p.addError("Only variable declarations are allowed in variable blocks")
					p.nextToken()
					return nil
				}
			}
		} else {
			_, isVariable := stmt.(*ast.VariableStatement)
			_, isVector := stmt.(*ast.VectorStatement)
			_, isAssign := stmt.(*ast.AssignmentStatement)

			if !isVariable && !isVector && !isAssign {
				p.addError("Only variable declarations are allowed in variable blocks")
				p.nextToken()
				return nil
			}
		}

		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}

	return block
}

func (p *Parser) parseFunctionStatement() ast.Statement {
	fuc := &ast.FunctionStatement{Token: p.curToken} // IDENT

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	fuc.Parameters = p.parseFunctionParameters()

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	fuc.ReturnType = p.parseFunctionReturnType()

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	fuc.Body = p.parseBlockStatement()

	return fuc
}

func (p *Parser) parseFunctionParameters() []*ast.Parameter {
	parameters := []*ast.Parameter{}

	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return parameters
	}
	p.nextToken()

	topParam := &ast.Parameter{
		Token: p.curToken,
		Type:  p.parseType(),
	}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	topParam.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if p.peekTokenIs(token.LBRACKET) {
		p.nextToken()
		p.nextToken()

		if p.curTokenIs(token.RBRACKET) {
			topParam.Size = &ast.IntegerLiteral{Token: token.Token{Type: token.INT, Literal: "1"}, Value: 1}
		} else {
			topParam.Size = p.parseExpression(LOWEST)

			if !p.expectPeek(token.RBRACKET) {
				return nil
			}
		}
		topParam.IsVector = true
	}

	parameters = append(parameters, topParam)

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()

		param := &ast.Parameter{
			Token: p.curToken,
		}

		if !p.peekTokenIs(token.IDENT) {
			param.Type = topParam.Type // set the top attribute type as default
			param.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		} else {
			param.Type = p.parseType()

			if !p.expectPeek(token.IDENT) {
				return nil
			}

			param.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		}

		// IS ARRAY
		if p.peekTokenIs(token.LBRACKET) {
			p.nextToken()
			p.nextToken()

			if p.curTokenIs(token.RBRACKET) {
				param.Size = &ast.IntegerLiteral{Token: token.Token{Type: token.INT, Literal: "1"}, Value: 1}
			} else {
				param.Size = p.parseExpression(LOWEST)

				if !p.expectPeek(token.RBRACKET) {
					return nil
				}
			}

			param.IsVector = true
		}

		parameters = append(parameters, param)
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return parameters
}

func (p *Parser) parseFunctionReturnType() *ast.ReturnType {
	returnType := &ast.ReturnType{
		Token: p.curToken,
		Type:  p.parseType(),
	}

	if p.peekTokenIs(token.LBRACKET) {
		p.nextToken()
		p.nextToken()

		if p.curTokenIs(token.RBRACKET) {
			returnType.Size = &ast.IntegerLiteral{Token: token.Token{Type: token.INT, Literal: "1"}, Value: 1}
		} else {
			returnType.Size = p.parseExpression(LOWEST)

			if !p.expectPeek(token.RBRACKET) {
				return nil
			}
		}
		returnType.IsVector = true
	}

	return returnType
}

func (p *Parser) parseGlobalStatement() ast.Statement {
	gl := &ast.GlobalStatement{
		Token: p.curToken,
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	gl.Body = p.parseVariableBlockStatement()

	return gl
}

func (p *Parser) parseConstStatement() ast.Statement {
	cs := &ast.ConstStatement{
		Token: p.curToken,
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	cs.Body = p.parseVariableBlockStatement()

	return cs
}

func (p *Parser) parseLocalStatement() ast.Statement {
	ls := &ast.LocalStatement{
		Token: p.curToken,
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	ls.Body = p.parseVariableBlockStatement()

	return ls
}

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

func (p *Parser) parseIndexExpression(left ast.Expression) ast.Expression {
	exp := &ast.IndexExpression{Token: p.curToken, Left: left}

	p.nextToken()
	exp.Index = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RBRACKET) {
		return nil
	}

	return exp
}

func (p *Parser) parseAccessorExpression(left ast.Expression) ast.Expression {
	exp := &ast.AccessorExpression{Token: p.curToken, Left: left}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	accessors := []ast.Expression{}
	accessors = append(accessors, &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal})

	for p.peekTokenIs(token.ACCESSOR) {
		p.nextToken()
		p.nextToken()
		accessors = append(accessors, &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal})
	}

	exp.Index = accessors

	return exp
}

func (p *Parser) parseAssignmentStatement(exp ast.Expression) *ast.AssignmentStatement {
	// check if expression is of type accessor or index
	_, isAccessorExpression := exp.(*ast.AccessorExpression)
	_, isIndexExpression := exp.(*ast.IndexExpression)

	if !isAccessorExpression && !isIndexExpression {
		msg := "illegal assignment target"
		p.errors = append(p.errors, msg)
		return nil
	}

	as := &ast.AssignmentStatement{Token: p.curToken, Left: exp}

	p.nextToken()
	as.Value = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return as
}

func (p *Parser) parseWhileStatement() *ast.WhileStatement {
	ws := &ast.WhileStatement{Token: p.curToken}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	p.nextToken()
	ws.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	ws.Body = p.parseBlockStatement()

	return ws
}

func (p *Parser) parseForStatement() *ast.ForStatement {
	fs := &ast.ForStatement{Token: p.curToken}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	fs.Var = p.parseIdentifier()

	if !p.expectPeek(token.COMMA) {
		return nil
	}
	p.nextToken()

	fs.Start = p.parseExpression(LOWEST)

	if !p.expectPeek(token.COMMA) {
		return nil
	}
	p.nextToken()

	fs.End = p.parseExpression(LOWEST)

	if !p.expectPeek(token.COMMA) {
		return nil
	}
	p.nextToken()

	fs.Increment = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	fs.Body = p.parseBlockStatement()

	return fs
}

func (p *Parser) parseIncrementStatement(target ast.Expression) *ast.IncrementStatement {
	stmt := &ast.IncrementStatement{Token: p.curToken, Var: target}
	p.nextToken()
	return stmt
}

func (p *Parser) parseDecrementStatement(target ast.Expression) *ast.DecrementStatement {
	stmt := &ast.DecrementStatement{Token: p.curToken, Var: target}
	p.nextToken()
	return stmt
}

func (p *Parser) parsePlusEqualsStatement(target ast.Expression) *ast.PlusEqualsStatement {
	stmt := &ast.PlusEqualsStatement{Token: p.curToken, Var: target}
	p.nextToken()
	stmt.Quantity = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) parseMultEqualsStatement(target ast.Expression) *ast.MultEqualsStatement {
	stmt := &ast.MultEqualsStatement{Token: p.curToken, Var: target}
	p.nextToken()
	stmt.Quantity = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) parseMinusEqualsStatement(target ast.Expression) *ast.MinusEqualsStatement {
	stmt := &ast.MinusEqualsStatement{Token: p.curToken, Var: target}
	p.nextToken()
	stmt.Quantity = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
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

func (p *Parser) parseType() *ast.Identifier {
	switch p.curToken.Literal {
	case "int":
		return &ast.Identifier{Token: p.curToken, Value: "int"}
	case "float":
		return &ast.Identifier{Token: p.curToken, Value: "float"}
	case "bool":
		return &ast.Identifier{Token: p.curToken, Value: "bool"}
	default:
		if p.peekTokenIs(token.ASSIGN) {
			return &ast.Identifier{Token: p.curToken, Value: "<unknown>"}
		}
		return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	}
}

func (p *Parser) isNextTokenFunctionStatement() bool {
	// Save the current lexer state
	backupPosition := p.l.Position
	backupReadPosition := p.l.ReadPosition
	backupCh := p.l.Ch
	backupCurToken := p.curToken
	backupPeekToken := p.peekToken

	counter := 10000
	for !p.curTokenIs(token.RPAREN) {
		p.nextToken()
		counter--
		if counter == 0 {
			panic("Oops, expected ) but never found it...")
		}
	}

	isFunction := p.peekTokenIs(token.IDENT)

	// Restore the lexer state
	p.l.Position = backupPosition
	p.l.ReadPosition = backupReadPosition
	p.l.Ch = backupCh
	p.curToken = backupCurToken
	p.peekToken = backupPeekToken

	return isFunction
}

// ---------------------- errors ----------------------
func (p *Parser) addError(errorMsg string) {
	msg := "Line %d: %s"
	p.errors = append(p.errors, fmt.Sprintf(msg, p.l.Line, errorMsg))
}

func (p *Parser) peekError(t token.TokenType) {
	msg := "Line %d: Expected next token to be %s, got %s instead"
	p.errors = append(p.errors, fmt.Sprintf(msg, p.l.Line, t, p.peekToken.Type))
}

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("Line %d: No prefix parse function for %s found", p.l.Line, t)
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
