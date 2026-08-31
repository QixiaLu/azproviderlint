// Package AZG007 defines an analyzer that reports redundant nil assignments to pointer fields
// in struct literals, where the field should be omitted instead.
package AZG007

import (
	"fmt"
	"go/ast"
	"go/types"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// Analyzer checks for pointer fields explicitly initialised to nil in a struct literal —
// `Selector: nil` — which is redundant because an omitted pointer field already takes its zero
// value (nil). Only pointer fields are flagged; slices, maps, and interfaces are left alone
// because an explicit nil there can be a deliberate, readable signal. Test files are skipped,
// since a nil entry in a test table is often semantically meaningful.
var Analyzer = &analysis.Analyzer{
	Name:     "AZG007",
	Doc:      "check for redundant nil assignments to pointer fields in struct literals that should be omitted",
	URL:      "https://github.com/katbyte/azproviderlint/blob/main/checks/AZG/AZG007_redundant_nil_pointer_field/README.md",
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

		// nil in a test table is often semantically meaningful, so skip test files.
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

			if !isNilIdent(pass, kv.Value) {
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

			// Only flag pointer fields; slices/maps/interfaces default to nil meaningfully.
			if _, isPointer := field.Type().Underlying().(*types.Pointer); !isPointer {
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
				Message: fmt.Sprintf("redundant nil assignment to pointer field %q - omit the field", keyIdent.Name),
				SuggestedFixes: []analysis.SuggestedFix{{
					Message:   "Remove the redundant nil field",
					TextEdits: []analysis.TextEdit{{Pos: kv.Pos(), End: end}},
				}},
			})
		}
	})

	return nil, nil
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
