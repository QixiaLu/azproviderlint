package azg008notests

type properties struct {
	Count *int
}

// Not flagged: tests=false skips _test.go files.
func derefInSkippedTest(props properties) int {
	return *props.Count
}
