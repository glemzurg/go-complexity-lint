package params_test

import (
	"testing"

	"github.com/glemzurg/go-complexity-lint/pkg/analyzer/params"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestParams(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, params.Analyzer, "params")
}
