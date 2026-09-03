# AZG007

The AZG007 analyzer reports struct literal fields explicitly initialised to their zero value — a pointer set to `nil` (`Selector: nil`), a string set to `""`, a numeric set to `0`, or a bool set to `false` — where the field should simply be omitted. An omitted field already takes its zero value, so the assignment adds noise without changing behaviour.

Slices, maps, and interfaces are deliberately left alone: an explicit `nil` there can be a readable, intentional signal, and omitting it is not always equivalent in intent. Zero values written as a named constant (`Mode: ModeNone`) are also left alone, since the name documents intent and omitting the field would lose that.

Only literal zero values are flagged, so non-idiomatic forms such as `-0.0` or `'\x00'` are not caught; this can be revisited if a real case comes up.

The report carries a suggested fix, so `azproviderlint -AZG007 -fix` (or an editor applying the suggested fix) removes the redundant field — along with its trailing comma and any trailing comment — automatically, leaving gofmt to tidy up the surrounding whitespace.

## Flagged Code

```go
return &profiles.ProfileLogScrubbing{
	State:    &policyDisabled,
	Selector: nil, // Selector is *string — redundant
}

// string, numeric, and bool zero values are redundant too
return KubeConfigModel{
	Host:              cluster.Server,
	Username:          name,
	Password:          "", // redundant
	ClientCertificate: "", // redundant
	ClientKey:         "", // redundant
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

// left alone: a named constant documents intent even when it is zero
return &Settings{
	Mode: ModeNone,
}
```

## Flags

`-AZG007.tests` also reports inside `_test.go` files, which are skipped by default because a zero entry in a test table is often a deliberate, meaningful row. Via the golangci-lint plugin the flag is set through settings:

```yaml
linters:
  settings:
    custom:
      azproviderlint:
        settings:
          AZG007:
            tests: true
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
