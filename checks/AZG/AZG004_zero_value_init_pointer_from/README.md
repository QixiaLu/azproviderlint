# AZG004

The AZG004 analyzer reports the manual `y := <zero>; if x != nil { y = *x }` idiom — a zero-value initialization immediately followed by a nil check that dereferences a pointer — where the generic `pointer.From(x)` helper from [go-azure-helpers](https://github.com/hashicorp/go-azure-helpers) should be used instead.

`pointer.From` returns the dereferenced value, or the type's zero value when the pointer is nil, so the whole init-and-nil-check dance collapses to a single expression. The check only fires when the variable is initialized to a zero value (`false`, `0`, `""`, `nil`), the `if` has no `else` branch and a single `x != nil` condition, and its body is exactly one `y = *x` assignment whose dereferenced expression matches the nil-checked one. Both the short-declaration form (`y := <zero>`) and the var-declaration form (`var y T`, `var y T = <zero>`, `var y = <zero>`) are matched.

## Flagged Code

```go
enabled := false
if props.Enabled != nil {
	enabled = *props.Enabled
}

// the var form is matched too
var enabled bool
if props.Enabled != nil {
	enabled = *props.Enabled
}
```

## Passing Code

```go
enabled := pointer.From(props.Enabled)

// left alone: an else branch, a non-zero initial value, a mismatched or
// re-declared variable, or non-adjacent statements
enabled := true
if props.Enabled != nil {
	enabled = *props.Enabled
}
```

## Ignoring Reports

When run via golangci-lint, reports can be ignored with a `//nolint:azproviderlint` Go code comment at the end of the offending line or on the line immediately preceding it:

```go
enabled := false //nolint:azproviderlint
```

To ignore only this check on a line — leaving any other azproviderlint checks active — use a `//azignore:AZG004` comment instead, in the same positions:

```go
enabled := false //azignore:AZG004
```
