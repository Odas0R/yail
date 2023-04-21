package ast

import (
	"testing"

	"github.com/odas0r/yail/token"
)

func TestVarDeclarationString(t *testing.T) {
	program := &Program{
		Statements: []Statement{
			&VarDeclaration{
				Token: token.Token{Type: token.INT, Literal: "int"},
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

func TestVectorDeclarationString(t *testing.T) {
	program := &Program{
		Statements: []Statement{
			&VectorDeclaration{
				Token: token.Token{Type: token.INT, Literal: "int"},
				Name: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "myVec"},
					Value: "myVec",
				},
				Size: &IntegerLiteral{
					Token: token.Token{Type: token.INT, Literal: "5"},
					Value: 5,
				},
				Values: []Expression{
					&IntegerLiteral{
						Token: token.Token{Type: token.INT, Literal: "1"},
						Value: 1,
					},
					&IntegerLiteral{
						Token: token.Token{Type: token.INT, Literal: "2"},
						Value: 2,
					},
					&IntegerLiteral{
						Token: token.Token{Type: token.INT, Literal: "3"},
						Value: 3,
					},
					&IntegerLiteral{
						Token: token.Token{Type: token.INT, Literal: "4"},
						Value: 4,
					},
					&IntegerLiteral{
						Token: token.Token{Type: token.INT, Literal: "5"},
						Value: 5,
					},
				},
			},
		},
	}

	if program.String() != "int myVec[5] = {1, 2, 3, 4, 5};" {
		t.Errorf("program.String() wrong. got=%q", program.String())
	}
}
