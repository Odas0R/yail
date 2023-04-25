package lexer

import (
	"testing"

	"github.com/odas0r/yail/token"
)

func TestNextToken(t *testing.T) {
	input := `int five = 5;
int ten = 10;
float pi = 3.14;
bool z = true;

add(int x, int y) {
  add = x + y;
};

result = add(five, ten);

int x[1] = {1};

!-/*5;

5 < 10 > 5;

if (5 < 10) {
	true;
} else {
	false;
}

10 == 10;
10 != 9;

"foobar"
"foo bar"

global { }
local { }
const { }

structs {
	point2D { float x, y; };
	pointND { float x[5]; };
}
`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.IDENT, "int"},
		{token.IDENT, "five"},
		{token.ASSIGN, "="},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},

		{token.IDENT, "int"},
		{token.IDENT, "ten"},
		{token.ASSIGN, "="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},

		{token.IDENT, "float"},
		{token.IDENT, "pi"},
		{token.ASSIGN, "="},
		{token.FLOAT, "3.14"},
		{token.SEMICOLON, ";"},

		{token.IDENT, "bool"},
		{token.IDENT, "z"},
		{token.ASSIGN, "="},
		{token.TRUE, "true"},
		{token.SEMICOLON, ";"},

		{token.IDENT, "add"},
		{token.LPAREN, "("},
		{token.IDENT, "int"},
		{token.IDENT, "x"},
		{token.COMMA, ","},
		{token.IDENT, "int"},
		{token.IDENT, "y"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.IDENT, "add"},
		{token.ASSIGN, "="},
		{token.IDENT, "x"},
		{token.PLUS, "+"},
		{token.IDENT, "y"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.SEMICOLON, ";"},

		{token.IDENT, "result"},
		{token.ASSIGN, "="},
		{token.IDENT, "add"},
		{token.LPAREN, "("},
		{token.IDENT, "five"},
		{token.COMMA, ","},
		{token.IDENT, "ten"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},

		{token.IDENT, "int"},
		{token.IDENT, "x"},
		{token.LBRACKET, "["},
		{token.INT, "1"},
		{token.RBRACKET, "]"},
		{token.ASSIGN, "="},
		{token.LBRACE, "{"},
		{token.INT, "1"},
		{token.RBRACE, "}"},
		{token.SEMICOLON, ";"},

		{token.BANG, "!"},
		{token.MINUS, "-"},
		{token.SLASH, "/"},
		{token.ASTERISK, "*"},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},

		{token.INT, "5"},
		{token.LT, "<"},
		{token.INT, "10"},
		{token.GT, ">"},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},

		{token.IF, "if"},
		{token.LPAREN, "("},
		{token.INT, "5"},
		{token.LT, "<"},
		{token.INT, "10"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.TRUE, "true"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.ELSE, "else"},
		{token.LBRACE, "{"},
		{token.FALSE, "false"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},

		{token.INT, "10"},
		{token.EQ, "=="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},

		{token.INT, "10"},
		{token.NOT_EQ, "!="},
		{token.INT, "9"},
		{token.SEMICOLON, ";"},

		{token.STRING, "foobar"},
		{token.STRING, "foo bar"},

		{token.GLOBAL, "global"},
		{token.LBRACE, "{"},
		{token.RBRACE, "}"},

		{token.LOCAL, "local"},
		{token.LBRACE, "{"},
		{token.RBRACE, "}"},

		{token.CONST, "const"},
		{token.LBRACE, "{"},
		{token.RBRACE, "}"},

		{token.STRUCTS, "structs"},
		{token.LBRACE, "{"},
		{token.IDENT, "point2D"},
		{token.LBRACE, "{"},
		{token.IDENT, "float"},
		{token.IDENT, "x"},
		{token.COMMA, ","},
		{token.IDENT, "y"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.SEMICOLON, ";"},
		{token.IDENT, "pointND"},
		{token.LBRACE, "{"},
		{token.IDENT, "float"},
		{token.IDENT, "x"},
		{token.LBRACKET, "["},
		{token.INT, "5"},
		{token.RBRACKET, "]"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},

		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		// t.Log(tok.Type, tok.Literal)

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q", i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
		}
	}
}
