# AZR008 - flatten functions must return empty slices/maps, not nil

The AZR008 analyzer reports `flatten*` functions that return `nil` for a slice or map result, where an empty container (`[]T{}`, `map[K]V{}`) should be returned instead.

Flatten helpers feed their result straight into schema state, and a nil slice is not interchangeable with an empty one: it can surface as a spurious plan diff or trigger a nil-map/slice assignment downstream. The nil-input guard that opens most flatten functions should therefore return an empty slice rather than `nil`.

The check fires on any function whose name begins with `flatten` (case-insensitive) that declares one or more slice or map result types (named types included) and returns a provably nil value in one of those positions: a literal `nil` (including conversions such as `[]T(nil)`), a naked `return` whose named container result has not been assigned yet, or a variable that is still nil — a zero-value `var` declaration or named result with no assignment before the return and its address never taken. Non-container results (strings, pointers, `error`) are ignored: a nil pointer (`*T`, `*[]T`, `*map[K]V`) is a deliberate absent signal, `interface{}` results are out of scope since the container shape is not declared, and `expand*` functions are out of scope since returning `nil` there is idiomatic.

Returns on an error path are skipped: a `return nil, err` (or `return nil, fmt.Errorf(...)`) that may carry a real error legitimately returns `nil` for the container, so only returns whose error positions are provably nil — the empty/nil-input branch, including `var noErr error; return nil, noErr` — are reported.

The report carries a suggested fix, so `azproviderlint -AZR008 -fix` (or an editor applying the suggested fix) rewrites each offending value into an empty composite literal of the declared type (`[]T{}`, `map[K]V{}`) automatically — naked returns become explicit (`return []T{}, err`), and a returned nil variable's now-unused declaration is deleted when that is safe (the fix is withheld when it is not).

## Flagged Code

```go
func flattenNetworkACLs(input *NetworkRuleSet) []NetworkACLs {
	if input == nil {
		return nil
	}
	// ...
}
```

## Passing Code

```go
func flattenNetworkACLs(input *NetworkRuleSet) []NetworkACLs {
	if input == nil {
		return []NetworkACLs{}
	}
	// ...
}

// make is fine too
func flattenNetworkACLs(input *NetworkRuleSet) []NetworkACLs {
	if input == nil {
		return make([]NetworkACLs, 0)
	}
	// ...
}
```

## Ignoring Reports

When run via golangci-lint, reports can be ignored with a `//nolint:azproviderlint` Go code comment at the end of the offending line or on the line immediately preceding it:

```go
return nil //nolint:azproviderlint
```

To ignore only this check on a line — leaving any other azproviderlint checks active — use a `//azignore:AZR008 - <reason>` comment instead, in the same positions:

```go
return nil //azignore:AZR008 - <reason>
```
