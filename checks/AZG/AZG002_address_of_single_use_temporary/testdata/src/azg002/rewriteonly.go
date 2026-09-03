package azg002

import (
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
)

// Should be flagged: the only pointer reference is the rewritten call, so the fix also
// deletes the import.
func rewriteDropsImport(s string) *string {
	return pointer.To(strings.ToUpper(s)) // want `pointer\.To should be replaced with the new\(\.\.\.\) builtin`
}
