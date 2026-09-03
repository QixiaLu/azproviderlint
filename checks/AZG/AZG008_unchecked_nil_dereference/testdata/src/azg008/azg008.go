package azg008

import (
	"os"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
)

type Status string

type properties struct {
	Status *Status
	Count  *int
	Name   *string
}

type model struct {
	Properties *properties
}

type data struct{}

func (data) Set(string, interface{}) {}

func use(interface{}) {}

// Should be flagged: bare dereference of an optional field, converted to string — the fix is
// pointer.FromEnum.
func invalidEnumConversion(d data, props properties) {
	d.Set("status", string(*props.Status)) // want "dereference of possibly-nil `props.Status` may panic - add a nil check or use pointer.From"
}

// Should be flagged: bare dereference of a non-enum field — the fix is pointer.From.
func invalidPlainDeref(d data, props properties) {
	d.Set("count", *props.Count) // want "dereference of possibly-nil `props.Count` may panic - add a nil check or use pointer.From"
}

// Should be flagged: a deeper selector chain with no guard on the dereferenced link.
func invalidChain(d data, m model) {
	d.Set("name", *m.Properties.Name) // want "dereference of possibly-nil `m.Properties.Name` may panic - add a nil check or use pointer.From"
}

// Should be flagged: a guard on a prefix of the chain does not cover the dereferenced field.
func invalidPrefixGuardOnly(d data, m model) {
	if m.Properties != nil {
		d.Set("name", *m.Properties.Name) // want "dereference of possibly-nil `m.Properties.Name` may panic - add a nil check or use pointer.From"
	}
}

// Should NOT be flagged here: an assignment target needs the pointer itself — AZG009's.
func writeTargetIsAZG009(props properties) {
	*props.Count = 1
}

// Should NOT be flagged: enclosing if proves the field non-nil.
func validIfGuard(d data, props properties) {
	if props.Status != nil {
		d.Set("status", string(*props.Status))
	}
}

// Should NOT be flagged: the && left operand proves the field non-nil for the right.
func validShortCircuitAnd(props properties) bool {
	return props.Count != nil && *props.Count > 0
}

// Should NOT be flagged: the || left operand short-circuits when the field is nil.
func validShortCircuitOr(props properties) bool {
	return props.Count == nil || *props.Count == 0
}

// Should NOT be flagged: the else branch of a pure-|| nil condition.
func validElseBranch(d data, props properties) {
	if props.Status == nil || *props.Status == "" {
		d.Set("status", "")
	} else {
		d.Set("status", string(*props.Status))
	}
}

// Should NOT be flagged: an early return guards everything below it.
func validEarlyReturn(d data, props properties) {
	if props.Status == nil {
		return
	}
	d.Set("status", string(*props.Status))
}

// Should NOT be flagged: an early exit via os.Exit or panic also terminates.
func validEarlyExit(props properties) int {
	if props.Count == nil {
		os.Exit(1)
	}
	return *props.Count
}

// Should NOT be flagged: a pure-|| early return covers both fields.
func validOrEarlyReturn(props properties) int {
	if props.Count == nil || props.Name == nil {
		return 0
	}
	use(*props.Name)
	return *props.Count
}

// Should NOT be flagged: the variable only ever holds provably non-nil values.
func validNonNilSources() (int, string) {
	count := new(int)
	name := pointer.To("x")
	status := new(Status)
	_ = status
	return *count, *name
}

// Should be flagged: reassignment to an unknown value invalidates the earlier guard.
func invalidReassigned(props properties, next func() *int) int {
	if props.Count == nil {
		return 0
	}
	count := next()
	return *count // want "dereference of possibly-nil `count` may panic - add a nil check or use pointer.From"
}

// Should be flagged: reassignment inside the guarded body invalidates the enclosing guard.
func invalidReassignedInsideGuard(next func() *int) int {
	x := next()
	if x != nil {
		x = next()
		return *x // want "dereference of possibly-nil `x` may panic - add a nil check or use pointer.From"
	}
	return 0
}

// Should be flagged: reassignment in a nested block invalidates the earlier early return.
func invalidReassignedNestedBlock(props properties, next func() *int) int {
	if props.Count == nil {
		return 0
	}
	{
		props.Count = next()
		return *props.Count // want "dereference of possibly-nil `props.Count` may panic - add a nil check or use pointer.From"
	}
}

// Should NOT be flagged: the reassignment inside the guard is itself provably non-nil.
func validReassignedNonNil(next func() *int) int {
	x := next()
	if x != nil {
		x = new(int)
		return *x
	}
	return 0
}

// Should NOT be flagged: a switch case condition proves the field non-nil.
func validSwitchCase(props properties) int {
	switch {
	case props.Count != nil:
		return *props.Count
	}
	return 0
}

// Should NOT be flagged: a guard outside a closure still covers dereferences inside it.
func validGuardOutsideClosure(d data, props properties) {
	if props.Status != nil {
		f := func() { d.Set("status", string(*props.Status)) }
		f()
	}
}

// Should NOT be flagged: the alias's source chain is guarded by the enclosing if.
func validAliasedGuard(d data, props properties) {
	s := props.Status
	if props.Status != nil {
		d.Set("status", string(*s))
	}
}

// Should NOT be flagged: the alias itself is guarded.
func validGuardedAliasDirect(d data, props properties) {
	s := props.Status
	if s != nil {
		d.Set("status", string(*s))
	}
}

// Should NOT be flagged: for-loop condition proves non-nil inside the body.
type node struct{ next *node }

func validForCond(n *node) int {
	count := 0
	for n != nil {
		n = (*n).next
		count++
	}
	return count
}

// Should NOT be flagged: `x, err := f()` followed by an error exit — the Go contract makes
// the other results valid.
func validErrContract(parse func(string) (*properties, error)) int {
	props, err := parse("x")
	if err != nil {
		return 0
	}
	return len((*props).Name2())
}

// Should be flagged: the error is discarded, so nothing proves the pointer valid.
func invalidDiscardedErr(parse func(string) (*properties, error)) int {
	props, _ := parse("x")
	return len((*props).Name2()) // want "dereference of possibly-nil `props` may panic - add a nil check or use pointer.From"
}

func (properties) Name2() string { return "" }

// Should NOT be flagged: the alias's source chain was nil-checked before the assignment.
func validAliasOfGuardedChain(d data, m model) {
	if m.Properties == nil {
		return
	}
	payload := m.Properties
	d.Set("status", string((*payload).StatusString()))
}

// Should be flagged: only the alias's source prefix was checked, not the derefd field.
func invalidAliasFieldUnguarded(d data, m model) {
	if m.Properties == nil {
		return
	}
	payload := m.Properties
	d.Set("count", *payload.Count) // want "dereference of possibly-nil `payload.Count` may panic - add a nil check or use pointer.From"
}

// Should NOT be flagged: `x, ok := f()` followed by a !ok exit — the comma-ok contract.
func validOkContract(lookup func(string) (*int, bool)) int {
	suffix, ok := lookup("x")
	if !ok {
		return 0
	}
	return *suffix
}

func (properties) StatusString() string { return "" }

// Should NOT be flagged: a bare pointer parameter's nil contract belongs to its callers.
func validBareParam(id *int) int {
	return *id
}

// Should be flagged: a field dereference through a parameter is always in scope.
func invalidParamField(props *properties) int {
	return *props.Count // want "dereference of possibly-nil `props.Count` may panic - add a nil check or use pointer.From"
}
