package lexer

import (
	"unicode"

	"github.com/dmarro89/ts-go-compiler/token"
)

// Lexer is responsible for scanning the source code
type Lexer struct {
	input        string
	position     int  // current position in input (points to current character)
	readPosition int  // current reading position in input (after current character)
	ch           byte // current character under examination
	line         int  // current line
	column       int  // current column
}

// New creates a new Lexer
func New(input string) *Lexer {
	l := &Lexer{input: input, line: 1, column: 0}
	l.readChar()
	return l
}

// readChar reads the next character and advances the position in the input
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
		return
	}

	l.ch = l.input[l.readPosition]
	l.position = l.readPosition
	l.readPosition++
	l.column++
	if l.ch == '\n' {
		l.line++
		l.column = 0
	}
}

// peekChar returns the next character without advancing the position
func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

// NextToken returns the next token from the input
func (l *Lexer) NextToken() token.Token {
	// Skip whitespace
	l.skipWhitespace()

	// Save current position for the token
	var tok token.Token

	tok.Line = l.line
	tok.Column = l.column

	switch l.ch {
	case '=':
		tok = token.Token{Type: token.ASSIGN, Literal: string(l.ch)}
	case '+':
		tok = token.Token{Type: token.PLUS, Literal: string(l.ch)}
	case '-':
		tok = token.Token{Type: token.MINUS, Literal: string(l.ch)}
	case '!':
		tok = token.Token{Type: token.BANG, Literal: string(l.ch)}
	case '*':
		tok = token.Token{Type: token.ASTERISK, Literal: string(l.ch)}
	case '/':
		// Handle comments
		if l.peekChar() == '/' {
			l.skipLineComment()
			return l.NextToken()
		} else if l.peekChar() == '*' {
			l.skipBlockComment()
			return l.NextToken()
		}
		tok = token.Token{Type: token.SLASH, Literal: string(l.ch)}
	case '<':
		tok = token.Token{Type: token.LT, Literal: string(l.ch)}
	case '>':
		tok = token.Token{Type: token.GT, Literal: string(l.ch)}
	case ',':
		tok = token.Token{Type: token.COMMA, Literal: string(l.ch)}
	case ';':
		tok = token.Token{Type: token.SEMICOLON, Literal: string(l.ch)}
	case ':':
		tok = token.Token{Type: token.COLON, Literal: string(l.ch)}
	case '(':
		tok = token.Token{Type: token.LPAREN, Literal: string(l.ch)}
	case ')':
		tok = token.Token{Type: token.RPAREN, Literal: string(l.ch)}
	case '{':
		tok = token.Token{Type: token.LBRACE, Literal: string(l.ch)}
	case '}':
		tok = token.Token{Type: token.RBRACE, Literal: string(l.ch)}
	case '.':
		tok = token.Token{Type: token.DOT, Literal: string(l.ch)}
	case '"':
		tok.Type = token.STRING
		tok.Literal = l.readString('"')
	case '\'':
		tok.Type = token.STRING
		tok.Literal = l.readString('\'')
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = lookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			tok.Type = token.INT
			tok.Literal = l.readNumber()
			return tok
		} else {
			tok = token.Token{Type: token.ILLEGAL, Literal: string(l.ch)}
		}
	}

	l.readChar()
	return tok
}

func (l *Lexer) skipWhitespace() {
	for unicode.IsSpace(rune(l.ch)) {
		l.readChar()
	}
}

// skipLineComment skips line comments
func (l *Lexer) skipLineComment() {
	// Skip the second '/'
	l.readChar()

	// Continue reading until we encounter a newline or EOF
	for l.ch != '\n' && l.ch != 0 {
		l.readChar()
	}
}

// skipBlockComment skips block comments
func (l *Lexer) skipBlockComment() {
	// Skip '/*'
	l.readChar()
	l.readChar()

	for {
		if l.ch == 0 {
			// End of file inside a block comment
			break
		}
		if l.ch == '*' && l.peekChar() == '/' {
			// End of block comment
			l.readChar() // skip '*'
			l.readChar() // skip '/'
			break
		}
		l.readChar()
	}
}

// readIdentifier reads an identifier
func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

// readNumber reads a number
func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

// readString reads a string delimited by quotes
func (l *Lexer) readString(quote byte) string {
	l.readChar() // skip the initial quote
	position := l.position

	for l.ch != quote && l.ch != 0 {
		l.readChar()
	}

	result := l.input[position:l.position]
	return result
}

// isLetter checks if a character is a letter
func isLetter(ch byte) bool {
	return unicode.IsLetter(rune(ch)) || ch == '_'
}

// isDigit checks if a character is a digit
func isDigit(ch byte) bool {
	return unicode.IsDigit(rune(ch))
}

// lookupIdent checks if an identifier is a keyword
func lookupIdent(ident string) token.TokenType {
	switch ident {
	case "function":
		return token.FUNCTION
	case "let":
		return token.LET
	case "const":
		return token.CONST
	case "var":
		return token.VAR
	case "return":
		return token.RETURN
	case "if":
		return token.IF
	case "else":
		return token.ELSE
	case "console":
		return token.CONSOLE
	case "log":
		return token.LOG
	case ".":
		return token.DOT
	default:
		return token.IDENT
	}
}
