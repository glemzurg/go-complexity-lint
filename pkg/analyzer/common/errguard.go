package common

import (
	"go/ast"
	"go/token"
)

// IsErrGuard reports whether an if statement is an idiomatic Go error guard clause.
// An if statement is an error guard when ALL of these are true:
//  1. The condition is `err != nil` (init statement can be anything)
//  2. The body contains exactly one statement
//  3. That statement is a return
//  4. All return values except the last are zero-value expressions
//  5. The last return value is either `err` or a function call
func IsErrGuard(ifStmt *ast.IfStmt) bool {
	if !isErrNotNil(ifStmt.Cond) {
		return false
	}
	if ifStmt.Else != nil {
		return false
	}
	if len(ifStmt.Body.List) != 1 {
		return false
	}
	retStmt, ok := ifStmt.Body.List[0].(*ast.ReturnStmt)
	if !ok {
		return false
	}
	if len(retStmt.Results) == 0 {
		return false
	}
	return isErrReturn(retStmt.Results)
}

// isErrNotNil checks if the condition is `err != nil` (or `nil != err`).
func isErrNotNil(cond ast.Expr) bool {
	binExpr, ok := cond.(*ast.BinaryExpr)
	if !ok || binExpr.Op != token.NEQ {
		return false
	}

	xIdent, xIsIdent := binExpr.X.(*ast.Ident)
	yIdent, yIsIdent := binExpr.Y.(*ast.Ident)

	// err != nil
	if xIsIdent && xIdent.Name == "err" && yIsIdent && yIdent.Name == "nil" {
		return true
	}
	// nil != err
	if xIsIdent && xIdent.Name == "nil" && yIsIdent && yIdent.Name == "err" {
		return true
	}
	return false
}

// isErrReturn checks that all return values except the last are zero-value
// expressions, and the last is either `err` or a function call.
func isErrReturn(results []ast.Expr) bool {
	last := results[len(results)-1]

	// Check the last return value is `err` or a function call.
	switch v := last.(type) {
	case *ast.Ident:
		if v.Name != "err" {
			return false
		}
	case *ast.CallExpr:
		// Any function call is allowed (fmt.Errorf, errors.New, etc.)
	default:
		return false
	}

	// Check all preceding return values are zero-value expressions.
	for _, expr := range results[:len(results)-1] {
		if !isZeroValue(expr) {
			return false
		}
	}

	return true
}

// isZeroValue reports whether an expression is a zero-value literal:
// nil, 0, "", false, or Type{}.
func isZeroValue(expr ast.Expr) bool {
	switch v := expr.(type) {
	case *ast.Ident:
		return v.Name == "nil" || v.Name == "false"
	case *ast.BasicLit:
		switch v.Kind {
		case token.INT:
			return v.Value == "0"
		case token.FLOAT:
			return v.Value == "0.0" || v.Value == "0."
		case token.STRING:
			return v.Value == `""` || v.Value == "``"
		}
	case *ast.CompositeLit:
		// Type{} with no elements
		return len(v.Elts) == 0
	}
	return false
}
