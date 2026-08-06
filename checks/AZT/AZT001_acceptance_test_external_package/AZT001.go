// Package AZT001 defines an analyzer that reports resource and data source acceptance test
// files that do not use an external _test package.
package AZT001

import (
	"path/filepath"
	"regexp"
	"strings"

	"golang.org/x/tools/go/analysis"
)

// Analyzer checks that acceptance test files for resources and data sources use an external
// test package (`package foo_test` rather than `package foo`) to prevent a circular
// dependency between the service package and the acceptance test framework.
var Analyzer = &analysis.Analyzer{
	Name: "AZT001",
	Doc:  "check that resource and data source acceptance test files use an external _test package to prevent circular dependencies",
	URL:  "https://github.com/katbyte/azproviderlint/blob/main/checks/AZT/AZT001_acceptance_test_external_package/README.md",
	Run:  run,
}

var acceptanceTestFile = regexp.MustCompile(`(resource|data_source).*_test\.go$`)

func run(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		filename := filepath.Base(pass.Fset.Position(file.Pos()).Filename)
		if !acceptanceTestFile.MatchString(filename) {
			continue
		}

		if strings.HasSuffix(file.Name.Name, "_test") {
			continue
		}

		pass.Reportf(file.Name.Pos(),
			"acceptance test files for resources and data sources must use a _test package to prevent circular dependencies")
	}

	return nil, nil
}
