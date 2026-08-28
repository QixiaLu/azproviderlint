# AZS004

The AZS004 analyzer reports `validation.StringInSlice([]string{...}, ...)` calls whose hand-written list references an SDK enum's constants. A partial list silently rejects values the API accepts, and even a complete list goes stale the moment the SDK adds a new value — when the SDK ships a possible-values helper, validation should use it. Incomplete lists are reported with the missing values named; complete lists are reported with a suggestion to switch to the helper; lists carrying values that are not part of the enum at all (typos, or deliberate legacy extras) are reported with the extra values named and advice to append any deliberate extras to the helper's result rather than swapping them away.

A type only counts as a closed enum when its package exports a possible-values helper (`PossibleValuesFor<Enum>()` in go-azure-sdk, `Possible<Enum>Values()` in older SDKs) returning either `[]string` or a slice of the enum type itself, so ordinary named string types with a few convenience constants are not reported. For the typed-slice track-1 form the helper cannot be passed to `StringInSlice` directly, so the advice routes through go-azure-helpers' generic conversion instead: `pointer.FromEnumSlice(pointer.To(cdn.PossibleTransformValues()))`. The `StringInSlice` callee is resolved through the type checker and matched by name, signature and a package path of/ending in `validation`, so both the plugin SDK's `helper/validation` and provider-internal wrappers of it (e.g. azurerm's `internal/tf/validation`) are recognised. Raw string literals in the list count towards coverage; lists containing non-constant elements or constants of more than one enum type are skipped, since a computed list or a deliberate union cannot be proven incomplete statically.

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

## Flags

Two flags suppress one reporting class each, for providers where subsets or supersets are policy rather than drift. A list that is exactly the enum is always reported, since switching to the helper there is a pure win.

- `-AZS004.allow-missing-values`: do not report in-place validation arrays that are missing enum values (deliberate subsets)
- `-AZS004.allow-extra-values`: do not report in-place validation arrays containing values that are not part of the enum (deliberate supersets)

Via the golangci-lint plugin the flags are set through a rule-name key in the settings:

```yaml
linters:
  settings:
    custom:
      azproviderlint:
        settings:
          AZS004:
            allow-missing-values: true
            allow-extra-values: true
```

## Ignoring Reports

A deliberately unsupported subset of an enum is a legitimate reason to suppress this check. When run via golangci-lint, reports can be ignored with a `//nolint:azproviderlint` Go code comment at the end of the offending line or on the line immediately preceding it:

```go
ValidateFunc: validation.StringInSlice([]string{ //nolint:azproviderlint
```

To ignore only this check on a line — leaving any other azproviderlint checks active — use a `//azignore:AZS004 - <reason>` comment instead, in the same positions:

```go
ValidateFunc: validation.StringInSlice([]string{ //azignore:AZS004 - <reason>
```
