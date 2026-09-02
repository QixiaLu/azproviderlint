# AZR008 - flatten functions must return empty slices/maps, not nil

The AZR008 analyzer reports `flatten*` functions that return `nil` for a slice or map result instead of an empty container (`[]T{}`, `map[K]V{}`).

Flatten results feed straight into schema state, where a nil container is not interchangeable with an empty one: it can surface as a spurious plan diff or a nil assignment downstream.

The check covers any `flatten*` function (case-insensitive) with slice or map results, named types included, and reports a position that is provably nil:

- a literal `nil`, including conversions like `[]T(nil)`
- a naked `return` whose named container result has not been assigned yet
- a variable that is still nil: a zero-value `var` or named result, unassigned before the return, address never taken

Error paths are skipped — `return nil, err` is legitimate — but only when the error can actually be non-nil, so `var noErr error; return nil, noErr` is still reported. Out of scope: pointers (`*T`, `*[]T`, `*map[K]V` — nil means absent), `interface{}` results (container shape undeclared), and `expand*` functions (nil is idiomatic there).

Reports carry a suggested fix applied via `-fix`: nils become empty literals, naked returns become explicit (`return []T{}, err`), and a returned variable's now-unused declaration is deleted when safe.

## Flagged Code

```go
func flattenNetworkACLs(input *NetworkRuleSet) []NetworkACLs {
	if input == nil {
		return nil
	}
	// ...
}

func flattenTags(input *Resource) map[string]interface{} {
	if input == nil {
		return nil
	}
	// ...
}

// naked return: ret is still nil here
func flattenReplicaSets(input *[]ReplicaSet) (ret []interface{}) {
	if input == nil {
		return
	}
	// ...
}

// out is provably still nil
func flattenRules(input *RuleSet) []Rule {
	var out []Rule
	if input == nil {
		return out
	}
	// ...
}
```

## Passing Code

```go
func flattenNetworkACLs(input *NetworkRuleSet) []NetworkACLs {
	if input == nil {
		return []NetworkACLs{} // or make([]NetworkACLs, 0)
	}
	// ...
}

// error path: nil container is legitimate
func flattenSku(input *Sku) ([]interface{}, error) {
	if input.Name == nil {
		return nil, fmt.Errorf("`name` was nil")
	}
	// ...
}

// nil pointer means absent, not empty
func flattenOptionalTags(input *Resource) *map[string]string {
	if input == nil {
		return nil
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
