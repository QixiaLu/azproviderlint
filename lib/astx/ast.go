// Package astx defines AST helper functions that are reused in multiple checks.
package astx

import (
	"go/ast"
	"go/constant"
	"go/token"
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

// UnsafeToMovePast reports whether evaluating expr after stmts could yield a different value:
// a statement writes to or takes the address of a variable expr reads — directly (`y = v`) or
// as the root of an index, field, slice, or dereference target (`column[y] = v` writes column)
// — or declares a name expr references, so moving expr past it would rebind the name to the
// shadowing variable.
func UnsafeToMovePast(pass *analysis.Pass, expr ast.Expr, stmts []ast.Stmt) bool {
	operands := map[types.Object]bool{}
	names := map[string]bool{}
	ast.Inspect(expr, func(n ast.Node) bool {
		if id, ok := n.(*ast.Ident); ok {
			if v, ok := pass.TypesInfo.Uses[id].(*types.Var); ok {
				operands[v] = true
				names[id.Name] = true
			}
		}
		return true
	})
	if len(operands) == 0 {
		return false
	}

	writes := func(target ast.Expr) bool {
		id := rootIdent(target)
		return id != nil && operands[pass.TypesInfo.Uses[id]]
	}

	unsafe := false
	for _, stmt := range stmts {
		// only a declaration in this same statement list can shadow a name at the moved
		// expression's new position; declarations in nested blocks end with their block
		if shadows(stmt, names) {
			return true
		}
		ast.Inspect(stmt, func(n ast.Node) bool {
			switch x := n.(type) {
			case *ast.AssignStmt:
				for _, lhs := range x.Lhs {
					unsafe = unsafe || writes(lhs)
				}
			case *ast.IncDecStmt:
				unsafe = unsafe || writes(x.X)
			case *ast.UnaryExpr:
				// a taken address may be written through by whatever receives it
				if x.Op == token.AND {
					unsafe = unsafe || writes(x.X)
				}
			case *ast.RangeStmt:
				if x.Tok == token.ASSIGN {
					unsafe = unsafe || writes(x.Key) || writes(x.Value)
				}
			}
			return !unsafe
		})
		if unsafe {
			return true
		}
	}
	return false
}

// shadows reports whether stmt itself declares any of names — a `x := ...` or `var x ...`
// whose new variable would capture references in an expression moved below it.
func shadows(stmt ast.Stmt, names map[string]bool) bool {
	switch d := stmt.(type) {
	case *ast.AssignStmt:
		if d.Tok == token.DEFINE {
			for _, lhs := range d.Lhs {
				if id, ok := lhs.(*ast.Ident); ok && names[id.Name] {
					return true
				}
			}
		}
	case *ast.DeclStmt:
		if gd, ok := d.Decl.(*ast.GenDecl); ok && gd.Tok == token.VAR {
			for _, spec := range gd.Specs {
				if vs, ok := spec.(*ast.ValueSpec); ok {
					for _, id := range vs.Names {
						if names[id.Name] {
							return true
						}
					}
				}
			}
		}
	}
	return false
}

// rootIdent returns the base identifier an lvalue writes through: column for `column[y]`,
// out for `out.Format`, p for `*p`; nil when the base is not an identifier.
func rootIdent(e ast.Expr) *ast.Ident {
	for {
		switch x := e.(type) {
		case *ast.Ident:
			return x
		case *ast.IndexExpr:
			e = x.X
		case *ast.IndexListExpr:
			e = x.X
		case *ast.SelectorExpr:
			e = x.X
		case *ast.SliceExpr:
			e = x.X
		case *ast.StarExpr:
			e = x.X
		case *ast.ParenExpr:
			e = x.X
		default:
			return nil
		}
	}
}

// EnclosingFile returns the *ast.File containing pos.
func EnclosingFile(pass *analysis.Pass, pos token.Pos) *ast.File {
	for _, file := range pass.Files {
		if file.FileStart <= pos && pos < file.FileEnd {
			return file
		}
	}
	return nil
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
