// Package AZG006 defines an analyzer that reports flatten* functions returning nil for a slice
// result, where an empty slice should be returned instead.
package AZG006

import (
	"fmt"
	"go/ast"
	"go/types"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// Analyzer checks that flatten* functions returning a slice type return an empty slice rather
// than nil. Callers of flatten helpers routinely range over or set the result straight into
// schema state, where a nil slice and an empty slice diverge — nil can surface as a spurious
// diff or a nil-map assignment — so the nil-input guard should return `[]T{}` (or
// `make([]T, 0)`) instead of nil.
var Analyzer = &analysis.Analyzer{
	Name:     "AZG006",
	Doc:      "check that flatten functions returning slices return an empty slice instead of nil",
	URL:      "https://github.com/katbyte/azproviderlint/blob/main/checks/AZG/AZG006_flatten_returns_nil_slice/README.md",
	Requires: []*analysis.Analyzer{inspect.Analyzer},
	Run:      run,
}

func run(pass *analysis.Pass) (any, error) {
	insp, ok := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	if !ok {
		return nil, nil
	}

	errorInterface, ok := types.Universe.Lookup("error").Type().Underlying().(*types.Interface)
	if !ok {
		return nil, nil
	}

	nodeFilter := []ast.Node{(*ast.FuncDecl)(nil)}

	insp.Preorder(nodeFilter, func(n ast.Node) {
		funcDecl, ok := n.(*ast.FuncDecl)
		if !ok || funcDecl.Name == nil {
			return
		}

		if !strings.HasPrefix(strings.ToLower(funcDecl.Name.Name), "flatten") {
			return
		}

		if funcDecl.Type.Results == nil || funcDecl.Body == nil {
			return
		}

		// Record which result positions are slice types; only those may be flagged.
		type sliceResult struct {
			index int
			typ   *ast.ArrayType
		}
		var sliceResults []sliceResult
		for i, result := range funcDecl.Type.Results.List {
			if arr, ok := result.Type.(*ast.ArrayType); ok {
				sliceResults = append(sliceResults, sliceResult{index: i, typ: arr})
			}
		}
		if len(sliceResults) == 0 {
			return
		}

		ast.Inspect(funcDecl.Body, func(node ast.Node) bool {
			retStmt, ok := node.(*ast.ReturnStmt)
			if !ok || len(retStmt.Results) == 0 {
				return true
			}

			// An error path (`return nil, err`) legitimately returns nil for the slice; only
			// the empty/nil-input branch, whose error results are all nil (or absent), is a
			// real finding, so skip any return carrying a non-nil error.
			for _, res := range retStmt.Results {
				if isNonNilError(pass, res, errorInterface) {
					return true
				}
			}

			// Collect every slice position that returns nil so a single report rewrites them
			// all — `return nil, nil` in an `([]T, []U)` flatten becomes `[]T{}, []U{}`.
			var edits []analysis.TextEdit
			for _, sr := range sliceResults {
				if sr.index >= len(retStmt.Results) {
					continue
				}
				expr := retStmt.Results[sr.index]
				if isNilIdent(pass, expr) {
					edits = append(edits, analysis.TextEdit{
						Pos:     expr.Pos(),
						End:     expr.End(),
						NewText: []byte(types.ExprString(sr.typ) + "{}"),
					})
				}
			}

			if len(edits) == 0 {
				return true
			}

			pass.Report(analysis.Diagnostic{
				Pos:     retStmt.Pos(),
				Message: fmt.Sprintf("flatten function %q should return an empty slice instead of nil", funcDecl.Name.Name),
				SuggestedFixes: []analysis.SuggestedFix{{
					Message:   "Return an empty slice instead of nil",
					TextEdits: edits,
				}},
			})
			return true
		})
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

// isNonNilError reports whether expr is a non-nil value implementing the error interface —
// `err`, `fmt.Errorf(...)`, `&myError{}` — marking the return as an error path rather than the
// empty/nil-input branch this check targets.
func isNonNilError(pass *analysis.Pass, expr ast.Expr, errorInterface *types.Interface) bool {
	if isNilIdent(pass, expr) {
		return false
	}
	t := pass.TypesInfo.TypeOf(expr)
	if t == nil {
		return false
	}
	return types.Implements(t, errorInterface)
}
