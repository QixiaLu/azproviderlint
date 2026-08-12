package azs004

import (
	tfvalidation "github.com/example/provider/tf/validation"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/virtualmachines"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// Should NOT be flagged: the SDK's possible-values helper always covers every value.
func validPossibleValuesHelper() {
	_ = validation.StringInSlice(virtualmachines.PossibleValuesForVirtualMachinePriorityTypes(), false)
}

// Should be flagged: a complete hand-written list still goes stale when the SDK adds a new
// value — the possible-values helper exists, so it should be used.
func invalidCompleteManualList() {
	_ = validation.StringInSlice([]string{ // want `enum validation for virtualmachines\.VirtualMachinePriorityTypes lists every value manually; use virtualmachines\.PossibleValuesForVirtualMachinePriorityTypes\(\) so new values are picked up automatically`
		string(virtualmachines.VirtualMachinePriorityTypesLow),
		string(virtualmachines.VirtualMachinePriorityTypesRegular),
		string(virtualmachines.VirtualMachinePriorityTypesSpot),
	}, false)
}

// Should be flagged: raw string literals count towards coverage, so this list is complete —
// but complete manual lists are still reported in favour of the helper.
func invalidMixedLiteralAndConstantsComplete() {
	_ = validation.StringInSlice([]string{ // want `lists every value manually`
		string(virtualmachines.VirtualMachinePriorityTypesLow),
		"Regular",
		string(virtualmachines.VirtualMachinePriorityTypesSpot),
	}, false)
}

// Should be flagged: the list references the enum but misses two of its values.
func invalidMissingValues() {
	_ = validation.StringInSlice([]string{ // want `enum validation for virtualmachines\.VirtualMachinePriorityTypes is missing VirtualMachinePriorityTypesRegular \("Regular"\), VirtualMachinePriorityTypesSpot \("Spot"\); use virtualmachines\.PossibleValuesForVirtualMachinePriorityTypes\(\)`
		string(virtualmachines.VirtualMachinePriorityTypesLow),
	}, false)
}

// Should be flagged: raw literals cover one value but another is still missing.
func invalidMixedStillMissing() {
	_ = validation.StringInSlice([]string{ // want `missing VirtualMachinePriorityTypesSpot \("Spot"\)`
		string(virtualmachines.VirtualMachinePriorityTypesLow),
		"Regular",
	}, false)
}

// Should be flagged: case-insensitive validation still has to cover every value.
func invalidIgnoreCase() {
	_ = validation.StringInSlice([]string{ // want `missing VirtualMachinePriorityTypesSpot \("Spot"\)`
		string(virtualmachines.VirtualMachinePriorityTypesLow),
		string(virtualmachines.VirtualMachinePriorityTypesRegular),
	}, true)
}

// Should be flagged: provider-internal wrappers of the plugin SDK helper resolve the same.
func invalidThroughWrapper() {
	_ = tfvalidation.StringInSlice([]string{ // want `missing VirtualMachinePriorityTypesSpot \("Spot"\)`
		string(virtualmachines.VirtualMachinePriorityTypesLow),
		string(virtualmachines.VirtualMachinePriorityTypesRegular),
	}, false)
}

// Should NOT be flagged: a non-constant element can contribute any value, so coverage of the
// list cannot be proven statically.
func validNonConstantElement(extra string) {
	_ = validation.StringInSlice([]string{
		string(virtualmachines.VirtualMachinePriorityTypesLow),
		extra,
	}, false)
}

// Should NOT be flagged: constants of two different enums form a deliberate union, which is
// out of scope.
func validTwoEnumUnion() {
	_ = validation.StringInSlice([]string{
		string(virtualmachines.VirtualMachinePriorityTypesLow),
		string(virtualmachines.CachingTypesNone),
	}, false)
}

// Should NOT be flagged: StorageTier has no possible-values helper, so it is not a closed
// enum and a subset of its constants is fine.
func validNotAClosedEnum() {
	_ = validation.StringInSlice([]string{
		string(virtualmachines.StorageTierHot),
	}, false)
}

// Should NOT be flagged: plain string lists with no enum constants are out of scope.
func validPlainStrings() {
	_ = validation.StringInSlice([]string{"a", "b"}, false)
}
