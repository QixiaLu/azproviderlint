# AZT001 - acceptance tests must use a _test package

The AZT001 analyzer reports acceptance test files that use the service package instead of an external `_test` package. A file counts as an acceptance test file when its name ends in `_resource_test.go`, `_data_source_test.go`, `_action_test.go` or `_ephemeral_test.go` — including the `_list`, `_identity` and `_gen` (generated) variants such as `_resource_list_test.go` and `_resource_identity_gen_test.go`. The suffix must match exactly: unit tests that merely contain `resource` in their name (`storage_queue_resource_manager_id_test.go`, `parse/resource_group_assignment_test.go`) are not acceptance tests and are ignored.

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

To ignore only this check on a line — leaving any other azproviderlint checks active — use a `//azignore:AZT001 - <reason>` comment instead, in the same positions:

```go
package compute //azignore:AZT001 - <reason>
```
