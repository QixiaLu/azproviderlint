package azg006maxargs

func consume2(k string, v interface{})            {}
func consume3(k string, v interface{}, extra int) {}

func flattenThing(v string) []string { return []string{v} }

// Should be flagged with maximum-arguments=2: the call has two arguments.
func invalidTwoArgs(v string) {
	apns := flattenThing(v) // want `"apns" is only used by the following call and should be inlined`
	consume2("k", apns)
}

// Should NOT be flagged with maximum-arguments=2: the call has three arguments.
func validThreeArgs(v string) {
	apns := flattenThing(v)
	consume3("k", apns, 1)
}
