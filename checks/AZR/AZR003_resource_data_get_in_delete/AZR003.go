// Package AZR003 defines an analyzer that reports ResourceData.Get being used inside a
// resource's Delete function, where it does not work as expected.
package AZR003

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// Analyzer checks Delete functions for schema data reads. During deletion the state may be
// partial, so `d.Get(...)` (untyped resources, found via the `Delete:` registration) and
// `metadata.ResourceData.Get(...)` (typed resources, inside `Delete() sdk.ResourceFunc`)
// do not behave as expected and should not be used.
var Analyzer = &analysis.Analyzer{
	Name:     "AZR003",
	Doc:      "check for ResourceData.Get being used inside a resource's Delete function where it does not work as expected",
	URL:      "https://github.com/katbyte/azproviderlint/blob/main/checks/AZR/AZR003_resource_data_get_in_delete/README.md",
	Requires: []*analysis.Analyzer{inspect.Analyzer},
	Run:      run,
}

func run(pass *analysis.Pass) (any, error) {
	insp, ok := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	if !ok {
		return nil, nil
	}

	// collect the functions registered as `Delete:` in resource definitions
	deleteFuncs := map[string]bool{}
	insp.Preorder([]ast.Node{(*ast.KeyValueExpr)(nil)}, func(n ast.Node) {
		kv, ok := n.(*ast.KeyValueExpr)
		if !ok {
			return
		}

		key, ok := kv.Key.(*ast.Ident)
		if !ok || key.Name != "Delete" {
			return
		}

		switch v := kv.Value.(type) {
		case *ast.Ident:
			deleteFuncs[v.Name] = true
		case *ast.SelectorExpr:
			deleteFuncs[v.Sel.Name] = true
		}
	})

	insp.Preorder([]ast.Node{(*ast.FuncDecl)(nil)}, func(n ast.Node) {
		fn, ok := n.(*ast.FuncDecl)
		if !ok || fn.Body == nil {
			return
		}

		// untyped resources: the function registered via `Delete: resourceFooDelete,`
		if fn.Recv == nil && deleteFuncs[fn.Name.Name] {
			checkUntypedDelete(pass, fn)
			return
		}

		// typed resources: the `Delete() sdk.ResourceFunc` method
		if fn.Recv != nil && fn.Name.Name == "Delete" && returnsResourceFunc(fn) {
			checkTypedDelete(pass, fn)
		}
	})

	return nil, nil
}

// checkUntypedDelete reports `d.Get(...)` calls, where `d` is the function's first parameter.
func checkUntypedDelete(pass *analysis.Pass, fn *ast.FuncDecl) {
	dataParam := firstParamName(fn)
	if dataParam == "" {
		return
	}

	ast.Inspect(fn.Body, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}

		sel, ok := call.Fun.(*ast.SelectorExpr)
		if !ok || sel.Sel.Name != "Get" {
			return true
		}

		recv, ok := sel.X.(*ast.Ident)
		if !ok || recv.Name != dataParam {
			return true
		}

		pass.Reportf(call.Pos(),
			"%s.Get should not be used within a Delete function as it does not work as expected during deletion", dataParam)
		return true
	})
}

// checkTypedDelete reports `metadata.ResourceData.Get(...)` calls.
func checkTypedDelete(pass *analysis.Pass, fn *ast.FuncDecl) {
	ast.Inspect(fn.Body, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}

		sel, ok := call.Fun.(*ast.SelectorExpr)
		if !ok || sel.Sel.Name != "Get" {
			return true
		}

		inner, ok := sel.X.(*ast.SelectorExpr)
		if !ok || inner.Sel.Name != "ResourceData" {
			return true
		}

		pass.Reportf(call.Pos(),
			"ResourceData.Get should not be used within a Delete function as it does not work as expected during deletion")
		return true
	})
}

// returnsResourceFunc reports whether the function's single result type is (sdk.)ResourceFunc.
func returnsResourceFunc(fn *ast.FuncDecl) bool {
	if fn.Type.Results == nil || len(fn.Type.Results.List) != 1 {
		return false
	}

	switch t := fn.Type.Results.List[0].Type.(type) {
	case *ast.Ident:
		return t.Name == "ResourceFunc"
	case *ast.SelectorExpr:
		return t.Sel.Name == "ResourceFunc"
	}
	return false
}

func firstParamName(fn *ast.FuncDecl) string {
	if fn.Type.Params == nil || len(fn.Type.Params.List) == 0 || len(fn.Type.Params.List[0].Names) == 0 {
		return ""
	}
	return fn.Type.Params.List[0].Names[0].Name
}
