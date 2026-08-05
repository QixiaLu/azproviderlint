// Package AZD collects the data source checks.
package AZD

import (
	"golang.org/x/tools/go/analysis"

	AZD001 "github.com/katbyte/azproviderlint/checks/AZD/AZD001_data_source_empty_set_id"
	AZD002 "github.com/katbyte/azproviderlint/checks/AZD/AZD002_data_source_mark_as_gone"
)

// Checks contains all AZD (data source) analyzers.
var Checks = []*analysis.Analyzer{
	AZD001.Analyzer,
	AZD002.Analyzer,
}
