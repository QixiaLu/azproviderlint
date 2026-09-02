package azs007

import (
	"azs007/pluginsdk"
	"azs007/schema"
	"azs007/sdk"
)

func flaggedLiterals() {
	// Missing comment entirely
	_ = &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
		Computed: true, // want `schema field has both Optional and Computed but is missing a '// Note: O\+C because \.\.\.' comment between the two fields`
	}

	// Has a comment, but it does not match the required pattern
	_ = &schema.Schema{
		Type:     schema.TypeInt,
		Optional: true,
		// this needs to be computed for reasons
		Computed: true, // want `schema field has both Optional and Computed but is missing a '// Note: O\+C because \.\.\.' comment between the two fields`
	}

	// Comment exists but is on the wrong side (after Computed, not between Optional and Computed)
	_ = &schema.Schema{
		Type:     schema.TypeBool,
		Optional: true,
		Computed: true, // want `schema field has both Optional and Computed but is missing a '// Note: O\+C because \.\.\.' comment between the two fields`
		// Note: O+C because this is in the wrong place
	}

	// pluginsdk alias - same underlying type, must also be flagged
	_ = &pluginsdk.Schema{
		Type:     pluginsdk.TypeString,
		Optional: true,
		Computed: true, // want `schema field has both Optional and Computed but is missing a '// Note: O\+C because \.\.\.' comment between the two fields`
	}

	// map literal (no explicit &pluginsdk.Schema{} type, inferred from map value type)
	_ = map[string]*pluginsdk.Schema{
		"size": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
			Computed: true, // want `schema field has both Optional and Computed but is missing a '// Note: O\+C because \.\.\.' comment between the two fields`
		},
	}
}

func flaggedTyped() {
	// Same checks apply inside a typed resource's Arguments() map
	_ = map[string]*schema.Schema{
		"authentication_methods": {
			Type:     schema.TypeSet,
			Optional: true,
			Computed: true, // want `schema field has both Optional and Computed but is missing a '// Note: O\+C because \.\.\.' comment between the two fields`
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
	}
}

// ---- Passing: Optional+Computed WITH a valid O+C comment between them ----

func passing() {
	// Standard "// NOTE: O+C" comment between Optional and Computed
	_ = &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
		// NOTE: O+C the API sets a default value if omitted
		Computed: true,
	}

	// Lowercase "// Note: O+C because" variant
	_ = &schema.Schema{
		Type:     schema.TypeInt,
		Optional: true,
		// Note: O+C because Azure returns a calculated value when omitted
		Computed: true,
	}

	// Comment with separator style used in some resources
	_ = &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
		// NOTE: O+C - Azure returns a generated value if azuread_authentication_only_enabled is true
		Computed: true,
	}

	// pluginsdk alias with correct comment
	_ = &pluginsdk.Schema{
		Type:     pluginsdk.TypeInt,
		Optional: true,
		// NOTE: O+C this gets a variable default based on the sku
		Computed:     true,
		ValidateFunc: func(v interface{}, k string) ([]string, []error) { return nil, nil },
	}

	// map literal with correct comment
	_ = map[string]*pluginsdk.Schema{
		"max_message_size_in_kilobytes": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
			// NOTE: O+C this gets a variable default based on the sku and can be updated without issues
			Computed: true,
		},
	}

	// Multi-line comment block between Optional and Computed - only first line needs to match
	_ = &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
		// Note: O+C because Azure returns a generated value if
		// azure_active_directory_administrator.azuread_authentication_only_enabled is true
		Computed: true,
	}

	// Other fields between Optional and Computed with a O+C comment
	_ = &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
		// NOTE: O+C creating a hub that has encryption enabled will set this property to true
		Computed:  true,
		Sensitive: true,
	}
}

// ---- Passing/Flagged: Computed declared before Optional (C+O order) ----

func computedBeforeOptional() {
	// Flagged: C+O order, no comment between the two fields
	_ = &schema.Schema{
		Type:     schema.TypeString,
		Computed: true, // want `schema field has both Optional and Computed but is missing a '// Note: O\+C because \.\.\.' comment between the two fields`
		Optional: true,
	}

	// Flagged: C+O order with unrelated comment between the fields
	_ = &schema.Schema{
		Type:     schema.TypeInt,
		Computed: true, // want `schema field has both Optional and Computed but is missing a '// Note: O\+C because \.\.\.' comment between the two fields`
		// this needs to be computed for reasons
		Optional: true,
	}

	// Passing: C+O order with valid O+C comment between Computed and Optional
	_ = &schema.Schema{
		Type:     schema.TypeString,
		Computed: true,
		// NOTE: O+C the API sets a default value if omitted
		Optional: true,
	}

	// Passing: C+O order, pluginsdk alias, valid comment
	_ = &pluginsdk.Schema{
		Type:     pluginsdk.TypeInt,
		Computed: true,
		// Note: O+C because Azure returns a calculated value when omitted
		Optional: true,
	}

	// Passing: C+O order, map literal, valid comment
	_ = map[string]*schema.Schema{
		"dns_label": {
			Type:     schema.TypeString,
			Computed: true,
			// NOTE: O+C - Azure assigns a value when omitted
			Optional: true,
		},
	}
}

// ---- Out of scope: neither Optional+Computed, so no check needed ----

func outOfScope() {
	// Only Optional - no check needed
	_ = &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
	}

	// Only Computed - no check needed
	_ = &schema.Schema{
		Type:     schema.TypeString,
		Computed: true,
	}

	// Required - no check needed
	_ = &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
	}

	// No Optional or Computed at all
	_ = &schema.Schema{
		Type:    schema.TypeString,
		Default: "default",
	}
}

// ---- Typed resource implementation (Arguments / Attributes pattern) ----

type ExampleTypedResource struct{}

var _ sdk.Resource = ExampleTypedResource{}

func (r ExampleTypedResource) ResourceType() string     { return "azurerm_example" }
func (r ExampleTypedResource) ModelObject() interface{} { return nil }
func (r ExampleTypedResource) IDValidationFunc() func(interface{}, string) ([]string, []error) {
	return nil
}
func (r ExampleTypedResource) Create() sdk.ResourceFunc { return nil }
func (r ExampleTypedResource) Read() sdk.ResourceFunc   { return nil }
func (r ExampleTypedResource) Delete() sdk.ResourceFunc { return nil }

// Arguments contains user-settable fields including O+C fields.
func (r ExampleTypedResource) Arguments() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Required: true,
		},

		// Flagged: O+C field without comment in a typed resource's Arguments()
		"storage_iops": {
			Type:     schema.TypeInt,
			Optional: true,
			Computed: true, // want `schema field has both Optional and Computed but is missing a '// Note: O\+C because \.\.\.' comment between the two fields`
		},

		// Passing: O+C field with correct comment in Arguments()
		"administrator_login": {
			Type:     schema.TypeString,
			Optional: true,
			// Note: O+C because Azure returns a generated value if azuread_authentication_only_enabled is true
			Computed: true,
			ForceNew: true,
		},
	}
}

// Attributes contains Computed-only (read-only) fields - these are never O+C.
func (r ExampleTypedResource) Attributes() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"fqdn": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}
