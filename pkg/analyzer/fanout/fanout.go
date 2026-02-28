package fanout

import (
	"fmt"
	"go/ast"
	"go/types"
	"strings"

	"github.com/glemzurg/go-complexity-lint/pkg/analyzer/common"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var Analyzer = &analysis.Analyzer{
	Name: "fanout",
	Doc: "reports functions with too many distinct function calls (fan out)\n\n" +
		"Counts unique non-builtin, non-stdlib function/method calls in a function. " +
		"The same function called multiple times counts as 1.",
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

var (
	warnAt int
	failAt int
)

func init() {
	Analyzer.Flags.IntVar(&warnAt, "warn", 6,
		"fan out count above this triggers a warning (yellow zone)")
	Analyzer.Flags.IntVar(&failAt, "fail", 9,
		"fan out count above this triggers a failure (red zone)")
}

func run(pass *analysis.Pass) (any, error) {
	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	defaults := common.Thresholds{WarnAt: warnAt, FailAt: failAt}
	if err := defaults.Validate("fanout"); err != nil {
		return nil, err
	}

	nodeFilter := []ast.Node{(*ast.FuncDecl)(nil)}

	insp.Preorder(nodeFilter, func(n ast.Node) {
		funcDecl := n.(*ast.FuncDecl)
		if funcDecl.Body == nil {
			return
		}

		funcName := common.FuncName(funcDecl)
		thresholds := common.ParseOverrides(funcDecl, "fanout", defaults)

		distinctCalls := countDistinctCalls(pass, funcDecl.Body)
		zone := thresholds.Classify(distinctCalls)

		if zone == common.ZoneGreen {
			return
		}

		pass.Report(analysis.Diagnostic{
			Pos:      funcDecl.Pos(),
			Category: zone.Category(),
			Message: fmt.Sprintf(
				"function %s has fan out of %d (warn: >%d, fail: >%d) [%s]",
				funcName, distinctCalls, thresholds.WarnAt, thresholds.FailAt,
				zone.Category()),
		})
	})

	return nil, nil
}

// countDistinctCalls counts the number of distinct non-builtin, non-stdlib
// function/method calls in a function body.
func countDistinctCalls(pass *analysis.Pass, body *ast.BlockStmt) int {
	seen := make(map[types.Object]bool)

	ast.Inspect(body, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}

		var obj types.Object

		switch fn := call.Fun.(type) {
		case *ast.Ident:
			obj = pass.TypesInfo.ObjectOf(fn)
		case *ast.SelectorExpr:
			obj = pass.TypesInfo.ObjectOf(fn.Sel)
		default:
			return true
		}

		if obj == nil {
			return true
		}

		// Exclude builtins (len, cap, make, etc.).
		if _, isBuiltin := obj.(*types.Builtin); isBuiltin {
			return true
		}

		// Exclude type conversions (int(x), string(b), etc.).
		if _, isTypeName := obj.(*types.TypeName); isTypeName {
			return true
		}

		// Exclude standard library functions.
		if pkg := obj.Pkg(); pkg != nil && isStdlib(pkg.Path()) {
			return true
		}

		seen[obj] = true
		return true
	})

	return len(seen)
}

// isStdlib reports whether a package path belongs to the Go standard library.
// Standard library packages have no dots in their path.
func isStdlib(pkgPath string) bool {
	return !strings.Contains(pkgPath, ".")
}
