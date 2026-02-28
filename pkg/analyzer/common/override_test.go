package common

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

func parseFuncDecl(t *testing.T, src string) *ast.FuncDecl {
	t.Helper()
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "test.go", "package p\n"+src, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}
	for _, decl := range f.Decls {
		if fn, ok := decl.(*ast.FuncDecl); ok {
			return fn
		}
	}
	t.Fatal("no FuncDecl found")
	return nil
}

func TestParseOverrides(t *testing.T) {
	defaults := Thresholds{WarnAt: 4, FailAt: 6}

	tests := []struct {
		name   string
		src    string
		metric string
		want   Thresholds
	}{
		{
			name:   "no override",
			src:    "func Foo() {}",
			metric: "nestdepth",
			want:   defaults,
		},
		{
			name:   "both warn and fail",
			src:    "//complexity:nestdepth:warn=8,fail=10\nfunc Foo() {}",
			metric: "nestdepth",
			want:   Thresholds{WarnAt: 8, FailAt: 10},
		},
		{
			name:   "warn only",
			src:    "//complexity:cyclo:warn=20\nfunc Foo() {}",
			metric: "cyclo",
			want:   Thresholds{WarnAt: 20, FailAt: 6},
		},
		{
			name:   "fail only",
			src:    "//complexity:params:fail=12\nfunc Foo() {}",
			metric: "params",
			want:   Thresholds{WarnAt: 4, FailAt: 12},
		},
		{
			name:   "wrong metric name",
			src:    "//complexity:cyclo:warn=20\nfunc Foo() {}",
			metric: "nestdepth",
			want:   defaults,
		},
		{
			name:   "malformed value",
			src:    "//complexity:nestdepth:warn=abc\nfunc Foo() {}",
			metric: "nestdepth",
			want:   defaults,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fn := parseFuncDecl(t, tt.src)
			got := ParseOverrides(fn, tt.metric, defaults)
			if got != tt.want {
				t.Errorf("ParseOverrides() = %+v, want %+v", got, tt.want)
			}
		})
	}
}
