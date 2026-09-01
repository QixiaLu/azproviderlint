// Package AZG collects the general Go style & readability checks.
package AZG

import (
	"golang.org/x/tools/go/analysis"

	AZG001 "github.com/katbyte/azproviderlint/checks/AZG/AZG001_combine_err_assignment_and_check"
	AZG002 "github.com/katbyte/azproviderlint/checks/AZG/AZG002_error_should_describe_expected_format"
	AZG003 "github.com/katbyte/azproviderlint/checks/AZG/AZG003_pointer_to_enum_conversion"
	AZG004 "github.com/katbyte/azproviderlint/checks/AZG/AZG004_zero_value_init_pointer_from"
	AZG005 "github.com/katbyte/azproviderlint/checks/AZG/AZG005_single_use_temporary"
)

// Checks contains all AZG (general Go style & readability) analyzers.
var Checks = []*analysis.Analyzer{
	AZG001.Analyzer,
	AZG002.Analyzer,
	AZG003.Analyzer,
	AZG004.Analyzer,
	AZG005.Analyzer,
}
