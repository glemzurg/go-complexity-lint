package common

import (
	"testing"
)

func TestIsExcluded(t *testing.T) {
	tests := []struct {
		name     string
		patterns string
		filename string
		want     bool
	}{
		{
			name:     "empty patterns never excludes",
			patterns: "",
			filename: "foo.go",
			want:     false,
		},
		{
			name:     "single pattern match",
			patterns: "*_gen.go",
			filename: "model_gen.go",
			want:     true,
		},
		{
			name:     "single pattern no match",
			patterns: "*_gen.go",
			filename: "model.go",
			want:     false,
		},
		{
			name:     "multiple patterns first matches",
			patterns: "*_gen.go,mock_*.go",
			filename: "model_gen.go",
			want:     true,
		},
		{
			name:     "multiple patterns second matches",
			patterns: "*_gen.go,mock_*.go",
			filename: "mock_service.go",
			want:     true,
		},
		{
			name:     "multiple patterns none match",
			patterns: "*_gen.go,mock_*.go",
			filename: "service.go",
			want:     false,
		},
		{
			name:     "whitespace around commas",
			patterns: " *_gen.go , mock_*.go ",
			filename: "mock_service.go",
			want:     true,
		},
		{
			name:     "full path matches base name",
			patterns: "*_gen.go",
			filename: "/home/user/project/pkg/model_gen.go",
			want:     true,
		},
		{
			name:     "full path no match",
			patterns: "*_gen.go",
			filename: "/home/user/project/pkg/model.go",
			want:     false,
		},
		{
			name:     "invalid glob pattern does not panic",
			patterns: "[",
			filename: "foo.go",
			want:     false,
		},
		{
			name:     "trailing comma ignored",
			patterns: "*_gen.go,",
			filename: "model_gen.go",
			want:     true,
		},
		{
			name:     "exact filename match",
			patterns: "generated.go",
			filename: "generated.go",
			want:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ExcludePatterns = tt.patterns
			t.Cleanup(func() { ExcludePatterns = "" })

			got := IsExcluded(tt.filename)
			if got != tt.want {
				t.Errorf("IsExcluded(%q) with patterns %q = %v, want %v",
					tt.filename, tt.patterns, got, tt.want)
			}
		})
	}
}
