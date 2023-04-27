package ast

import (
	"testing"

	"github.com/odas0r/yail/token"
)

func TestVarDeclarationString(t *testing.T) {
	program := &Program{
		Statements: []Statement{
			&VariableStatement{
				Token: token.Token{Type: token.IDENT, Literal: "int"},
				Type: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "int"},
					Value: "int",
				},
				Name: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "myVar"},
					Value: "myVar",
				},
				Value: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "anotherVar"},
					Value: "anotherVar",
				},
			},
		},
	}

	if program.String() != "int myVar = anotherVar;" {
		t.Errorf("program.String() wrong. got=%q", program.String())
	}
}
