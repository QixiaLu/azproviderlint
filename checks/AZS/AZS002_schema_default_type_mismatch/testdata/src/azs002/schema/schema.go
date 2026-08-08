// Package schema is a minimal stub of the plugin SDK's helper/schema for the analyzer tests.
package schema

type ValueType int

const (
	TypeInvalid ValueType = iota
	TypeBool
	TypeInt
	TypeFloat
	TypeString
	TypeList
	TypeMap
	TypeSet
)

type SchemaDefaultFunc func() (any, error)

type Schema struct {
	Type          ValueType
	Optional      bool
	Required      bool
	Computed      bool
	ForceNew      bool
	Default       any
	DefaultFunc   SchemaDefaultFunc
	Description   string
	MaxItems      int
	MinItems      int
	AtLeastOneOf  []string
	ExactlyOneOf  []string
	ConflictsWith []string
	Elem          any
}

type Resource struct {
	Schema map[string]*Schema
}
