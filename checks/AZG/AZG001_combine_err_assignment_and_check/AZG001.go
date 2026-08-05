// Package AZG001 defines an analyzer that reports '_, err := SomeFunc()' assignments
// that should be combined with the following 'if err != nil' into a single 'if' init statement.
package AZG001

import (
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// Analyzer checks for `_, err := SomeFunc()` followed immediately by `if err != nil`
// and reports that they should be combined into a single `if` init statement.
var Analyzer = &analysis.Analyzer{
	Name:     "AZG001",
	Doc:      "check for '_, err := SomeFunc()' followed by 'if err != nil' that should be combined into a single 'if' init statement",
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

		// Must be := or = with exactly 2 LHS values
		if assignStmt.Tok != token.DEFINE && assignStmt.Tok != token.ASSIGN {
			continue
		}
		if len(assignStmt.Lhs) != 2 {
			continue
		}

		// First LHS must be blank identifier (_)
		firstIdent, ok := assignStmt.Lhs[0].(*ast.Ident)
		if !ok || firstIdent.Name != "_" {
			continue
		}

		// Second LHS must be "err"
		secondIdent, ok := assignStmt.Lhs[1].(*ast.Ident)
		if !ok || secondIdent.Name != "err" {
			continue
		}

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

		pass.Reportf(assignStmt.Pos(),
			"'_, err' assignment should be combined with the following 'if err != nil' into a single 'if' init statement")
	}
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
