package common

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

func parseFunc(t *testing.T, src string) *ast.FuncDecl {
	t.Helper()

	file, err := parser.ParseFile(token.NewFileSet(), "test.go", "package p\n"+src, 0)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	return file.Decls[0].(*ast.FuncDecl)
}

func TestIsExemptContextParam(t *testing.T) {
	tests := []struct {
		name  string
		param string
		typ   string
		want  bool
	}{
		{name: "ctx context.Context", param: "ctx", typ: "context.Context", want: true},
		{name: "wrong name", param: "c", typ: "context.Context", want: false},
		{name: "wrong type selector", param: "ctx", typ: "context.CancelFunc", want: false},
		{name: "unqualified Context", param: "ctx", typ: "Context", want: false},
		{name: "pointer to context.Context", param: "ctx", typ: "*context.Context", want: false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			fn := parseFunc(t, "func f("+tc.param+" "+tc.typ+") {}")
			field := fn.Type.Params.List[0]
			got := IsExemptContextParam(field.Names[0].Name, field.Type)
			if got != tc.want {
				t.Fatalf("IsExemptContextParam(%q, %q) = %v, want %v", tc.param, tc.typ, got, tc.want)
			}
		})
	}
}

func TestIsExemptContextParamGrouped(t *testing.T) {
	fn := parseFunc(t, "func f(ctx, cancel context.Context) {}")
	field := fn.Type.Params.List[0]

	if !IsExemptContextParam(field.Names[0].Name, field.Type) {
		t.Fatal("expected ctx to be exempt")
	}
	if IsExemptContextParam(field.Names[1].Name, field.Type) {
		t.Fatal("expected cancel not to be exempt")
	}
}
