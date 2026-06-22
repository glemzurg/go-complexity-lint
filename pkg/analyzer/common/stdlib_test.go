package common

import "testing"

func TestIsStdlib(t *testing.T) {
	tests := []struct {
		name    string
		pkgPath string
		want    bool
	}{
		{name: "fmt", pkgPath: "fmt", want: true},
		{name: "errors", pkgPath: "errors", want: true},
		{name: "context", pkgPath: "context", want: true},
		{name: "crypto subpackage", pkgPath: "crypto/rand", want: true},
		{name: "empty", pkgPath: "", want: false},
		{name: "unsafe", pkgPath: "unsafe", want: true},
		{name: "analyzer common package", pkgPath: "github.com/glemzurg/go-complexity-lint/pkg/analyzer/common", want: false},
		{name: "analyzer fanout testdata package", pkgPath: "github.com/glemzurg/go-complexity-lint/pkg/analyzer/fanout/testdata/src/fanout", want: false},
		{name: "external module", pkgPath: "github.com/glemzurg/go-complexity-lint", want: false},
		{name: "external dotted first segment", pkgPath: "ext.pkg/dep", want: false},
		{name: "connect", pkgPath: "connectrpc.com/connect", want: false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := IsStdlib(tc.pkgPath); got != tc.want {
				t.Fatalf("IsStdlib(%q) = %v, want %v", tc.pkgPath, got, tc.want)
			}
		})
	}
}