package common

import (
	"go/ast"
	"testing"
)

func TestFuncName(t *testing.T) {
	tests := []struct {
		name string
		src  string
		want string
	}{
		{
			name: "plain function",
			src:  "func Foo() {}",
			want: "Foo",
		},
		{
			name: "pointer receiver method",
			src:  "type T struct{}\nfunc (t *T) Bar() {}",
			want: "*T.Bar",
		},
		{
			name: "value receiver method",
			src:  "type T struct{}\nfunc (t T) Baz() {}",
			want: "T.Baz",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fn := parseFuncDecl(t, tt.src)
			got := FuncName(fn)
			if got != tt.want {
				t.Errorf("FuncName() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestExprName(t *testing.T) {
	tests := []struct {
		name string
		expr ast.Expr
		want string
	}{
		{
			name: "ident",
			expr: &ast.Ident{Name: "Foo"},
			want: "Foo",
		},
		{
			name: "star expr",
			expr: &ast.StarExpr{X: &ast.Ident{Name: "Foo"}},
			want: "*Foo",
		},
		{
			name: "index expr (generic single type param)",
			expr: &ast.IndexExpr{X: &ast.Ident{Name: "Foo"}},
			want: "Foo",
		},
		{
			name: "index list expr (generic multiple type params)",
			expr: &ast.IndexListExpr{X: &ast.Ident{Name: "Foo"}},
			want: "Foo",
		},
		{
			name: "star wrapping index expr",
			expr: &ast.StarExpr{X: &ast.IndexExpr{X: &ast.Ident{Name: "Foo"}}},
			want: "*Foo",
		},
		{
			name: "unknown expr type",
			expr: &ast.ArrayType{},
			want: "?",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ExprName(tt.expr)
			if got != tt.want {
				t.Errorf("ExprName() = %q, want %q", got, tt.want)
			}
		})
	}
}
