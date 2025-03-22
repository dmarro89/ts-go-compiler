package typecheck

import (
	"fmt"

	"github.com/dmarro89/ts-go-compiler/ast"
)

// TypeScript type
type Type interface {
	String() string
}

// BasicType represents primitive types like number, string, etc.
type BasicType struct {
	Name string
}

func (t *BasicType) String() string {
	return t.Name
}

// TypeChecker performs type checking on the AST
type TypeChecker struct {
	errors []string
	env    *TypeEnvironment
}

// TypeEnvironment stores variable types
type TypeEnvironment struct {
	store map[string]Type
	outer *TypeEnvironment
}

// New creates a new TypeChecker
func New() *TypeChecker {
	return &TypeChecker{
		errors: []string{},
		env:    NewTypeEnvironment(),
	}
}

// NewTypeEnvironment creates a new type environment
func NewTypeEnvironment() *TypeEnvironment {
	return &TypeEnvironment{
		store: make(map[string]Type),
	}
}

func (tc *TypeChecker) Check(program *ast.Program) []string {
	for _, stmt := range program.Statements {
		tc.checkStatement(stmt)
	}
	return tc.errors
}

func (tc *TypeChecker) checkStatement(stmt ast.Statement) Type {
	switch s := stmt.(type) {
	case *ast.LetStatement:
		return tc.checkLetStatement(s)
	case *ast.ReturnStatement:
		return tc.checkReturnStatement(s)
	case *ast.ExpressionStatement:
		return tc.checkExpression(s.Expression)
	case *ast.BlockStatement:
		return tc.checkBlockStatement(s)
	default:
		return &BasicType{Name: "void"}
	}
}

func (tc *TypeChecker) checkBlockStatement(block *ast.BlockStatement) Type {
	var lastType Type
	for _, stmt := range block.Statements {
		lastType = tc.checkStatement(stmt)
	}
	return lastType
}

func (tc *TypeChecker) checkLetStatement(stmt *ast.LetStatement) Type {
	if stmt.Value == nil {
		// Default to any type if no value is assigned
		tc.env.Set(stmt.Name.Value, &BasicType{Name: "any"})
		return &BasicType{Name: "any"}
	}

	valueType := tc.checkExpression(stmt.Value)
	tc.env.Set(stmt.Name.Value, valueType)
	return valueType
}

func (tc *TypeChecker) checkReturnStatement(stmt *ast.ReturnStatement) Type {
	if stmt.ReturnValue == nil {
		return &BasicType{Name: "void"}
	}
	return tc.checkExpression(stmt.ReturnValue)
}

func (tc *TypeChecker) checkExpression(expr ast.Expression) Type {
	switch e := expr.(type) {
	case *ast.IntegerLiteral:
		return &BasicType{Name: "number"}
	case *ast.StringLiteral:
		return &BasicType{Name: "string"}
	case *ast.Identifier:
		return tc.checkIdentifier(e)
	case *ast.FunctionLiteral:
		return &BasicType{Name: "function"}
	default:
		tc.addError(fmt.Sprintf("unknown expression type: %T", expr))
		return &BasicType{Name: "unknown"}
	}
}

func (tc *TypeChecker) checkIdentifier(ident *ast.Identifier) Type {
	if val, ok := tc.env.Get(ident.Value); ok {
		return val
	}
	tc.addError(fmt.Sprintf("undefined variable: %s", ident.Value))
	return &BasicType{Name: "unknown"}
}

func (tc *TypeChecker) addError(msg string) {
	tc.errors = append(tc.errors, msg)
}

// Get retrieves a type from the environment
func (env *TypeEnvironment) Get(name string) (Type, bool) {
	obj, ok := env.store[name]
	if !ok && env.outer != nil {
		obj, ok = env.outer.Get(name)
	}
	return obj, ok
}

// Set adds a type to the environment
func (env *TypeEnvironment) Set(name string, val Type) Type {
	env.store[name] = val
	return val
}
