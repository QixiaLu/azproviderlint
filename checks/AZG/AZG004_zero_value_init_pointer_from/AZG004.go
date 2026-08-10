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

			assignStmt, varName := zeroValueAssignment(pass, stmts[i])
			if assignStmt == nil {
				continue
			}

			if !matchingNilCheckAssignment(ifStmt, varName) {
				continue
			}

			pass.Report(analysis.Diagnostic{
				Pos:     assignStmt.Pos(),
				Message: "zero-value initialization followed by a nil check and pointer dereference can be simplified with pointer.From",
			})
		}
	})

	return nil, nil
}

// zeroValueAssignment reports whether stmt is a `varName := <zero-value>` short declaration,
// returning the assignment and the declared variable name when it matches.
func zeroValueAssignment(pass *analysis.Pass, stmt ast.Stmt) (assign *ast.AssignStmt, varName string) {
	assignStmt, ok := stmt.(*ast.AssignStmt)
	if !ok || assignStmt.Tok != token.DEFINE {
		return nil, ""
	}

	if len(assignStmt.Lhs) != 1 || len(assignStmt.Rhs) != 1 {
		return nil, ""
	}

	lhsIdent, ok := assignStmt.Lhs[0].(*ast.Ident)
	if !ok {
		return nil, ""
	}

	if !isZeroValue(pass, assignStmt.Rhs[0]) {
		return nil, ""
	}

	return assignStmt, lhsIdent.Name
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

	return astExprEqual(checkedExpr, starExpr.X)
}

// isNilIdent reports whether expr is the nil identifier.
func isNilIdent(expr ast.Expr) bool {
	ident, ok := expr.(*ast.Ident)
	return ok && ident.Name == "nil"
}

// astExprEqual reports whether two expressions are structurally equal by comparing their
// source rendering, which handles selectors such as props.Name.
func astExprEqual(a, b ast.Expr) bool {
	return types.ExprString(a) == types.ExprString(b)
}
