package common

import "go/ast"

// IsExemptContextParam reports whether a parameter is the idiomatic Go context
// argument: named ctx with type context.Context. Such parameters are omitted
// from params complexity counts because they are boilerplate, not decision load.
func IsExemptContextParam(name string, typ ast.Expr) bool {
	return name == "ctx" && isContextContextType(typ)
}

func isContextContextType(typ ast.Expr) bool {
	sel, ok := typ.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	pkg, ok := sel.X.(*ast.Ident)
	if !ok || pkg.Name != "context" {
		return false
	}
	return sel.Sel.Name == "Context"
}
