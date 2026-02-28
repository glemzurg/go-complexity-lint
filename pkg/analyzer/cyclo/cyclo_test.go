package cyclo_test

import (
	"testing"

	"github.com/glemzurg/go-complexity-lint/pkg/analyzer/cyclo"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestCyclo(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, cyclo.Analyzer, "cyclo")
}
