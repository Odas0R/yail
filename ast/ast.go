package ast

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"github.com/odas0r/yail/token"
)

type Node interface {
	TokenLiteral() string
	String() string
	StringAST(level int) string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

func (p *Program) String() string {
	var out bytes.Buffer

	for _, s := range p.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

// Program
// |
// | Statement: VariableStatement
// | | Expression: Identifier (int)
// | | Expression: Identifier (myVar)
// | | Expression: IntegerLiteral (5)
// ...
func (p *Program) PrintAST() string {
	var out bytes.Buffer

	out.WriteString("Program\n")

	for _, s := range p.Statements {
		out.WriteString(s.StringAST(1))
	}

	return out.String()
}

type Identifier struct {
	Token token.Token // the token.IDENT token
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string       { return i.Value }
func (i *Identifier) StringAST(level int) string {
	var out bytes.Buffer

	out.WriteString(strings.Repeat("| ", level))
	out.WriteString("Expression: Identifier (")
	out.WriteString(i.Value)
	out.WriteString(")\n")

	return out.String()
}

type VariableStatement struct {
	Token token.Token // The token.TYPE_INT, token.TYPE_FLOAT, or token.TYPE_BOOL token
	Type  *Identifier // The type of the variable (e.g., int, float, or bool)
	Name  *Identifier // The variable name (e.g., x, y, or z)
	Value Expression  // The value assigned to the variable, can be nil
}

func (vs *VariableStatement) statementNode()       {}
func (vs *VariableStatement) TokenLiteral() string { return vs.Token.Literal }
func (vs *VariableStatement) String() string {
	var out bytes.Buffer

	out.WriteString(vs.Type.String())
	if vs.Type.String() != "" {
		out.WriteString(" ")
	}
	out.WriteString(vs.Name.String())
	out.WriteString(" = ")

	if vs.Value != nil {
		out.WriteString(vs.Value.String())
	}

	out.WriteString(";")
	return out.String()
}
func (vs *VariableStatement) StringAST(level int) string {
	var out strings.Builder

	out.WriteString(strings.Repeat("| ", level))
	out.WriteString("Statement: VariableStatement\n")
	out.WriteString(strings.Repeat("| ", level+1))
	out.WriteString("Token: " + vs.TokenLiteral() + "\n")

	out.WriteString(strings.Repeat("| ", level+1))
	out.WriteString("Expression(Type): Indentifier (")
	out.WriteString(vs.Type.String())
	out.WriteString(")\n")

	out.WriteString(strings.Repeat("| ", level+1))
	out.WriteString("Expression(Name): Identifier (")
	out.WriteString(vs.Name.String())
	out.WriteString(")\n")

	out.WriteString(strings.Repeat("| ", level+1))
	out.WriteString("Expression(Value):")

	if vs.Value != nil {
		out.WriteString("\n")
		out.WriteString(vs.Value.StringAST(level + 2))
	} else {
		out.WriteString(" <nil>\n")
	}

	return out.String()
}

type ExpressionStatement struct {
	Token      token.Token // the first token of the expression
	Expression Expression
}

func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}

	return ""
}
func (es *ExpressionStatement) StringAST(indent int) string {
	var out strings.Builder

	out.WriteString(strings.Repeat("| ", indent))
	out.WriteString("Statement: ExpressionStatement\n")

	if es.Expression != nil {
		out.WriteString(es.Expression.StringAST(indent + 1))
	}

	return out.String()
}

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteral) String() string       { return il.Token.Literal }
func (il *IntegerLiteral) StringAST(indent int) string {
	var out strings.Builder

	out.WriteString(strings.Repeat("| ", indent))
	out.WriteString("Expression: IntegerLiteral (")
	out.WriteString(il.Token.Literal)
	out.WriteString(")\n")

	return out.String()
}

type FloatLiteral struct {
	Token token.Token
	Value float64
}

func (fl *FloatLiteral) expressionNode()      {}
func (fl *FloatLiteral) TokenLiteral() string { return fl.Token.Literal }
func (fl *FloatLiteral) String() string       { return fl.Token.Literal }
func (fl *FloatLiteral) StringAST(indent int) string {
	var out strings.Builder

	out.WriteString(strings.Repeat("| ", indent))
	out.WriteString("Expression: FloatLiteral (")
	out.WriteString(fl.Token.Literal)
	out.WriteString(")\n")

	return out.String()
}

type VectorStatement struct {
	Token  token.Token  // The token.TYPE_INT, token.TYPE_FLOAT or token.TYPE_BOOL token
	Size   Expression   // The size of the vector, can be integer or expression
	Type   *Identifier  // The type of the vector (e.g., int, float, or bool)
	Name   *Identifier  // The variable name (e.g., x, y, or z)
	Values []Expression // The values assigned to the vector, can be nil or an array of expressions
}

func (vl *VectorStatement) statementNode()       {}
func (vs *VectorStatement) TokenLiteral() string { return vs.Token.Literal }
func (vs *VectorStatement) String() string {
	var out bytes.Buffer

	out.WriteString(vs.Type.String() + " ")
	out.WriteString(vs.Name.String())
	out.WriteString("[")

	if vs.Size != nil {
		out.WriteString(vs.Size.String())
	}

	out.WriteString("]")

	if len(vs.Values) > 0 {
		out.WriteString(" = {")
		for i, value := range vs.Values {
			out.WriteString(value.String())
			if i < len(vs.Values)-1 {
				out.WriteString(", ")
			}
		}
		out.WriteString("}")
	}

	out.WriteString(";")
	return out.String()
}
func (vs *VectorStatement) StringAST(indent int) string {
	var out strings.Builder

	out.WriteString(strings.Repeat("| ", indent))
	out.WriteString("Statement: VectorStatement\n")

	indent++

	out.WriteString(strings.Repeat("| ", indent))
	out.WriteString("Expression(Type): Indentifier (")
	out.WriteString(vs.Type.String())
	out.WriteString(")\n")

	out.WriteString(strings.Repeat("| ", indent))
	out.WriteString("Expression(Name): Identifier (")
	out.WriteString(vs.Name.String())
	out.WriteString(")\n")

	out.WriteString(strings.Repeat("| ", indent))
	out.WriteString(fmt.Sprintf("Expression(Size): %T (", vs.Size))
	if vs.Size != nil {
		out.WriteString(vs.Size.String())
	}
	out.WriteString(")\n")

	out.WriteString(strings.Repeat("| ", indent))
	out.WriteString("Expression(Values):\n")

	if len(vs.Values) > 0 {
		for _, value := range vs.Values {
			out.WriteString(value.StringAST(indent + 1))
		}
	}

	return out.String()
}

type Boolean struct {
	Token token.Token
	Value bool
}

func (b *Boolean) expressionNode()      {}
func (b *Boolean) TokenLiteral() string { return b.Token.Literal }
func (b *Boolean) String() string       { return b.Token.Literal }
func (b *Boolean) StringAST(indent int) string {
	var out strings.Builder

	out.WriteString(strings.Repeat("| ", indent))
	out.WriteString("Expression: Boolean (")
	out.WriteString(b.Token.Literal)
	out.WriteString(")\n")

	return out.String()
}

type PrefixExpression struct {
	Token    token.Token // the prefix token, e.g. !
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode()      {}
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }
func (pe *PrefixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")

	return out.String()
}
func (pe *PrefixExpression) StringAST(indent int) string {
	var out strings.Builder

	out.WriteString(strings.Repeat("| ", indent))
	out.WriteString("Expression: PrefixExpression (")
	out.WriteString(pe.Operator)
	out.WriteString(")\n")

	out.WriteString(pe.Right.StringAST(indent + 1))

	return out.String()
}

type InfixExpression struct {
	Token    token.Token // the operator token, e.g. +,-,/, *
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) expressionNode()      {}
func (ie *InfixExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *InfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(" " + ie.Operator + " ")
	out.WriteString(ie.Right.String())
	out.WriteString(")")

	return out.String()
}
func (ie *InfixExpression) StringAST(indent int) string {
	var out strings.Builder

	out.WriteString(strings.Repeat("| ", indent))
	out.WriteString("Expression: InfixExpression (")
	out.WriteString(ie.Operator)
	out.WriteString(")\n")

	out.WriteString(ie.Left.StringAST(indent + 1))
	out.WriteString(ie.Right.StringAST(indent + 1))

	return out.String()
}

type BlockStatement struct {
	Token      token.Token // the '{' token
	Statements []Statement
}

func (bs *BlockStatement) statementNode()       {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BlockStatement) String() string {
	var out bytes.Buffer

	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}
func (bs *BlockStatement) StringAST(indent int) string {
	var out strings.Builder

	out.WriteString(strings.Repeat("| ", indent))
	out.WriteString("Statement: BlockStatement\n")

	for _, s := range bs.Statements {
		out.WriteString(s.StringAST(indent + 1))
	}

	return out.String()
}

type IfExpression struct {
	Token       token.Token // the 'if' token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ie *IfExpression) expressionNode()      {}
func (ie *IfExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IfExpression) String() string {
	var out bytes.Buffer

	out.WriteString("if")
	out.WriteString(ie.Condition.String())
	out.WriteString(" ")
	out.WriteString(ie.Consequence.String())

	if ie.Alternative != nil {
		out.WriteString("else ")
		out.WriteString(ie.Alternative.String())
	}

	return out.String()
}
func (ie *IfExpression) StringAST(indent int) string {
	var out strings.Builder

	out.WriteString(strings.Repeat("| ", indent))
	out.WriteString("Expression: IfExpression\n")

	out.WriteString(strings.Repeat("| ", indent+1))
	out.WriteString("Condition:\n")
	out.WriteString(ie.Condition.StringAST(indent + 2))

	out.WriteString(strings.Repeat("| ", indent+1))
	out.WriteString("Consequence:\n")
	out.WriteString(ie.Consequence.StringAST(indent + 2))

	if ie.Alternative != nil {
		out.WriteString(strings.Repeat("| ", indent+1))
		out.WriteString("Alternative:\n")
		out.WriteString(ie.Alternative.StringAST(indent + 2))
	}

	return out.String()
}

type WhileStatement struct {
	Token     token.Token // the 'if' token
	Condition Expression
	Body      *BlockStatement
}

func (ws *WhileStatement) statementNode()       {}
func (ws *WhileStatement) TokenLiteral() string { return ws.Token.Literal }
func (ws *WhileStatement) String() string {
	var out bytes.Buffer

	out.WriteString("while")
	out.WriteString(" ")
	out.WriteString(ws.Condition.String())
	out.WriteString(" ")
	out.WriteString(ws.Body.String())

	return out.String()
}
func (ws *WhileStatement) StringAST(indent int) string {
	var out strings.Builder

	out.WriteString(strings.Repeat("| ", indent))
	out.WriteString("Statement: WhileStatement\n")

	out.WriteString(strings.Repeat("| ", indent+1))
	out.WriteString("Condition:\n")
	out.WriteString(ws.Condition.StringAST(indent + 2))

	out.WriteString(strings.Repeat("| ", indent+1))
	out.WriteString("Body:\n")
	out.WriteString(ws.Body.StringAST(indent + 2))

	return out.String()
}

type ForStatement struct {
	Token     token.Token // the 'if' token
	Var       Expression
	Start     Expression
	End       Expression
	Increment Expression
	Body      *BlockStatement
}

func (fs *ForStatement) statementNode()       {}
func (fs *ForStatement) TokenLiteral() string { return fs.Token.Literal }
func (fs *ForStatement) String() string {
	var out bytes.Buffer

	out.WriteString("for")
	out.WriteString(" (")
	out.WriteString(fs.Var.String())
	out.WriteString(", ")
	out.WriteString(fs.Start.String())
	out.WriteString(", ")
	out.WriteString(fs.End.String())
	out.WriteString(", ")
	out.WriteString(fs.Increment.String())
	out.WriteString(") ")
	out.WriteString(fs.Body.String())

	return out.String()
}
func (fs *ForStatement) StringAST(indent int) string {
	var out strings.Builder

	out.WriteString(strings.Repeat("| ", indent))
	out.WriteString("Statement: ForStatement\n")

	out.WriteString(strings.Repeat("| ", indent+1))
	out.WriteString("Var:\n")
	out.WriteString(fs.Start.StringAST(indent + 2))

	out.WriteString(strings.Repeat("| ", indent+1))
	out.WriteString("Initial:\n")
	out.WriteString(fs.Start.StringAST(indent + 2))

	out.WriteString(strings.Repeat("| ", indent+1))
	out.WriteString("End:\n")
	out.WriteString(fs.End.StringAST(indent + 2))

	out.WriteString(strings.Repeat("| ", indent+1))
	out.WriteString("Increment:\n")
	out.WriteString(fs.Increment.StringAST(indent + 2))

	out.WriteString(strings.Repeat("| ", indent+1))
	out.WriteString("Body:\n")
	out.WriteString(fs.Body.StringAST(indent + 2))

	return out.String()
}

type Attribute struct {
	Token    token.Token
	Name     *Identifier
	Type     *Identifier // Type of the parameter
	IsVector bool
	Size     Expression // Can be nil if the size is not specified
	Value    Expression // Can be nil if the value is not specified
}

func (a *Attribute) expressionNode()      {}
func (a *Attribute) TokenLiteral() string { return a.Token.Literal }
func (a *Attribute) String() string {
	var out bytes.Buffer

	out.WriteString(a.Type.String())
	out.WriteString(" ")
	out.WriteString(a.Name.String())
	if a.IsVector {
		out.WriteString("[")
		if a.Size != nil {
			out.WriteString(a.Size.String())
		}
		out.WriteString("]")
	}

	return out.String()
}
func (a *Attribute) StringAST(indent int) string {
	var out strings.Builder

	out.WriteString(strings.Repeat("| ", indent))
	out.WriteString("Expression: Attribute\n")

	out.WriteString(strings.Repeat("| ", indent+1))
	out.WriteString("Vector: " + strconv.FormatBool(a.IsVector) + "\n")

	out.WriteString(strings.Repeat("| ", indent+1))
	out.WriteString("Expression(Type): Identifier (" + a.Type.String() + ")\n")

	out.WriteString(strings.Repeat("| ", indent+1))
	out.WriteString("Expression(Name): Identifier (" + a.Name.String() + ")\n")

	if a.IsVector {
		if a.Size != nil {
			out.WriteString(strings.Repeat("| ", indent+1))
			out.WriteString("Expression(Size):\n")
			out.WriteString(a.Size.StringAST(indent + 2))
		}
	}

	if a.Value != nil {
		out.WriteString(strings.Repeat("| ", indent+1))
		out.WriteString("Expression(Value):\n")
		out.WriteString(a.Value.StringAST(indent + 2))
	}

	return out.String()
}

type ReturnType struct {
	Token    token.Token
	Type     *Identifier // Type of the parameter
	IsVector bool
	Size     Expression
}

func (rt *ReturnType) expressionNode()      {}
func (rt *ReturnType) TokenLiteral() string { return rt.Token.Literal }
func (rt *ReturnType) String() string {
	var out bytes.Buffer

	out.WriteString(rt.Type.String())
	if rt.IsVector {
		out.WriteString("[")
		if rt.Size != nil {
			out.WriteString(rt.Size.String())
		}
		out.WriteString("]")
	}

	return out.String()
}
func (rt *ReturnType) StringAST(indent int) string {
	var out strings.Builder

	out.WriteString(strings.Repeat("| ", indent+1))
	out.WriteString("Vector: " + strconv.FormatBool(rt.IsVector) + "\n")

	out.WriteString(strings.Repeat("| ", indent))
	out.WriteString("Expression: ReturnType\n")

	out.WriteString(strings.Repeat("| ", indent+1))
	out.WriteString("Expression(Type): Identifier (" + rt.Type.String() + ")\n")

	if rt.IsVector {
		if rt.Size != nil {
			out.WriteString(strings.Repeat("| ", indent+1))
			out.WriteString("Size:\n")
			out.WriteString(rt.Size.StringAST(indent + 2))
		}
	}

	return out.String()
}

type Parameter struct {
	Token    token.Token
	Name     *Identifier
	Type     *Identifier // Type of the parameter
	IsVector bool
	Size     Expression // Can be nil if the size is not specified
}

func (p *Parameter) expressionNode()      {}
func (p *Parameter) TokenLiteral() string { return p.Token.Literal }
func (p *Parameter) String() string {
	var out bytes.Buffer

	out.WriteString(p.Type.String())
	out.WriteString(" ")
	out.WriteString(p.Name.String())
	if p.IsVector {
		out.WriteString("[")
		if p.Size != nil {
			out.WriteString(p.Size.String())
		}
		out.WriteString("]")
	}

	return out.String()
}
func (p *Parameter) StringAST(indent int) string {
	var out bytes.Buffer

	out.WriteString(strings.Repeat("| ", indent))
	out.WriteString("Expression: Parameter\n")

	out.WriteString(strings.Repeat("| ", indent+1))
	out.WriteString("Vector: " + strconv.FormatBool(p.IsVector) + "\n")

	out.WriteString(strings.Repeat("| ", indent+1))
	out.WriteString("Expression(Type): Identifier (" + p.Type.String() + ")\n")

	out.WriteString(strings.Repeat("| ", indent+1))
	out.WriteString("Expression(Name): Identifier (" + p.Name.String() + ")\n")

	if p.Size != nil {
		out.WriteString(strings.Repeat("| ", indent+1))
		out.WriteString("Expression(Size):\n")
		out.WriteString(p.Size.StringAST(indent + 2))
	}

	return out.String()
}

type FunctionStatement struct {
	Token      token.Token     // The function name token
	Parameters []*Parameter    // The function parameters
	ReturnType *ReturnType     // The function return type
	Body       *BlockStatement // The function body
}

func (fs *FunctionStatement) statementNode()       {}
func (fs *FunctionStatement) TokenLiteral() string { return fs.Token.Literal }
func (fs *FunctionStatement) String() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range fs.Parameters {
		params = append(params, p.String())
	}

	out.WriteString(fs.TokenLiteral())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")
	out.WriteString(fs.ReturnType.String())
	out.WriteString(" ")
	out.WriteString(fs.Body.String())

	return out.String()
}
func (fs *FunctionStatement) StringAST(indent int) string {
	var out bytes.Buffer

	out.WriteString(strings.Repeat("| ", indent))
	out.WriteString("Statement: FunctionStatement\n")
	out.WriteString(strings.Repeat("| ", indent+1))
	out.WriteString("Token: " + fs.TokenLiteral() + "\n")
	out.WriteString(strings.Repeat("| ", indent+1))
	out.WriteString("Parameters:\n")
	for _, p := range fs.Parameters {
		out.WriteString(p.StringAST(indent + 2))
	}

	out.WriteString(strings.Repeat("| ", indent+1))
	out.WriteString("ReturnType:\n")
	out.WriteString(fs.ReturnType.StringAST(indent + 2))

	out.WriteString(strings.Repeat("| ", indent+1))
	out.WriteString("Body:\n")
	out.WriteString(fs.Body.StringAST(indent + 2))

	return out.String()
}

type CallExpression struct {
	Token     token.Token // the '(' token
	Function  Expression  // Identifier or FunctionLiteral
	Arguments []Expression
}

func (ce *CallExpression) expressionNode()      {}
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Literal }
func (ce *CallExpression) String() string {
	var out bytes.Buffer

	args := []string{}
	for _, a := range ce.Arguments {
		args = append(args, a.String())
	}

	out.WriteString(ce.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")

	return out.String()
}
func (ce *CallExpression) StringAST(indent int) string {
	var out bytes.Buffer

	out.WriteString(strings.Repeat("| ", indent))
	out.WriteString("Expression: CallExpression\n")
	out.WriteString(strings.Repeat("| ", indent+1))
	out.WriteString("Function:\n")
	out.WriteString(ce.Function.StringAST(indent + 2))

	out.WriteString(strings.Repeat("| ", indent+1))
	out.WriteString("Arguments:\n")
	for _, a := range ce.Arguments {
		out.WriteString(a.StringAST(indent + 2))
	}

	return out.String()
}

type StringLiteral struct {
	Token token.Token
	Value string
}

func (sl *StringLiteral) expressionNode()      {}
func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Literal }
func (sl *StringLiteral) String() string       { return sl.Token.Literal }
func (sl *StringLiteral) StringAST(indent int) string {
	var out bytes.Buffer

	out.WriteString(strings.Repeat("| ", indent))
	out.WriteString("Expression: StringLiteral\n")

	out.WriteString(strings.Repeat("| ", indent+1))
	out.WriteString("Token: " + sl.TokenLiteral() + "\n")

	out.WriteString(strings.Repeat("| ", indent+1))
	out.WriteString("Value: " + sl.Value + "\n")

	return out.String()
}

// structs {
//   <struct_name> { <struct_definition> };
//	 <struct_name> { <struct_definition> };
//   ...
// }

type StructsStatement struct {
	Token   token.Token
	Structs []*Struct
}

func (sd *StructsStatement) statementNode()       {}
func (sd *StructsStatement) TokenLiteral() string { return sd.Token.Literal }
func (sd *StructsStatement) String() string {
	var out bytes.Buffer

	out.WriteString("structs { ")
	for _, str := range sd.Structs {
		out.WriteString(str.String())
	}
	out.WriteString(" }")

	return out.String()
}
func (sd *StructsStatement) StringAST(indent int) string {
	var out bytes.Buffer

	out.WriteString(strings.Repeat("| ", indent))
	out.WriteString("Statement: StructsStatement\n")

	out.WriteString(strings.Repeat("| ", indent+1))
	out.WriteString("Token: " + sd.TokenLiteral() + "\n")

	out.WriteString(strings.Repeat("| ", indent+1))
	out.WriteString("Structs:\n")

	for _, str := range sd.Structs {
		out.WriteString(str.StringAST(indent + 2))
	}

	return out.String()
}

type Struct struct {
	Token      token.Token
	Attributes []*Attribute
}

func (s *Struct) expressionNode()      {}
func (s *Struct) TokenLiteral() string { return s.Token.Literal }
func (s *Struct) String() string {
	var out bytes.Buffer

	out.WriteString(s.TokenLiteral())
	out.WriteString(" { ")
	for i, attr := range s.Attributes {
		out.WriteString(attr.String())
		if i < len(s.Attributes)-1 {
			out.WriteString(", ")
		} else {
			out.WriteString("; };")
		}
	}

	return out.String()
}
func (s *Struct) StringAST(indent int) string {
	var out bytes.Buffer

	out.WriteString(strings.Repeat("| ", indent))
	out.WriteString("Expression: Struct\n")

	out.WriteString(strings.Repeat("| ", indent+1))
	out.WriteString("Token: " + s.TokenLiteral() + "\n")

	out.WriteString(strings.Repeat("| ", indent+1))
	out.WriteString("Attributes:\n")

	for _, attr := range s.Attributes {
		out.WriteString(attr.StringAST(indent + 2))
	}

	return out.String()
}

type GlobalStatement struct {
	Token token.Token
	Body  *BlockStatement
}

func (gs *GlobalStatement) statementNode()       {}
func (gs *GlobalStatement) TokenLiteral() string { return gs.Token.Literal }
func (gs *GlobalStatement) String() string {
	var out bytes.Buffer

	out.WriteString("global { ")
	for _, v := range gs.Body.Statements {
		out.WriteString(v.String())
	}
	out.WriteString(" }")

	return out.String()
}
func (gs *GlobalStatement) StringAST(indent int) string {
	var out bytes.Buffer

	out.WriteString(strings.Repeat("| ", indent))
	out.WriteString("Statement: GlobalStatement\n")

	out.WriteString(strings.Repeat("| ", indent+1))
	out.WriteString("Token: " + gs.TokenLiteral() + "\n")

	out.WriteString(strings.Repeat("| ", indent+1))
	out.WriteString("Body:\n")
	for _, v := range gs.Body.Statements {
		out.WriteString(v.StringAST(indent + 2))
	}

	return out.String()
}

type ConstStatement struct {
	Token token.Token
	Body  *BlockStatement
}

func (cs *ConstStatement) statementNode()       {}
func (cs *ConstStatement) TokenLiteral() string { return cs.Token.Literal }
func (cs *ConstStatement) String() string {
	var out bytes.Buffer

	out.WriteString("const { ")
	for _, v := range cs.Body.Statements {
		out.WriteString(v.String())
	}
	out.WriteString(" }")

	return out.String()
}
func (cs *ConstStatement) StringAST(indent int) string {
	var out bytes.Buffer

	out.WriteString(strings.Repeat("| ", indent))
	out.WriteString("Statement: ConstStatement\n")

	out.WriteString(strings.Repeat("| ", indent+1))
	out.WriteString("Token: " + cs.TokenLiteral() + "\n")

	out.WriteString(strings.Repeat("| ", indent+1))
	out.WriteString("Body:\n")

	for _, v := range cs.Body.Statements {
		out.WriteString(v.StringAST(indent + 2))
	}

	return out.String()
}

type LocalStatement struct {
	Token token.Token
	Body  *BlockStatement
}

func (ls *LocalStatement) statementNode()       {}
func (ls *LocalStatement) TokenLiteral() string { return ls.Token.Literal }
func (ls *LocalStatement) String() string {
	var out bytes.Buffer

	out.WriteString("local { ")
	for _, v := range ls.Body.Statements {
		out.WriteString(v.String())
	}
	out.WriteString(" }")

	return out.String()
}
func (ls *LocalStatement) StringAST(indent int) string {
	var out bytes.Buffer

	out.WriteString(strings.Repeat("| ", indent))
	out.WriteString("Statement: LocalStatement\n")

	out.WriteString(strings.Repeat("| ", indent+1))
	out.WriteString("Token: " + ls.TokenLiteral() + "\n")

	out.WriteString(strings.Repeat("| ", indent+1))
	out.WriteString("Body:\n")

	for _, v := range ls.Body.Statements {
		out.WriteString(v.StringAST(indent + 2))
	}

	return out.String()
}

type IndexExpression struct {
	Token token.Token
	Left  Expression // The object being accessed
	Index Expression // The index being accessed
}

func (ie *IndexExpression) expressionNode()      {}
func (ie *IndexExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IndexExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString("[")
	out.WriteString(ie.Index.String())
	out.WriteString("])")

	return out.String()
}
func (ie *IndexExpression) StringAST(indent int) string {
	var out bytes.Buffer

	out.WriteString(strings.Repeat("| ", indent))
	out.WriteString("Expression: IndexExpression\n")

	out.WriteString(strings.Repeat("| ", indent+1))
	out.WriteString("Token: " + ie.TokenLiteral() + "\n")

	out.WriteString(strings.Repeat("| ", indent+1))
	out.WriteString("Expression(Left):\n")
	out.WriteString(ie.Left.StringAST(indent + 2))

	out.WriteString(strings.Repeat("| ", indent+1))
	out.WriteString("Expression(Index):\n")

	out.WriteString(ie.Index.StringAST(indent + 2))

	return out.String()
}

type AccessorExpression struct {
	Token token.Token
	Left  Expression   // The object being accessed
	Index []Expression // The index being accessed
}

func (ae *AccessorExpression) expressionNode()      {}
func (ae *AccessorExpression) TokenLiteral() string { return ae.Token.Literal }
func (ae *AccessorExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ae.Left.String())
	out.WriteString(".")
	out.WriteString(ae.Index[0].String())
	for _, i := range ae.Index[1:] {
		out.WriteString(".")
		out.WriteString(i.String())
	}
	out.WriteString(")")

	return out.String()
}
func (ae *AccessorExpression) StringAST(indent int) string {
	var out bytes.Buffer

	out.WriteString(strings.Repeat("| ", indent))
	out.WriteString("Expression: AccessorExpression\n")
	out.WriteString(strings.Repeat("| ", indent+1))
	out.WriteString("Token: " + ae.TokenLiteral() + "\n")
	out.WriteString(strings.Repeat("| ", indent+1))
	out.WriteString("Expression(Left):\n")
	out.WriteString(ae.Left.StringAST(indent + 2))
	out.WriteString(strings.Repeat("| ", indent+1))
	out.WriteString("Expression(Index):\n")
	for _, i := range ae.Index {
		out.WriteString(i.StringAST(indent + 2))
	}

	return out.String()
}

type AssignmentStatement struct {
	Token token.Token
	Left  Expression
	Value Expression
}

func (ae *AssignmentStatement) statementNode()       {}
func (ae *AssignmentStatement) TokenLiteral() string { return ae.Token.Literal }
func (ae *AssignmentStatement) String() string {
	var out bytes.Buffer

	out.WriteString(ae.Left.String())
	out.WriteString(" = ")
	out.WriteString(ae.Value.String())
	out.WriteString(";")

	return out.String()
}
func (ae *AssignmentStatement) StringAST(indent int) string {
	var out bytes.Buffer

	out.WriteString(strings.Repeat("| ", indent))
	out.WriteString("Statement: AssignmentStatement\n")

	out.WriteString(strings.Repeat("| ", indent+1))
	out.WriteString("Token: " + ae.TokenLiteral() + "\n")

	out.WriteString(strings.Repeat("| ", indent+1))
	out.WriteString("Expression(Left):\n")

	out.WriteString(ae.Left.StringAST(indent + 2))

	out.WriteString(strings.Repeat("| ", indent+1))
	out.WriteString("Expression(Value):\n")

	out.WriteString(ae.Value.StringAST(indent + 2))

	return out.String()
}

type IncrementStatement struct {
	Token token.Token
	Var   Expression
}

func (is *IncrementStatement) statementNode()       {}
func (is *IncrementStatement) TokenLiteral() string { return is.Token.Literal }
func (is *IncrementStatement) String() string {
	var out bytes.Buffer

	out.WriteString(is.Var.String())
	out.WriteString("++")
	out.WriteString(";")

	return out.String()
}
func (is *IncrementStatement) StringAST(indent int) string {
	var out bytes.Buffer

	out.WriteString(strings.Repeat("| ", indent))
	out.WriteString("Statement: IncrementStatement\n")
	out.WriteString(strings.Repeat("| ", indent+1))
	out.WriteString("Token: " + is.TokenLiteral() + "\n")
	out.WriteString(strings.Repeat("| ", indent+1))
	out.WriteString("Expression(Var):\n")
	out.WriteString(is.Var.StringAST(indent + 2))

	return out.String()
}

type DecrementStatement struct {
	Token token.Token
	Var   Expression
}

func (ds *DecrementStatement) statementNode()       {}
func (ds *DecrementStatement) TokenLiteral() string { return ds.Token.Literal }
func (ds *DecrementStatement) String() string {
	var out bytes.Buffer

	out.WriteString(ds.Var.String())
	out.WriteString("--")
	out.WriteString(";")

	return out.String()
}
func (ds *DecrementStatement) StringAST(indent int) string {
	var out bytes.Buffer

	out.WriteString(strings.Repeat("| ", indent))
	out.WriteString("Statement: DecrementStatement\n")
	out.WriteString(strings.Repeat("| ", indent+1))
	out.WriteString("Token: " + ds.TokenLiteral() + "\n")
	out.WriteString(strings.Repeat("| ", indent+1))
	out.WriteString("Expression(Var):\n")
	out.WriteString(ds.Var.StringAST(indent + 2))

	return out.String()
}

type PlusEqualsStatement struct {
	Token    token.Token
	Var      Expression
	Quantity Expression
}

func (ps *PlusEqualsStatement) statementNode()       {}
func (ps *PlusEqualsStatement) TokenLiteral() string { return ps.Token.Literal }
func (ps *PlusEqualsStatement) String() string {
	var out bytes.Buffer

	out.WriteString(ps.Var.String())
	out.WriteString(" += ")
	out.WriteString(ps.Quantity.String())
	out.WriteString(";")

	return out.String()
}
func (ps *PlusEqualsStatement) StringAST(indent int) string {
	var out bytes.Buffer

	out.WriteString(strings.Repeat("| ", indent))
	out.WriteString("Statement: PlusEqualsStatement\n")
	out.WriteString(strings.Repeat("| ", indent+1))
	out.WriteString("Token: " + ps.TokenLiteral() + "\n")
	out.WriteString(strings.Repeat("| ", indent+1))
	out.WriteString("Expression(Var):\n")
	out.WriteString(ps.Var.StringAST(indent + 2))
	out.WriteString(strings.Repeat("| ", indent+1))
	out.WriteString("Expression(Quantity):\n")
	out.WriteString(ps.Quantity.StringAST(indent + 2))

	return out.String()
}

type MinusEqualsStatement struct {
	Token    token.Token
	Var      Expression
	Quantity Expression
}

func (ms *MinusEqualsStatement) statementNode()       {}
func (ms *MinusEqualsStatement) TokenLiteral() string { return ms.Token.Literal }
func (ms *MinusEqualsStatement) String() string {
	var out bytes.Buffer

	out.WriteString(ms.Var.String())
	out.WriteString(" -= ")
	out.WriteString(ms.Quantity.String())
	out.WriteString(";")

	return out.String()
}
func (ms *MinusEqualsStatement) StringAST(indent int) string {
	var out bytes.Buffer

	out.WriteString(strings.Repeat("| ", indent))
	out.WriteString("Statement: MinusEqualsStatement\n")
	out.WriteString(strings.Repeat("| ", indent+1))
	out.WriteString("Token: " + ms.TokenLiteral() + "\n")
	out.WriteString(strings.Repeat("| ", indent+1))
	out.WriteString("Expression(Var):\n")
	out.WriteString(ms.Var.StringAST(indent + 2))
	out.WriteString(strings.Repeat("| ", indent+1))
	out.WriteString("Expression(Quantity):\n")
	out.WriteString(ms.Quantity.StringAST(indent + 2))

	return out.String()
}
