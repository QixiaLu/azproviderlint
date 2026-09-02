// Package validation is a minimal stand-in for the terraform-plugin-sdk helper/validation
// package used only by the AZS004 analysistest fixtures.
package validation

// StringInSlice is a stub matching the real helper's signature.
func StringInSlice(valid []string, ignoreCase bool) func(interface{}, string) ([]string, []error) {
	return func(i interface{}, k string) ([]string, []error) {
		return nil, nil
	}
}
