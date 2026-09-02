# AZS007 - optional+computed fields must have a Note: O+C comment

The AZS007 analyzer reports `schema.Schema` fields that have both `Optional: true` and `Computed: true` (O+C) without a `// Note: O+C because ...` comment between the two field declarations.

Failing to document *why* a field is O+C makes it hard to review and maintain. The required comment forces authors to explain the API behaviour driving the decision.

azurerm's `pluginsdk` type aliases and the typed SDK resource pattern (`Arguments()` / `Attributes()` methods) are both recognised.

## Flagged Code

```go
"max_message_size_in_kilobytes": {
	Type:     schema.TypeInt,
	Optional: true,
	Computed: true,
},
```

```go
"max_message_size_in_kilobytes": {
	Type:     schema.TypeInt,
	Optional: true,
	// this needs to be computed
	Computed: true,
},
```

## Passing Code

```go
"max_message_size_in_kilobytes": {
	Type:     schema.TypeInt,
	Optional: true,
	// NOTE: O+C this gets a variable default based on the sku and can be updated without issues
	Computed: true,
},
```

```go
"administrator_login": {
	Type:     schema.TypeString,
	Optional: true,
	// Note: O+C because Azure returns a generated value if
	// azure_active_directory_administrator.azuread_authentication_only_enabled is true
	Computed: true,
	ForceNew: true,
},
```

The comment must match `// Note: O+C` (case-insensitive) and appear on a line strictly between `Optional:` and `Computed:` in the source. Multi-line O+C comments are supported, only the first line must match the pattern.

## Options

| Option | Default | Effect |
|---|---|---|
| `exclude-packages` | (empty) | comma-separated package names to skip entirely (e.g. state-migration snapshot packages) |

Set via `-AZS007.<option>` on the CLI or a rule-name key in the plugin's golangci settings.

## Ignoring Reports

When run via golangci-lint, reports can be ignored with a `//nolint:azproviderlint` Go code comment at the end of the offending line or on the line immediately preceding it:

```go
Computed: true, //nolint:azproviderlint
```

To ignore only this check — leaving any other azproviderlint checks active — use a `//azignore:AZS007 - <reason>` comment instead:

```go
Computed: true, //azignore:AZS007 - <reason>
```
