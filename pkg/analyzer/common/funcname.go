package common

import "go/ast"

// FuncName returns the qualified name of a function declaration,
// including the receiver type for methods.
func FuncName(funcDecl *ast.FuncDecl) string {
	name := funcDecl.Name.Name
	if funcDecl.Recv != nil && len(funcDecl.Recv.List) > 0 {
		recvType := ExprName(funcDecl.Recv.List[0].Type)
		name = recvType + "." + name
	}
	return name
}

// ExprName extracts a human-readable name from a type expression.
func ExprName(expr ast.Expr) string {
	switch e := expr.(type) {
	case *ast.Ident:
		return e.Name
	case *ast.StarExpr:
		return "*" + ExprName(e.X)
	case *ast.IndexExpr:
		return ExprName(e.X)
	case *ast.IndexListExpr:
		return ExprName(e.X)
	default:
		return "?"
	}
}
