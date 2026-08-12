// Package AZR007 defines an analyzer that reports pluginsdk.StateChangeConf usage, which
// should be replaced with a custom poller implementing the pollers.PollerType interface.
package AZR007

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// pluginSDKPkgPath is the import path of the provider's pluginsdk helper package that declares
// StateChangeConf.
const pluginSDKPkgPath = "github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"

// Analyzer checks for `pluginsdk.StateChangeConf{...}` composite literals. Going forward the
// provider prefers custom pollers that implement the go-azure-sdk `pollers.PollerType`
// interface and are driven via `pollers.NewPoller(...).PollUntilDone(ctx)`.
//
// Reference: https://github.com/hashicorp/terraform-provider-azurerm/pull/30066
var Analyzer = &analysis.Analyzer{
	Name:     "AZR007",
	Doc:      "check for pluginsdk.StateChangeConf usage that should use a Custom Poller instead",
	URL:      "https://github.com/katbyte/azproviderlint/blob/main/checks/AZR/AZR007_wait_for_state_context/README.md",
	Requires: []*analysis.Analyzer{inspect.Analyzer},
	Run:      run,
}

func run(pass *analysis.Pass) (any, error) {
	insp, ok := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	if !ok {
		return nil, nil
	}

	nodeFilter := []ast.Node{
		(*ast.CompositeLit)(nil),
	}

	insp.Preorder(nodeFilter, func(n ast.Node) {
		lit, ok := n.(*ast.CompositeLit)
		if !ok {
			return
		}

		sel, ok := lit.Type.(*ast.SelectorExpr)
		if !ok || sel.Sel.Name != "StateChangeConf" {
			return
		}

		ident, ok := sel.X.(*ast.Ident)
		if !ok || ident.Name != "pluginsdk" {
			return
		}

		pkgName, ok := pass.TypesInfo.Uses[ident].(*types.PkgName)
		if !ok || pkgName.Imported().Path() != pluginSDKPkgPath {
			return
		}

		pass.Reportf(lit.Pos(),
			"prefer a Custom Poller over pluginsdk.StateChangeConf")
	})

	return nil, nil
}
