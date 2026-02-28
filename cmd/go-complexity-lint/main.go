package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/glemzurg/go-complexity-lint/pkg/analyzer/cyclo"
	"github.com/glemzurg/go-complexity-lint/pkg/analyzer/fanout"
	"github.com/glemzurg/go-complexity-lint/pkg/analyzer/nestdepth"
	"github.com/glemzurg/go-complexity-lint/pkg/analyzer/params"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/checker"
	"golang.org/x/tools/go/packages"
)

func main() {
	progname := filepath.Base(os.Args[0])
	log.SetFlags(0)
	log.SetPrefix(progname + ": ")

	analyzers := []*analysis.Analyzer{
		nestdepth.Analyzer,
		cyclo.Analyzer,
		params.Analyzer,
		fanout.Analyzer,
	}

	if err := analysis.Validate(analyzers); err != nil {
		log.Fatal(err)
	}

	// Register analyzer flags with namespace prefix (e.g., cyclo.warn).
	for _, a := range analyzers {
		a.Flags.VisitAll(func(f *flag.Flag) {
			name := a.Name + "." + f.Name
			flag.Var(f.Value, name, f.Usage)
		})
	}
	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, `%[1]s is a complexity linter for Go programs.

Usage: %[1]s [-flag] [package]

Analyzers:
  nestdepth   reports functions with deep nesting
  cyclo       reports functions with high cyclomatic complexity
  params      reports functions with too many parameters
  fanout      reports functions with high fan-out

Flags are namespaced by analyzer, e.g.:
  -nestdepth.warn=4  -nestdepth.fail=6
  -cyclo.warn=9      -cyclo.fail=14
  -params.warn=4     -params.fail=6
  -fanout.warn=6     -fanout.fail=9

Exit codes:
  0  no red-zone violations (warnings may be present)
  1  red-zone violations found or analysis error
`, progname)
		os.Exit(1)
	}

	// Load packages.
	cfg := &packages.Config{
		Mode: packages.LoadAllSyntax,
	}
	pkgs, err := packages.Load(cfg, args...)
	if err != nil {
		log.Fatal(err)
	}
	if n := packages.PrintErrors(pkgs); n > 0 {
		os.Exit(1)
	}

	// Run analysis.
	graph, err := checker.Analyze(analyzers, pkgs, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Print diagnostics.
	if err := graph.PrintText(os.Stderr, -1); err != nil {
		log.Fatal(err)
	}

	// Exit 1 only if any diagnostic has "error" category (red zone).
	for act := range graph.All() {
		if act.Err != nil {
			os.Exit(1)
		}
		if act.IsRoot {
			for _, diag := range act.Diagnostics {
				if diag.Category == "error" {
					os.Exit(1)
				}
			}
		}
	}
}
