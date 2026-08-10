package azg003

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/virtualmachines"
)

// Should NOT be flagged: uses the generic pointer.ToEnum helper.
func validToEnum() *virtualmachines.VirtualMachinePriorityTypes {
	priority := "Spot"
	return pointer.ToEnum(virtualmachines.VirtualMachinePriorityTypes(priority))
}

// Should NOT be flagged: uses the generic pointer.ToEnum helper with a literal.
func validToEnumLiteral() *virtualmachines.OperatingSystemTypes {
	return pointer.ToEnum(virtualmachines.OperatingSystemTypes("Linux"))
}

// Should NOT be flagged: pointer.To of a plain string is fine.
func validPlainString() *string {
	return pointer.To("regular string")
}

// Should NOT be flagged: pointer.To of a plain int is fine.
func validPlainInt() *int {
	return pointer.To(42)
}

// Should be flagged: pointer.To with an explicit enum conversion of a variable.
func invalidPriority() *virtualmachines.VirtualMachinePriorityTypes {
	priority := "Spot"
	return pointer.To(virtualmachines.VirtualMachinePriorityTypes(priority)) // want `pointer\.To with an explicit go-azure-sdk enum conversion should use pointer\.ToEnum\[VirtualMachinePriorityTypes\] instead`
}

// Should be flagged: pointer.To with an explicit enum conversion of a literal.
func invalidOSType() *virtualmachines.OperatingSystemTypes {
	return pointer.To(virtualmachines.OperatingSystemTypes("Linux")) // want `pointer\.To with an explicit go-azure-sdk enum conversion should use pointer\.ToEnum\[OperatingSystemTypes\] instead`
}

// Should be flagged: pointer.To with an explicit enum conversion of a map lookup.
func invalidFromMap(config map[string]interface{}) *virtualmachines.VirtualMachinePriorityTypes {
	return pointer.To(virtualmachines.VirtualMachinePriorityTypes(config["priority"].(string))) // want `pointer\.To with an explicit go-azure-sdk enum conversion should use pointer\.ToEnum\[VirtualMachinePriorityTypes\] instead`
}
