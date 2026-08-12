// Package pointer is a minimal stand-in for github.com/hashicorp/go-azure-helpers/lang/pointer
// used only by the AZG003 analysistest fixtures.
package pointer

// To returns a pointer to v.
func To[T any](v T) *T { return &v }

// ToEnum returns a pointer to the enum value v.
func ToEnum[T ~string](v T) *T { return &v }
