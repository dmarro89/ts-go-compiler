package compiler

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCompile(t *testing.T) {
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
			`let z = 10;
			return z;`,
			`let z = 10;
return z;`,
		},
	}

	compiler := New()

	for _, tt := range tests {
		output, err := compiler.Compile(tt.input)
		if err != nil {
			t.Errorf("compile error: %s", err)
			continue
		}

		// Normalize whitespace for comparison
		output = strings.TrimSpace(output)
		expected := strings.TrimSpace(tt.expected)

		if output != expected {
			t.Errorf("expected=%q, got=%q", expected, output)
		}
	}
}

func TestCompileFile(t *testing.T) {
	// Create temporary directory
	tempDir, err := os.MkdirTemp("", "compiler_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test input file
	inputContent := `let x = 5;
let y = "hello";
let z = x;
return z;`

	inputFile := filepath.Join(tempDir, "test.ts")
	err = os.WriteFile(inputFile, []byte(inputContent), 0644)
	if err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	// Set output file
	outputFile := filepath.Join(tempDir, "test.js")

	// Compile the file
	compiler := New()
	err = compiler.CompileFile(inputFile, outputFile)
	if err != nil {
		t.Fatalf("failed to compile file: %v", err)
	}

	// Read output and verify
	outputContent, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("failed to read output file: %v", err)
	}

	expected := `let x = 5;
let y = "hello";
let z = x;
return z;`

	if strings.TrimSpace(string(outputContent)) != strings.TrimSpace(expected) {
		t.Errorf("expected=%q, got=%q", expected, string(outputContent))
	}
}

func TestCompileWithErrors(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			"Parse error",
			`let x = ;`,
		},
		{
			"Type error",
			`let x = y;`,
		},
	}

	compiler := New()

	for _, tt := range tests {
		_, err := compiler.Compile(tt.input)
		if err == nil {
			t.Errorf("%s: expected error but got none", tt.name)
		}
	}
}

func TestCompileSimpleFile(t *testing.T) {
	// Set up the test file
	testFile := "../testdata/test.ts"

	// Create a compiler instance
	compiler := New()

	// Compile the file
	err := compiler.CompileFile(testFile, "simple_file.ts")
	if err != nil {
		// Print parser errors if available
		if perrs, ok := err.(interface{ Errors() []string }); ok {
			t.Log("Parser errors:")
			for _, msg := range perrs.Errors() {
				t.Logf("  - %s", msg)
			}
		}
		t.Fatalf("Error during the compilation of %s: %v", testFile, err)
	}

	if _, err := os.Stat("simple_file.ts"); os.IsNotExist(err) {
		t.Fatalf("Output file %s was not created", "simple_file.ts")
	}

	os.Remove("simple_file.ts")
}

func TestCompileNonExistentFile(t *testing.T) {
	compiler := New()
	err := compiler.CompileFile("testdata/non_existent.ts", "non_existent.ts")

	if err == nil {
		t.Fatal("Error was expected for non-existent file, but none was generated")
	}
}
