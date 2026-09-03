// Package AZV collects the validation checks.
package AZV

import (
	"golang.org/x/tools/go/analysis"

	AZV001 "github.com/katbyte/azproviderlint/checks/AZV/AZV001_error_should_describe_expected_format"
)

// Checks contains all AZV (validation) analyzers.
var Checks = []*analysis.Analyzer{
	AZV001.Analyzer,
}
