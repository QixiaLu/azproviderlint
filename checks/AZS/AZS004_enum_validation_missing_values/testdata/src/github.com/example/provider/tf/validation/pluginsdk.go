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
