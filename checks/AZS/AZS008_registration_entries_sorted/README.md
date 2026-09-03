# AZS008

The AZS008 analyzer reports `Registration` methods containing map or slice entries that are not sorted alphabetically (case-insensitively). Each method receives at most one report, even when multiple sections are unsorted.

`Registration` methods (`SupportedResources`, `SupportedDataSources`, `Resources`, `DataSources`, and friends) return the terraform types a service exposes. Keeping those map keys and slice elements sorted alphabetically keeps registrations easy to scan, keeps diffs small, and avoids merge conflicts when several PRs add entries at once. Map and slice literals matching a method's declared result type are checked whether they are returned directly or assigned to a local variable. Unrelated literals with other types are ignored.

Both `registration.go` and the generated `registration_gen.go` (with its `autoRegistration` receiver) are in scope by default: an unsorted generated file means the generator's input or template needs fixing — apply the fix there, not to the generated output. Set `generated` to false to skip generated files.

## Options

| Option | Default | Effect |
|---|---|---|
| `generated` | true | check generated `registration_gen.go` files |

Set via `-AZS008.<option>` on the CLI or a rule-name key in the plugin's golangci settings.

When entries are grouped into sections separated by blank lines or headings, each section is validated independently rather than across the whole literal, so intentionally grouped registrations are not forced into one global ordering. A heading comment starts a section when it has a blank line before it. Other comments are treated as attached to the following entry.

The report carries a suggested fix, so `azproviderlint -AZS008 -fix` (or an editor applying the suggested fix) reorders the entries automatically. Only safely rewritable unsorted sections are changed; sections with entries sharing a line or crossed by a multiline comment are left alone. Each entry is moved as whole source lines, so its attached comments travel with it and section headings stay in place.

## Flagged Code

```go
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_managed_disk":     nil,
		"azurerm_availability_set": nil, // should come first alphabetically
	}
}

func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		WorkspaceResource{},
		ApiManagementResource{}, // should come first alphabetically
	}
}
```

## Passing Code

```go
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_availability_set": nil,
		"azurerm_managed_disk":     nil,
	}
}

func (r Registration) SectionedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		// CDN
		"azurerm_cdn_profile": nil,

		// FrontDoor
		"azurerm_cdn_frontdoor_custom_domain": nil,
		"azurerm_cdn_frontdoor_profile":       nil,
	}
}

func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		ApiManagementResource{},
		WorkspaceResource{},
	}
}
```

## Ignoring Reports

AZS008 reports on the registration method. When run via golangci-lint, the report can be ignored with a `//nolint:azproviderlint` Go code comment on the method declaration or on the line immediately preceding it:

```go
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource { //nolint:azproviderlint
	return map[string]*pluginsdk.Resource{
		"azurerm_managed_disk":     nil,
		"azurerm_availability_set": nil,
	}
}
```

To ignore only this check on the method — leaving any other azproviderlint checks active — use a `//azignore:AZS008 - <reason>` comment instead, in the same positions:

```go

func (r Registration) SupportedResources() map[string]*pluginsdk.Resource { //azignore:AZS008 - intentional ordering
	return map[string]*pluginsdk.Resource{
		"azurerm_managed_disk":     nil,
		"azurerm_availability_set": nil,
	}
}
```
