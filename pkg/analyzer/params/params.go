package params

import (
	"fmt"
	"go/ast"

	"github.com/glemzurg/go-complexity-lint/pkg/analyzer/common"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var Analyzer = &analysis.Analyzer{
	Name: "params",
	Doc: "reports functions with too many parameters\n\n" +
		"Counts the number of parameters in a function signature, " +
		"properly handling grouped parameters like func(a, b int). " +
		"A leading ctx context.Context parameter is not counted.",
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

var (
	warnAt int
	failAt int
)

func init() {
	Analyzer.Flags.IntVar(&warnAt, "warn", 5,
		"parameter count at or above this triggers a warning (yellow zone)")
	Analyzer.Flags.IntVar(&failAt, "fail", 7,
		"parameter count at or above this triggers a failure (red zone)")
	Analyzer.Flags.StringVar(&common.ExcludePatterns, "exclude", "",
		"comma-separated filename glob patterns to skip (e.g. *_gen.go)")
}

func run(pass *analysis.Pass) (any, error) {
	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	defaults := common.Thresholds{WarnAt: warnAt, FailAt: failAt}
	if err := defaults.Validate("params"); err != nil {
		return nil, err
	}

	nodeFilter := []ast.Node{(*ast.FuncDecl)(nil)}

	insp.Preorder(nodeFilter, func(n ast.Node) {
		funcDecl := n.(*ast.FuncDecl)
		if funcDecl.Type == nil {
			return
		}
		if common.IsExcluded(pass.Fset.Position(funcDecl.Pos()).Filename) {
			return
		}

		funcName := common.FuncName(funcDecl)
		thresholds := common.ParseOverrides(funcDecl, "params", defaults)

		paramCount := countParams(funcDecl.Type)
		zone := thresholds.Classify(paramCount)

		if zone == common.ZoneGreen {
			return
		}

		pass.Report(analysis.Diagnostic{
			Pos:      funcDecl.Pos(),
			Category: zone.Category(),
			Message: fmt.Sprintf(
				"function %s has %d parameters (warn: >=%d, fail: >=%d) [%s] "+
					"(reduce by grouping coherently related subsets of parameters into structs — do not simply wrap all params into a single struct)",
				funcName, paramCount, thresholds.WarnAt, thresholds.FailAt,
				zone.Category()),
		})
	})

	return nil, nil
}

// countParams counts the total number of parameters, handling grouped params.
// func(a, b int, c string) has 3 params despite 2 field entries.
// Idiomatic ctx context.Context parameters are excluded from the count.
func countParams(funcType *ast.FuncType) int {
	if funcType.Params == nil {
		return 0
	}
	count := 0
	for _, field := range funcType.Params.List {
		if len(field.Names) == 0 {
			// Unnamed parameter (e.g., in interface method signatures).
			count++
			continue
		}
		for _, name := range field.Names {
			if common.IsExemptContextParam(name.Name, field.Type) {
				continue
			}
			count++
		}
	}
	return count
}
