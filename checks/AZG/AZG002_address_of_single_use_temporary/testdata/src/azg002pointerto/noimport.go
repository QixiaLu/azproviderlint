package azg002pointerto

import (
	"strings"
)

type widget struct {
	name *string
}

// Should be flagged: the fix must add the pointer import, opening a new section after the
// standard-library imports.
func noImportAddressOf(s string) *widget {
	upper := strings.ToUpper(s)
	name := upper + "!" // want `"name" is only used as an address by the following statement and should be inlined with pointer\.To`
	return &widget{name: &name}
}
