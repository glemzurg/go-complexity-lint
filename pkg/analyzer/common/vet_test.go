package common

import (
	"flag"
	"testing"

	"golang.org/x/tools/go/analysis"
)

func TestConfigureRedZoneOnly(t *testing.T) {
	var warnAt, failAt int

	analyzer := &analysis.Analyzer{Name: "testmetric"}
	analyzer.Flags.Init("testmetric", flag.ExitOnError)
	analyzer.Flags.IntVar(&warnAt, "warn", 5, "warning threshold")
	analyzer.Flags.IntVar(&failAt, "fail", 7, "failure threshold")

	ConfigureRedZoneOnly([]*analysis.Analyzer{analyzer})

	if warnAt != 7 {
		t.Fatalf("warnAt = %d, want 7", warnAt)
	}
	if failAt != 7 {
		t.Fatalf("failAt = %d, want 7", failAt)
	}
}
