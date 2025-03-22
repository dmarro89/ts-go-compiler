package lexer

import (
	"testing"

	"github.com/dmarro89/ts-go-compiler/token"
)

func TestNextToken(t *testing.T) {
	input := `let five = 5;
	let ten = 10;

	let add = function(x, y) {
		return x + y;
	};

	let result = add(five, ten);
	"hello world";
	'hello world';
	console.log("hello world");
	`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.LET, "let"},
		{token.IDENT, "five"},
		{token.ASSIGN, "="},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},

		{token.LET, "let"},
		{token.IDENT, "ten"},
		{token.ASSIGN, "="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},

		{token.LET, "let"},
		{token.IDENT, "add"},
		{token.ASSIGN, "="},
		{token.FUNCTION, "function"},
		{token.LPAREN, "("},
		{token.IDENT, "x"},
		{token.COMMA, ","},
		{token.IDENT, "y"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.IDENT, "x"},
		{token.PLUS, "+"},
		{token.IDENT, "y"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.SEMICOLON, ";"},

		{token.LET, "let"},
		{token.IDENT, "result"},
		{token.ASSIGN, "="},
		{token.IDENT, "add"},
		{token.LPAREN, "("},
		{token.IDENT, "five"},
		{token.COMMA, ","},
		{token.IDENT, "ten"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},

		{token.STRING, "hello world"},
		{token.SEMICOLON, ";"},

		{token.STRING, "hello world"},
		{token.SEMICOLON, ";"},

		{token.CONSOLE, "console"},
		{token.DOT, "."},
		{token.LOG, "log"},
		{token.LPAREN, "("},
		{token.STRING, "hello world"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},

		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%s, got=%s",
				i, getTokenTypeString(tt.expectedType), getTokenTypeString(tok.Type))
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%s, got=%s",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

// utility method to get the string value of token type
func getTokenTypeString(t token.TokenType) string {
	switch t {
	case token.ILLEGAL:
		return "ILLEGAL"
	case token.EOF:
		return "EOF"
	case token.IDENT:
		return "IDENT"
	case token.INT:
		return "INT"
	case token.STRING:
		return "STRING"
	case token.ASSIGN:
		return "ASSIGN"
	case token.PLUS:
		return "PLUS"
	case token.COMMA:
		return "COMMA"
	case token.SEMICOLON:
		return "SEMICOLON"
	case token.LPAREN:
		return "LPAREN"
	case token.RPAREN:
		return "RPAREN"
	case token.LBRACE:
		return "LBRACE"
	case token.RBRACE:
		return "RBRACE"
	case token.FUNCTION:
		return "FUNCTION"
	case token.LET:
		return "LET"
	case token.RETURN:
		return "RETURN"
	case token.BANG:
		return "BANG"
	case token.MINUS:
		return "MINUS"
	case token.ASTERISK:
		return "ASTERISK"
	case token.SLASH:
		return "SLASH"
	case token.LT:
		return "LT"
	case token.GT:
		return "GT"
	case token.COLON:
		return "COLON"
	case token.DOT:
		return "DOT"
	case token.CONSOLE:
		return "CONSOLE"
	case token.LOG:
		return "LOG"
	default:
		return "UNKNOWN"
	}
}

func TestNextTokenOperators(t *testing.T) {
	input := `let a = 1 + 2 - 3 * 4 / 5;
	let b = true && false || !true;
	let c = 3 > 2;
	let d = 2 < 3;
	let e = a == b;
	let f = a != b;
	`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.LET, "let"},
		{token.IDENT, "a"},
		{token.ASSIGN, "="},
		{token.INT, "1"},
		{token.PLUS, "+"},
		{token.INT, "2"},
		{token.MINUS, "-"},
		{token.INT, "3"},
		{token.ASTERISK, "*"},
		{token.INT, "4"},
		{token.SLASH, "/"},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},

		// Additional operators can be tested here
		{token.LET, "let"},
		{token.IDENT, "b"},
		{token.ASSIGN, "="},
		{token.IDENT, "true"},
		// ... continue with other tokens
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestComment(t *testing.T) {
	input := `
	// This is a line comment
	let a = 5; // This is another comment
	/* 
	   This is a 
	   block comment
	*/
	let b = 10;
	`

	l := New(input)
	var tokens []token.Token

	for {
		tok := l.NextToken()
		tokens = append(tokens, tok)
		if tok.Type == token.EOF {
			break
		}
	}

	// After processing comments, we should only have the meaningful tokens
	expected := []token.TokenType{token.LET, token.IDENT, token.ASSIGN, token.INT, token.SEMICOLON, token.LET, token.IDENT, token.ASSIGN, token.INT, token.SEMICOLON, token.EOF}

	if len(tokens) != len(expected) {
		t.Fatalf("Expected %d tokens, got %d", len(expected), len(tokens))
	}

	for i, tok := range tokens {
		if i < len(expected) && tok.Type != expected[i] {
			t.Fatalf("Expected token %d to be %q, got %q", i, expected[i], tok.Type)
		}
	}
}
