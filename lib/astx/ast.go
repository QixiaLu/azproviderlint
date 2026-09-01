// Package astx defines AST helper functions that are reused in multiple checks.
package astx

import (
	"go/ast"
	"go/constant"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

func IsTrueConstant(pass *analysis.Pass, e ast.Expr) bool {
	if e == nil {
		return false
	}
	value := pass.TypesInfo.Types[e].Value
	return value != nil && value.Kind() == constant.Bool && constant.BoolVal(value)
}

// ContainsCallExpr reports whether expr contains any function call.
func ContainsCallExpr(expr ast.Expr) bool {
	found := false
	ast.Inspect(expr, func(n ast.Node) bool {
		if _, ok := n.(*ast.CallExpr); ok {
			found = true
			return false
		}
		return true
	})
	return found
}

// IsNilValue reports whether expr denotes a nil value: the predeclared nil identifier,
// possibly wrapped in parentheses or type conversions (`error(nil)`, `[]T(nil)`).
func IsNilValue(pass *analysis.Pass, expr ast.Expr) bool {
	switch e := ast.Unparen(expr).(type) {
	case *ast.Ident:
		_, isNil := pass.TypesInfo.Uses[e].(*types.Nil)
		return isNil
	case *ast.CallExpr:
		if len(e.Args) == 1 && pass.TypesInfo.Types[e.Fun].IsType() {
			return IsNilValue(pass, e.Args[0])
		}
	}
	return false
}
