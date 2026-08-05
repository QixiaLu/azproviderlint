// Package AZT002 defines an analyzer that reports acceptance tests reading provider
// credentials from the environment instead of provisioning their own identity.
package AZT002

import (
	"go/ast"
	"go/token"
	"strconv"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// Analyzer checks for `os.Getenv("ARM_CLIENT_ID")` / `os.Getenv("ARM_CLIENT_SECRET")` /
// `os.Getenv("ARM_CLIENT_SECRET_ALT")`. Test configurations should not reuse the testing
// credentials - instead create an azurerm_user_assigned_identity as part of the test
// configuration with as minimal permissions as possible, which is cleaned up with the test.
var Analyzer = &analysis.Analyzer{
	Name:     "AZT002",
	Doc:      "check for acceptance tests reading provider credentials from the environment instead of creating a user assigned identity",
	Requires: []*analysis.Analyzer{inspect.Analyzer},
	Run:      run,
}

var bannedVariables = map[string]bool{
	"ARM_CLIENT_ID":         true,
	"ARM_CLIENT_SECRET":     true,
	"ARM_CLIENT_SECRET_ALT": true,
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
		if !ok || sel.Sel.Name != "Getenv" {
			return
		}

		pkg, ok := sel.X.(*ast.Ident)
		if !ok || pkg.Name != "os" {
			return
		}

		lit, ok := call.Args[0].(*ast.BasicLit)
		if !ok || lit.Kind != token.STRING {
			return
		}

		name, err := strconv.Unquote(lit.Value)
		if err != nil || !bannedVariables[name] {
			return
		}

		pass.Reportf(call.Pos(),
			"tests should not obtain credentials via os.Getenv(%q), create an azurerm_user_assigned_identity with minimal permissions as part of the test configuration instead", name)
	})

	return nil, nil
}
