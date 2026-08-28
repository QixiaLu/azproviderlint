// Package validation is a minimal stand-in for a provider-internal wrapper around the plugin
// SDK's helper/validation package (e.g. azurerm's internal/tf/validation), used only by the
// AZS004 analysistest fixtures.
package validation

import (
	sdkvalidation "github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// StringInSlice is a stub matching azurerm's wrapper of the real helper.
func StringInSlice(valid []string, ignoreCase bool) func(interface{}, string) ([]string, []error) {
	return sdkvalidation.StringInSlice(valid, ignoreCase)
}

// StringInEnumSlice is a stub matching azurerm's generic wrapper accepting a track-1 style
// typed enum slice directly; its presence switches AZS004's typed-helper advice to name it.
func StringInEnumSlice[T ~string](valid []T, ignoreCase bool) func(interface{}, string) ([]string, []error) {
	strs := make([]string, len(valid))
	for i, v := range valid {
		strs[i] = string(v)
	}
	return sdkvalidation.StringInSlice(strs, ignoreCase)
}
