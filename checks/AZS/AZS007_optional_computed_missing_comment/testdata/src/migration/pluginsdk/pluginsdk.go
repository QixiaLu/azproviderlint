// Package pluginsdk mirrors terraform-provider-azurerm's internal/tf/pluginsdk type aliases.
package pluginsdk

import "migration/schema"

type (
	Schema   = schema.Schema
	Resource = schema.Resource
)

const (
	TypeBool   = schema.TypeBool
	TypeInt    = schema.TypeInt
	TypeFloat  = schema.TypeFloat
	TypeString = schema.TypeString
	TypeList   = schema.TypeList
	TypeMap    = schema.TypeMap
	TypeSet    = schema.TypeSet
)
