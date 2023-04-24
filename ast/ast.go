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

type VarStatement struct {
	Token token.Token // The token.TYPE_INT, token.TYPE_FLOAT, or token.TYPE_BOOL token
	Name  *Identifier // The variable name (e.g., x, y, or z)
	Value Expression  // The value assigned to the variable, can be nil
}

func (vs *VarStatement) statementNode()       {}
func (vs *VarStatement) TokenLiteral() string { return vs.Token.Literal }
func (vs *VarStatement) String() string {
	var out bytes.Buffer

	out.WriteString(vs.TokenLiteral())
	out.WriteString(" ")
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
	Name   *Identifier  // The variable name (e.g., x, y, or z)
	Values []Expression // The values assigned to the vector, can be nil or an array of expressions
}

func (vl *VectorStatement) statementNode()       {}
func (vs *VectorStatement) TokenLiteral() string { return vs.Token.Literal }
func (vs *VectorStatement) String() string {
	var out bytes.Buffer

	out.WriteString(vs.TokenLiteral() + " ")
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

type Attribute struct {
	Token    token.Token
	Name     *Identifier
	IsVector bool
	Size     Expression // Can be nil if the size is not specified
}

func (p *Attribute) expressionNode()      {}
func (p *Attribute) TokenLiteral() string { return p.Token.Literal }
func (p *Attribute) String() string {
	var out bytes.Buffer

	out.WriteString(p.TokenLiteral())
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

type ReturnType struct {
	Token    token.Token
	IsVector bool
	Size     Expression
}

func (ti *ReturnType) expressionNode()      {}
func (ti *ReturnType) TokenLiteral() string { return ti.Token.Literal }
func (ti *ReturnType) String() string {
	var out bytes.Buffer

	out.WriteString(ti.TokenLiteral())
	if ti.IsVector {
		out.WriteString("[")
		if ti.Size != nil {
			out.WriteString(ti.Size.String())
		}
		out.WriteString("]")
	}

	return out.String()
}

type FunctionLiteral struct {
	Token      token.Token     // The function name token
	Parameters []*Attribute    // The function parameters
	ReturnType *ReturnType     // The function return type
	Body       *BlockStatement // The function body
}

func (fl *FunctionLiteral) expressionNode()      {}
func (fl *FunctionLiteral) TokenLiteral() string { return fl.Token.Literal }
func (fl *FunctionLiteral) String() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range fl.Parameters {
		params = append(params, p.String())
	}

	out.WriteString(fl.TokenLiteral())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")
	out.WriteString(fl.ReturnType.String())
	out.WriteString(fl.Body.String())

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

type StructsDefinition struct {
	Token   token.Token
	Structs []*StructLiteral
}

func (sd *StructsDefinition) statementNode()       {}
func (sd *StructsDefinition) TokenLiteral() string { return sd.Token.Literal }
func (sd *StructsDefinition) String() string {
	var out bytes.Buffer

	out.WriteString("structs {")
	for _, str := range sd.Structs {
		out.WriteString("\n\t")
		out.WriteString(str.String())
	}
	out.WriteString("\n}")

	return out.String()
}

type StructLiteral struct {
	Token      token.Token
	Attributes []*Attribute
}

func (sl *StructLiteral) expressionNode()      {}
func (sl *StructLiteral) TokenLiteral() string { return sl.Token.Literal }
func (sl *StructLiteral) String() string {
	var out bytes.Buffer

	out.WriteString(sl.TokenLiteral())
	out.WriteString(" { ")
	for i, attr := range sl.Attributes {
		out.WriteString(attr.String())
		if i < len(sl.Attributes)-1 {
			out.WriteString(", ")
		} else {
			out.WriteString("; };")
		}
	}

	return out.String()
}
