// Package checks exposes all azproviderlint analyzers, grouped by category.
package checks

import (
	"slices"

	"github.com/katbyte/azproviderlint/checks/AZC"
	"github.com/katbyte/azproviderlint/checks/AZD"
	"github.com/katbyte/azproviderlint/checks/AZG"
	"github.com/katbyte/azproviderlint/checks/AZR"
	"github.com/katbyte/azproviderlint/checks/AZS"
	"github.com/katbyte/azproviderlint/checks/AZT"
	"github.com/katbyte/azproviderlint/checks/azignore"
)

// All contains every azproviderlint analyzer across all categories, each honouring
// '//azignore:<Name>' comment directives.
var All = azignore.Wrap(slices.Concat(
	AZC.Checks,
	AZD.Checks,
	AZG.Checks,
	AZR.Checks,
	AZS.Checks,
	AZT.Checks,
))
