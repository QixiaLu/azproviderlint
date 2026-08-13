// Package AZS004 defines an analyzer that reports enum validations built from a hand-written
// []string that does not cover every possible value of the SDK enum being validated.
package AZS004

import (
	"fmt"
	"go/ast"
	"go/constant"
	"go/types"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// Analyzer checks `validation.StringInSlice([]string{...}, ...)` calls whose slice references
// an SDK enum's constants. The only accepted form is the SDK's own possible-values helper
// (`PossibleValuesFor<Enum>()` / `Possible<Enum>Values()`): a partial hand-written list
// silently rejects values the API accepts, and even a complete one goes stale the moment the
// SDK adds a new value, so both are reported — with the missing values named when the list is
// incomplete.
//
// A type only counts as an enum when its package exports such a helper, so ordinary named
// string types with a few convenience constants are not reported. Slices containing
// non-constant elements or constants of more than one enum type are skipped, since coverage
// of a deliberate union or a computed list cannot be proven statically.
var Analyzer = &analysis.Analyzer{
	Name:     "AZS004",
	Doc:      "check that enum validation uses the SDK's possible-values helper or lists every possible value",
	URL:      "https://github.com/katbyte/azproviderlint/blob/main/checks/AZS/AZS004_enum_validation_missing_values/README.md",
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
		if !ok || !isStringInSliceCall(pass, call) {
			return
		}

		// only a hand-written []string literal can be proven complete or incomplete; calls
		// (e.g. PossibleValuesFor<Enum>()) and variables are out of scope
		lit, ok := ast.Unparen(call.Args[0]).(*ast.CompositeLit)
		if !ok {
			return
		}

		covered := map[string]bool{}
		var enums []*types.Named
		for _, elem := range lit.Elts {
			tv, ok := pass.TypesInfo.Types[elem]
			if !ok || tv.Value == nil || tv.Value.Kind() != constant.String {
				// a non-constant element may contribute any value, so coverage is unprovable
				return
			}
			covered[constant.StringVal(tv.Value)] = true
			enums = appendEnumConstTypes(pass, enums, elem)
		}

		// no enum constants referenced (plain strings), or a deliberate union of two enums
		if len(enums) != 1 {
			return
		}

		enum := enums[0]
		helper, values := enumValues(enum)
		if helper == "" {
			return
		}

		var missing []string
		for _, v := range values {
			if !covered[v.value] {
				missing = append(missing, fmt.Sprintf("%s (%q)", v.name, v.value))
			}
		}

		pkgName := enum.Obj().Pkg().Name()
		if len(missing) == 0 {
			pass.Reportf(lit.Pos(),
				"enum validation for %s.%s lists every value manually; use %s.%s() so new values are picked up automatically",
				pkgName, enum.Obj().Name(), pkgName, helper)
			return
		}

		pass.Reportf(lit.Pos(),
			"enum validation for %s.%s is missing %s; use %s.%s()",
			pkgName, enum.Obj().Name(), strings.Join(missing, ", "), pkgName, helper)
	})

	return nil, nil
}

// isStringInSliceCall reports whether call resolves to a StringInSlice helper: a function
// named StringInSlice with a ([]string, bool) parameter list, declared in a package whose
// import path is or ends in "validation". This matches the plugin SDK's helper/validation
// package and provider-internal wrappers of it (e.g. azurerm's internal/tf/validation)
// without depending on what the wrapper package is called locally.
func isStringInSliceCall(pass *analysis.Pass, call *ast.CallExpr) bool {
	if len(call.Args) < 1 {
		return false
	}

	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}

	obj := pass.TypesInfo.ObjectOf(sel.Sel)
	if obj == nil || obj.Name() != "StringInSlice" {
		return false
	}

	pkg := obj.Pkg()
	if pkg == nil || (pkg.Path() != "validation" && !strings.HasSuffix(pkg.Path(), "/validation")) {
		return false
	}

	sig, ok := obj.Type().(*types.Signature)
	if !ok || sig.Params().Len() != 2 {
		return false
	}

	slice, ok := sig.Params().At(0).Type().(*types.Slice)
	if !ok {
		return false
	}

	elem, ok := slice.Elem().(*types.Basic)
	if !ok || elem.Kind() != types.String {
		return false
	}

	basic, ok := sig.Params().At(1).Type().(*types.Basic)
	return ok && basic.Kind() == types.Bool
}

// appendEnumConstTypes walks elem and appends the named string type of every enum constant it
// references (deduplicated), handling `string(pkg.Const)` conversions, parenthesised forms and
// bare constant references uniformly.
func appendEnumConstTypes(pass *analysis.Pass, enums []*types.Named, elem ast.Expr) []*types.Named {
	ast.Inspect(elem, func(n ast.Node) bool {
		id, ok := n.(*ast.Ident)
		if !ok {
			return true
		}

		c, ok := pass.TypesInfo.Uses[id].(*types.Const)
		if !ok {
			return true
		}

		named, ok := types.Unalias(c.Type()).(*types.Named)
		if !ok {
			return true
		}

		basic, ok := named.Underlying().(*types.Basic)
		if !ok || basic.Info()&types.IsString == 0 {
			return true
		}

		for _, seen := range enums {
			if seen.Obj() == named.Obj() {
				return true
			}
		}
		enums = append(enums, named)
		return true
	})

	return enums
}

type enumValue struct {
	name  string
	value string
}

// enumValues returns the name of the possible-values helper for the enum and every constant of
// the enum type declared in its package, in declaration-scope order. The helper's presence is
// what marks the type as a closed enum: the generated go-azure-sdk packages emit
// `PossibleValuesFor<Name>()` and the older track-1 style SDKs emit `Possible<Name>Values()`.
// A named string type without either helper returns ("", nil) and is not treated as an enum.
func enumValues(named *types.Named) (string, []enumValue) {
	obj := named.Obj()
	pkg := obj.Pkg()
	if pkg == nil {
		return "", nil
	}

	var helper string
	for _, candidate := range []string{"PossibleValuesFor" + obj.Name(), "Possible" + obj.Name() + "Values"} {
		if _, ok := pkg.Scope().Lookup(candidate).(*types.Func); ok {
			helper = candidate
			break
		}
	}
	if helper == "" {
		return "", nil
	}

	var values []enumValue
	for _, name := range pkg.Scope().Names() {
		c, ok := pkg.Scope().Lookup(name).(*types.Const)
		if !ok || !types.Identical(types.Unalias(c.Type()), named) {
			continue
		}
		values = append(values, enumValue{name: name, value: constant.StringVal(c.Val())})
	}

	return helper, values
}
