package common

import "testing"

func TestWarningsModeSet(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    WarningsMode
		wantErr bool
	}{
		{name: "default", input: "default", want: WarningsDefault},
		{name: "none", input: "none", want: WarningsNone},
		{name: "error", input: "error", want: WarningsError},
		{name: "invalid", input: "quiet", wantErr: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var mode WarningsMode
			err := mode.Set(tc.input)
			if (err != nil) != tc.wantErr {
				t.Fatalf("Set(%q) error = %v, wantErr %v", tc.input, err, tc.wantErr)
			}
			if !tc.wantErr && mode != tc.want {
				t.Fatalf("Set(%q) = %v, want %v", tc.input, mode, tc.want)
			}
		})
	}
}

func TestWarningsModeString(t *testing.T) {
	tests := []struct {
		mode WarningsMode
		want string
	}{
		{WarningsDefault, "default"},
		{WarningsNone, "none"},
		{WarningsError, "error"},
	}

	for _, tc := range tests {
		if got := tc.mode.String(); got != tc.want {
			t.Fatalf("WarningsMode(%d).String() = %q, want %q", tc.mode, got, tc.want)
		}
	}
}

func TestWarningsModeReportDiagnostic(t *testing.T) {
	tests := []struct {
		name     string
		mode     WarningsMode
		category string
		want     bool
	}{
		{name: "default prints warning", mode: WarningsDefault, category: "warning", want: true},
		{name: "default prints error", mode: WarningsDefault, category: "error", want: true},
		{name: "none hides warning", mode: WarningsNone, category: "warning", want: false},
		{name: "none prints error", mode: WarningsNone, category: "error", want: true},
		{name: "error prints warning", mode: WarningsError, category: "warning", want: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.mode.ReportDiagnostic(tc.category); got != tc.want {
				t.Fatalf("ReportDiagnostic(%q) = %v, want %v", tc.category, got, tc.want)
			}
		})
	}
}

func TestWarningsModeDiagnosticFails(t *testing.T) {
	tests := []struct {
		name     string
		mode     WarningsMode
		category string
		want     bool
	}{
		{name: "default ignores warning", mode: WarningsDefault, category: "warning", want: false},
		{name: "default fails on error", mode: WarningsDefault, category: "error", want: true},
		{name: "none ignores warning", mode: WarningsNone, category: "warning", want: false},
		{name: "none fails on error", mode: WarningsNone, category: "error", want: true},
		{name: "error fails on warning", mode: WarningsError, category: "warning", want: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.mode.DiagnosticFails(tc.category); got != tc.want {
				t.Fatalf("DiagnosticFails(%q) = %v, want %v", tc.category, got, tc.want)
			}
		})
	}
}
