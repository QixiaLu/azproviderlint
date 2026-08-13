# AZS004

The AZS004 analyzer reports `validation.StringInSlice([]string{...}, ...)` calls whose hand-written list references an SDK enum's constants. A partial list silently rejects values the API accepts, and even a complete list goes stale the moment the SDK adds a new value — when the SDK ships a possible-values helper, validation should use it. Incomplete lists are reported with the missing values named; complete lists are reported with a suggestion to switch to the helper.

A type only counts as a closed enum when its package exports a possible-values helper (`PossibleValuesFor<Enum>()` in go-azure-sdk, `Possible<Enum>Values()` in older SDKs), so ordinary named string types with a few convenience constants are not reported. The `StringInSlice` callee is resolved through the type checker and matched by name, signature and a package path of/ending in `validation`, so both the plugin SDK's `helper/validation` and provider-internal wrappers of it (e.g. azurerm's `internal/tf/validation`) are recognised. Raw string literals in the list count towards coverage; lists containing non-constant elements or constants of more than one enum type are skipped, since a computed list or a deliberate union cannot be proven incomplete statically.

## Flagged Code

```go
validation.StringInSlice([]string{
	string(virtualmachines.VirtualMachinePriorityTypesLow),
	string(virtualmachines.VirtualMachinePriorityTypesRegular),
	// missing VirtualMachinePriorityTypesSpot ("Spot")
}, false)
```

## Passing Code

```go
validation.StringInSlice(virtualmachines.PossibleValuesForVirtualMachinePriorityTypes(), false)
```

## Ignoring Reports

A deliberately unsupported subset of an enum is a legitimate reason to suppress this check. When run via golangci-lint, reports can be ignored with a `//nolint:azproviderlint` Go code comment at the end of the offending line or on the line immediately preceding it:

```go
ValidateFunc: validation.StringInSlice([]string{ //nolint:azproviderlint
```

To ignore only this check on a line — leaving any other azproviderlint checks active — use a `//azignore:AZS004` comment instead, in the same positions:

```go
ValidateFunc: validation.StringInSlice([]string{ //azignore:AZS004
```
