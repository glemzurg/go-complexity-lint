package common

import (
	"go/ast"
	"go/token"
)

// IsErrGuard reports whether an if statement is an idiomatic Go error guard clause.
// An if statement is an error guard when ALL of these are true:
//  1. The condition is `<ident> != nil` (or `nil != <ident>`)
//  2. There is no else clause
//  3. The body contains exactly one statement
//  4. That statement is a return
//  5. All return values except the last are zero-value expressions
//  6. The last return value is either the same identifier or a function call
func IsErrGuard(ifStmt *ast.IfStmt) bool {
	errName := identNotNil(ifStmt.Cond)
	if errName == "" {
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
	return isErrReturn(retStmt.Results, errName)
}

// identNotNil checks if the condition is `<ident> != nil` (or `nil != <ident>`).
// Returns the identifier name, or "" if the condition doesn't match.
func identNotNil(cond ast.Expr) string {
	binExpr, ok := cond.(*ast.BinaryExpr)
	if !ok || binExpr.Op != token.NEQ {
		return ""
	}

	xIdent, xIsIdent := binExpr.X.(*ast.Ident)
	yIdent, yIsIdent := binExpr.Y.(*ast.Ident)

	// <ident> != nil
	if xIsIdent && yIsIdent && yIdent.Name == "nil" {
		return xIdent.Name
	}
	// nil != <ident>
	if xIsIdent && xIdent.Name == "nil" && yIsIdent {
		return yIdent.Name
	}
	return ""
}

// isErrReturn checks that all return values except the last are zero-value
// expressions, and the last is either the error identifier or a function call.
func isErrReturn(results []ast.Expr, errName string) bool {
	last := results[len(results)-1]

	// Check the last return value is the error variable or a function call.
	switch v := last.(type) {
	case *ast.Ident:
		if v.Name != errName {
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
