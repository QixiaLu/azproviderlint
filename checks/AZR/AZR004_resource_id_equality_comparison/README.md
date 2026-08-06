# AZR004

The AZR004 analyzer reports Resource IDs being compared with `==` or `!=` (`a.ID() == b.ID()`).

Azure Resource IDs contain user-specified segments that are compared case-insensitively by the API, so string equality is unreliable. Use `resourceids.Match(a, b)` from `github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids` instead (and `!resourceids.Match(a, b)` in place of `!=`).

## Flagged Code

```go
if subnetId.ID() == other.ID() {
	// ...
}
```

## Passing Code

```go
if resourceids.Match(subnetId, other) {
	// ...
}
```

## Ignoring Reports

When run via golangci-lint, reports can be ignored with a `//nolint:azproviderlint` Go code comment at the end of the offending line or on the line immediately preceding it:

```go
if subnetId.ID() == other.ID() { //nolint:azproviderlint
```

To ignore only this check on a line — leaving any other azproviderlint checks active — use a `//azignore:AZR004` comment instead, in the same positions:

```go
if subnetId.ID() == other.ID() { //azignore:AZR004
```
