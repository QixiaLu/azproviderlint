package azs004

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/virtualmachines"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// Should be flagged: the list covers every enum value but also carries an extra value that is
// not part of the enum — a plain swap to the helper would drop it, so the advice says to
// append deliberate extras.
func invalidSupersetList() {
	_ = validation.StringInSlice([]string{ // want `enum validation for virtualmachines\.VirtualMachinePriorityTypes has extra values not in the enum: "Legacy"; use virtualmachines\.PossibleValuesForVirtualMachinePriorityTypes\(\), appending any deliberate extras`
		string(virtualmachines.VirtualMachinePriorityTypesLow),
		string(virtualmachines.VirtualMachinePriorityTypesRegular),
		string(virtualmachines.VirtualMachinePriorityTypesSpot),
		"Legacy",
	}, false)
}

// Should be flagged: missing an enum value and carrying an extra one — both clauses reported.
// The extra here is the kind of typo (wrong case) the extras clause is designed to surface.
func invalidMissingAndExtra() {
	_ = validation.StringInSlice([]string{ // want `enum validation for virtualmachines\.VirtualMachinePriorityTypes is missing VirtualMachinePriorityTypesRegular \("Regular"\), VirtualMachinePriorityTypesSpot \("Spot"\) and has extra values not in the enum: "low"; use virtualmachines\.PossibleValuesForVirtualMachinePriorityTypes\(\), appending any deliberate extras`
		string(virtualmachines.VirtualMachinePriorityTypesLow),
		"low",
	}, false)
}
