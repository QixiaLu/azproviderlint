# AZS006

The AZS006 analyzer pairs every registered data source with the same-named registered resource and reports resource schema properties the data source does not expose anywhere. Data sources are meant to mirror their resource, and a property added to the resource but never plumbed into the data source is the usual way they drift apart.

All registration flavours contribute to the pairing (untyped plugin SDK maps, typed SDK slices, framework wrapped slices — correlated via map keys and `ResourceType()` methods, exactly as in AZS005), and schemas are resolved per flavour: the untyped registration value's function body, typed `Arguments()` + `Attributes()` methods, or the framework `Schema()` method. Property names are collected recursively — string-keyed map literals whose value type is a schema type (`*Schema`, framework `Attribute`/`Block`) contribute their constant keys, and calls into same-package schema helper functions are followed.

Matching is by property name across the whole schema (top-level and nested): a resource property is "covered" if its name appears anywhere in the data source's schema. This deliberately under-reports when schemas are restructured (a nested name matching an unrelated data source property masks the gap) but never false-positives from restructuring. A missing block is reported once by dotted path with its entire subtree suppressed — once `backup` is called out there is no value in also listing `backup.retention_days`. Pairs where either side yields no properties, or where the data source contains a non-constant schema key, are skipped. Write-only arguments (`foo_wo`, `foo_wo_version`) are never readable and are exempt.

## Flagged Code

```go
func resourceExample() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Schema: map[string]*pluginsdk.Schema{
			"name": {Type: pluginsdk.TypeString, Required: true},
			"zone": {Type: pluginsdk.TypeString, Optional: true},
		},
	}
}

func dataSourceExample() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Schema: map[string]*pluginsdk.Schema{
			"name": {Type: pluginsdk.TypeString, Required: true},
			// "zone" is not exposed
		},
	}
}
```

## Passing Code

```go
func dataSourceExample() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Schema: map[string]*pluginsdk.Schema{
			"name": {Type: pluginsdk.TypeString, Required: true},
			"zone": {Type: pluginsdk.TypeString, Computed: true},
		},
	}
}
```

## Ignoring Reports

Some resource properties have no meaningful data source form (write-only credentials, create-only inputs), so suppressions are expected. When run via golangci-lint, reports can be ignored with a `//nolint:azproviderlint` Go code comment at the end of the offending line or on the line immediately preceding it:

```go
"azurerm_example": dataSourceExample(), //nolint:azproviderlint
```

To ignore only this check on a line — leaving any other azproviderlint checks active — use a `//azignore:AZS006` comment instead, in the same positions:

```go
"azurerm_example": dataSourceExample(), //azignore:AZS006
```
