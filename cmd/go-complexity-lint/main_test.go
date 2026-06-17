package main

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestWarningsModesCLI(t *testing.T) {
	bin := buildBinary(t)
	pkg := filepath.Join("..", "..", "pkg", "analyzer", "params", "testdata", "src", "params")

	tests := []struct {
		name       string
		args       []string
		wantExit   int
		wantSubstr string
		wantAbsent string
	}{
		{
			name:       "default prints warnings and exits 0",
			args:       []string{"-warnings=default", "-params.fail=100", pkg},
			wantExit:   0,
			wantSubstr: "FiveParams",
		},
		{
			name:       "none suppresses warnings and exits 0",
			args:       []string{"-warnings=none", "-params.fail=100", pkg},
			wantExit:   0,
			wantAbsent: "FiveParams",
		},
		{
			name:       "error fails on warnings",
			args:       []string{"-warnings=error", "-params.fail=100", pkg},
			wantExit:   1,
			wantSubstr: "FiveParams",
		},
		{
			name:       "default still fails on red-zone violations",
			args:       []string{"-warnings=default", pkg},
			wantExit:   1,
			wantSubstr: "SevenParams",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cmd := exec.Command(bin, tc.args...)
			var stderr bytes.Buffer
			cmd.Stderr = &stderr
			err := cmd.Run()

			exit := 0
			if err != nil {
				exitErr, ok := err.(*exec.ExitError)
				if !ok {
					t.Fatalf("Run() error = %v", err)
				}
				exit = exitErr.ExitCode()
			}
			if exit != tc.wantExit {
				t.Fatalf("exit code = %d, want %d; stderr:\n%s", exit, tc.wantExit, stderr.String())
			}

			out := stderr.String()
			if tc.wantSubstr != "" && !strings.Contains(out, tc.wantSubstr) {
				t.Fatalf("stderr missing %q:\n%s", tc.wantSubstr, out)
			}
			if tc.wantAbsent != "" && strings.Contains(out, tc.wantAbsent) {
				t.Fatalf("stderr unexpectedly contains %q:\n%s", tc.wantAbsent, out)
			}
		})
	}
}

func buildBinary(t *testing.T) string {
	t.Helper()

	tmp := t.TempDir()
	bin := filepath.Join(tmp, "go-complexity-lint")
	cmd := exec.Command("go", "build", "-o", bin, ".")
	cmd.Dir = "."
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("build binary: %v\n%s", err, out)
	}
	return bin
}

func TestWarningsModeInvalidCLI(t *testing.T) {
	bin := buildBinary(t)
	pkg := filepath.Join("..", "..", "pkg", "analyzer", "params", "testdata", "src", "params")

	cmd := exec.Command(bin, "-warnings=quiet", pkg)
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err == nil {
		t.Fatal("expected invalid -warnings value to fail")
	}
}
