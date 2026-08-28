// Package AZS003 defines an analyzer that reports optional or required list blocks whose
// properties are all optional with no defaults, so an empty block (`foo {}`) is valid
// configuration that can crash expand functions or produce spurious diffs.
package AZS003

import (
	"go/ast"
	"go/token"
	"go/types"

	"github.com/katbyte/azproviderlint/helpers"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// Analyzer checks for `Type: schema.TypeList` blocks (Optional or Required) whose nested
// Resource schema contains only optional properties with no Default/DefaultFunc and no
// AtLeastOneOf/ExactlyOneOf constraint. Such blocks accept `foo {}`, which yields a nil
// element that expand functions commonly crash on (`raw[0].(map[string]interface{})`), or
// a permanent diff. A single Required property, Default, or AtLeastOneOf/ExactlyOneOf
// constraint on any property makes the block safe.
var Analyzer = &analysis.Analyzer{
	Name:     "AZS003",
	Doc:      "check for list blocks that allow empty blocks because every property is optional with no default or constraint",
	URL:      "https://github.com/katbyte/azproviderlint/blob/main/checks/AZS/AZS003_schema_allows_empty_block/README.md",
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
		cl, ok := n.(*ast.CompositeLit)
		if !ok || !isSchemaHelperType(pass, cl, "Schema") {
			return
		}

		fields := compositeLitFields(cl)

		var schemaType string
		switch v := fields["Type"].(type) {
		case *ast.SelectorExpr:
			schemaType = v.Sel.Name
		case *ast.Ident:
			schemaType = v.Name
		}
		if schemaType != "TypeList" {
			return
		}

		// computed-only blocks cannot be set in configuration
		if !helpers.IsTrueConstant(pass, fields["Optional"]) && !helpers.IsTrueConstant(pass, fields["Required"]) {
			return
		}

		ref, ok := fields["Elem"].(*ast.UnaryExpr)
		if !ok || ref.Op != token.AND {
			return
		}
		resource, ok := ref.X.(*ast.CompositeLit)
		if !ok || !isSchemaHelperType(pass, resource, "Resource") {
			return
		}
		properties, ok := compositeLitFields(resource)["Schema"].(*ast.CompositeLit)
		if !ok || len(properties.Elts) == 0 {
			return
		}

		for _, elt := range properties.Elts {
			kv, ok := elt.(*ast.KeyValueExpr)
			if !ok {
				return
			}
			value := kv.Value
			if addr, isAddr := value.(*ast.UnaryExpr); isAddr && addr.Op == token.AND {
				value = addr.X
			}
			// a property defined elsewhere (variable/function) cannot be inspected - stay quiet
			property, ok := value.(*ast.CompositeLit)
			if !ok {
				return
			}
			propertyFields := compositeLitFields(property)
			if helpers.IsTrueConstant(pass, propertyFields["Required"]) ||
				isSet(propertyFields["Default"]) || isSet(propertyFields["DefaultFunc"]) ||
				isSet(propertyFields["AtLeastOneOf"]) || isSet(propertyFields["ExactlyOneOf"]) {
				return
			}
		}

		pass.Reportf(cl.Pos(),
			"schema allows an empty block as every property is optional with no default - add AtLeastOneOf/ExactlyOneOf constraints, a Required property, or a Default so empty blocks cannot crash expand functions or cause spurious diffs")
	})

	return nil, nil
}

// isSchemaHelperType reports whether the composite literal's type is the named type from a
// package named "schema", which azurerm's pluginsdk aliases also resolve to.
func isSchemaHelperType(pass *analysis.Pass, cl *ast.CompositeLit, name string) bool {
	t := pass.TypesInfo.TypeOf(cl)
	if t == nil {
		return false
	}
	if ptr, ok := types.Unalias(t).(*types.Pointer); ok {
		t = ptr.Elem()
	}
	named, ok := types.Unalias(t).(*types.Named)
	if !ok {
		return false
	}
	obj := named.Obj()
	return obj.Name() == name && obj.Pkg() != nil && obj.Pkg().Name() == "schema"
}

func compositeLitFields(cl *ast.CompositeLit) map[string]ast.Expr {
	fields := make(map[string]ast.Expr, len(cl.Elts))
	for _, elt := range cl.Elts {
		if kv, ok := elt.(*ast.KeyValueExpr); ok {
			if key, ok := kv.Key.(*ast.Ident); ok {
				fields[key.Name] = kv.Value
			}
		}
	}
	return fields
}

// isSet reports whether the field is present and not literally nil.
func isSet(e ast.Expr) bool {
	if e == nil {
		return false
	}
	ident, ok := e.(*ast.Ident)
	return !ok || ident.Name != "nil"
}
