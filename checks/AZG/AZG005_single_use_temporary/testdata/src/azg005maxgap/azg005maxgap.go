package azg005maxgap

// Should be flagged with max-gap=3: the consumer is two lines below the declaration.
func withinGap(v string) string {
	s := v + "!" // want `"s" is only used by the statement on line \d+ and should be inlined`
	_ = len(v)
	return s
}

// Should NOT be flagged with max-gap=3: the consumer is five lines below the declaration.
func beyondGap(v string) string {
	s := v + "!"
	_ = len(v)
	// filler
	// filler
	// filler
	return s
}
