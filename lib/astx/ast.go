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

// IsUseOf reports whether expr is (modulo parentheses) an identifier resolving to obj.
func IsUseOf(pass *analysis.Pass, expr ast.Expr, obj types.Object) bool {
	id, ok := ast.Unparen(expr).(*ast.Ident)
	return ok && pass.TypesInfo.Uses[id] == obj
}

// UseCount counts references to obj within node.
func UseCount(pass *analysis.Pass, node ast.Node, obj types.Object) int {
	count := 0
	ast.Inspect(node, func(n ast.Node) bool {
		if id, ok := n.(*ast.Ident); ok && pass.TypesInfo.Uses[id] == obj {
			count++
		}
		return true
	})
	return count
}

// SourceText returns the raw source bytes of node.
func SourceText(pass *analysis.Pass, node ast.Node) ([]byte, bool) {
	tf := pass.Fset.File(node.Pos())
	if tf == nil {
		return nil, false
	}
	content, err := pass.ReadFile(tf.Name())
	if err != nil {
		return nil, false
	}
	start, end := tf.Offset(node.Pos()), tf.Offset(node.End())
	if start < 0 || end > len(content) || start > end {
		return nil, false
	}
	return content[start:end], true
}
