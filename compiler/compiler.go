package compiler

import (
	"errors"
	"os"
	"strings"

	"github.com/dmarro89/ts-go-compiler/codegen"
	"github.com/dmarro89/ts-go-compiler/lexer"
	"github.com/dmarro89/ts-go-compiler/parser"
	"github.com/dmarro89/ts-go-compiler/typecheck"
)

// Compiler handles the compilation process
type Compiler struct {
}

// New creates a new compiler
func New() *Compiler {
	return &Compiler{}
}

// CompileFile compiles a TypeScript file
func (c *Compiler) CompileFile(filename string, outputFile string) error {
	// Read input file
	input, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	// Compile source code
	output, err := c.Compile(string(input))
	if err != nil {
		return err
	}

	// Write output file
	return os.WriteFile(outputFile, []byte(output), 0644)
}

// Compile compiles TypeScript source code
func (c *Compiler) Compile(input string) (string, error) {
	// Initialize lexer
	l := lexer.New(input)

	// Parse input
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		return "", errors.New("parse errors: " + strings.Join(p.Errors(), ", "))
	}

	// Type check
	tc := typecheck.New()
	errs := tc.Check(program)
	if len(errs) > 0 {
		return "", errors.New("type error: " + errs[0])
	}

	// Generate code
	generator := codegen.New()
	output := generator.GenerateJavaScript(program)

	return output, nil
}
