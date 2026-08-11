// Package AZG005 defines an analyzer that reports fmt.Errorf calls with a single constant
// string and no format placeholders, which should use errors.New instead.
package AZG005

import (
	"go/ast"
	"go/token"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// Analyzer checks for `fmt.Errorf("fixed message")` calls whose only argument is a string
// literal containing no `%` format placeholders. Such calls don't format anything, so the
// simpler and cheaper `errors.New("fixed message")` should be used instead.
var Analyzer = &analysis.Analyzer{
	Name:     "AZG005",
	Doc:      "check for fmt.Errorf with a fixed string and no format placeholders that should use errors.New instead",
	URL:      "https://github.com/katbyte/azproviderlint/blob/main/checks/AZG/AZG005_fixed_error_string_errorf/README.md",
	Requires: []*analysis.Analyzer{inspect.Analyzer},
	Run:      run,
}

func run(pass *analysis.Pass) (any, error) {
	insp, ok := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	if !ok {
		return nil, nil
	}

	nodeFilter := []ast.Node{
		(*ast.CallExpr)(nil),
	}

	insp.Preorder(nodeFilter, func(n ast.Node) {
		call, ok := n.(*ast.CallExpr)
		if !ok || len(call.Args) != 1 {
			return
		}

		sel, ok := call.Fun.(*ast.SelectorExpr)
		if !ok || sel.Sel.Name != "Errorf" {
			return
		}

		ident, ok := sel.X.(*ast.Ident)
		if !ok || ident.Name != "fmt" {
			return
		}

		lit, ok := call.Args[0].(*ast.BasicLit)
		if !ok || lit.Kind != token.STRING {
			return
		}

		if strings.Contains(lit.Value, "%") {
			return
		}

		pass.Reportf(call.Pos(),
			"fixed error strings should use errors.New instead of fmt.Errorf")
	})

	return nil, nil
}
