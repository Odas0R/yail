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
			&VectorStatement{
				Token: token.Token{Type: token.IDENT, Literal: "int"},
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

func TestStructDefinitionString(t *testing.T) {
	program := &Program{
		Statements: []Statement{
			&StructsStatement{
				Token: token.Token{Type: token.STRUCTS, Literal: "structs"},
				Structs: []*Struct{
					{
						Token: token.Token{Type: token.IDENT, Literal: "point2D"},
						Attributes: []*Attribute{
							{
								Token: token.Token{Type: token.IDENT, Literal: "float"},
								Name: &Identifier{
									Token: token.Token{Type: token.IDENT, Literal: "x"},
									Value: "x",
								},
							},
							{
								Token: token.Token{Type: token.IDENT, Literal: "float"},
								Name: &Identifier{
									Token: token.Token{Type: token.IDENT, Literal: "y"},
									Value: "y",
								},
							},
						},
					},
					{
						Token: token.Token{Type: token.IDENT, Literal: "point3D"},
						Attributes: []*Attribute{
							{
								Token: token.Token{Type: token.IDENT, Literal: "float"},
								Name: &Identifier{
									Token: token.Token{Type: token.IDENT, Literal: "x"},
									Value: "x",
								},
							},
							{
								Token: token.Token{Type: token.IDENT, Literal: "float"},
								Name: &Identifier{
									Token: token.Token{Type: token.IDENT, Literal: "y"},
									Value: "y",
								},
							},
							{
								Token: token.Token{Type: token.IDENT, Literal: "float"},
								Name: &Identifier{
									Token: token.Token{Type: token.IDENT, Literal: "z"},
									Value: "z",
								},
							},
						},
					},
					{
						Token: token.Token{Type: token.IDENT, Literal: "pointND"},
						Attributes: []*Attribute{
							{
								Token: token.Token{Type: token.IDENT, Literal: "float"},
								Name: &Identifier{
									Token: token.Token{Type: token.IDENT, Literal: "x"},
									Value: "x",
								},
								IsVector: true,
							},
						},
					},
					{
						Token: token.Token{Type: token.IDENT, Literal: "pointNDSize"},
						Attributes: []*Attribute{
							{
								Token: token.Token{Type: token.IDENT, Literal: "float"},
								Name: &Identifier{
									Token: token.Token{Type: token.IDENT, Literal: "x"},
									Value: "x",
								},
								IsVector: true,
								Size: &IntegerLiteral{
									Token: token.Token{Type: token.INT, Literal: "5"},
									Value: 5,
								},
							},
						},
					},
				},
			},
		},
	}

	want := `structs {
	point2D { float x, float y; };
	point3D { float x, float y, float z; };
	pointND { float x[]; };
	pointNDSize { float x[5]; };
}`

	if program.String() != want {
		t.Errorf("program.String() wrong. got=%q", program.String())
	}
}
