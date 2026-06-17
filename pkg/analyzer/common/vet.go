package common

import "golang.org/x/tools/go/analysis"

// ConfigureRedZoneOnly sets each analyzer's warn threshold to its fail
// threshold so only red-zone violations are reported. go vet invokes the
// tool through unitchecker, which has no -warnings flag; red-zone-only
// is the vet default. Explicit -metric.warn flags still override.
func ConfigureRedZoneOnly(analyzers []*analysis.Analyzer) {
	for _, a := range analyzers {
		fail := a.Flags.Lookup("fail")
		if fail == nil {
			continue
		}
		_ = a.Flags.Set("warn", fail.Value.String())
	}
}
