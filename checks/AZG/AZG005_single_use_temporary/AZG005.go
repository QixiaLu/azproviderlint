// Package AZG005 defines an analyzer that reports single-use temporaries immediately consumed
// by the next statement — `x := <expr>` followed by `y = x` or `return x` where x has no other
// use — which should be inlined.
package AZG005

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// Analyzer checks for a short variable declaration whose variable is used exactly once in the
// entire function, by the immediately following statement, as the bare right-hand side of a
// plain assignment or as the sole return value. Such a temporary adds a name without adding
// information — `output.Format = pointer.From(input.Format)` reads as well as the two-line
// form — so it should be inlined.
//
// The consuming statement must be `y = x` (single pair, plain `=`) or `return x`; call
// arguments are deliberately out of scope since naming an argument is usually intentional
// documentation. Assignments whose left-hand side contains a function call are skipped, since
// inlining would reorder the call relative to the temporary's initializer.
var Analyzer = &analysis.Analyzer{
	Name:     "AZG005",
	Doc:      "check for single-use temporaries immediately consumed by the next statement that should be inlined",
	URL:      "https://github.com/katbyte/azproviderlint/blob/main/checks/AZG/AZG005_single_use_temporary/README.md",
	Requires: []*analysis.Analyzer{inspect.Analyzer},
	Run:      run,
}

func run(pass *analysis.Pass) (any, error) {
	insp, ok := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	if !ok {
		return nil, nil
	}

	nodeFilter := []ast.Node{
		(*ast.FuncDecl)(nil),
		(*ast.FuncLit)(nil),
	}

	insp.Preorder(nodeFilter, func(n ast.Node) {
		var body *ast.BlockStmt
		switch fn := n.(type) {
		case *ast.FuncDecl:
			body = fn.Body
		case *ast.FuncLit:
			body = fn.Body
		}
		if body == nil {
			return
		}

		// examine every statement list in this function, but not those of nested function
		// literals — they are visited as their own FuncLit node
		ast.Inspect(body, func(x ast.Node) bool {
			if _, isLit := x.(*ast.FuncLit); isLit && x != n {
				return false
			}

			var stmts []ast.Stmt
			switch node := x.(type) {
			case *ast.BlockStmt:
				stmts = node.List
			case *ast.CaseClause:
				stmts = node.Body
			case *ast.CommClause:
				stmts = node.Body
			default:
				return true
			}

			for i := range len(stmts) - 1 {
				checkPair(pass, body, stmts[i], stmts[i+1])
			}
			return true
		})
	})

	return nil, nil
}

// checkPair reports first when it declares a single-use temporary that second consumes.
func checkPair(pass *analysis.Pass, body *ast.BlockStmt, first, second ast.Stmt) {
	assign, ok := first.(*ast.AssignStmt)
	if !ok || assign.Tok != token.DEFINE || len(assign.Lhs) != 1 || len(assign.Rhs) != 1 {
		return
	}

	ident, ok := assign.Lhs[0].(*ast.Ident)
	if !ok || ident.Name == "_" {
		return
	}

	obj := pass.TypesInfo.Defs[ident]
	if obj == nil {
		return
	}

	useExpr, ok := consumesAsBareValue(pass, second, obj)
	if !ok {
		return
	}

	if useCount(pass, body, obj) != 1 {
		return
	}

	pass.Report(analysis.Diagnostic{
		Pos:            assign.Pos(),
		Message:        fmt.Sprintf("%q is only used by the following statement and should be inlined", ident.Name),
		SuggestedFixes: suggestedFixes(pass, assign, second, useExpr),
	})
}

// consumesAsBareValue reports whether stmt consumes obj as the bare right-hand side of a plain
// single assignment or as the sole return value, returning the consuming expression.
func consumesAsBareValue(pass *analysis.Pass, stmt ast.Stmt, obj types.Object) (ast.Expr, bool) {
	switch s := stmt.(type) {
	case *ast.AssignStmt:
		if s.Tok != token.ASSIGN || len(s.Lhs) != 1 || len(s.Rhs) != 1 {
			return nil, false
		}
		if !isUseOf(pass, s.Rhs[0], obj) {
			return nil, false
		}
		// a blank assignment is a discard, not a consumption worth inlining into
		if lhs, ok := s.Lhs[0].(*ast.Ident); ok && lhs.Name == "_" {
			return nil, false
		}
		// inlining moves the initializer after the left-hand side's operands in evaluation
		// order, so a call on the left-hand side could observe the swap
		if containsCallExpr(s.Lhs[0]) {
			return nil, false
		}
		return s.Rhs[0], true
	case *ast.ReturnStmt:
		if len(s.Results) == 1 && isUseOf(pass, s.Results[0], obj) {
			return s.Results[0], true
		}
	}
	return nil, false
}

// suggestedFixes deletes the temporary's declaration and replaces the consuming reference with
// the declaration's initializer, spliced as raw source text so multi-line initializers keep
// their exact original formatting.
func suggestedFixes(pass *analysis.Pass, assign *ast.AssignStmt, consumer ast.Stmt, useExpr ast.Expr) []analysis.SuggestedFix {
	exprSrc, ok := sourceText(pass, assign.Rhs[0])
	if !ok {
		return nil
	}

	return []analysis.SuggestedFix{{
		Message: "Inline the temporary into the consuming statement",
		TextEdits: []analysis.TextEdit{
			{Pos: assign.Pos(), End: consumer.Pos()},
			{Pos: useExpr.Pos(), End: useExpr.End(), NewText: exprSrc},
		},
	}}
}

// sourceText returns the raw source bytes of node.
func sourceText(pass *analysis.Pass, node ast.Node) ([]byte, bool) {
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

// isUseOf reports whether expr is (modulo parentheses) an identifier resolving to obj.
func isUseOf(pass *analysis.Pass, expr ast.Expr, obj types.Object) bool {
	id, ok := ast.Unparen(expr).(*ast.Ident)
	return ok && pass.TypesInfo.Uses[id] == obj
}

// useCount counts references to obj within body.
func useCount(pass *analysis.Pass, body *ast.BlockStmt, obj types.Object) int {
	count := 0
	ast.Inspect(body, func(n ast.Node) bool {
		if id, ok := n.(*ast.Ident); ok && pass.TypesInfo.Uses[id] == obj {
			count++
		}
		return true
	})
	return count
}

// containsCallExpr reports whether expr contains any function call.
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
