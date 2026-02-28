package nestdepth

import (
	"fmt"
	"go/ast"
	"go/token"

	"github.com/glemzurg/go-complexity-lint/pkg/analyzer/common"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var Analyzer = &analysis.Analyzer{
	Name: "nestdepth",
	Doc: "reports functions with excessive nesting depth\n\n" +
		"Nesting depth is incremented by control-flow constructs: " +
		"if/else, for/range, switch, type switch, select, case/default, " +
		"and anonymous function literals. Error guard clauses are exempt.",
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

var (
	warnAt int
	failAt int
)

func init() {
	Analyzer.Flags.IntVar(&warnAt, "warn", 4,
		"nesting depth above this triggers a warning (yellow zone)")
	Analyzer.Flags.IntVar(&failAt, "fail", 6,
		"nesting depth above this triggers a failure (red zone)")
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
		thresholds := common.ParseOverrides(funcDecl, "nestdepth", defaults)

		deepestPos, deepestDepth := walkBody(funcDecl.Body, 0)
		zone := thresholds.Classify(deepestDepth)

		if zone == common.ZoneGreen {
			return
		}

		pass.Report(analysis.Diagnostic{
			Pos:      deepestPos,
			Category: zone.Category(),
			Message: fmt.Sprintf(
				"function %s has a nesting depth of %d (warn: >%d, fail: >%d) [%s]",
				funcName, deepestDepth, thresholds.WarnAt, thresholds.FailAt,
				zone.Category()),
		})
	})

	return nil, nil
}

func walkBody(block *ast.BlockStmt, currentDepth int) (token.Pos, int) {
	deepestDepth := currentDepth
	deepestPos := block.Lbrace

	for _, stmt := range block.List {
		pos, depth := walkStmt(stmt, currentDepth)
		if depth > deepestDepth {
			deepestDepth = depth
			deepestPos = pos
		}
	}
	return deepestPos, deepestDepth
}

func walkStmt(stmt ast.Stmt, currentDepth int) (token.Pos, int) {
	deepestDepth := currentDepth
	deepestPos := stmt.Pos()

	switch s := stmt.(type) {
	case *ast.IfStmt:
		// Error guard clauses don't count as nesting.
		if common.IsErrGuard(s) {
			return deepestPos, deepestDepth
		}

		if s.Init != nil {
			pos, depth := walkNodeForFuncLit(s.Init, currentDepth)
			if depth > deepestDepth {
				deepestDepth = depth
				deepestPos = pos
			}
		}

		newDepth := currentDepth + 1
		pos, depth := walkBody(s.Body, newDepth)
		if depth > deepestDepth {
			deepestDepth = depth
			deepestPos = pos
		}

		if s.Else != nil {
			pos, depth := walkElse(s.Else, currentDepth)
			if depth > deepestDepth {
				deepestDepth = depth
				deepestPos = pos
			}
		}

	case *ast.ForStmt:
		pos, depth := walkBody(s.Body, currentDepth+1)
		if depth > deepestDepth {
			deepestDepth = depth
			deepestPos = pos
		}

	case *ast.RangeStmt:
		pos, depth := walkBody(s.Body, currentDepth+1)
		if depth > deepestDepth {
			deepestDepth = depth
			deepestPos = pos
		}

	case *ast.SwitchStmt:
		pos, depth := walkBody(s.Body, currentDepth+1)
		if depth > deepestDepth {
			deepestDepth = depth
			deepestPos = pos
		}

	case *ast.TypeSwitchStmt:
		pos, depth := walkBody(s.Body, currentDepth+1)
		if depth > deepestDepth {
			deepestDepth = depth
			deepestPos = pos
		}

	case *ast.SelectStmt:
		pos, depth := walkBody(s.Body, currentDepth+1)
		if depth > deepestDepth {
			deepestDepth = depth
			deepestPos = pos
		}

	case *ast.CaseClause:
		newDepth := currentDepth + 1
		for _, bodyStmt := range s.Body {
			pos, depth := walkStmt(bodyStmt, newDepth)
			if depth > deepestDepth {
				deepestDepth = depth
				deepestPos = pos
			}
		}

	case *ast.CommClause:
		newDepth := currentDepth + 1
		for _, bodyStmt := range s.Body {
			pos, depth := walkStmt(bodyStmt, newDepth)
			if depth > deepestDepth {
				deepestDepth = depth
				deepestPos = pos
			}
		}

	case *ast.BlockStmt:
		for _, bodyStmt := range s.List {
			pos, depth := walkStmt(bodyStmt, currentDepth)
			if depth > deepestDepth {
				deepestDepth = depth
				deepestPos = pos
			}
		}

	case *ast.LabeledStmt:
		pos, depth := walkStmt(s.Stmt, currentDepth)
		if depth > deepestDepth {
			deepestDepth = depth
			deepestPos = pos
		}

	default:
		pos, depth := walkNodeForFuncLit(stmt, currentDepth)
		if depth > deepestDepth {
			deepestDepth = depth
			deepestPos = pos
		}
	}

	return deepestPos, deepestDepth
}

func walkElse(elseNode ast.Stmt, currentDepth int) (token.Pos, int) {
	switch e := elseNode.(type) {
	case *ast.BlockStmt:
		return walkBody(e, currentDepth+1)
	case *ast.IfStmt:
		return walkStmt(e, currentDepth)
	default:
		return elseNode.Pos(), currentDepth
	}
}

func walkNodeForFuncLit(node ast.Node, currentDepth int) (token.Pos, int) {
	deepestPos := node.Pos()
	deepestDepth := currentDepth

	ast.Inspect(node, func(n ast.Node) bool {
		fl, ok := n.(*ast.FuncLit)
		if !ok {
			return true
		}
		if fl.Body != nil {
			pos, depth := walkBody(fl.Body, currentDepth+1)
			if depth > deepestDepth {
				deepestDepth = depth
				deepestPos = pos
			}
		}
		return false
	})

	return deepestPos, deepestDepth
}
