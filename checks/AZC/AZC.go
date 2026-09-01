// Package AZC collects the client & SDK usage checks.
package AZC

import (
	"golang.org/x/tools/go/analysis"

	AZC001 "github.com/katbyte/azproviderlint/checks/AZC/AZC001_client_missing_resource_manager_endpoint"
)

// Checks contains all AZC (client & SDK usage) analyzers.
var Checks = []*analysis.Analyzer{
	AZC001.Analyzer,
}
