package azg006

import "fmt"

type NetworkACLs struct {
	Name string
}

type NetworkRuleSet struct {
	Rules []string
}

type ACLList []NetworkACLs

func validate(input *NetworkRuleSet) error {
	return nil
}

// Should NOT be flagged: returns an empty slice on the nil-input guard
func flattenReturnsEmptySlice(input *NetworkRuleSet) []NetworkACLs {
	if input == nil {
		return []NetworkACLs{}
	}
	return []NetworkACLs{{Name: "test"}}
}

// Should NOT be flagged: expand functions are allowed to return nil
func expandCanReturnNil(input []NetworkACLs) *NetworkRuleSet {
	if len(input) == 0 {
		return nil
	}
	return &NetworkRuleSet{}
}

// Should NOT be flagged: the result is not a slice
func flattenToString(input *NetworkRuleSet) string {
	if input == nil {
		return ""
	}
	return "result"
}

// Should NOT be flagged: every slice position returns an empty slice
func flattenMultipleSlicesValid(input *NetworkRuleSet) ([]NetworkACLs, []any, error) {
	if input == nil {
		return []NetworkACLs{}, []any{}, nil
	}
	return []NetworkACLs{{Name: "test"}}, []any{"a"}, nil
}

// Should be flagged: returns nil for a slice result
func flattenReturnsNil(input *NetworkRuleSet) []NetworkACLs {
	if input == nil {
		return nil // want `flatten function "flattenReturnsNil" should return an empty slice instead of nil`
	}
	return []NetworkACLs{{Name: "test"}}
}

// Should be flagged: returns nil in the slice position of a multi-value return with a nil error
func flattenWithErrorNil(input *NetworkRuleSet) ([]NetworkACLs, error) {
	if input == nil {
		return nil, nil // want `flatten function "flattenWithErrorNil" should return an empty slice instead of nil`
	}
	return []NetworkACLs{{Name: "test"}}, nil
}

// Should be flagged: one of several slice positions returns nil with a nil error
func flattenMultipleSlicesOneNil(input *NetworkRuleSet) ([]NetworkACLs, []any, error) {
	if input == nil {
		return []NetworkACLs{}, nil, nil // want `flatten function "flattenMultipleSlicesOneNil" should return an empty slice instead of nil`
	}
	return []NetworkACLs{{Name: "test"}}, []any{"a"}, nil
}

// Should NOT be flagged on the error path (return nil, err); only the nil-input branch is flagged
func flattenWithErrorHandling(input *NetworkRuleSet) ([]NetworkACLs, error) {
	if input == nil {
		return nil, nil // want `flatten function "flattenWithErrorHandling" should return an empty slice instead of nil`
	}
	if err := validate(input); err != nil {
		return nil, err
	}
	return []NetworkACLs{{Name: "test"}}, nil
}

// Should NOT be flagged: the nil-input branch itself returns an error
func flattenNilInputReturnsError(input *NetworkRuleSet) ([]NetworkACLs, error) {
	if input == nil {
		return nil, fmt.Errorf("input is nil")
	}
	return []NetworkACLs{{Name: "test"}}, nil
}

// Should NOT be flagged: error path in a three-value return (return nil, nil, err)
func flattenThreeWithError(input *NetworkRuleSet) ([]NetworkACLs, []any, error) {
	if err := validate(input); err != nil {
		return nil, nil, err
	}
	return []NetworkACLs{{Name: "test"}}, []any{"a"}, nil
}

// Should be flagged: the flatten prefix match is case-insensitive
func FlattenUpperCase(input *NetworkRuleSet) []NetworkACLs {
	if input == nil {
		return nil // want `flatten function "FlattenUpperCase" should return an empty slice instead of nil`
	}
	return []NetworkACLs{{Name: "test"}}
}

// Should be flagged: a multi-name slice field expands to every return position it covers
func flattenMultiName(input *NetworkRuleSet) (a, b []NetworkACLs, err error) {
	if input == nil {
		return nil, nil, nil // want `flatten function "flattenMultiName" should return an empty slice instead of nil`
	}
	return []NetworkACLs{{Name: "test"}}, []NetworkACLs{{Name: "test2"}}, nil
}

// Should be flagged: named slice types are detected via the type checker, not the syntax
func flattenNamedSliceType(input *NetworkRuleSet) ACLList {
	if input == nil {
		return nil // want `flatten function "flattenNamedSliceType" should return an empty slice instead of nil`
	}
	return ACLList{{Name: "test"}}
}

// Should NOT be flagged: a nil return inside a closure answers to the closure, not the flatten function
func flattenWithClosure(input *NetworkRuleSet) []NetworkACLs {
	build := func() []string {
		return nil
	}
	_ = build
	if input == nil {
		return []NetworkACLs{}
	}
	return []NetworkACLs{{Name: "test"}}
}

// Should NOT be flagged: naked returns are intentionally out of scope
func flattenNakedReturn(input *NetworkRuleSet) (ret []NetworkACLs) {
	if input == nil {
		return
	}
	ret = []NetworkACLs{{Name: "test"}}
	return
}
