# AZS001

The AZS001 analyzer reports typed SDK model fields (struct fields tagged `tfschema`) using non-64-bit numeric types: `int`, `int16` or `int32` instead of `int64`, and `float32` instead of `float64`.

The typed SDK's `Encode`/`Decode` work with `int64` and `float64`; models using other widths fail at runtime. The check covers slices, maps and pointers of these types, and resolves named types and aliases through the type checker (matching what `reflect.Kind()` sees at runtime), so `type Capacity int` is also flagged.

## Flagged Code

```go
type ServerModel struct {
	Capacity     int       `tfschema:"capacity"`
	Priority     int32     `tfschema:"priority"`
	CpuThreshold float32   `tfschema:"cpu_threshold"`
	AllowedPorts []int     `tfschema:"allowed_ports"`
}
```

## Passing Code

```go
type ServerModel struct {
	Capacity     int64     `tfschema:"capacity"`
	Priority     int64     `tfschema:"priority"`
	CpuThreshold float64   `tfschema:"cpu_threshold"`
	AllowedPorts []int64   `tfschema:"allowed_ports"`
}
```

## Ignoring Reports

When run via golangci-lint, reports can be ignored with a `//nolint:azproviderlint` Go code comment at the end of the offending line or on the line immediately preceding it:

```go
Capacity int `tfschema:"capacity"` //nolint:azproviderlint
```

To ignore only this check on a line — leaving any other azproviderlint checks active — use a `//azignore:AZS001` comment instead, in the same positions:

```go
Capacity int `tfschema:"capacity"` //azignore:AZS001
```
