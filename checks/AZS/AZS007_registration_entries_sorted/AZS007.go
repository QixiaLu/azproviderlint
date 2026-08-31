// Package AZS007 defines an analyzer that reports registration.go map and slice entries
// that are not sorted alphabetically.
package AZS007

import (
	"bytes"
	"go/ast"
	"go/token"
	"go/types"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// Analyzer checks that map and slice entries returned directly by `Registration` methods in
// registration.go files are sorted alphabetically. Entries grouped into sections separated by
// blank or comment lines are validated within each section independently.
var Analyzer = &analysis.Analyzer{
	Name:     "AZS007",
	Doc:      "check that registration.go map and slice entries are sorted alphabetically",
	URL:      "https://github.com/katbyte/azproviderlint/blob/main/checks/AZS/AZS007_registration_entries_sorted/README.md",
	Requires: []*analysis.Analyzer{inspect.Analyzer},
	Run:      run,
}

func run(pass *analysis.Pass) (any, error) {
	insp, ok := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	if !ok {
		return nil, nil
	}

	insp.Preorder([]ast.Node{(*ast.FuncDecl)(nil)}, func(n ast.Node) {
		funcDecl, ok := n.(*ast.FuncDecl)
		if !ok {
			return
		}

		if filepath.Base(pass.Fset.Position(funcDecl.Pos()).Filename) != "registration.go" {
			return
		}

		if !hasRegistrationReceiver(funcDecl) {
			return
		}

		analyzeRegistrationMethod(pass, funcDecl)
	})

	return nil, nil
}

// hasRegistrationReceiver reports whether the function has a Registration receiver.
func hasRegistrationReceiver(funcDecl *ast.FuncDecl) bool {
	if funcDecl.Recv == nil || len(funcDecl.Recv.List) == 0 {
		return false
	}

	recv := funcDecl.Recv.List[0]
	var typeName string

	switch t := recv.Type.(type) {
	case *ast.Ident:
		typeName = t.Name
	case *ast.StarExpr:
		if ident, ok := t.X.(*ast.Ident); ok {
			typeName = ident.Name
		}
	}

	return typeName == "Registration"
}

// analyzeRegistrationMethod validates the map or slice literals a registration method returns
// directly. Literals assigned to a local variable before being returned are left alone, since the
// value that reaches the return can be reassigned or appended to.
func analyzeRegistrationMethod(pass *analysis.Pass, funcDecl *ast.FuncDecl) {
	if funcDecl.Body == nil {
		return
	}

	ast.Inspect(funcDecl.Body, func(n ast.Node) bool {
		returnStmt, ok := n.(*ast.ReturnStmt)
		if !ok {
			return true
		}
		for _, result := range returnStmt.Results {
			if compositeLit, ok := result.(*ast.CompositeLit); ok {
				validateSorting(pass, compositeLit)
			}
		}
		return true
	})
}

// splitIntoSections groups composite-literal elements into sections separated by blank or comment
// lines: a gap of more than one source line between adjacent elements starts a new section.
func splitIntoSections(fset *token.FileSet, elts []ast.Expr) [][]ast.Expr {
	if len(elts) == 0 {
		return nil
	}

	var sections [][]ast.Expr
	current := []ast.Expr{elts[0]}

	for i := 1; i < len(elts); i++ {
		prevEnd := fset.Position(elts[i-1].End()).Line
		currStart := fset.Position(elts[i].Pos()).Line
		if currStart-prevEnd > 1 {
			sections = append(sections, current)
			current = []ast.Expr{elts[i]}
		} else {
			current = append(current, elts[i])
		}
	}

	return append(sections, current)
}

// validateSorting reports a composite literal whose entries are unsorted, checking each
// blank- or comment-separated section independently.
func validateSorting(pass *analysis.Pass, compositeLit *ast.CompositeLit) {
	if compositeLit.Type == nil {
		return
	}

	typ := pass.TypesInfo.TypeOf(compositeLit)
	if typ == nil {
		return
	}

	var isMap bool
	switch typ.Underlying().(type) {
	case *types.Map:
		isMap = true
	case *types.Slice:
	default:
		return
	}

	sections := splitIntoSections(pass.Fset, compositeLit.Elts)

	var unsorted [][]ast.Expr
	for _, section := range sections {
		names := make([]string, 0, len(section))
		resolvable := true
		for _, elt := range section {
			name, ok := sortKey(elt, isMap)
			if !ok {
				resolvable = false
				break
			}
			names = append(names, name)
		}

		// Skip sections with an unresolvable key rather than judge a partial subset.
		if !resolvable {
			continue
		}

		if !sort.SliceIsSorted(names, func(i, j int) bool { return lessFold(names[i], names[j]) }) {
			unsorted = append(unsorted, section)
		}
	}

	if len(unsorted) == 0 {
		return
	}

	pass.Report(analysis.Diagnostic{
		Pos:            compositeLit.Pos(),
		Message:        "registration entries should be sorted alphabetically",
		SuggestedFixes: sortFixes(pass, compositeLit, unsorted, isMap),
	})
}

// sortFixes builds a single suggested fix reordering every unsorted section alphabetically. It
// returns nil when entries cannot be rewritten safely (unresolvable key, entries sharing a line,
// or a brace sharing an entry's line) so the diagnostic is still reported without a broken fix.
func sortFixes(pass *analysis.Pass, compositeLit *ast.CompositeLit, sections [][]ast.Expr, isMap bool) []analysis.SuggestedFix {
	tf := pass.Fset.File(compositeLit.Pos())
	if tf == nil {
		return nil
	}

	// Bail out when a brace shares an entry's line: entries are copied as whole lines, so moving
	// one would drag the brace with it (and a trailing brace makes LineStart(endLine+1) panic).
	if elts := compositeLit.Elts; len(elts) > 0 {
		if tf.Line(compositeLit.Lbrace) == tf.Line(elts[0].Pos()) ||
			tf.Line(compositeLit.Rbrace) == tf.Line(elts[len(elts)-1].End()) {
			return nil
		}
	}

	content, err := pass.ReadFile(tf.Name())
	if err != nil {
		return nil
	}

	var edits []analysis.TextEdit
	for _, section := range sections {
		edit, ok := sortSectionEdit(tf, content, section, isMap)
		if !ok {
			return nil
		}
		edits = append(edits, edit)
	}

	return []analysis.SuggestedFix{{
		Message:   "Sort registration entries alphabetically",
		TextEdits: edits,
	}}
}

// sortSectionEdit produces a text edit reordering one section's entries alphabetically. Each
// entry is copied as the whole source lines it spans, so its trailing comment moves with it.
func sortSectionEdit(tf *token.File, content []byte, section []ast.Expr, isMap bool) (analysis.TextEdit, bool) {
	type block struct {
		key  string
		text []byte
	}

	blocks := make([]block, len(section))
	var editStart, editEnd token.Pos
	prevEndLine := 0
	for i, elt := range section {
		key, ok := sortKey(elt, isMap)
		if !ok {
			return analysis.TextEdit{}, false
		}

		startLine := tf.Line(elt.Pos())
		endLine := tf.Line(elt.End())
		if i > 0 && startLine <= prevEndLine {
			return analysis.TextEdit{}, false // entries share a line; not safely reorderable
		}
		prevEndLine = endLine

		lineStart := tf.LineStart(startLine)
		lineEnd := tf.LineStart(endLine+1) - 1
		if i == 0 {
			editStart = lineStart
		}
		editEnd = lineEnd
		blocks[i] = block{key: key, text: content[tf.Offset(lineStart):tf.Offset(lineEnd)]}
	}

	sort.SliceStable(blocks, func(i, j int) bool {
		return lessFold(blocks[i].key, blocks[j].key)
	})

	var buf bytes.Buffer
	for i, b := range blocks {
		if i > 0 {
			buf.WriteByte('\n')
		}
		buf.Write(b.text)
	}

	return analysis.TextEdit{Pos: editStart, End: editEnd, NewText: buf.Bytes()}, true
}

// lessFold orders registration keys case-insensitively.
func lessFold(a, b string) bool {
	return strings.ToLower(a) < strings.ToLower(b)
}

// sortKey returns the alphabetical sort key for a registration entry: the map key for map
// literals, the resource struct or constructor name for slice literals.
func sortKey(elt ast.Expr, isMap bool) (string, bool) {
	if isMap {
		kv, ok := elt.(*ast.KeyValueExpr)
		if !ok {
			return "", false
		}
		lit, ok := kv.Key.(*ast.BasicLit)
		if !ok || lit.Kind != token.STRING {
			return "", false
		}
		key, err := strconv.Unquote(lit.Value)
		if err != nil {
			return "", false
		}
		return key, true
	}

	return sliceEntryKey(elt)
}

// sliceEntryKey extracts the sort key from a typed or framework slice entry: the struct type
// name of a composite literal (`FooResource{}`, `&FooResource{}`, `pkg.FooResource{}`), the
// constructor identifier of a framework entry (`newFooResource`, `pkg.NewFooResource`), or the
// unquoted value of a string literal (`[]string` lists). Any package
// qualifier is dropped so `pkg.FooResource` and `FooResource` sort identically.
func sliceEntryKey(elt ast.Expr) (string, bool) {
	switch e := elt.(type) {
	case *ast.UnaryExpr:
		if e.Op == token.AND {
			return sliceEntryKey(e.X)
		}
	case *ast.CompositeLit:
		return sliceEntryKey(e.Type)
	case *ast.Ident:
		return e.Name, true
	case *ast.SelectorExpr:
		return e.Sel.Name, true
	case *ast.BasicLit:
		if e.Kind == token.STRING {
			if s, err := strconv.Unquote(e.Value); err == nil {
				return s, true
			}
		}
	}
	return "", false
}
