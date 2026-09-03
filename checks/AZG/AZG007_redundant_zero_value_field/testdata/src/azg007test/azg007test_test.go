package azg007test

// With the tests flag on, redundant zero values in test files ARE reported.
var _ = []row{
	{
		name:     "empty",
		expected: "",  // want `redundant zero-value assignment to field "expected" - omit the field`
		count:    0,   // want `redundant zero-value assignment to field "count" - omit the field`
		ptr:      nil, // want `redundant nil assignment to pointer field "ptr" - omit the field`
	},
}
