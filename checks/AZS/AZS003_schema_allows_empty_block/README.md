# AZS003

The AZS003 analyzer reports `Type: schema.TypeList` blocks (Optional or Required) whose nested `Resource` schema contains only optional properties with no `Default`/`DefaultFunc` and no `AtLeastOneOf`/`ExactlyOneOf` constraint. Such schemas accept an empty block:

```hcl
resource "azurerm_example" "example" {
  settings {}
}
```

which yields a `nil` list element that expand functions commonly crash on (`raw[0].(map[string]interface{})` — see [azurerm #11426](https://github.com/hashicorp/terraform-provider-azurerm/issues/11426)), or a permanent diff. A single `Required` property, `Default`, or `AtLeastOneOf`/`ExactlyOneOf` constraint on any property makes the block safe. azurerm's `pluginsdk` type aliases are recognised.

Note this is a heuristic: blocks whose expand functions guard against `nil` elements are safe in practice but still reported — suppress those with `//azignore:AZS003`.

Ports [tfproviderlint PR #236 (XS003)](https://github.com/bflad/tfproviderlint/pull/236).

## Flagged Code

```go
"settings": {
	Type:     schema.TypeList,
	Optional: true,
	MaxItems: 1,
	Elem: &schema.Resource{
		Schema: map[string]*schema.Schema{
			"foo": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	},
},
```

## Passing Code

```go
"settings": {
	Type:     schema.TypeList,
	Optional: true,
	MaxItems: 1,
	Elem: &schema.Resource{
		Schema: map[string]*schema.Schema{
			"foo": {
				Type:         schema.TypeString,
				Optional:     true,
				AtLeastOneOf: []string{"settings.0.foo"},
			},
		},
	},
},
```

## Ignoring Reports

When run via golangci-lint, reports can be ignored with a `//nolint:azproviderlint` Go code comment at the end of the offending line or on the line immediately preceding it:

```go
"settings": { //nolint:azproviderlint
```

To ignore only this check on a line — leaving any other azproviderlint checks active — use a `//azignore:AZS003` comment instead, in the same positions:

```go
"settings": { //azignore:AZS003
```
