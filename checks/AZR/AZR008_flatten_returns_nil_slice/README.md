# AZR008 - flatten functions must return empty slices, not nil

The AZR008 analyzer reports `flatten*` functions that return `nil` for a slice result, where an empty slice (`[]T{}` or `make([]T, 0)`) should be returned instead.

Flatten helpers feed their result straight into schema state, and a nil slice is not interchangeable with an empty one: it can surface as a spurious plan diff or trigger a nil-map/slice assignment downstream. The nil-input guard that opens most flatten functions should therefore return an empty slice rather than `nil`.

The check fires on any function whose name begins with `flatten` (case-insensitive) that declares one or more slice result types and returns `nil` in one of those slice positions — including the slice position of a multi-value return such as `([]T, error)`. Non-slice results (strings, pointers, `error`) are ignored, and `expand*` functions are out of scope since returning `nil` there is idiomatic.

Returns on an error path are skipped: a `return nil, err` (or `return nil, nil, err`, `return nil, fmt.Errorf(...)`) that carries a non-nil error value legitimately returns `nil` for the slice, so only returns whose error results are all `nil` (or absent) — the empty/nil-input branch — are reported.

The report carries a suggested fix, so `azproviderlint -AZR008 -fix` (or an editor applying the suggested fix) rewrites each offending `nil` into an empty composite literal of the declared slice type (`[]T{}`) automatically — every nil slice position in a single return is fixed at once.

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
