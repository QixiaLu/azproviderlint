// Package AZS007 defines an analyzer that reports schema declarations that have both
// Optional: true and Computed: true without a "// Note: O+C because ..." comment between
// the two fields, explaining why the field is Optional+Computed.
package AZS007

import (
	"go/ast"
	"go/token"
	"go/types"
	"regexp"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// Analyzer checks that schema fields with both Optional: true and Computed: true have a
// "// Note: O+C because ..." comment between the Optional and Computed field declarations,
// documenting why the field must be Optional+Computed rather than just Optional, or Optional with a Default.
var Analyzer = &analysis.Analyzer{
	Name:     "AZS007",
	Doc:      "check that schema fields with both Optional and Computed have a '// Note: O+C because ...' comment between them",
	URL:      "https://github.com/katbyte/azproviderlint/blob/main/checks/AZS/AZS007_schema_optional_computed_missing_comment/README.md",
	Requires: []*analysis.Analyzer{inspect.Analyzer},
	Run:      run,
}

const failureMessage = "schema field has both Optional and Computed but is missing a '// Note: O+C because ...' comment between the two fields"

// ocCommentRe matches the required documentation comment.
// The pattern is case-insensitive "note:" followed by "O+C" and at least one non-space character.
// Both "// Note: O+C because ..." and "// NOTE: O+C ..." variants used in the codebase are matched.
var ocCommentRe = regexp.MustCompile(`(?i)//\s*note:\s*o\+c\b`)

func run(pass *analysis.Pass) (any, error) {
	insp, ok := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	if !ok {
		return nil, nil
	}

	// State migration schemas are historical snapshots of a resource's schema at a point in
	// time. They are stripped of comments, validation, defaults and other metadata by
	// convention, so O+C fields legitimately lack a reason comment. Skip the entire package.
	if pass.Pkg.Name() == "migration" {
		return nil, nil
	}

	// Build a per-file map of comment lines for fast lookup.
	// commentLines[file][line] = comment text on that line.
	type fileComments = map[int]string
	commentsByFile := make(map[*token.File]fileComments)
	for _, f := range pass.Files {
		fc := make(fileComments)
		for _, cg := range f.Comments {
			for _, c := range cg.List {
				line := pass.Fset.Position(c.Slash).Line
				fc[line] = c.Text
			}
		}
		commentsByFile[pass.Fset.File(f.Pos())] = fc
	}

	nodeFilter := []ast.Node{
		(*ast.CompositeLit)(nil),
	}

	insp.Preorder(nodeFilter, func(n ast.Node) {
		cl, ok := n.(*ast.CompositeLit)
		if !ok {
			return
		}

		if !isSchemaType(pass, cl) {
			return
		}

		optionalPos, computedPos, hasOptional, hasComputed := findOptionalAndComputedPositions(pass, cl)
		if !hasOptional || !hasComputed {
			return
		}

		// Determine the line range to search for a O+C comment.
		// The comment must appear on a line strictly between the Optional and Computed field lines,
		tf := pass.Fset.File(cl.Pos())
		if tf == nil {
			return
		}

		optionalLine := pass.Fset.Position(optionalPos).Line
		computedLine := pass.Fset.Position(computedPos).Line

		// Normalise so low is always the earlier line regardless of field order in the literal.
		low, high := optionalLine, computedLine
		if computedLine < optionalLine {
			low, high = computedLine, optionalLine
		}

		fc, hasFC := commentsByFile[tf]
		if !hasFC {
			pass.Reportf(computedPos, failureMessage)
			return
		}

		for line := low + 1; line < high; line++ {
			text, exists := fc[line]
			if !exists {
				continue
			}
			if ocCommentRe.MatchString(strings.TrimSpace(text)) {
				return
			}
		}

		pass.Reportf(computedPos, failureMessage)
	})

	return nil, nil
}

// isSchemaType reports whether the composite literal represents a schema.Schema value,
// including through pluginsdk's type alias.
func isSchemaType(pass *analysis.Pass, cl *ast.CompositeLit) bool {
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
	return obj.Name() == "Schema" && obj.Pkg() != nil && obj.Pkg().Name() == "schema"
}

// findOptionalAndComputedPositions returns the source positions of the Optional: true and
// Computed: true field assignments within a schema.Schema composite literal, and whether
// each was found.
func findOptionalAndComputedPositions(pass *analysis.Pass, cl *ast.CompositeLit) (optionalPos, computedPos token.Pos, hasOptional, hasComputed bool) {
	for _, elt := range cl.Elts {
		kv, ok := elt.(*ast.KeyValueExpr)
		if !ok {
			continue
		}

		key, ok := kv.Key.(*ast.Ident)
		if !ok {
			continue
		}

		switch key.Name {
		case "Optional":
			if isTrueBoolValue(pass, kv.Value) {
				hasOptional = true
				optionalPos = kv.Key.Pos()
			}
		case "Computed":
			if isTrueBoolValue(pass, kv.Value) {
				hasComputed = true
				computedPos = kv.Key.Pos()
			}
		}
	}

	return
}

// isTrueBoolValue reports whether the expression is the constant `true`.
func isTrueBoolValue(pass *analysis.Pass, e ast.Expr) bool {
	if e == nil {
		return false
	}

	if ident, ok := e.(*ast.Ident); ok && ident.Name == "true" {
		return true
	}

	tv := pass.TypesInfo.Types[e]
	return tv.IsValue() && tv.Value != nil && tv.Value.String() == "true"
}
