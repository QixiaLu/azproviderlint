package azg001

import "fmt"

// Should be flagged: _, err assignment followed by if err != nil
func badExample() {
	_, err := fmt.Println("hello") // want `'_, err' assignment should be combined with the following 'if err != nil' into a single 'if' init statement`
	if err != nil {
		panic(err)
	}
}

// Should be flagged: using = instead of :=
func badExampleAssign() {
	var err error
	_, err = fmt.Println("hello") // want `'_, err' assignment should be combined with the following 'if err != nil' into a single 'if' init statement`
	if err != nil {
		panic(err)
	}
}

// Should NOT be flagged: already combined
func goodCombined() {
	if _, err := fmt.Println("hello"); err != nil {
		panic(err)
	}
}

// Should NOT be flagged: not followed by if err != nil
func goodNoIf() {
	_, err := fmt.Println("hello")
	_ = err
}

// Should NOT be flagged: the error is used between assignment and if
func goodIntervening() {
	_, err := fmt.Println("hello")
	fmt.Println("something else")
	if err != nil {
		panic(err)
	}
}

// Should NOT be flagged: more than 2 return values
func goodMultiReturn() {
	a, b, err := multiReturn()
	if err != nil {
		panic(err)
	}
	_ = a
	_ = b
}

// Should NOT be flagged: first value is not blank
func goodNamedReturn() {
	n, err := fmt.Println("hello")
	if err != nil {
		panic(err)
	}
	_ = n
}

// Should NOT be flagged: if has a different condition
func goodDifferentCondition() {
	_, err := fmt.Println("hello")
	if err == nil {
		fmt.Println("ok")
	}
}

// Should NOT be flagged: if already has an init statement
func goodIfWithInit() {
	_, err := fmt.Println("hello")
	if x := 1; err != nil {
		panic(err)
		_ = x
	}
}

func multiReturn() (int, int, error) {
	return 1, 2, nil
}
