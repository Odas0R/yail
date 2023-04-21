package token

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// Identifiers + literals
	IDENT  = "IDENT" // add, foobar, x, y, ...
	INT    = "INT"   // 1343456
	FLOAT  = "FLOAT" // 1.343456
	BOOL   = "BOOL"  // bool z
	STRING = "STRING"

	TYPE_INT   = "_INT"
	TYPE_FLOAT = "_FLOAT"
	TYPE_BOOL  = "_BOOL"

	STRUCTS  = "STRUCTS"
	ACCESSOR = "."

	// Operators
	ASSIGN = "="
	PLUS   = "+"

	MINUS    = "-"
	BANG     = "!"
	ASTERISK = "*"
	SLASH    = "/"

	LT = "<"
	GT = ">"

	EQ     = "=="
	NOT_EQ = "!="

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
	TRUE  = "TRUE"
	FALSE = "FALSE"
	IF    = "IF"
	ELSE  = "ELSE"
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

	// Types
	"int":   TYPE_INT,
	"float": TYPE_FLOAT,
	"bool":  TYPE_BOOL,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
