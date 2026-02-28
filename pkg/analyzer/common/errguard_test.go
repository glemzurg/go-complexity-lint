package common

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

func parseIfStmt(t *testing.T, bodyCode string) *ast.IfStmt {
	t.Helper()
	src := "package p\nfunc f() error {\n" + bodyCode + "\nreturn nil\n}"
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "test.go", src, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}
	fn := f.Decls[0].(*ast.FuncDecl)
	for _, stmt := range fn.Body.List {
		if ifStmt, ok := stmt.(*ast.IfStmt); ok {
			return ifStmt
		}
	}
	t.Fatal("no IfStmt found")
	return nil
}

func TestIsErrGuard(t *testing.T) {
	tests := []struct {
		name string
		code string
		want bool
	}{
		{
			name: "simple err return",
			code: `if err != nil { return err }`,
			want: true,
		},
		{
			name: "nil and err return",
			code: `if err != nil { return nil, err }`,
			want: true,
		},
		{
			name: "with init statement",
			code: `if err := doSomething(); err != nil { return nil, err }`,
			want: true,
		},
		{
			name: "with fmt.Errorf",
			code: `if err != nil { return nil, fmt.Errorf("failed: %w", err) }`,
			want: true,
		},
		{
			name: "zero value struct and err",
			code: `if err != nil { return MyStruct{}, err }`,
			want: true,
		},
		{
			name: "zero int and err",
			code: `if err != nil { return 0, err }`,
			want: true,
		},
		{
			name: "empty string and err",
			code: `if err != nil { return "", err }`,
			want: true,
		},
		{
			name: "false and err",
			code: `if err != nil { return false, err }`,
			want: true,
		},
		{
			name: "two statements in body",
			code: `if err != nil { log.Print(err); return err }`,
			want: false,
		},
		{
			name: "condition is not err != nil",
			code: `if x > 0 { return nil, err }`,
			want: false,
		},
		{
			name: "non-zero return value",
			code: `if err != nil { return result, err }`,
			want: false,
		},
		{
			name: "has else clause",
			code: `if err != nil { return err } else { return nil }`,
			want: false,
		},
		{
			name: "body is not return",
			code: `if err != nil { panic(err) }`,
			want: false,
		},
		{
			name: "no return values",
			code: `if err != nil { return }`,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Need to parse with enough context for the various types
			src := "package p\nimport \"fmt\"\nimport \"log\"\nvar result int\nvar err error\nvar x int\ntype MyStruct struct{}\nfunc doSomething() error { return nil }\nfunc f() (int, error) {\n" + tt.code + "\nreturn 0, nil\n}"
			fset := token.NewFileSet()
			f, err2 := parser.ParseFile(fset, "test.go", src, parser.ParseComments)
			if err2 != nil {
				t.Fatal(err2)
			}

			var ifStmt *ast.IfStmt
			ast.Inspect(f, func(n ast.Node) bool {
				if is, ok := n.(*ast.IfStmt); ok && ifStmt == nil {
					ifStmt = is
					return false
				}
				return true
			})
			if ifStmt == nil {
				t.Fatal("no IfStmt found")
			}

			got := IsErrGuard(ifStmt)
			if got != tt.want {
				t.Errorf("IsErrGuard() = %v, want %v", got, tt.want)
			}
		})
	}
}
