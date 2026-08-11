# AZG005

The AZG005 analyzer reports `fmt.Errorf("fixed message")` calls whose only argument is a string literal containing no `%` format placeholders. Such calls don't format anything, so the simpler and cheaper `errors.New("fixed message")` should be used instead.

The check only fires on a single string-literal argument with no `%`. Calls with format placeholders (`%s`, `%d`, `%v`, `%w`, …) or additional arguments are left alone.

## Flagged Code

```go
fmt.Errorf("something went wrong")
fmt.Errorf("invalid input")
```

## Passing Code

```go
errors.New("something went wrong")

// has a placeholder, so fmt.Errorf is appropriate
fmt.Errorf("value %s is invalid", value)
fmt.Errorf("wrapped: %w", err)
```

## Ignoring Reports

When run via golangci-lint, reports can be ignored with a `//nolint:azproviderlint` Go code comment at the end of the offending line or on the line immediately preceding it:

```go
return fmt.Errorf("something went wrong") //nolint:azproviderlint
```

To ignore only this check on a line — leaving any other azproviderlint checks active — use a `//azignore:AZG005` comment instead, in the same positions:

```go
return fmt.Errorf("something went wrong") //azignore:AZG005
```
