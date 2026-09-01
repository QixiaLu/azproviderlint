package azg006maxgap

func consume(k string, v interface{}) {}

func flattenThing(v string) []string { return []string{v} }

// Should be flagged with max-gap=3: the consuming call is two lines below the declaration.
func withinGap(v string, n *int) {
	apns := flattenThing(v) // want `"apns" is only used by the call on line \d+ and should be inlined`
	*n = 1
	consume("k", apns)
}

// Should NOT be flagged with max-gap=3: the consuming call is five lines below the declaration.
func beyondGap(v string, n *int) {
	apns := flattenThing(v)
	*n = 1
	// filler
	// filler
	// filler
	consume("k", apns)
}
