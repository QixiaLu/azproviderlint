# AZD002

The AZD002 analyzer reports data sources calling `metadata.MarkAsGone(...)`.

`MarkAsGone` is for resources, where a deleted remote object should be removed from state. A data source that cannot find what the user asked for should return an error instead, so the user learns why their configuration cannot be applied.

This check only applies to files whose name contains `data_source`.

## Flagged Code

```go
if response.WasNotFound(resp.HttpResponse) {
	return metadata.MarkAsGone(id)
}
```

## Passing Code

```go
if response.WasNotFound(resp.HttpResponse) {
	return fmt.Errorf("%s was not found", id)
}
```

## Ignoring Reports

When run via golangci-lint, reports can be ignored with a `//nolint:azproviderlint` Go code comment at the end of the offending line or on the line immediately preceding it:

```go
return metadata.MarkAsGone(id) //nolint:azproviderlint
```

To ignore only this check on a line — leaving any other azproviderlint checks active — use a `//azignore:AZD002 - <reason>` comment instead, in the same positions:

```go
return metadata.MarkAsGone(id) //azignore:AZD002 - <reason>
```
