package common

import "fmt"

// WarningsMode controls how yellow-zone diagnostics are reported by the
// standalone CLI. Red-zone diagnostics are always printed and always fail.
type WarningsMode int

const (
	// WarningsDefault prints warnings and exits 0 when only warnings are present.
	WarningsDefault WarningsMode = iota
	// WarningsNone suppresses warning output and does not fail on warnings.
	WarningsNone
	// WarningsError prints warnings and exits 1 when any warning is present.
	WarningsError
)

// String returns the flag value for m.
func (m WarningsMode) String() string {
	switch m {
	case WarningsNone:
		return "none"
	case WarningsError:
		return "error"
	default:
		return "default"
	}
}

// Set parses a -warnings flag value.
func (m *WarningsMode) Set(value string) error {
	switch value {
	case "default":
		*m = WarningsDefault
	case "none":
		*m = WarningsNone
	case "error":
		*m = WarningsError
	default:
		return fmt.Errorf("invalid warnings mode %q (want default, none, or error)", value)
	}
	return nil
}

// ReportDiagnostic reports whether a diagnostic with the given category
// should be printed. Empty category (green zone) is never emitted by analyzers.
func (m WarningsMode) ReportDiagnostic(category string) bool {
	if category == "warning" {
		return m != WarningsNone
	}
	return category == "error"
}

// DiagnosticFails reports whether a diagnostic with the given category
// should contribute to a non-zero exit code.
func (m WarningsMode) DiagnosticFails(category string) bool {
	if category == "error" {
		return true
	}
	return category == "warning" && m == WarningsError
}
