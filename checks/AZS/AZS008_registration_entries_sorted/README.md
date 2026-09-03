# AZS008

The AZS008 analyzer reports map and slice entries in `registration.go` files that are not sorted alphabetically (case-insensitively).

`Registration` methods (`SupportedResources`, `SupportedDataSources`, `Resources`, `DataSources`, and friends) return the terraform types a service exposes. Keeping those map keys and slice elements sorted alphabetically keeps registrations easy to scan, keeps diffs small, and avoids merge conflicts when several PRs add entries at once. Map and slice literals matching a method's declared result type are checked whether they are returned directly or assigned to a local variable. Unrelated literals with other types are ignored.

When entries are grouped into sections separated by blank lines or comment lines (typically a comment header above a group), each section is validated independently rather than across the whole literal, so intentionally grouped registrations are not forced into one global ordering.

The report carries a suggested fix, so `azproviderlint -AZS008 -fix` (or an editor applying the suggested fix) reorders the entries automatically. Only safely rewritable unsorted sections are changed; sections with entries sharing a line or crossed by a multiline comment are left alone. Each entry is moved as whole source lines, so its attached comments travel with it. Comment lines act as section boundaries, so section headings stay in place. A comment directly between two entries with no blank lines is treated as attached to the following entry and moves with it.

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

When run via golangci-lint, reports can be ignored with a `//nolint:azproviderlint` Go code comment at the end of the offending line or on the line immediately preceding it:

```go
return map[string]*pluginsdk.Resource{ //nolint:azproviderlint
```

To ignore only this check on a line — leaving any other azproviderlint checks active — use a `//azignore:AZS008 - <reason>` comment instead, in the same positions:

```go
return map[string]*pluginsdk.Resource{ //azignore:AZS008 - <reason>
```
