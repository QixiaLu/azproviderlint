# AZG007

The AZG007 analyzer reports pointer fields explicitly initialised to `nil` in a struct literal — `Selector: nil` — where the field should simply be omitted. An omitted pointer field already takes its zero value (`nil`), so the assignment adds noise without changing behaviour.

Only pointer fields are flagged. Slices, maps, and interfaces are deliberately left alone: an explicit `nil` there can be a readable, intentional signal, and omitting it is not always equivalent in intent. Test files are skipped as well, since a `nil` entry in a test table is often semantically meaningful.

The report carries a suggested fix, so `azproviderlint -AZG007 -fix` (or an editor applying the suggested fix) removes the redundant field — along with its trailing comma and any trailing comment — automatically, leaving gofmt to tidy up the surrounding whitespace.

## Flagged Code

```go
return &profiles.ProfileLogScrubbing{
	State:    &policyDisabled,
	Selector: nil, // Selector is *string — redundant
}
```

## Passing Code

```go
return &profiles.ProfileLogScrubbing{
	State: &policyDisabled,
}

// left alone: slices, maps, and interfaces
return &Config{
	Items: nil,
	Data:  nil,
}
```

## Ignoring Reports

When run via golangci-lint, reports can be ignored with a `//nolint:azproviderlint` Go code comment at the end of the offending line or on the line immediately preceding it:

```go
Selector: nil, //nolint:azproviderlint
```

To ignore only this check on a line — leaving any other azproviderlint checks active — use a `//azignore:AZG007 - <reason>` comment instead, in the same positions:

```go
Selector: nil, //azignore:AZG007 - <reason>
```
