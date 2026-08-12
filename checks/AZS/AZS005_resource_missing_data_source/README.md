# AZS005

The AZS005 analyzer reports registered resources that have no data source of the same terraform type name registered in the service package. Every registration flavour contributes to both sides of the comparison, so a resource registered one way is covered by a data source registered any other way:

- untyped plugin SDK: `SupportedResources()` / `SupportedDataSources()` map keys
- typed SDK: `Resources()` / `DataSources()` elements, correlated via their `ResourceType()` methods
- framework: `FrameworkResources()` / `FrameworkDataSources()` elements, correlated via their `ResourceType()` methods

Conditionally registered entries (`resources["azurerm_x"] = ...` assignments and `append(out, FooResource{})` calls behind feature flags) are collected too, and generated auto-registration delegation (`append(out, r.autoRegistration.Resources()...)`) is followed — the delegated-to methods' own literal entries are collected when their declarations are visited. Type names declared as package-level vars (`var FooResourceName = "azurerm_foo"`, common for sharing with locks helpers) resolve through their initializer. Invoke-style "action" resources — ones that perform an operation or mint a credential rather than manage a durable object, currently recognised by the name suffixes `_run_command` and `_sas_token` — have no meaningful data source form and are never reported. Registration methods are recognised by name plus return shape (`map[string]*Resource`, `[]Resource`, `[]DataSource`, `[]FrameworkWrappedResource`, `[]FrameworkWrappedDataSource`), and terraform type names are resolved through the type checker, so named constants work. If any data source entry's name cannot be resolved statically the package is skipped entirely, so an unresolvable data source can never produce a false "missing data source" report; unresolvable resource entries are skipped individually.

## Flagged Code

```go
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_example": resourceExample(), // no data source named azurerm_example registered
	}
}
```

## Passing Code

```go
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_example": resourceExample(),
	}
}

func (r Registration) SupportedDataSources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_example": dataSourceExample(),
	}
}
```

## Ignoring Reports

Not every resource needs a data source, so suppressions are expected — the check exists to make the gap a deliberate decision rather than an accident. When run via golangci-lint, reports can be ignored with a `//nolint:azproviderlint` Go code comment at the end of the offending line or on the line immediately preceding it:

```go
"azurerm_example": resourceExample(), //nolint:azproviderlint
```

To ignore only this check on a line — leaving any other azproviderlint checks active — use a `//azignore:AZS005` comment instead, in the same positions:

```go
"azurerm_example": resourceExample(), //azignore:AZS005
```
