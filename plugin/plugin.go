package plugin

import (
	"fmt"

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

type Settings struct {
	NestdepthWarn *int `json:"nestdepth-warn"`
	NestdepthFail *int `json:"nestdepth-fail"`
	CycloWarn     *int `json:"cyclo-warn"`
	CycloFail     *int `json:"cyclo-fail"`
	ParamsWarn    *int `json:"params-warn"`
	ParamsFail    *int `json:"params-fail"`
	FanoutWarn    *int `json:"fanout-warn"`
	FanoutFail    *int `json:"fanout-fail"`
}

func New(conf any) (register.LinterPlugin, error) {
	s, err := register.DecodeSettings[Settings](conf)
	if err != nil {
		return nil, err
	}
	return &complexityPlugin{settings: s}, nil
}

type complexityPlugin struct {
	settings Settings
}

func (p *complexityPlugin) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	analyzers := []*analysis.Analyzer{
		nestdepth.Analyzer,
		cyclo.Analyzer,
		params.Analyzer,
		fanout.Analyzer,
	}

	flagOverrides := []struct {
		analyzer *analysis.Analyzer
		warn     *int
		fail     *int
	}{
		{nestdepth.Analyzer, p.settings.NestdepthWarn, p.settings.NestdepthFail},
		{cyclo.Analyzer, p.settings.CycloWarn, p.settings.CycloFail},
		{params.Analyzer, p.settings.ParamsWarn, p.settings.ParamsFail},
		{fanout.Analyzer, p.settings.FanoutWarn, p.settings.FanoutFail},
	}

	for _, o := range flagOverrides {
		if o.warn != nil {
			if err := o.analyzer.Flags.Set("warn", fmt.Sprint(*o.warn)); err != nil {
				return nil, fmt.Errorf("setting %s.warn: %w", o.analyzer.Name, err)
			}
		}
		if o.fail != nil {
			if err := o.analyzer.Flags.Set("fail", fmt.Sprint(*o.fail)); err != nil {
				return nil, fmt.Errorf("setting %s.fail: %w", o.analyzer.Name, err)
			}
		}
	}

	return analyzers, nil
}

func (p *complexityPlugin) GetLoadMode() string {
	return register.LoadModeTypesInfo
}
