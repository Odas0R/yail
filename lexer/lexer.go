package lexer

import (
	"strings"

	"github.com/odas0r/yail/token"
)

type Lexer struct {
	input        string
	Position     int // current position in input (points to current char)
	ReadPosition int // current reading position in input (after current char)

	// current char under examination; In order to support UTF-8, Unicode we need
	// to use a rune instead of a byte.
	Ch   byte
	Line int // current line number
}

func New(input string) *Lexer {
	l := &Lexer{input: input, Line: 1}
	l.readChar()
	return l
}

// readChar() reads the next character from the input and advances the position
// and readPosition by one. It also sets the ch field to the character read.
//
// Example:
//
//	input := "let five = 5;"
//	l := New(input)
//	l.readChar()
//	l.ch == 'l'
//	l.readChar()
//	l.ch == 'e'
//	...
func (l *Lexer) readChar() {
	// Check if we've reached the end of the input text
	if l.ReadPosition >= len(l.input) {
		// If so, set the character to NUL (ASCII code 0) to indicate EOF
		l.Ch = 0
	} else {
		// Otherwise, read the next character from the input text
		l.Ch = l.input[l.ReadPosition]
	}

	// Increment the line number if the current character is a newline
	if l.Ch == '\n' {
		l.Line++
	}

	// Update the Lexer's position and readPosition fields to reflect the new character
	l.Position = l.ReadPosition
	l.ReadPosition += 1
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhiteSpace()

	switch l.Ch {
	case '#':
		l.readComment()
		return l.NextToken()
	case '=':
		if l.peekChar() == '=' {
			ch := l.Ch
			l.readChar()
			literal := string(ch) + string(l.Ch)
			tok = token.Token{Type: token.EQ, Literal: literal}
		} else {
			tok = newToken(token.ASSIGN, l.Ch)
		}
	case '!':
		if l.peekChar() == '=' {
			ch := l.Ch
			l.readChar()
			literal := string(ch) + string(l.Ch)
			tok = token.Token{Type: token.NOT_EQ, Literal: literal}
		} else {
			tok = newToken(token.BANG, l.Ch)
		}
	case '"':
		tok.Type = token.STRING
		tok.Literal = l.readString()
	case '(':
		tok = newToken(token.LPAREN, l.Ch)
	case ')':
		tok = newToken(token.RPAREN, l.Ch)
	case '.':
		tok = newToken(token.ACCESSOR, l.Ch)
	case '+':
		if l.peekChar() == '+' {
			ch := l.Ch
			l.readChar()
			literal := string(ch) + string(l.Ch)
			tok = token.Token{Type: token.INCREMENT, Literal: literal}
		} else if l.peekChar() == '=' {
			ch := l.Ch
			l.readChar()
			literal := string(ch) + string(l.Ch)
			tok = token.Token{Type: token.PLUS_EQ, Literal: literal}
		} else {
			tok = newToken(token.PLUS, l.Ch)
		}
	case '-':
		if l.peekChar() == '-' {
			ch := l.Ch
			l.readChar()
			literal := string(ch) + string(l.Ch)
			tok = token.Token{Type: token.DECREMENT, Literal: literal}
		} else if l.peekChar() == '=' {
			ch := l.Ch
			l.readChar()
			literal := string(ch) + string(l.Ch)
			tok = token.Token{Type: token.MINUS_EQ, Literal: literal}
		} else {
			tok = newToken(token.MINUS, l.Ch)
		}
	case '/':
		tok = newToken(token.SLASH, l.Ch)
	case '*':
		if l.peekChar() == '=' {
			ch := l.Ch
			l.readChar()
			literal := string(ch) + string(l.Ch)
			tok = token.Token{Type: token.MULT_EQ, Literal: literal}
		} else {
			tok = newToken(token.ASTERISK, l.Ch)
		}
	case '<':
		if l.peekChar() == '=' {
			ch := l.Ch
			l.readChar()
			literal := string(ch) + string(l.Ch)
			tok = token.Token{Type: token.LTE, Literal: literal}
		} else {
			tok = newToken(token.LT, l.Ch)
		}
	case '>':
		if l.peekChar() == '=' {
			ch := l.Ch
			l.readChar()
			literal := string(ch) + string(l.Ch)
			tok = token.Token{Type: token.GTE, Literal: literal}
		} else {
			tok = newToken(token.GT, l.Ch)
		}
	case ';':
		tok = newToken(token.SEMICOLON, l.Ch)
	case ',':
		tok = newToken(token.COMMA, l.Ch)
	case '{':
		tok = newToken(token.LBRACE, l.Ch)
	case '}':
		tok = newToken(token.RBRACE, l.Ch)
	case '[':
		tok = newToken(token.LBRACKET, l.Ch)
	case ']':
		tok = newToken(token.RBRACKET, l.Ch)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isLetter(l.Ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.Ch) {
			tok.Literal = l.readNumber()
			if strings.Contains(tok.Literal, ".") {
				tok.Type = token.FLOAT
			} else {
				tok.Type = token.INT
			}
			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.Ch)
		}
	}

	l.readChar()
	return tok
}

func (l *Lexer) readIdentifier() string {
	position := l.Position
	for isIdentifierChar(l.Ch) {
		l.readChar()
	}
	return l.input[position:l.Position]
}

func (l *Lexer) readNumber() string {
	position := l.Position
	hasDecimal := false
	for isDigit(l.Ch) || (l.Ch == '.' && !hasDecimal) {
		if l.Ch == '.' {
			hasDecimal = true
		}
		l.readChar()
	}
	return l.input[position:l.Position]
}

func (l *Lexer) readString() string {
	position := l.Position + 1
	for {
		l.readChar()
		if l.Ch == '"' || l.Ch == 0 {
			break
		}
	}
	return l.input[position:l.Position]
}

func (l *Lexer) readComment() {
	for l.Ch != '\n' && l.Ch != 0 {
		l.readChar()
	}
}

// ------------------ helpers ------------------

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func isDigitFloat(ch byte) bool {
	return '0' <= ch && ch <= '9' || ch == '.'
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isIdentifierChar(ch byte) bool {
	return isLetter(ch) || isDigit(ch)
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

func (l *Lexer) peekChar() byte {
	if l.ReadPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.ReadPosition]
	}
}

func (l *Lexer) skipWhiteSpace() {
	for l.Ch == ' ' || l.Ch == '\t' || l.Ch == '\n' || l.Ch == '\r' {
		l.readChar()
	}
}
