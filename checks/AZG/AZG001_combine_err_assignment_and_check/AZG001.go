// Package AZG001 defines an analyzer that reports 'err := SomeFunc()' and '_, err := SomeFunc()'
// assignments that should be combined with the following 'if err != nil' into a single 'if' init statement.
package AZG001

import (
	"go/ast"
	"go/token"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// Analyzer checks for `err := SomeFunc()` or `_, err := SomeFunc()` followed immediately by
// `if err != nil` and reports that they should be combined into a single `if` init statement.
var Analyzer = &analysis.Analyzer{
	Name:     "AZG001",
	Doc:      "check for 'err := SomeFunc()' or '_, err := SomeFunc()' followed by 'if err != nil' that should be combined into a single 'if' init statement",
	URL:      "https://github.com/katbyte/azproviderlint/blob/main/checks/AZG/AZG001_combine_err_assignment_and_check/README.md",
	Requires: []*analysis.Analyzer{inspect.Analyzer},
	Run:      run,
}

func run(pass *analysis.Pass) (any, error) {
	insp, ok := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	if !ok {
		return nil, nil
	}

	nodeFilter := []ast.Node{
		(*ast.BlockStmt)(nil),
	}

	insp.Preorder(nodeFilter, func(n ast.Node) {
		block, ok := n.(*ast.BlockStmt)
		if !ok {
			return
		}
		checkBlock(pass, block.List)
	})

	return nil, nil
}

func checkBlock(pass *analysis.Pass, stmts []ast.Stmt) {
	for i := range len(stmts) - 1 {
		assignStmt, ok := stmts[i].(*ast.AssignStmt)
		if !ok {
			continue
		}

		if assignStmt.Tok != token.DEFINE && assignStmt.Tok != token.ASSIGN {
			continue
		}

		// The last LHS value must be "err" and any values before it must be blank identifiers (_)
		lhsNames := make([]string, 0, len(assignStmt.Lhs))
		for _, lhs := range assignStmt.Lhs {
			ident, ok := lhs.(*ast.Ident)
			if !ok {
				break
			}
			lhsNames = append(lhsNames, ident.Name)
		}
		if len(lhsNames) != len(assignStmt.Lhs) {
			continue
		}
		if lhsNames[len(lhsNames)-1] != "err" {
			continue
		}
		allBlankPrefix := true
		for _, name := range lhsNames[:len(lhsNames)-1] {
			if name != "_" {
				allBlankPrefix = false
				break
			}
		}
		if !allBlankPrefix {
			continue
		}
		errIdent := assignStmt.Lhs[len(assignStmt.Lhs)-1].(*ast.Ident)

		// Next statement must be `if err != nil`
		ifStmt, ok := stmts[i+1].(*ast.IfStmt)
		if !ok {
			continue
		}

		// The if must have no init statement (otherwise it's already combined)
		if ifStmt.Init != nil {
			continue
		}

		// The condition must be `err != nil`
		if !isErrNotNil(ifStmt.Cond) {
			continue
		}

		// Combining a := declaration moves err into the if statement's scope, so it must
		// not be used again after the if statement
		if assignStmt.Tok == token.DEFINE && usedAfter(pass, errIdent, stmts[i+2:]) {
			continue
		}

		pass.Reportf(assignStmt.Pos(),
			"'%s' assignment should be combined with the following 'if err != nil' into a single 'if' init statement",
			strings.Join(lhsNames, ", "))
	}
}

// usedAfter reports whether the variable declared by ident is referenced in any of stmts.
func usedAfter(pass *analysis.Pass, ident *ast.Ident, stmts []ast.Stmt) bool {
	obj := pass.TypesInfo.ObjectOf(ident)
	if obj == nil {
		return true // can't resolve, assume it is used
	}

	for _, stmt := range stmts {
		used := false
		ast.Inspect(stmt, func(n ast.Node) bool {
			id, ok := n.(*ast.Ident)
			if !ok {
				return true
			}
			if pass.TypesInfo.ObjectOf(id) == obj {
				used = true
			}
			return !used
		})
		if used {
			return true
		}
	}

	return false
}

// isErrNotNil checks if the expression is `err != nil`.
func isErrNotNil(expr ast.Expr) bool {
	binExpr, ok := expr.(*ast.BinaryExpr)
	if !ok {
		return false
	}

	if binExpr.Op != token.NEQ {
		return false
	}

	xIdent, ok := binExpr.X.(*ast.Ident)
	if !ok || xIdent.Name != "err" {
		return false
	}

	yIdent, ok := binExpr.Y.(*ast.Ident)
	if !ok || yIdent.Name != "nil" {
		return false
	}

	return true
}
