// Package AZG collects the general Go style & readability checks.
package AZG

import (
	"golang.org/x/tools/go/analysis"

	AZG001 "github.com/katbyte/azproviderlint/checks/AZG/AZG001_combine_err_assignment_and_check"
	AZG002 "github.com/katbyte/azproviderlint/checks/AZG/AZG002_error_should_describe_expected_format"
)

// Checks contains all AZG (general Go style & readability) analyzers.
var Checks = []*analysis.Analyzer{
	AZG001.Analyzer,
	AZG002.Analyzer,
}
