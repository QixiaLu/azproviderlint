package azr004

type ThingId struct{}

func (ThingId) ID() string { return "" }

func Match(a, b interface{}) bool { return true }

// Should be flagged: comparing Resource IDs with ==
func badEqual(a, b ThingId) bool {
	return a.ID() == b.ID() // want `Resource IDs should not be compared with == or !=, use resourceids\.Match instead`
}

// Should be flagged: comparing Resource IDs with !=
func badNotEqual(a, b ThingId) bool {
	return a.ID() != b.ID() // want `Resource IDs should not be compared with == or !=, use resourceids\.Match instead`
}

// Should be flagged: one side is an ID() call
func badOneSide(a ThingId, s string) bool {
	return a.ID() == s // want `Resource IDs should not be compared with == or !=, use resourceids\.Match instead`
}

// Should NOT be flagged: using a Match helper
func goodMatch(a, b ThingId) bool {
	return Match(a, b)
}

// Should NOT be flagged: comparing plain strings
func goodStrings(a, b string) bool {
	return a == b
}
