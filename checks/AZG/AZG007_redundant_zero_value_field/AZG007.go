// Package AZG007 defines an analyzer that reports redundant zero-value assignments to struct
// literal fields, where the field should be omitted instead.
package AZG007

import (
	"fmt"
	"go/ast"
	"go/constant"
	"go/token"
	"go/types"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// Analyzer checks for fields explicitly initialised to their zero value in a struct literal —
// a pointer set to `nil`, a string set to `""`, a numeric set to `0`, or a bool set to `false` —
// which is redundant because an omitted field already takes its zero value. Only pointer and
// basic (string/numeric/bool) fields are flagged; slices, maps, and interfaces are left alone
// because an explicit nil there can be a deliberate, readable signal. Zero values written as a
// named constant (`Type: TypeNone`) are also left alone, since the name conveys intent. Test
// files are skipped, since a zero entry in a test table is often semantically meaningful.
var Analyzer = &analysis.Analyzer{
	Name:     "AZG007",
	Doc:      "check for redundant zero-value assignments to struct literal fields that should be omitted",
	URL:      "https://github.com/katbyte/azproviderlint/blob/main/checks/AZG/AZG007_redundant_zero_value_field/README.md",
	Requires: []*analysis.Analyzer{inspect.Analyzer},
	Run:      run,
}

func run(pass *analysis.Pass) (any, error) {
	insp, ok := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	if !ok {
		return nil, nil
	}

	nodeFilter := []ast.Node{(*ast.CompositeLit)(nil)}

	// All comment groups across the package's files, used to avoid deleting a comment that
	// leads the field following the one being removed. Each file's comments are position-sorted.
	var commentGroups []*ast.CommentGroup
	for _, f := range pass.Files {
		commentGroups = append(commentGroups, f.Comments...)
	}

	insp.Preorder(nodeFilter, func(n ast.Node) {
		compositeLit, ok := n.(*ast.CompositeLit)
		if !ok {
			return
		}

		// A zero entry in a test table is often semantically meaningful, so skip test files.
		if strings.HasSuffix(pass.Fset.Position(compositeLit.Pos()).Filename, "_test.go") {
			return
		}

		structType := structTypeOf(pass.TypesInfo.TypeOf(compositeLit))
		if structType == nil {
			return
		}

		for i, elt := range compositeLit.Elts {
			kv, ok := elt.(*ast.KeyValueExpr)
			if !ok {
				continue
			}

			keyIdent, ok := kv.Key.(*ast.Ident)
			if !ok {
				continue
			}

			fieldObj, _, _ := types.LookupFieldOrMethod(structType, true, pass.Pkg, keyIdent.Name)
			field, ok := fieldObj.(*types.Var)
			if !ok {
				continue
			}

			message := redundantZeroMessage(pass, field.Type(), kv.Value, keyIdent.Name)
			if message == "" {
				continue
			}

			// Delete the element up to the start of the next one (or the closing brace for the
			// last element), which takes the trailing comma and any trailing comment with it;
			// gofmt collapses the leftover whitespace afterwards.
			end := compositeLit.Rbrace
			if i+1 < len(compositeLit.Elts) {
				end = compositeLit.Elts[i+1].Pos()
			}

			// A comment on a later line than this field leads the following field, so stop the
			// deletion before it; a same-line trailing comment stays inside the removed span.
			eltLine := pass.Fset.Position(elt.End()).Line
			for _, cg := range commentGroups {
				if cg.Pos() > elt.End() && cg.Pos() < end && pass.Fset.Position(cg.Pos()).Line > eltLine {
					end = cg.Pos()
					break
				}
			}

			pass.Report(analysis.Diagnostic{
				Pos:     kv.Pos(),
				Message: message,
				SuggestedFixes: []analysis.SuggestedFix{{
					Message:   "Remove the redundant field",
					TextEdits: []analysis.TextEdit{{Pos: kv.Pos(), End: end}},
				}},
			})
		}
	})

	return nil, nil
}

// redundantZeroMessage returns a diagnostic message when value is a redundant zero-value
// assignment for a field of fieldType, or an empty string otherwise. Pointer fields are
// redundant when set to `nil`; string/numeric/bool fields when set to a literal `""`, `0`, or
// `false`. Slices, maps, interfaces, and named-constant zero values are left alone.
func redundantZeroMessage(pass *analysis.Pass, fieldType types.Type, value ast.Expr, name string) string {
	switch underlying := fieldType.Underlying().(type) {
	case *types.Pointer:
		if isNilIdent(pass, value) {
			return fmt.Sprintf("redundant nil assignment to pointer field %q - omit the field", name)
		}
	case *types.Basic:
		if isBasicZeroLit(pass, underlying, value) {
			return fmt.Sprintf("redundant zero-value assignment to field %q - omit the field", name)
		}
	}
	return ""
}

// isNilIdent reports whether expr is the predeclared nil identifier.
func isNilIdent(pass *analysis.Pass, expr ast.Expr) bool {
	ident, ok := expr.(*ast.Ident)
	if !ok {
		return false
	}
	_, isNil := pass.TypesInfo.Uses[ident].(*types.Nil)
	return isNil
}

// isBasicZeroLit reports whether value is a literal zero value for the given basic type: the
// predeclared `false` for booleans, an empty string literal for strings, and a `0` numeric
// literal for numeric types. Named constants that happen to be zero (`Type: TypeNone`) return
// false, since the name documents intent and omitting the field would lose that.
func isBasicZeroLit(pass *analysis.Pass, basic *types.Basic, value ast.Expr) bool {
	switch info := basic.Info(); {
	case info&types.IsBoolean != 0:
		ident, ok := value.(*ast.Ident)
		return ok && pass.TypesInfo.Uses[ident] == types.Universe.Lookup("false")
	case info&types.IsString != 0:
		lit, ok := value.(*ast.BasicLit)
		if !ok || lit.Kind != token.STRING {
			return false
		}
		cv := pass.TypesInfo.Types[value].Value
		return cv != nil && constant.StringVal(cv) == ""
	case info&types.IsNumeric != 0:
		lit, ok := value.(*ast.BasicLit)
		if !ok || (lit.Kind != token.INT && lit.Kind != token.FLOAT && lit.Kind != token.IMAG) {
			return false
		}
		cv := pass.TypesInfo.Types[value].Value
		return cv != nil && constant.Sign(cv) == 0
	}
	return false
}

// structTypeOf returns the struct type behind t, dereferencing a pointer if needed, or nil when
// t is not a struct.
func structTypeOf(t types.Type) types.Type {
	if t == nil {
		return nil
	}
	if ptr, ok := t.Underlying().(*types.Pointer); ok {
		t = ptr.Elem()
	}
	if _, ok := t.Underlying().(*types.Struct); ok {
		return t
	}
	return nil
}
