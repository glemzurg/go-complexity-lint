package common

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

func TestErrGuardCallExprs(t *testing.T) {
	src := `package p
import "ext/wrap"
var err error
func f() error {
	_ = wrap.A()
	if err != nil {
		return nil, wrap.B(err)
	}
	if x := 1; x > 0 {
		return wrap.C()
	}
	return nil
}`
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "test.go", src, 0)
	if err != nil {
		t.Fatal(err)
	}
	var fn *ast.FuncDecl
	for _, decl := range f.Decls {
		if fd, ok := decl.(*ast.FuncDecl); ok && fd.Name.Name == "f" {
			fn = fd
			break
		}
	}
	if fn == nil {
		t.Fatal("function f not found")
	}

	excluded := ErrGuardCallExprs(fn.Body)
	var excludedNames []string
	ast.Inspect(fn.Body, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}
		if _, ok := excluded[call]; !ok {
			return true
		}
		switch fun := call.Fun.(type) {
		case *ast.Ident:
			excludedNames = append(excludedNames, fun.Name)
		case *ast.SelectorExpr:
			excludedNames = append(excludedNames, fun.Sel.Name)
		}
		return true
	})

	if len(excludedNames) != 1 || excludedNames[0] != "B" {
		t.Fatalf("excluded calls = %v, want [B]", excludedNames)
	}
}