package codegen

import (
	"bytes"
	"fmt"

	"github.com/dmarro89/ts-go-compiler/ast"
)

// Generator generates code from an AST
type Generator struct {
}

// New creates a new code generator
func New() *Generator {
	return &Generator{}
}

// GenerateJavaScript generates JavaScript code
func (g *Generator) GenerateJavaScript(program *ast.Program) string {
	var out bytes.Buffer

	for _, stmt := range program.Statements {
		out.WriteString(g.generateJSStatement(stmt))
		out.WriteString("\n")
	}

	return out.String()
}

func (g *Generator) generateJSStatement(stmt ast.Statement) string {
	switch s := stmt.(type) {
	case *ast.LetStatement:
		return fmt.Sprintf("let %s = %s;", s.Name.Value, g.generateJSExpression(s.Value))
	case *ast.ReturnStatement:
		return fmt.Sprintf("return %s;", g.generateJSExpression(s.ReturnValue))
	case *ast.ExpressionStatement:
		return g.generateJSExpression(s.Expression) + ";"
	default:
		return ""
	}
}

func (g *Generator) generateJSExpression(expr ast.Expression) string {
	if expr == nil {
		return ""
	}

	switch e := expr.(type) {
	case *ast.IntegerLiteral:
		return e.Token.Literal
	case *ast.StringLiteral:
		return fmt.Sprintf("\"%s\"", e.Value)
	case *ast.Identifier:
		return e.Value
	default:
		return ""
	}
}
