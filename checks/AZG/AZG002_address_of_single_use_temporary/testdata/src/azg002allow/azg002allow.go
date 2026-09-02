package azg002allow

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
)

// Should NOT be flagged: allow: pointer.To tolerates existing helper calls.
func validAllowed(v int) *int {
	return pointer.To(v)
}

// Should be flagged: the temporary is still inlined, with new().
func invalidTemp(v int) *int {
	count := v + 1 // want `"count" is only used as an address by the following statement and should be inlined with new\(\)`
	return &count
}
