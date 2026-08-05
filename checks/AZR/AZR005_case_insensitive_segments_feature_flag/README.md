# AZR005

The AZR005 analyzer reports assignments to the
`features.TreatUserSpecifiedSegmentsAsCaseInsensitive` feature flag.

The case-aware comparisons feature is not ready for use: there is a substantial number of
unresolved dependencies required for it to not cause more problems than it solves. Until
the rollout is completed it must not be configured or exposed in any form.

## Flagged Code

```go
features.TreatUserSpecifiedSegmentsAsCaseInsensitive = true
```

## Passing Code

```go
// remove the assignment entirely - the flag must not be set
```

## Ignoring Reports

When run via golangci-lint, reports can be ignored with a `//nolint:azproviderlint` Go code
comment at the end of the offending line or on the line immediately preceding it:

```go
features.TreatUserSpecifiedSegmentsAsCaseInsensitive = true //nolint:azproviderlint
```
