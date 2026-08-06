// Package azignore implements '//azignore:AZX001' comment directives, letting individual
// checks be suppressed per line without disabling every azproviderlint check on that line the
// way '//nolint:azproviderlint' does. Directives work under any driver (golangci-lint plugin,
// standalone binary or go vet) as filtering happens inside the analyzers themselves.
package azignore

import (
	"strings"

	"golang.org/x/tools/go/analysis"
)

const prefix = "azignore:"

// Wrap replaces each analyzer's Run with one that drops diagnostics on lines carrying a
// '//azignore:<Name>' comment, either at the end of the line or on the line immediately
// preceding it. Multiple checks can be listed: '//azignore:AZG001,AZR001'.
func Wrap(analyzers []*analysis.Analyzer) []*analysis.Analyzer {
	for _, a := range analyzers {
		wrap(a)
	}
	return analyzers
}

func wrap(a *analysis.Analyzer) {
	run := a.Run
	name := a.Name
	a.Run = func(pass *analysis.Pass) (any, error) {
		ignored := ignoredLines(pass, name)
		report := pass.Report
		pass.Report = func(d analysis.Diagnostic) {
			pos := pass.Fset.Position(d.Pos)
			if lines, ok := ignored[pos.Filename]; ok && lines[pos.Line] {
				return
			}
			report(d)
		}
		return run(pass)
	}
}

// ignoredLines collects, per filename, the lines on which diagnostics from the named
// analyzer are suppressed: the line of each matching directive and the line below it.
func ignoredLines(pass *analysis.Pass, name string) map[string]map[int]bool {
	ignored := map[string]map[int]bool{}

	for _, file := range pass.Files {
		for _, group := range file.Comments {
			for _, comment := range group.List {
				text := strings.TrimSpace(strings.TrimPrefix(comment.Text, "//"))
				if !strings.HasPrefix(text, prefix) {
					continue
				}

				matched := false
				for n := range strings.SplitSeq(strings.TrimPrefix(text, prefix), ",") {
					if strings.TrimSpace(n) == name {
						matched = true
						break
					}
				}
				if !matched {
					continue
				}

				pos := pass.Fset.Position(comment.Pos())
				if ignored[pos.Filename] == nil {
					ignored[pos.Filename] = map[int]bool{}
				}
				ignored[pos.Filename][pos.Line] = true
				ignored[pos.Filename][pos.Line+1] = true
			}
		}
	}

	return ignored
}
