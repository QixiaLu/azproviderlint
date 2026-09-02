package azs004allowextra

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/virtualmachines"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// With -AZS004.allow-extra-values, a deliberate superset is not reported.
func validDeliberateSuperset() {
	_ = validation.StringInSlice([]string{
		string(virtualmachines.VirtualMachinePriorityTypesLow),
		string(virtualmachines.VirtualMachinePriorityTypesRegular),
		string(virtualmachines.VirtualMachinePriorityTypesSpot),
		"Legacy",
	}, false)
}

// Missing values are still reported; the suppressed extras clause is simply absent, and the
// appending advice stays because the list does carry extras.
func invalidMissingStillReported() {
	_ = validation.StringInSlice([]string{ // want `enum validation for virtualmachines\.VirtualMachinePriorityTypes is missing VirtualMachinePriorityTypesRegular \("Regular"\), VirtualMachinePriorityTypesSpot \("Spot"\); use virtualmachines\.PossibleValuesForVirtualMachinePriorityTypes\(\), appending any deliberate extras`
		string(virtualmachines.VirtualMachinePriorityTypesLow),
		"Legacy",
	}, false)
}
