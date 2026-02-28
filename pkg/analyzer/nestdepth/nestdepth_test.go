package nestdepth_test

import (
	"testing"

	"github.com/glemzurg/go-complexity-lint/pkg/analyzer/nestdepth"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestNestDepth(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, nestdepth.Analyzer, "nestdepth")
}
