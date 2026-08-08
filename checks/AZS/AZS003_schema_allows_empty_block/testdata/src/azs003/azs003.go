package azs003

import (
	"azs003/pluginsdk"
	"azs003/schema"
)

func someDefault() (any, error) { return nil, nil }

func sharedProperty() *schema.Schema { return &schema.Schema{Type: schema.TypeString, Optional: true} }

// Should be flagged: list blocks whose properties are all optional with no defaults
func flagged() {
	_ = map[string]*schema.Schema{
		"settings": { // want "schema allows an empty block as every property is optional with no default"
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"foo": {
						Type:     schema.TypeString,
						Optional: true,
					},
					"bar": {
						Type:     schema.TypeInt,
						Optional: true,
					},
				},
			},
		},
	}

	_ = &pluginsdk.Schema{ // want "schema allows an empty block as every property is optional with no default"
		Type:     pluginsdk.TypeList,
		Optional: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"foo": {
					Type:     pluginsdk.TypeString,
					Optional: true,
				},
			},
		},
	}

	// required blocks can be empty too
	_ = &schema.Schema{ // want "schema allows an empty block as every property is optional with no default"
		Type:     schema.TypeList,
		Required: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"foo": {Type: schema.TypeString, Optional: true},
			},
		},
	}
}

// Should NOT be flagged: a constraint, required property or default on any property makes
// empty blocks impossible or harmless
func constrained() {
	_ = &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"foo": {Type: schema.TypeString, Optional: true, AtLeastOneOf: []string{"settings.0.foo", "settings.0.bar"}},
				"bar": {Type: schema.TypeInt, Optional: true, AtLeastOneOf: []string{"settings.0.foo", "settings.0.bar"}},
			},
		},
	}

	_ = &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"foo": {Type: schema.TypeString, Optional: true, ExactlyOneOf: []string{"settings.0.foo"}},
			},
		},
	}

	_ = &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {Type: schema.TypeString, Required: true},
				"foo":  {Type: schema.TypeString, Optional: true},
			},
		},
	}

	_ = &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"foo": {Type: schema.TypeString, Optional: true, Default: "bar"},
			},
		},
	}

	_ = &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"foo": {Type: pluginsdk.TypeString, Optional: true, DefaultFunc: someDefault},
			},
		},
	}
}

// Should NOT be flagged: out of scope shapes
func outOfScope() {
	// computed-only blocks cannot be set in configuration
	_ = &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"foo": {Type: schema.TypeString, Computed: true},
			},
		},
	}

	// sets are out of scope
	_ = &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"foo": {Type: schema.TypeString, Optional: true},
			},
		},
	}

	// primitive element lists have no nested block properties
	_ = &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Elem:     &schema.Schema{Type: schema.TypeString},
	}

	// properties defined elsewhere cannot be inspected
	_ = &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"foo": sharedProperty(),
			},
		},
	}

	// empty nested schema maps have nothing to constrain
	_ = &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Elem:     &schema.Resource{Schema: map[string]*schema.Schema{}},
	}
}
