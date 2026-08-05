// Package AZD002_mark_as_gone defines an analyzer that reports data sources using MarkAsGone instead of
// returning an error when the resource cannot be found.
package AZD002

import (
	"go/ast"
	"path/filepath"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// Analyzer checks for `metadata.MarkAsGone(...)` in data source files. Data Sources should
// return an error when a resource cannot be found rather than marking the resource as gone.
var Analyzer = &analysis.Analyzer{
	Name:     "AZD002_mark_as_gone",
	Doc:      "check for data sources using MarkAsGone instead of returning an error when the resource cannot be found",
	Requires: []*analysis.Analyzer{inspect.Analyzer},
	Run:      run,
}

func run(pass *analysis.Pass) (any, error) {
	insp, ok := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	if !ok {
		return nil, nil
	}

	nodeFilter := []ast.Node{
		(*ast.SelectorExpr)(nil),
	}

	insp.Preorder(nodeFilter, func(n ast.Node) {
		if !inDataSourceFile(pass, n) {
			return
		}

		sel, ok := n.(*ast.SelectorExpr)
		if !ok || sel.Sel.Name != "MarkAsGone" {
			return
		}

		pass.Reportf(sel.Pos(),
			"data sources should return an error when a resource cannot be found instead of calling MarkAsGone")
	})

	return nil, nil
}

func inDataSourceFile(pass *analysis.Pass, n ast.Node) bool {
	filename := filepath.Base(pass.Fset.Position(n.Pos()).Filename)
	return strings.Contains(filename, "data_source")
}
