package cyclo

import (
	"fmt"
	"go/ast"

	"github.com/glemzurg/go-complexity-lint/pkg/analyzer/common"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var Analyzer = &analysis.Analyzer{
	Name: "cyclo",
	Doc: "reports functions with high cyclomatic complexity\n\n" +
		"Cyclomatic complexity is 1 + 1 for each branching/looping decision " +
		"(if, for, range, case). Else clauses, default, and boolean operators " +
		"do not count. Error guard clauses are exempt.",
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

var (
	warnAt int
	failAt int
)

func init() {
	Analyzer.Flags.IntVar(&warnAt, "warn", 9,
		"cyclomatic complexity above this triggers a warning (yellow zone)")
	Analyzer.Flags.IntVar(&failAt, "fail", 14,
		"cyclomatic complexity above this triggers a failure (red zone)")
}

func run(pass *analysis.Pass) (any, error) {
	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	defaults := common.Thresholds{WarnAt: warnAt, FailAt: failAt}

	nodeFilter := []ast.Node{(*ast.FuncDecl)(nil)}

	insp.Preorder(nodeFilter, func(n ast.Node) {
		funcDecl := n.(*ast.FuncDecl)
		if funcDecl.Body == nil {
			return
		}

		funcName := common.FuncName(funcDecl)
		thresholds := common.ParseOverrides(funcDecl, "cyclo", defaults)

		complexity := calcComplexity(funcDecl.Body)
		zone := thresholds.Classify(complexity)

		if zone == common.ZoneGreen {
			return
		}

		pass.Report(analysis.Diagnostic{
			Pos:      funcDecl.Pos(),
			Category: zone.Category(),
			Message: fmt.Sprintf(
				"function %s has cyclomatic complexity of %d (warn: >%d, fail: >%d) [%s]",
				funcName, complexity, thresholds.WarnAt, thresholds.FailAt,
				zone.Category()),
		})
	})

	return nil, nil
}

// calcComplexity computes the cyclomatic complexity of a function body.
// Base complexity is 1. Each branching/looping decision adds 1.
func calcComplexity(body *ast.BlockStmt) int {
	complexity := 1

	ast.Inspect(body, func(n ast.Node) bool {
		switch s := n.(type) {
		case *ast.IfStmt:
			// Error guard clauses are exempt.
			if common.IsErrGuard(s) {
				return false
			}
			complexity++
		case *ast.ForStmt:
			complexity++
		case *ast.RangeStmt:
			complexity++
		case *ast.CaseClause:
			// Each case is a decision (not default).
			if s.List != nil {
				complexity++
			}
		case *ast.CommClause:
			// Each select case is a decision (not default).
			if s.Comm != nil {
				complexity++
			}
		}
		return true
	})

	return complexity
}
