# AZT001

The AZT001 analyzer reports resource and data source acceptance test files (`*resource*_test.go`, `*data_source*_test.go`) that use the service package instead of an external `_test` package.

Acceptance tests must live in a `_test` package to prevent a circular dependency between the service package and the acceptance test framework.

## Flagged Code

```go
// in example_resource_test.go
package compute
```

## Passing Code

```go
// in example_resource_test.go
package compute_test
```

## Ignoring Reports

When run via golangci-lint, reports can be ignored with a `//nolint:azproviderlint` Go code comment at the end of the offending line or on the line immediately preceding it:

```go
package compute //nolint:azproviderlint
```

To ignore only this check on a line — leaving any other azproviderlint checks active — use a `//azignore:AZT001` comment instead, in the same positions:

```go
package compute //azignore:AZT001
```
