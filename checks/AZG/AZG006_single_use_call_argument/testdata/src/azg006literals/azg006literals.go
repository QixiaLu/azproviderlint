package azg006literals

type resourceData struct{}

func (resourceData) Set(k string, v interface{}) error { return nil }

func flattenThing(v string) []string { return []string{v} }

// Should be flagged with only-when-literals: the sibling is a literal.
func invalidLiteralSibling(d resourceData, v string) {
	apns := flattenThing(v) // want `"apns" is only used by the following call and should be inlined`
	d.Set("apns_credential", apns)
}

// Should NOT be flagged with only-when-literals: the sibling is a plain identifier.
func validIdentSibling(d resourceData, v, key string) {
	apns := flattenThing(v)
	d.Set(key, apns)
}
