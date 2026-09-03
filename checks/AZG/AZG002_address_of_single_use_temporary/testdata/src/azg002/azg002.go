package azg002

type Enum string

type props struct {
	Name    *string
	Count   *int
	Enabled *bool
}

func use(interface{}) {}

func newName() string { return "x" }

// Should be flagged: the temporary's only use is its address in the next statement's return.
func invalidReturn() *string {
	name := "hello" // want `"name" is only used as an address by the following statement and should be inlined with new\(\)`
	return &name
}

// Should be flagged: the address is a composite literal field value.
func invalidField() props {
	name := "hello" // want `"name" is only used as an address by the following statement and should be inlined with new\(\)`
	return props{Name: &name}
}

// Should be flagged: the address is assigned to a struct field.
func invalidAssign(p *props) {
	count := 3 // want `"count" is only used as an address by the following statement and should be inlined with new\(\)`
	p.Count = &count
}

// Should be flagged: the address is a call argument.
func invalidCallArg() {
	enabled := true // want `"enabled" is only used as an address by the following statement and should be inlined with new\(\)`
	use(&enabled)
}

// Should be flagged: the address is taken two statements later, within max-gap.
func invalidLater(p *props) {
	name := "x" // want `"name" is only used as an address by the statement on line \d+ and should be inlined with new\(\)`
	use(p)
	p.Name = &name
}

// Should be flagged: a conversion initializer is side-effect free.
func invalidConversion(s string) *Enum {
	e := Enum(s) // want `"e" is only used as an address by the following statement and should be inlined with new\(\)`
	return &e
}

// Should NOT be flagged: the initializer calls a function, so inlining would move the call
// past the intervening statements' effects.
func validCallInitializer() *string {
	name := newName()
	return &name
}

// Should NOT be flagged: the variable is used beyond the address-of.
func validSecondUse(p *props) {
	name := "x"
	use(name)
	p.Name = &name
}

// Should NOT be flagged: the variable is reassigned before its address is taken.
func validReassigned(p *props) {
	name := "x"
	name = "y"
	p.Name = &name
}

// Should NOT be flagged: an intervening statement writes to an operand of the initializer.
func validInterveningWrite(p *props, s string) {
	name := s + "!"
	s = "changed"
	p.Name = &name
}

// Should NOT be flagged: the address is captured by a function literal, which defers the
// initializer's evaluation to whenever the closure runs.
func validClosureCapture() {
	name := "x"
	defer func() { use(&name) }()
}

// Should NOT be flagged: multi-line initializers are out of scope.
func validMultiLine() *[]string {
	parts := []string{
		"a",
	}
	return &parts
}

// Should NOT be flagged: a channel receive in the initializer is a side effect.
func validReceive(ch chan string, p *props) {
	name := <-ch
	p.Name = &name
}
