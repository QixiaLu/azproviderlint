// Package pointer is a minimal stand-in for github.com/hashicorp/go-azure-helpers/lang/pointer
// used only by the AZG004 analysistest fixtures.
package pointer

// From returns the value pointed to by v, or the zero value of T when v is nil.
func From[T any](v *T) T {
	if v == nil {
		var zero T
		return zero
	}
	return *v
}
