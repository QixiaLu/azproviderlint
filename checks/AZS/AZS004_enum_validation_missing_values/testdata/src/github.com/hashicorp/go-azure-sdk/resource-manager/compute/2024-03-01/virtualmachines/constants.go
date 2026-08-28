// Package virtualmachines is a minimal stand-in for a generated go-azure-sdk package used
// only by the AZS004 analysistest fixtures.
package virtualmachines

type VirtualMachinePriorityTypes string

const (
	VirtualMachinePriorityTypesLow     VirtualMachinePriorityTypes = "Low"
	VirtualMachinePriorityTypesRegular VirtualMachinePriorityTypes = "Regular"
	VirtualMachinePriorityTypesSpot    VirtualMachinePriorityTypes = "Spot"
)

func PossibleValuesForVirtualMachinePriorityTypes() []string {
	return []string{
		string(VirtualMachinePriorityTypesLow),
		string(VirtualMachinePriorityTypesRegular),
		string(VirtualMachinePriorityTypesSpot),
	}
}

type CachingTypes string

const (
	CachingTypesNone      CachingTypes = "None"
	CachingTypesReadOnly  CachingTypes = "ReadOnly"
	CachingTypesReadWrite CachingTypes = "ReadWrite"
)

func PossibleValuesForCachingTypes() []string {
	return []string{
		string(CachingTypesNone),
		string(CachingTypesReadOnly),
		string(CachingTypesReadWrite),
	}
}

// StorageTier is a named string type with constants but no possible-values helper, so AZS004
// must not treat it as a closed enum.
type StorageTier string

const (
	StorageTierHot  StorageTier = "Hot"
	StorageTierCool StorageTier = "Cool"
)

// DiskCreateOptionTypes mimics a track-1 style enum: its possible-values helper returns a
// typed slice rather than []string, so AZS004 must not treat it as a closed enum — the
// "use the helper" advice would not compile inside StringInSlice.
type DiskCreateOptionTypes string

const (
	DiskCreateOptionTypesAttach    DiskCreateOptionTypes = "Attach"
	DiskCreateOptionTypesEmpty     DiskCreateOptionTypes = "Empty"
	DiskCreateOptionTypesFromImage DiskCreateOptionTypes = "FromImage"
)

func PossibleDiskCreateOptionTypesValues() []DiskCreateOptionTypes {
	return []DiskCreateOptionTypes{
		DiskCreateOptionTypesAttach,
		DiskCreateOptionTypesEmpty,
		DiskCreateOptionTypesFromImage,
	}
}
