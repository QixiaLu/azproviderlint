# AZS002 - schema defaults must match the declared type

The AZS002 analyzer reports `schema.Schema` declarations whose `Default` value's type does not match the declared `Type` — e.g. a `bool` default on a `schema.TypeInt` schema. The plugin SDK's `InternalValidate` does not type-check `Default`, so a mismatch only surfaces as an error at plan time.

Constant values are resolved through the type checker, so named constants (`Default: SkuStandard`) are checked too, and azurerm's `pluginsdk` type aliases are recognised. Non-constant defaults are skipped, `TypeFloat` accepts both int and float constants, and list/set/map schema types (which cannot have literal defaults) are out of scope.

Ports [tfproviderlint PR #329 (S038)](https://github.com/bflad/tfproviderlint/pull/329) with direct constant-kind comparison.

## Flagged Code

```go
&schema.Schema{
	Type:     schema.TypeInt,
	Optional: true,
	Default:  true,
}
```

## Passing Code

```go
&schema.Schema{
	Type:     schema.TypeBool,
	Optional: true,
	Default:  true,
}
```

## Ignoring Reports

When run via golangci-lint, reports can be ignored with a `//nolint:azproviderlint` Go code comment at the end of the offending line or on the line immediately preceding it:

```go
Default: true, //nolint:azproviderlint
```

To ignore only this check on a line — leaving any other azproviderlint checks active — use a `//azignore:AZS002 - <reason>` comment instead, in the same positions:

```go
Default: true, //azignore:AZS002 - <reason>
```
