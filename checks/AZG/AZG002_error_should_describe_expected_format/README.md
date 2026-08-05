# AZG002

The AZG002 analyzer reports error messages containing `invalid format of ...`. These
messages aren't descriptive - they tell the user something is wrong without telling them
how to fix it. Error messages should describe the expected format instead.

## Flagged Code

```go
return fmt.Errorf("invalid format of %q", name)
```

## Passing Code

```go
return fmt.Errorf("%q must start with a letter, may contain letters and numbers, and must end with a letter", name)
```

## Ignoring Reports

When run via golangci-lint, reports can be ignored with a `//nolint:azproviderlint` Go code
comment at the end of the offending line or on the line immediately preceding it:

```go
return fmt.Errorf("invalid format of %q", name) //nolint:azproviderlint
```
