# AZG001

The AZG001 analyzer reports `_, err := SomeFunc()` assignments immediately followed by
`if err != nil`, which should be combined into a single `if` init statement to keep the
error scoped to its check and the happy path unindented.

## Flagged Code

```go
_, err := client.Delete(ctx, id)
if err != nil {
	return fmt.Errorf("deleting %s: %+v", id, err)
}
```

## Passing Code

```go
if _, err := client.Delete(ctx, id); err != nil {
	return fmt.Errorf("deleting %s: %+v", id, err)
}
```

## Ignoring Reports

When run via golangci-lint, reports can be ignored with a `//nolint:azproviderlint` Go code
comment at the end of the offending line or on the line immediately preceding it:

```go
_, err := client.Delete(ctx, id) //nolint:azproviderlint
if err != nil {
	return err
}
```
