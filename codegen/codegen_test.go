package codegen

import (
	"strings"
	"testing"

	"github.com/dmarro89/ts-go-compiler/lexer"
	"github.com/dmarro89/ts-go-compiler/parser"
)

func TestJavaScriptGeneration(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			`let x = 5;`,
			`let x = 5;`,
		},
		{
			`let y = "hello";`,
			`let y = "hello";`,
		},
		{
			`let z = x;`,
			`let z = x;`,
		},
		{
			`return 10;`,
			`return 10;`,
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := parser.New(l)
		program := p.ParseProgram()

		generator := New()
		output := generator.GenerateJavaScript(program)

		// Trim whitespace for comparison
		output = strings.TrimSpace(output)
		expected := strings.TrimSpace(tt.expected)

		if output != expected {
			t.Errorf("expected=%q, got=%q", expected, output)
		}
	}
}

func TestJavaScriptMultiStatementGeneration(t *testing.T) {
	input := `
	let x = 5;
	let y = "hello";
	let z = x;
	return z;
	`

	expected := `let x = 5;
let y = "hello";
let z = x;
return z;`

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	generator := New()
	output := generator.GenerateJavaScript(program)

	// Normalize whitespace for comparison
	output = strings.TrimSpace(output)
	expected = strings.TrimSpace(expected)

	if output != expected {
		t.Errorf("expected=%q, got=%q", expected, output)
	}
}
