# AZG003

The AZG003 analyzer reports `pointer.To` calls that wrap an explicit [go-azure-sdk](https://github.com/hashicorp/go-azure-sdk) enum type conversion — `pointer.To(sdk.SomeEnum(v))` — where the generic `pointer.ToEnum[sdk.SomeEnum](v)` helper from [go-azure-helpers](https://github.com/hashicorp/go-azure-helpers) should be used instead.

`pointer.ToEnum` makes the intent explicit and keeps the enum type in one place, avoiding the redundant `sdk.SomeEnum(...)` conversion. The check only fires when the converted type is a go-azure-sdk enum: a named string type declared in a `github.com/hashicorp/go-azure-sdk` package that exposes the generated `PossibleValuesFor<Name>() []T` helper. Plain `pointer.To` calls on strings, ints, or non-SDK types are left alone.

`pointer.ToEnum` has the signature `func ToEnum[T ~string](input string) *T`, so the rewrite `pointer.ToEnum[sdk.SomeEnum](v)` only compiles when the converted value `v` is assignable to `string`. When `v` is itself a named string type (for example another enum), the report notes that the value must be wrapped in an explicit `string(...)` conversion — e.g. `pointer.ToEnum[sdk.SomeEnum](string(v))` — for the suggestion to compile.

The report carries a suggested fix, so `azproviderlint -AZG003 -fix` (or an editor applying the suggested fix) rewrites the call automatically — including the `string(...)` wrap when the converted value is not assignable to string.

## Flagged Code

```go
return pointer.To(virtualmachines.VirtualMachinePriorityTypes(priority))
return pointer.To(managedclusters.ArtifactSource(config["artifact_source"].(string)))
```

## Passing Code

```go
return pointer.ToEnum[virtualmachines.VirtualMachinePriorityTypes](priority)
return pointer.ToEnum[managedclusters.ArtifactSource](config["artifact_source"].(string))

// non-enum values are unaffected
return pointer.To("regular string")
return pointer.To(42)
```

## Ignoring Reports

When run via golangci-lint, reports can be ignored with a `//nolint:azproviderlint` Go code comment at the end of the offending line or on the line immediately preceding it:

```go
return pointer.To(virtualmachines.OperatingSystemTypes("Linux")) //nolint:azproviderlint
```

To ignore only this check on a line — leaving any other azproviderlint checks active — use a `//azignore:AZG003 - <reason>` comment instead, in the same positions:

```go
return pointer.To(virtualmachines.OperatingSystemTypes("Linux")) //azignore:AZG003 - <reason>
```
