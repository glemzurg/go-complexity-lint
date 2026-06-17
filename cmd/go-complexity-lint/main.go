package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/glemzurg/go-complexity-lint/pkg/analyzer/common"
	"github.com/glemzurg/go-complexity-lint/pkg/analyzer/cyclo"
	"github.com/glemzurg/go-complexity-lint/pkg/analyzer/fanout"
	"github.com/glemzurg/go-complexity-lint/pkg/analyzer/nestdepth"
	"github.com/glemzurg/go-complexity-lint/pkg/analyzer/params"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/checker"
	"golang.org/x/tools/go/analysis/unitchecker"
	"golang.org/x/tools/go/packages"
)

//complexity:cyclo:warn=20,fail=25
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

	// When invoked by "go vet -vettool", delegate to unitchecker
	// which handles the -flags, -V=full, and *.cfg protocol.
	if len(os.Args) > 1 {
		arg := os.Args[1]
		if arg == "-flags" || strings.HasPrefix(arg, "-V=") || strings.HasSuffix(arg, ".cfg") {
			// go vet has no -warnings flag; default to red-zone-only reporting.
			common.ConfigureRedZoneOnly(analyzers)
			unitchecker.Main(analyzers...)
		}
	}

	if err := analysis.Validate(analyzers); err != nil {
		log.Fatal(err)
	}

	var warningsMode common.WarningsMode

	// Register a global -exclude flag (shared across all analyzers).
	flag.StringVar(&common.ExcludePatterns, "exclude", "",
		"comma-separated filename glob patterns to skip (e.g. *_gen.go)")
	flag.Var(&warningsMode, "warnings",
		"warning handling: default (print, exit 0), none (suppress), error (print, exit 1)")

	// Register analyzer flags with namespace prefix (e.g., cyclo.warn).
	// Also register hyphen-separated aliases (e.g., cyclo-warn) for convenience.
	for _, a := range analyzers {
		a.Flags.VisitAll(func(f *flag.Flag) {
			if f.Name == "exclude" {
				return // covered by the global -exclude flag
			}
			name := a.Name + "." + f.Name
			flag.Var(f.Value, name, f.Usage)
			alias := a.Name + "-" + f.Name
			flag.Var(f.Value, alias, f.Usage)
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

Flags are namespaced by analyzer (dot or hyphen separator). The warn/fail
values are inclusive lower bounds (a value at or above the threshold triggers
the zone). Defaults shown:
  -nestdepth.warn=5  -nestdepth.fail=7
  -cyclo.warn=10     -cyclo.fail=15
  -params.warn=5     -params.fail=7
  -fanout.warn=7     -fanout.fail=10

Hyphen-separated aliases also work:
  -cyclo-warn=10     -cyclo-fail=15

  -exclude="*_gen.go,mock_*.go"  skip files matching glob patterns

  -warnings=default  print warnings, exit 0 when only warnings are present
  -warnings=none     suppress warning output, exit 0 when only warnings are present
  -warnings=error    print warnings, exit 1 when any warning is present

Exit codes:
  0  no failing diagnostics under the selected -warnings mode
  1  failing diagnostics found or analysis error
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

	failed, err := printDiagnostics(os.Stderr, graph, warningsMode)
	if err != nil {
		log.Fatal(err)
	}
	if failed {
		os.Exit(1)
	}
}
