package azs004allowmissing

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/virtualmachines"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// With -AZS004.allow-missing-values, a deliberate subset is not reported.
func validDeliberateSubset() {
	_ = validation.StringInSlice([]string{
		string(virtualmachines.VirtualMachinePriorityTypesLow),
	}, false)
}

// A complete manual list is still reported: switching to the helper is a pure win.
func invalidCompleteManualList() {
	_ = validation.StringInSlice([]string{ // want `lists every value manually`
		string(virtualmachines.VirtualMachinePriorityTypesLow),
		string(virtualmachines.VirtualMachinePriorityTypesRegular),
		string(virtualmachines.VirtualMachinePriorityTypesSpot),
	}, false)
}

// Extra values are still reported; the suppressed missing clause is simply absent.
func invalidExtraStillReported() {
	_ = validation.StringInSlice([]string{ // want `enum validation for virtualmachines\.VirtualMachinePriorityTypes has extra values not in the enum: "Legacy"; use virtualmachines\.PossibleValuesForVirtualMachinePriorityTypes\(\), appending any deliberate extras`
		string(virtualmachines.VirtualMachinePriorityTypesLow),
		"Legacy",
	}, false)
}
