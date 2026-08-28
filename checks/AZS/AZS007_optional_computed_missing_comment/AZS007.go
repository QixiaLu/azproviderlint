// Package AZS007 defines an analyzer that reports schema declarations that have both
// Optional: true and Computed: true without a "// Note: O+C because ..." comment between
// the two fields, explaining why the field is Optional+Computed.
package AZS007

import (
	"flag"
	"go/ast"
	"go/token"
	"go/types"
	"regexp"
	"strings"

	"github.com/katbyte/azproviderlint/helpers"
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
	URL:      "https://github.com/katbyte/azproviderlint/blob/main/checks/AZS/AZS007_optional_computed_missing_comment/README.md",
	Requires: []*analysis.Analyzer{inspect.Analyzer},
	Run:      run,
}

var excludePackages string

func init() {
	Analyzer.Flags.Init("AZS007", flag.ContinueOnError)
	Analyzer.Flags.StringVar(&excludePackages, "exclude-packages", "",
		"comma-separated list of package names to skip")
}

const failureMessage = "schema field has both Optional and Computed but is missing a '// Note: O+C because ...' comment between the two fields"

// ocCommentRe matches the required documentation comment.
// The pattern is case-insensitive "note:" followed by "O+C" and at least one non-space character.
// Both "// Note: O+C because ..." and "// NOTE: O+C ..." variants used in the codebase are matched.
var ocCommentRe = regexp.MustCompile(`(?i)//\s*note:\s*o\+c\s\S+`)

func run(pass *analysis.Pass) (any, error) {
	insp, ok := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	if !ok {
		return nil, nil
	}

	// Skip packages listed in exclude-packages (comma-separated).
	for _, pkg := range strings.Split(excludePackages, ",") {
		if strings.TrimSpace(pkg) == pass.Pkg.Name() {
			return nil, nil
		}
	}

	// a file's line to comment map is only built the first time a Schema literal in that file is found to have both Optional and Computed set.
	type fileComments = map[int]string
	commentsByFile := make(map[*token.File]fileComments)

	// astFileByTokenFile maps a *token.File back to its *ast.File so comments can be
	// indexed on demand without a second pass over pass.Files.
	astFileByTokenFile := make(map[*token.File]*ast.File, len(pass.Files))
	for _, f := range pass.Files {
		astFileByTokenFile[pass.Fset.File(f.Pos())] = f
	}

	// commentsFor returns the line to comment map for tf, building it on first access.
	commentsFor := func(tf *token.File) fileComments {
		if fc, ok := commentsByFile[tf]; ok {
			return fc
		}
		fc := make(fileComments)
		if f, ok := astFileByTokenFile[tf]; ok {
			for _, cg := range f.Comments {
				for _, c := range cg.List {
					fc[pass.Fset.Position(c.Slash).Line] = c.Text
				}
			}
		}
		commentsByFile[tf] = fc
		return fc
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
		// The comment must appear on a line between the Optional and Computed field lines.
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

		fc := commentsFor(tf)
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
			if helpers.IsTrueConstant(pass, kv.Value) {
				hasOptional = true
				optionalPos = kv.Key.Pos()
			}
		case "Computed":
			if helpers.IsTrueConstant(pass, kv.Value) {
				hasComputed = true
				computedPos = kv.Key.Pos()
			}
		}
	}

	return
}
