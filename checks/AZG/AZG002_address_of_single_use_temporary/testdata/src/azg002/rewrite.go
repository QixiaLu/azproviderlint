package azg002

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
)

type widget struct {
	name  *string
	count *int
}

// Should be flagged: pointer.To should be the new() builtin; the From call keeps the import.
func rewriteKeepsImport(s string) widget {
	w := widget{
		name:  pointer.To(s), // want `pointer\.To should be replaced with the new\(\.\.\.\) builtin`
		count: pointer.To(1), // want `pointer\.To should be replaced with the new\(\.\.\.\) builtin`
	}
	_ = pointer.From(w.name)
	return w
}
