package azs002

import (
	"azs002/pluginsdk"
	"azs002/schema"
)

type Tier string

const (
	TierFree        Tier = "free"
	defaultCapacity      = 4
	defaultSku           = "Standard"
	defaultEnabled       = true
)

var runtimeDefault = "computed-at-runtime"

func getDefault() any { return nil }

// Should be flagged: literal defaults of the wrong type
func flaggedLiterals() {
	_ = &schema.Schema{
		Type:     schema.TypeInt,
		Optional: true,
		Default:  true, // want "schema Default value type bool does not match the declared Type TypeInt"
	}

	_ = &schema.Schema{
		Type:     schema.TypeBool,
		Optional: true,
		Default:  "true", // want "schema Default value type string does not match the declared Type TypeBool"
	}

	_ = &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
		Default:  30, // want "schema Default value type int does not match the declared Type TypeString"
	}

	_ = &schema.Schema{
		Type:     schema.TypeInt,
		Optional: true,
		Default:  1.5, // want "schema Default value type float does not match the declared Type TypeInt"
	}
}

// Should be flagged: pluginsdk aliases, map-elided literals and named constants all resolve
// through the type checker
func flaggedResolved() {
	_ = &pluginsdk.Schema{
		Type:     pluginsdk.TypeInt,
		Optional: true,
		Default:  "5", // want "schema Default value type string does not match the declared Type TypeInt"
	}

	_ = map[string]*pluginsdk.Schema{
		"sku": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  false, // want "schema Default value type bool does not match the declared Type TypeString"
		},
	}

	_ = &schema.Schema{
		Type:     schema.TypeInt,
		Optional: true,
		Default:  TierFree, // want "schema Default value type string does not match the declared Type TypeInt"
	}
}

// Should NOT be flagged: matching defaults, including named constants and int-for-float
func matching() {
	_ = &schema.Schema{Type: schema.TypeBool, Optional: true, Default: defaultEnabled}
	_ = &schema.Schema{Type: schema.TypeInt, Optional: true, Default: defaultCapacity}
	_ = &schema.Schema{Type: schema.TypeString, Optional: true, Default: defaultSku}
	_ = &schema.Schema{Type: schema.TypeString, Optional: true, Default: string(TierFree)}
	_ = &schema.Schema{Type: schema.TypeFloat, Optional: true, Default: 1.5}
	_ = &schema.Schema{Type: schema.TypeFloat, Optional: true, Default: 1}
	_ = &pluginsdk.Schema{Type: pluginsdk.TypeInt, Optional: true, Default: -1}
}

// Should NOT be flagged: non-constant or nil defaults cannot be judged statically, and
// schema types without literal defaults are out of scope
func outOfScope() {
	_ = &schema.Schema{Type: schema.TypeString, Optional: true, Default: runtimeDefault}
	_ = &schema.Schema{Type: schema.TypeString, Optional: true, Default: getDefault()}
	_ = &schema.Schema{Type: schema.TypeString, Optional: true, Default: nil}
	_ = &schema.Schema{Type: schema.TypeString, Optional: true}
	_ = &schema.Schema{Optional: true, Default: "no declared Type"}
	_ = &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Default:  "lists cannot have literal defaults",
		Elem:     &schema.Schema{Type: schema.TypeString},
	}
}
