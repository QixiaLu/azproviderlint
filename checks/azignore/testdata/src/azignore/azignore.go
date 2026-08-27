package azignore

// Should NOT be flagged: directive at the end of the line
func ignoredSameLine() {
	err := doSomething() //azignore:AZG001
	if err != nil {
		panic(err)
	}
}

// Should NOT be flagged: directive on the line immediately preceding
func ignoredLineAbove() {
	//azignore:AZG001
	err := doSomething()
	if err != nil {
		panic(err)
	}
}

// Should NOT be flagged: directive listing multiple checks
func ignoredList() {
	err := doSomething() //azignore:AZR001, AZG001
	if err != nil {
		panic(err)
	}
}

// Should NOT be flagged: directive with a reason
func ignoredWithReason() {
	err := doSomething() //azignore:AZG001 - legacy pattern, refactor tracked in #42
	if err != nil {
		panic(err)
	}
}

// Should NOT be flagged: multi-rule directive with an em dash reason
func ignoredListWithReason() {
	err := doSomething() //azignore:AZR001, AZG001 — deliberate
	if err != nil {
		panic(err)
	}
}

// Should NOT be flagged: gofmt adds a space after '//' when the rule list starts with an
// uppercase letter — the spaced form must still suppress
func ignoredGofmtSpaced() {
	err := doSomething() // azignore:AZG001 - spaced by gofmt
	if err != nil {
		panic(err)
	}
}

// Should be flagged: directive names a different check
func notIgnoredOtherCheck() {
	err := doSomething() //azignore:AZR001 // want `'err' assignment should be combined with the following 'if err != nil' into a single 'if' init statement`
	if err != nil {
		panic(err)
	}
}

// Should be flagged: directive is too far above the offending line
func notIgnoredTooFar() {
	//azignore:AZG001

	err := doSomething() // want `'err' assignment should be combined with the following 'if err != nil' into a single 'if' init statement`
	if err != nil {
		panic(err)
	}
}

func doSomething() error {
	return nil
}
