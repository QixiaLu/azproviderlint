package azg008

// _test.go files are checked by default (tests=true); no import block, so no fix is offered.
func derefInTest(props properties) int {
	return *props.Count // want "dereference of possibly-nil `props.Count` may panic - add a nil check or use pointer.From"
}
