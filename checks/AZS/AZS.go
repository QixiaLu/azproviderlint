// Package AZS collects the schema & typed SDK model checks.
package AZS

import (
	"golang.org/x/tools/go/analysis"

	AZS001 "github.com/katbyte/azproviderlint/checks/AZS/AZS001_typed_sdk_model_64bit_types"
)

// Checks contains all AZS (schema & typed SDK model) analyzers.
var Checks = []*analysis.Analyzer{
	AZS001.Analyzer,
}
