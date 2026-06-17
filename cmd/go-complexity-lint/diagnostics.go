package main

import (
	"fmt"
	"io"

	"github.com/glemzurg/go-complexity-lint/pkg/analyzer/common"
	"go/token"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/checker"
)

// printDiagnostics emits filtered diagnostics and returns whether the run
// should exit with a non-zero status.
func printDiagnostics(w io.Writer, graph *checker.Graph, mode common.WarningsMode) (bool, error) {
	type key struct {
		pos token.Position
		end token.Position
		*analysis.Analyzer
		message string
	}
	seen := make(map[key]bool)

	failed := false

	for act := range graph.All() {
		if act.Err != nil {
			if _, err := fmt.Fprintf(w, "%s: %v\n", act.Analyzer.Name, act.Err); err != nil {
				return failed, err
			}
			failed = true
			continue
		}
		if !act.IsRoot {
			continue
		}

		for _, diag := range act.Diagnostics {
			if !mode.ReportDiagnostic(diag.Category) {
				continue
			}
			if mode.DiagnosticFails(diag.Category) {
				failed = true
			}

			posn := act.Package.Fset.Position(diag.Pos)
			end := act.Package.Fset.Position(diag.End)
			k := key{posn, end, act.Analyzer, diag.Message}
			if seen[k] {
				continue
			}
			seen[k] = true

			if _, err := fmt.Fprintf(w, "%s: %s\n", posn, diag.Message); err != nil {
				return failed, err
			}
		}
	}

	return failed, nil
}
