package azg004

import (
	"fmt"
	"strings"

	"github.com/example/first"
	"github.com/hashicorp/go-azure-helpers/lang/response"
)

type groupedModel struct {
	Name   *string
	Status int
}

// Should be flagged: the file groups its imports gci-style (standard library first, then
// external packages) — the inserted import must land in sorted position inside the external
// group, not between the standard-library imports.
func invalidGroupedImports(m *groupedModel) {
	name := "" // want `pointer\.From`
	if m.Name != nil {
		name = *m.Name
	}
	_ = strings.TrimSpace(name)
	_ = fmt.Sprintf("%s %s %v", name, first.Value(), response.WasNotFound(m.Status))
}
