package plugin

import (
	"github.com/glemzurg/go-complexity-lint/pkg/analyzer/cyclo"
	"github.com/glemzurg/go-complexity-lint/pkg/analyzer/fanout"
	"github.com/glemzurg/go-complexity-lint/pkg/analyzer/nestdepth"
	"github.com/glemzurg/go-complexity-lint/pkg/analyzer/params"
	"github.com/golangci/plugin-module-register/register"
	"golang.org/x/tools/go/analysis"
)

func init() {
	register.Plugin("go-complexity-lint", New)
}

func New(_ any) (register.LinterPlugin, error) {
	return &complexityPlugin{}, nil
}

type complexityPlugin struct{}

func (p *complexityPlugin) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{
		nestdepth.Analyzer,
		cyclo.Analyzer,
		params.Analyzer,
		fanout.Analyzer,
	}, nil
}

func (p *complexityPlugin) GetLoadMode() string {
	return register.LoadModeTypesInfo
}
