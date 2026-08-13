// Package AZG004 defines an analyzer that reports a zero-value initialization immediately
// followed by a nil check and pointer dereference — `y := <zero>; if x != nil { y = *x }` —
// where the generic pointer.From(x) helper should be used instead.
package AZG004

import (
	"go/ast"
	"go/constant"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// Analyzer checks for the manual `y := <zero>; if x != nil { y = *x }` idiom that reassigns a
// zero-initialised variable from a dereferenced pointer only when it is non-nil. The generic
// `pointer.From(x)` helper returns the dereferenced value or the type's zero value when the
// pointer is nil, so those call sites should use it instead.
var Analyzer = &analysis.Analyzer{
	Name:     "AZG004",
	Doc:      "check for zero-value initialization followed by a nil check and pointer dereference that should use pointer.From",
	URL:      "https://github.com/katbyte/azproviderlint/blob/main/checks/AZG/AZG004_zero_value_init_pointer_from/README.md",
	Requires: []*analysis.Analyzer{inspect.Analyzer},
	Run:      run,
}

func run(pass *analysis.Pass) (any, error) {
	insp, ok := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	if !ok {
		return nil, nil
	}

	// Blocks and case clauses both hold statement lists where the two-statement pattern
	// (zero-init assignment followed by an if) can appear.
	nodeFilter := []ast.Node{
		(*ast.BlockStmt)(nil),
		(*ast.CaseClause)(nil),
	}

	insp.Preorder(nodeFilter, func(n ast.Node) {
		var stmts []ast.Stmt
		switch node := n.(type) {
		case *ast.BlockStmt:
			stmts = node.List
		case *ast.CaseClause:
			stmts = node.Body
		default:
			return
		}

		// Look for an assignment immediately followed by a matching if statement.
		for i := range len(stmts) - 1 {
			ifStmt, ok := stmts[i+1].(*ast.IfStmt)
			if !ok {
				continue
			}

			assignPos, varName := zeroValueInit(pass, stmts[i])
			if varName == "" {
				continue
			}

			if !matchingNilCheckAssignment(ifStmt, varName) {
				continue
			}

			pass.Report(analysis.Diagnostic{
				Pos:     assignPos,
				Message: "zero-value initialization followed by a nil check and pointer dereference can be simplified with pointer.From",
			})
		}
	})

	return nil, nil
}

// zeroValueInit reports whether stmt initialises a single variable to its zero value, returning
// the initialization position and the variable name. It matches both the short-declaration form
// (`varName := <zero>`) and the var-declaration form, covering all three of its variants:
// `var varName T` (implicitly zero), `var varName T = <zero>`, and `var varName = <zero>`.
func zeroValueInit(pass *analysis.Pass, stmt ast.Stmt) (pos token.Pos, varName string) {
	switch s := stmt.(type) {
	case *ast.AssignStmt:
		if s.Tok != token.DEFINE || len(s.Lhs) != 1 || len(s.Rhs) != 1 {
			return token.NoPos, ""
		}

		lhsIdent, ok := s.Lhs[0].(*ast.Ident)
		if !ok {
			return token.NoPos, ""
		}

		if !isZeroValue(pass, s.Rhs[0]) {
			return token.NoPos, ""
		}

		return s.Pos(), lhsIdent.Name

	case *ast.DeclStmt:
		genDecl, ok := s.Decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.VAR || len(genDecl.Specs) != 1 {
			return token.NoPos, ""
		}

		valueSpec, ok := genDecl.Specs[0].(*ast.ValueSpec)
		if !ok || len(valueSpec.Names) != 1 {
			return token.NoPos, ""
		}

		switch len(valueSpec.Values) {
		case 0:
			// `var varName T` with no initializer is implicitly the zero value; a type is
			// required (`var varName` alone is not valid Go).
			if valueSpec.Type == nil {
				return token.NoPos, ""
			}
		case 1:
			// `var varName T = <expr>` / `var varName = <expr>` must initialise to a zero value.
			if !isZeroValue(pass, valueSpec.Values[0]) {
				return token.NoPos, ""
			}
		default:
			return token.NoPos, ""
		}

		return s.Pos(), valueSpec.Names[0].Name
	}

	return token.NoPos, ""
}

// isZeroValue reports whether expr is a zero value: false, 0, "", or nil.
func isZeroValue(pass *analysis.Pass, expr ast.Expr) bool {
	tv, ok := pass.TypesInfo.Types[expr]
	if !ok {
		return false
	}

	if tv.Value != nil {
		switch tv.Value.Kind() {
		case constant.Bool:
			return !constant.BoolVal(tv.Value)
		case constant.String:
			return constant.StringVal(tv.Value) == ""
		case constant.Int, constant.Float:
			return constant.Sign(tv.Value) == 0
		case constant.Unknown, constant.Complex:
			return false
		}
	}

	return isNilIdent(expr)
}

// matchingNilCheckAssignment reports whether ifStmt is `if x != nil { varName = *x }` with no
// else branch, where the dereferenced expression matches the nil-checked one.
func matchingNilCheckAssignment(ifStmt *ast.IfStmt, varName string) bool {
	if ifStmt.Else != nil {
		return false
	}

	binExpr, ok := ifStmt.Cond.(*ast.BinaryExpr)
	if !ok || binExpr.Op != token.NEQ {
		return false
	}

	if !isNilIdent(binExpr.Y) {
		return false
	}

	checkedExpr := binExpr.X
	if containsCallExpr(checkedExpr) {
		return false
	}

	if len(ifStmt.Body.List) != 1 {
		return false
	}

	assignStmt, ok := ifStmt.Body.List[0].(*ast.AssignStmt)
	if !ok || assignStmt.Tok != token.ASSIGN {
		return false
	}

	if len(assignStmt.Lhs) != 1 || len(assignStmt.Rhs) != 1 {
		return false
	}

	lhsIdent, ok := assignStmt.Lhs[0].(*ast.Ident)
	if !ok || lhsIdent.Name != varName {
		return false
	}

	starExpr, ok := assignStmt.Rhs[0].(*ast.StarExpr)
	if !ok {
		return false
	}

	return types.ExprString(checkedExpr) == types.ExprString(starExpr.X)
}

// isNilIdent reports whether expr is the nil identifier.
func isNilIdent(expr ast.Expr) bool {
	ident, ok := expr.(*ast.Ident)
	return ok && ident.Name == "nil"
}

// containsCallExpr reports whether expr contains any function call, which would make the
// pointer.From rewrite non-equivalent because the manual idiom evaluates the expression twice.
func containsCallExpr(expr ast.Expr) bool {
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
