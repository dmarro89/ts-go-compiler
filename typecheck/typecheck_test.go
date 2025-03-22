package typecheck

import (
	"testing"

	"github.com/dmarro89/ts-go-compiler/lexer"
	"github.com/dmarro89/ts-go-compiler/parser"
)

func TestBasicTypeChecking(t *testing.T) {
	input := `
	let x = 5;
	let y = "hello";
	let z = x;
	`

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		t.Fatalf("parser has %d errors", len(p.Errors()))
		for _, msg := range p.Errors() {
			t.Errorf("parser error: %q", msg)
		}
		t.FailNow()
	}

	tc := New()
	errors := tc.Check(program)

	if len(errors) > 0 {
		t.Fatalf("type checker has %d errors", len(errors))
		for _, msg := range errors {
			t.Errorf("type error: %q", msg)
		}
	}

	// Check the types in environment
	varX, ok := tc.env.Get("x")
	if !ok {
		t.Fatalf("variable x not found in environment")
	}
	if varX.String() != "number" {
		t.Errorf("expected x to be of type number, got %s", varX.String())
	}

	varY, ok := tc.env.Get("y")
	if !ok {
		t.Fatalf("variable y not found in environment")
	}
	if varY.String() != "string" {
		t.Errorf("expected y to be of type string, got %s", varY.String())
	}

	varZ, ok := tc.env.Get("z")
	if !ok {
		t.Fatalf("variable z not found in environment")
	}
	if varZ.String() != "number" {
		t.Errorf("expected z to be of type number, got %s", varZ.String())
	}
}

func TestUndefinedVariable(t *testing.T) {
	input := `
	let x = y;
	`

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		t.Fatalf("parser has %d errors", len(p.Errors()))
		for _, msg := range p.Errors() {
			t.Errorf("parser error: %q", msg)
		}
		t.FailNow()
	}

	tc := New()
	errors := tc.Check(program)

	if len(errors) == 0 {
		t.Fatalf("expected type errors but got none")
	}

	expectedError := "undefined variable: y"
	for _, err := range errors {
		if err == expectedError {
			return // Test passed
		}
	}

	t.Errorf("expected error message %q not found in errors: %v", expectedError, errors)
}
