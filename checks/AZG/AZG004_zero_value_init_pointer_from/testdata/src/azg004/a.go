package azg004

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
)

type Props struct {
	Enabled *bool
	Name    *string
	Count   *int
}

// Should NOT be flagged: already uses pointer.From.
func validFrom(props *Props) bool {
	return pointer.From(props.Enabled)
}

// Should NOT be flagged: assigns the result of pointer.From to a local.
func validFromLocal(props *Props) {
	enabled := pointer.From(props.Enabled)
	useBool(enabled)
}

// Should be flagged: zero-init bool followed by nil check + dereference.
func invalidBool(props *Props) {
	enabled := false // want `pointer\.From`
	if props.Enabled != nil {
		enabled = *props.Enabled
	}
	useBool(enabled)
}

// Should be flagged: zero-init string.
func invalidString(props *Props) {
	name := "" // want `pointer\.From`
	if props.Name != nil {
		name = *props.Name
	}
	useString(name)
}

// Should be flagged: zero-init int.
func invalidInt(props *Props) {
	count := 0 // want `pointer\.From`
	if props.Count != nil {
		count = *props.Count
	}
	useInt(count)
}

// Should be flagged: the variable is used more than once afterwards.
func invalidMultipleUses(props *Props) {
	enabled := false // want `pointer\.From`
	if props.Enabled != nil {
		enabled = *props.Enabled
	}
	useBool(enabled)
	useBool(enabled)
}

// Should NOT be flagged: has an else branch.
func edgeElseBranch(props *Props) {
	enabled := false
	if props.Enabled != nil {
		enabled = *props.Enabled
	} else {
		enabled = true
	}
	useBool(enabled)
}

// Should NOT be flagged: not initialized to a zero value.
func edgeNonZeroInit(props *Props) {
	enabled := true
	if props.Enabled != nil {
		enabled = *props.Enabled
	}
	useBool(enabled)
}

// Should NOT be flagged: the if body assigns a different variable.
func edgeDifferentVar(props *Props) {
	enabled := false
	if props.Enabled != nil {
		other := *props.Enabled
		useBool(other)
	}
	useBool(enabled)
}

// Should NOT be flagged: the assignment and if are not adjacent.
func edgeNotAdjacent(props *Props) {
	enabled := false
	doSomething()
	if props.Enabled != nil {
		enabled = *props.Enabled
	}
	useBool(enabled)
}

// Should NOT be flagged: the if body uses := instead of =.
func edgeShortDeclInBody(props *Props) {
	_ = false
	if props.Enabled != nil {
		enabled := *props.Enabled
		useBool(enabled)
	}
}

func useBool(b bool)     {}
func useString(s string) {}
func useInt(i int)       {}
func doSomething()       {}
