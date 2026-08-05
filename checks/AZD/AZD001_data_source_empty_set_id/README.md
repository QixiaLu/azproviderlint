# AZD001

The AZD001 analyzer reports data sources calling `d.SetId("")`.

Unlike resources (where clearing the ID removes a deleted resource from state), a data
source that cannot find what the user asked for should return an error - silently setting
an empty ID leaves the user with unexplained empty attributes and downstream diffs.

This check only applies to files whose name contains `data_source`.

## Flagged Code

```go
if response.WasNotFound(resp.HttpResponse) {
	d.SetId("")
	return nil
}
```

## Passing Code

```go
if response.WasNotFound(resp.HttpResponse) {
	return fmt.Errorf("%s was not found", id)
}
```

## Ignoring Reports

When run via golangci-lint, reports can be ignored with a `//nolint:azproviderlint` Go code
comment at the end of the offending line or on the line immediately preceding it:

```go
d.SetId("") //nolint:azproviderlint
```
