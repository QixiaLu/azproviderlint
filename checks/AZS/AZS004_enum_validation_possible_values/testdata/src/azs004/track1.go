package azs004

import (
	tfvalidation "github.com/example/provider/tf/validation"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/virtualmachines"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// Should be flagged: DiskCreateOptionTypes is a track-1 style enum whose possible-values
// helper returns a typed slice ([]DiskCreateOptionTypes), so the advice must route through
// go-azure-helpers' enum-slice conversion instead of passing the helper to StringInSlice
// directly.
func invalidTrackOneIncompleteList() {
	_ = validation.StringInSlice([]string{ // want `enum validation for virtualmachines\.DiskCreateOptionTypes is missing DiskCreateOptionTypesFromImage \("FromImage"\); use pointer\.FromEnumSlice\(pointer\.To\(virtualmachines\.PossibleDiskCreateOptionTypesValues\(\)\)\)`
		string(virtualmachines.DiskCreateOptionTypesAttach),
		string(virtualmachines.DiskCreateOptionTypesEmpty),
	}, false)
}

// Should be flagged: a complete list of a track-1 enum still goes stale, with the same
// conversion-flavoured advice.
func invalidTrackOneCompleteList() {
	_ = validation.StringInSlice([]string{ // want `enum validation for virtualmachines\.DiskCreateOptionTypes lists every value manually; use pointer\.FromEnumSlice\(pointer\.To\(virtualmachines\.PossibleDiskCreateOptionTypesValues\(\)\)\) so new values are picked up automatically`
		string(virtualmachines.DiskCreateOptionTypesAttach),
		string(virtualmachines.DiskCreateOptionTypesEmpty),
		string(virtualmachines.DiskCreateOptionTypesFromImage),
	}, false)
}

// Should be flagged: the wrapper package exports StringInEnumSlice, so the track-1 advice
// names it instead of the go-azure-helpers conversion — echoing the call's ignoreCase
// argument so the suggested replacement call is complete.
func invalidTrackOneWithWrapperHelper() {
	_ = tfvalidation.StringInSlice([]string{ // want `enum validation for virtualmachines\.DiskCreateOptionTypes is missing DiskCreateOptionTypesFromImage \("FromImage"\); use validation\.StringInEnumSlice\(virtualmachines\.PossibleDiskCreateOptionTypesValues\(\), true\)`
		string(virtualmachines.DiskCreateOptionTypesAttach),
		string(virtualmachines.DiskCreateOptionTypesEmpty),
	}, true)
}

// Should NOT be flagged: already using the wrapper with the typed possible-values helper.
func validTrackOneWrapperHelper() {
	_ = tfvalidation.StringInEnumSlice(virtualmachines.PossibleDiskCreateOptionTypesValues(), false)
}
