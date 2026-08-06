# AZG001

The AZG001 analyzer reports `err := SomeFunc()` and `_, err := SomeFunc()` assignments immediately followed by `if err != nil`, which should be combined into a single `if` init statement to keep the error scoped to its check and the happy path unindented.

Assignments are only reported when combining is safe: any non-`err` values must be blank identifiers (`_`), and an `err` declared with `:=` must not be used again after the `if` statement (combining would move it into the `if` statement's scope).

## Flagged Code

```go
_, err := client.Delete(ctx, id)
if err != nil {
	return fmt.Errorf("deleting %s: %+v", id, err)
}
```

```go
err := resourceGroupClient.WaitForDeletion(ctx, id)
if err != nil {
	return fmt.Errorf("waiting for deletion of %s: %+v", id, err)
}
```

## Passing Code

```go
if _, err := client.Delete(ctx, id); err != nil {
	return fmt.Errorf("deleting %s: %+v", id, err)
}
```

```go
if err := resourceGroupClient.WaitForDeletion(ctx, id); err != nil {
	return fmt.Errorf("waiting for deletion of %s: %+v", id, err)
}
```

## Ignoring Reports

When run via golangci-lint, reports can be ignored with a `//nolint:azproviderlint` Go code comment at the end of the offending line or on the line immediately preceding it:

```go
_, err := client.Delete(ctx, id) //nolint:azproviderlint
if err != nil {
	return err
}
```

To ignore only this check on a line — leaving any other azproviderlint checks active — use a `//azignore:AZG001` comment instead, in the same positions:

```go
_, err := client.Delete(ctx, id) //azignore:AZG001
if err != nil {
	return err
}
```
