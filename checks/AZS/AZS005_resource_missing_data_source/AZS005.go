// Package AZS005 defines an analyzer that reports registered resources that have no
// corresponding data source of the same name, across every registration flavour: untyped
// plugin SDK maps, typed sdk.Resource slices and framework wrapped resource slices.
package AZS005

import (
	"go/ast"
	"go/constant"
	"go/token"
	"go/types"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// Analyzer correlates a service package's registered resources and data sources by their
// terraform type name and reports every resource without a matching data source. All
// registration flavours contribute to both sides:
//
//   - untyped plugin SDK: `SupportedResources()` / `SupportedDataSources()` map keys
//   - typed SDK: `Resources()` / `DataSources()` elements' `ResourceType()` methods
//   - framework: `FrameworkResources()` / `FrameworkDataSources()` elements' `ResourceType()` methods
//
// Conditionally registered entries (`m["azurerm_x"] = ...`, `append(out, Foo{})` behind
// feature flags) are collected too. If any data source entry's name cannot be resolved the
// package is skipped entirely, so an unresolvable data source can never produce a false
// "missing data source" report.
var Analyzer = &analysis.Analyzer{
	Name:     "AZS005",
	Doc:      "check for registered resources that have no corresponding data source",
	URL:      "https://github.com/katbyte/azproviderlint/blob/main/checks/AZS/AZS005_resource_missing_data_source/README.md",
	Requires: []*analysis.Analyzer{inspect.Analyzer},
	Run:      run,
}

// registration method names for each side; the untyped map methods and the typed/framework
// slice methods are distinguished by their return shape when collecting entries.
var (
	resourceMethods   = map[string]bool{"SupportedResources": true, "Resources": true, "FrameworkResources": true}
	dataSourceMethods = map[string]bool{"SupportedDataSources": true, "DataSources": true, "FrameworkDataSources": true}
)

// actionStyleSuffixes marks invoke-style resources — ones that perform an operation or mint a
// credential rather than manage a durable object (the kind that will eventually become
// framework Actions). These have no meaningful data source form, so they are never reported.
// There is no structural marker for them in any registration flavour, so this is a name
// convention list.
var actionStyleSuffixes = []string{
	"_run_command",
	"_sas_token",
}

type entry struct {
	name string
	pos  token.Pos
}

func run(pass *analysis.Pass) (any, error) {
	insp, ok := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	if !ok {
		return nil, nil
	}

	var resources []entry
	dataSources := map[string]bool{}
	resolvable := true

	nodeFilter := []ast.Node{
		(*ast.FuncDecl)(nil),
	}

	insp.Preorder(nodeFilter, func(n ast.Node) {
		fn, ok := n.(*ast.FuncDecl)
		if !ok || fn.Recv == nil || fn.Body == nil {
			return
		}

		isResource := resourceMethods[fn.Name.Name]
		isDataSource := dataSourceMethods[fn.Name.Name]
		if !isResource && !isDataSource {
			return
		}

		if !registrationReturnShape(pass, fn) {
			return
		}

		entries, allResolved := collectEntries(pass, fn.Body)
		if isDataSource {
			if !allResolved {
				// an unresolvable data source name could hide a match — skip the package
				resolvable = false
				return
			}
			for _, e := range entries {
				dataSources[e.name] = true
			}
			return
		}

		// unresolvable resource entries are simply skipped: they can only under-report
		resources = append(resources, entries...)
	})

	if !resolvable {
		return nil, nil
	}

	reported := map[string]bool{}
	for _, r := range resources {
		if dataSources[r.name] || reported[r.name] || isActionStyle(r.name) {
			continue
		}
		reported[r.name] = true
		pass.Reportf(r.pos, "resource %q has no corresponding data source", r.name)
	}

	return nil, nil
}

// isActionStyle reports whether a resource name marks an invoke-style resource that has no
// meaningful data source form.
func isActionStyle(name string) bool {
	for _, suffix := range actionStyleSuffixes {
		if strings.HasSuffix(name, suffix) {
			return true
		}
	}
	return false
}

// registrationReturnShape reports whether fn returns either a map[string]*Resource (untyped
// registration) or a slice of a named Resource/DataSource/FrameworkWrapped* type
// (typed/framework registration), which guards against unrelated methods sharing the names.
func registrationReturnShape(pass *analysis.Pass, fn *ast.FuncDecl) bool {
	if fn.Type.Results == nil || len(fn.Type.Results.List) != 1 {
		return false
	}

	ret := pass.TypesInfo.TypeOf(fn.Type.Results.List[0].Type)
	switch t := types.Unalias(ret).(type) {
	case *types.Map:
		basic, ok := types.Unalias(t.Key()).(*types.Basic)
		if !ok || basic.Kind() != types.String {
			return false
		}
		ptr, ok := types.Unalias(t.Elem()).(*types.Pointer)
		if !ok {
			return false
		}
		named, ok := types.Unalias(ptr.Elem()).(*types.Named)
		return ok && named.Obj().Name() == "Resource"
	case *types.Slice:
		named, ok := types.Unalias(t.Elem()).(*types.Named)
		if !ok {
			return false
		}
		switch named.Obj().Name() {
		case "Resource", "DataSource", "FrameworkWrappedResource", "FrameworkWrappedDataSource":
			return true
		}
	}

	return false
}

// collectEntries walks a registration method body and collects the terraform type name of
// every registered entry: map literal keys, `m[key] = ...` assignments, slice literal
// elements and `append(...)` arguments. It reports whether every encountered entry resolved
// to a name.
func collectEntries(pass *analysis.Pass, body *ast.BlockStmt) ([]entry, bool) {
	var entries []entry
	allResolved := true

	add := func(name string, pos token.Pos, ok bool) {
		if !ok {
			allResolved = false
			return
		}
		entries = append(entries, entry{name: name, pos: pos})
	}

	ast.Inspect(body, func(n ast.Node) bool {
		switch node := n.(type) {
		case *ast.CompositeLit:
			t := types.Unalias(pass.TypesInfo.TypeOf(node))
			if _, isMap := t.(*types.Map); isMap {
				for _, elt := range node.Elts {
					kv, ok := elt.(*ast.KeyValueExpr)
					if !ok {
						continue
					}
					name, ok := constantString(pass, kv.Key)
					add(name, kv.Key.Pos(), ok)
				}
				return false
			}
			if _, isSlice := t.(*types.Slice); isSlice {
				for _, elt := range node.Elts {
					name, ok := resourceTypeOf(pass, elt)
					add(name, elt.Pos(), ok)
				}
				return false
			}
		case *ast.AssignStmt:
			// m["azurerm_x"] = resourceX() — conditional registration behind feature flags
			for _, lhs := range node.Lhs {
				idx, ok := lhs.(*ast.IndexExpr)
				if !ok {
					continue
				}
				if _, isMap := types.Unalias(pass.TypesInfo.TypeOf(idx.X)).(*types.Map); !isMap {
					continue
				}
				name, ok := constantString(pass, idx.Index)
				add(name, idx.Index.Pos(), ok)
			}
		case *ast.CallExpr:
			// out = append(out, FooResource{}) — conditional registration behind feature flags
			if fun, ok := node.Fun.(*ast.Ident); ok && fun.Name == "append" && len(node.Args) > 1 {
				for _, arg := range node.Args[1:] {
					name, ok := resourceTypeOf(pass, arg)
					add(name, arg.Pos(), ok)
				}
				return false
			}
		}
		return true
	})

	return entries, allResolved
}

// constantString resolves expr to a constant string via the type checker, so literals and
// named constants both work.
func constantString(pass *analysis.Pass, expr ast.Expr) (string, bool) {
	tv, ok := pass.TypesInfo.Types[expr]
	if !ok || tv.Value == nil || tv.Value.Kind() != constant.String {
		return "", false
	}
	return constant.StringVal(tv.Value), true
}

// resourceTypeOf resolves a typed/framework registration element (`FooResource{}` or
// `&FooResource{}`) to the constant string its ResourceType() method returns.
func resourceTypeOf(pass *analysis.Pass, elem ast.Expr) (string, bool) {
	t := types.Unalias(pass.TypesInfo.TypeOf(elem))
	if ptr, ok := t.(*types.Pointer); ok {
		t = types.Unalias(ptr.Elem())
	}
	named, ok := t.(*types.Named)
	if !ok {
		return "", false
	}

	decl := resourceTypeMethod(pass, named.Obj().Name())
	if decl == nil {
		return "", false
	}

	var name string
	found := false
	ast.Inspect(decl.Body, func(n ast.Node) bool {
		ret, ok := n.(*ast.ReturnStmt)
		if !ok || found || len(ret.Results) != 1 {
			return true
		}
		name, found = constantString(pass, ret.Results[0])
		return !found
	})

	return name, found
}

// resourceTypeMethod finds the `func (r TypeName) ResourceType() string` declaration for the
// named receiver type in the package being analyzed.
func resourceTypeMethod(pass *analysis.Pass, typeName string) *ast.FuncDecl {
	for _, file := range pass.Files {
		for _, decl := range file.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok || fn.Name.Name != "ResourceType" || fn.Recv == nil || len(fn.Recv.List) != 1 || fn.Body == nil {
				continue
			}
			recv := fn.Recv.List[0].Type
			if star, ok := recv.(*ast.StarExpr); ok {
				recv = star.X
			}
			if id, ok := recv.(*ast.Ident); ok && id.Name == typeName {
				return fn
			}
		}
	}
	return nil
}
