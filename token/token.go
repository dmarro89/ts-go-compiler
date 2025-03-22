package token

type TokenType int

const (
	ILLEGAL TokenType = iota
	EOF

	// Identifiers + literals
	IDENT  // variable names, functions, etc.
	INT    // integer numbers
	STRING // strings

	// Equals and not equals
	EQ     // ==
	NOT_EQ // !=

	// Operators
	ASSIGN   // =
	PLUS     // +
	MINUS    // -
	BANG     // !
	ASTERISK // *
	SLASH    // /

	LT // <
	GT // >

	// Delimiters
	COMMA     // ,
	SEMICOLON // ;
	COLON     // :

	LPAREN // (
	RPAREN // )
	LBRACE // {
	RBRACE // }

	// Keywords
	FUNCTION
	LET
	CONST
	VAR
	RETURN
	IF
	ELSE
	CONSOLE
	LOG
	DOT

	TRUE
	FALSE
)

// Token represents a token in our lexer
type Token struct {
	Type    TokenType
	Literal string
	Line    int
	Column  int
}

func NewToken(t TokenType, l string, line int, column int) *Token {
	return &Token{Type: t, Literal: l, Line: line, Column: column}
}
