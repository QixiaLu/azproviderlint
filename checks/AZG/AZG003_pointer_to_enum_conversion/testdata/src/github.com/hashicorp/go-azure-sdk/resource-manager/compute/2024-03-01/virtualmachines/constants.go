// Package virtualmachines is a minimal stand-in for a generated go-azure-sdk resource-manager
// package used only by the AZG003 analysistest fixtures. Enum types follow the generated SDK
// convention of a matching PossibleValuesFor<Name>() []string helper declared in constants.go.
package virtualmachines

type DiskCount int64

const (
	DiskCountOne DiskCount = 1
	DiskCountTwo DiskCount = 2
)

func PossibleValuesForDiskCount() []int64 {
	return []int64{
		int64(DiskCountOne),
		int64(DiskCountTwo),
	}
}

type OperatingSystemTypes string

const (
	OperatingSystemTypesLinux   OperatingSystemTypes = "Linux"
	OperatingSystemTypesWindows OperatingSystemTypes = "Windows"
)

func PossibleValuesForOperatingSystemTypes() []string {
	return []string{
		string(OperatingSystemTypesLinux),
		string(OperatingSystemTypesWindows),
	}
}

type OperatingSystemTypesList []string

type ResourceIdentifier string

type VirtualMachinePriorityTypes string

const (
	VirtualMachinePriorityTypesRegular VirtualMachinePriorityTypes = "Regular"
	VirtualMachinePriorityTypesSpot    VirtualMachinePriorityTypes = "Spot"
)

func PossibleValuesForVirtualMachinePriorityTypes() []string {
	return []string{
		string(VirtualMachinePriorityTypesRegular),
		string(VirtualMachinePriorityTypesSpot),
	}
}
