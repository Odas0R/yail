package ast

import (
	"bytes"
	"strings"

	"github.com/odas0r/yail/token"
)

type Node interface {
	TokenLiteral() string
	String() string
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

type Identifier struct {
	Token token.Token // the token.IDENT token
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string       { return i.Value }

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

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteral) String() string       { return il.Token.Literal }

type FloatLiteral struct {
	Token token.Token
	Value float64
}

func (fl *FloatLiteral) expressionNode()      {}
func (fl *FloatLiteral) TokenLiteral() string { return fl.Token.Literal }
func (fl *FloatLiteral) String() string       { return fl.Token.Literal }

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

type Boolean struct {
	Token token.Token
	Value bool
}

func (b *Boolean) expressionNode()      {}
func (b *Boolean) TokenLiteral() string { return b.Token.Literal }
func (b *Boolean) String() string       { return b.Token.Literal }

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

type StringLiteral struct {
	Token token.Token
	Value string
}

func (sl *StringLiteral) expressionNode()      {}
func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Literal }
func (sl *StringLiteral) String() string       { return sl.Token.Literal }

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

	out.WriteString("structs {")
	for _, str := range sd.Structs {
		out.WriteString("\n\t")
		out.WriteString(str.String())
	}
	out.WriteString("\n}")

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

type GlobalStatement struct {
	Token token.Token
	Body  *BlockStatement
}

func (gs *GlobalStatement) statementNode()       {}
func (gs *GlobalStatement) TokenLiteral() string { return gs.Token.Literal }
func (gs *GlobalStatement) String() string {
	var out bytes.Buffer

	out.WriteString("global {")
	for _, v := range gs.Body.Statements {
		out.WriteString("\n\t")
		out.WriteString(v.String())
	}
	out.WriteString("\n}")

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

	out.WriteString("const {")
	for _, v := range cs.Body.Statements {
		out.WriteString("\n\t")
		out.WriteString(v.String())
	}
	out.WriteString("\n}")

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

	out.WriteString("local {")
	for _, v := range ls.Body.Statements {
		out.WriteString("\n\t")
		out.WriteString(v.String())
	}
	out.WriteString("\n}")

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
