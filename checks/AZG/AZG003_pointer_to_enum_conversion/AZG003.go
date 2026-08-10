// Package AZG003 defines an analyzer that reports pointer.To being used with an explicit
// go-azure-sdk enum type conversion (pointer.To(sdk.Enum(v))) where the generic
// pointer.ToEnum[sdk.Enum](v) helper should be used instead.
package AZG003

import (
	"fmt"
	"go/ast"
	"go/types"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const (
	// pointerPkgPath is the import path of the go-azure-helpers pointer package.
	pointerPkgPath = "github.com/hashicorp/go-azure-helpers/lang/pointer"
	// goAzureSDKPath is a fragment of the go-azure-sdk import path used to identify SDK enum types.
	goAzureSDKPath = "github.com/hashicorp/go-azure-sdk"
)

// Analyzer checks for `pointer.To(sdk.SomeEnum(v))` calls that convert a go-azure-sdk enum
// type with an explicit conversion. The generic `pointer.ToEnum[sdk.SomeEnum](v)` helper is
// clearer and type-safe, so those call sites should use it instead.
var Analyzer = &analysis.Analyzer{
	Name:     "AZG003",
	Doc:      "check for pointer.To with an explicit go-azure-sdk enum conversion that should use pointer.ToEnum instead",
	URL:      "https://github.com/katbyte/azproviderlint/blob/main/checks/AZG/AZG003_pointer_to_enum_conversion/README.md",
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

		argCall, ok := call.Args[0].(*ast.CallExpr)
		if !ok || len(argCall.Args) != 1 {
			return
		}

		sel, ok := call.Fun.(*ast.SelectorExpr)
		if !ok || sel.Sel.Name != "To" {
			return
		}

		ident, ok := sel.X.(*ast.Ident)
		if !ok || ident.Name != "pointer" {
			return
		}

		pkgName, ok := pass.TypesInfo.Uses[ident].(*types.PkgName)
		if !ok || pkgName.Imported().Path() != pointerPkgPath {
			return
		}

		named, ok := pass.TypesInfo.TypeOf(argCall.Fun).(*types.Named)
		if !ok || !isAzureSDKEnumType(pass, named) {
			return
		}

		pass.Report(analysis.Diagnostic{
			Pos:     call.Pos(),
			Message: fmt.Sprintf("pointer.To with an explicit go-azure-sdk enum conversion should use pointer.ToEnum[%s] instead", named.Obj().Name()),
		})
	})

	return nil, nil
}

// isAzureSDKEnumType reports whether named is a go-azure-sdk enum type: a named type with a
// string or integer underlying type, declared in a go-azure-sdk package, that either exposes
// the generated `PossibleValuesFor<Name>() []T` helper or is declared in a constants.go file.
func isAzureSDKEnumType(pass *analysis.Pass, named *types.Named) bool {
	basic, ok := named.Underlying().(*types.Basic)
	if !ok {
		return false
	}

	if basic.Info()&(types.IsString|types.IsInteger) == 0 {
		return false
	}

	obj := named.Obj()
	pkg := obj.Pkg()
	if pkg == nil || !strings.Contains(pkg.Path(), goAzureSDKPath) {
		return false
	}

	// The generated SDK emits a PossibleValuesFor<TypeName>() []T helper for each enum.
	lookup := pkg.Scope().Lookup("PossibleValuesFor" + obj.Name())
	if lookup == nil {
		// Fall back to the convention that enum types are declared in constants.go.
		return strings.HasSuffix(pass.Fset.Position(obj.Pos()).Filename, "constants.go")
	}

	fn, ok := lookup.(*types.Func)
	if !ok {
		return false
	}

	sig, ok := fn.Type().(*types.Signature)
	if !ok || sig.Params().Len() != 0 || sig.Results().Len() != 1 {
		return false
	}

	slice, ok := sig.Results().At(0).Type().(*types.Slice)
	if !ok {
		return false
	}

	elem, ok := slice.Elem().(*types.Basic)
	return ok && elem.Kind() == basic.Kind()
}
