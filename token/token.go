package token

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"
	COMMENT = "#"

	// Identifiers + literals
	IDENT  = "IDENT" // add, foobar, x, y, int, float, bool, string
	INT    = "INT"   // 1343456
	FLOAT  = "FLOAT" // 1.343456
	BOOL   = "BOOL"  // bool z
	STRING = "STRING"

	ACCESSOR = "."

	// Operators
	ASSIGN = "="
	PLUS   = "+"

	MINUS    = "-"
	BANG     = "!"
	ASTERISK = "*"
	SLASH    = "/"

	AND = "AND"
	OR  = "OR"
	LT  = "<"
	LTE = "<="
	GT  = ">"
	GTE = ">="

	EQ     = "=="
	NOT_EQ = "!="

	INCREMENT = "++"
	DECREMENT = "--"
	PLUS_EQ   = "+="
	MINUS_EQ  = "-="

	// Delimiters
	COMMA     = ","
	SEMICOLON = ";"

	LPAREN   = "("
	RPAREN   = ")"
	LBRACE   = "{"
	RBRACE   = "}"
	LBRACKET = "["
	RBRACKET = "]"

	// Keywords
	TRUE    = "TRUE"
	FALSE   = "FALSE"
	IF      = "IF"
	ELSE    = "ELSE"
	STRUCTS = "STRUCTS"
	GLOBAL  = "GLOBAL"
	LOCAL   = "LOCAL"
	CONST   = "CONST"
	WHILE   = "WHILE"
	FOR		 = "FOR"
)

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

var keywords = map[string]TokenType{
	"true":    TRUE,
	"false":   FALSE,
	"if":      IF,
	"else":    ELSE,
	"structs": STRUCTS,
	"const":   CONST,
	"global":  GLOBAL,
	"local":   LOCAL,
	"while":   WHILE,
	"and":     AND,
	"or":      OR,
	"for":     FOR,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
