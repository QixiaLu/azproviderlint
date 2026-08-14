package azg004

import (
	"strings"
)

type conf struct {
	Description *string
}

// Should be flagged: this file does not import the pointer package, so the suggested fix also
// inserts the import in sorted position.
func invalidNoPointerImport(c *conf) {
	description := "" // want `pointer\.From`
	if c.Description != nil {
		description = *c.Description
	}
	_ = strings.TrimSpace(description)
}
