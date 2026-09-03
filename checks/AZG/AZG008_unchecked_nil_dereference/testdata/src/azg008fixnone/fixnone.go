package azg008fixnone

type properties struct {
	Count *int
}

// Still flagged with fix-with=none; the diagnostic just carries no suggested fix.
func deref(props properties) int {
	return *props.Count // want "dereference of possibly-nil `props.Count` may panic - add a nil check or use pointer.From"
}
