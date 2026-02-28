package fanout_test

import (
	"testing"

	"github.com/glemzurg/go-complexity-lint/pkg/analyzer/fanout"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestFanout(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, fanout.Analyzer, "fanout")
}
