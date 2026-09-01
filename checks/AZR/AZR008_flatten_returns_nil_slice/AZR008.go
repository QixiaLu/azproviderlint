// Package AZR008 reports flatten* functions that return nil slices.
package AZR008

import (
	"errors"
	"fmt"
	"go/ast"
	"go/types"
	"strings"

	"github.com/katbyte/azproviderlint/lib/astx"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// Analyzer checks that flatten* helpers return empty slices instead of nil.
var Analyzer = &analysis.Analyzer{
	Name:     "AZR008",
	Doc:      "check that flatten functions returning slices return an empty slice instead of nil",
	URL:      "https://github.com/katbyte/azproviderlint/blob/main/checks/AZR/AZR008_flatten_returns_nil_slice/README.md",
	Requires: []*analysis.Analyzer{inspect.Analyzer},
	Run:      run,
}

// resultInfo describes one logical return position.
type resultInfo struct {
	typeStr string
	isSlice bool
	isError bool
}

func run(pass *analysis.Pass) (any, error) {
	insp, ok := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	if !ok {
		return nil, nil
	}

	errorInterface, ok := types.Universe.Lookup("error").Type().Underlying().(*types.Interface)
	if !ok {
		return nil, errors.New("AZR008: could not resolve the built-in error interface type")
	}

	nodeFilter := []ast.Node{(*ast.FuncDecl)(nil)}

	insp.Preorder(nodeFilter, func(n ast.Node) {
		funcDecl, ok := n.(*ast.FuncDecl)
		if !ok {
			return
		}

		name := funcDecl.Name.Name
		if len(name) < 7 || !strings.EqualFold(name[:7], "flatten") {
			return
		}

		if funcDecl.Type.Results == nil || funcDecl.Body == nil {
			return
		}

		// Keep one entry per return position, including multi-name fields.
		var results []resultInfo
		haveSlice := false
		for _, field := range funcDecl.Type.Results.List {
			info := resultInfo{typeStr: types.ExprString(field.Type)}
			if t := pass.TypesInfo.TypeOf(field.Type); t != nil {
				if _, ok := t.Underlying().(*types.Slice); ok {
					info.isSlice = true
				}
				info.isError = types.Implements(t, errorInterface)
			}
			haveSlice = haveSlice || info.isSlice
			positions := len(field.Names)
			if positions == 0 {
				positions = 1
			}
			for range positions {
				results = append(results, info)
			}
		}
		if !haveSlice {
			return
		}

		ast.Inspect(funcDecl.Body, func(node ast.Node) bool {
			if _, ok := node.(*ast.FuncLit); ok {
				return false
			}

			retStmt, ok := node.(*ast.ReturnStmt)
			if !ok || len(retStmt.Results) == 0 {
				return true
			}

			// An error path (`return nil, err`) legitimately returns nil for the slice; only
			// the nil-input branch, whose error positions are all nil, is a real finding. Any
			// value in a declared error position is an error path unless it is nil — the type
			// checker already guarantees it implements error.
			for i, res := range retStmt.Results {
				if i < len(results) && results[i].isError && !astx.IsNilValue(pass, res) {
					return true
				}
			}

			var edits []analysis.TextEdit
			for i, res := range retStmt.Results {
				if i >= len(results) || !results[i].isSlice {
					continue
				}
				if astx.IsNilValue(pass, res) {
					edits = append(edits, analysis.TextEdit{
						Pos:     res.Pos(),
						End:     res.End(),
						NewText: []byte(results[i].typeStr + "{}"),
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
