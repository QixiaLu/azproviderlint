# AZR006 - use timeouts wrappers, not StopContext

The AZR006 analyzer reports `ctx` being assigned directly from the provider meta object (`ctx := meta.(*clients.Client).StopContext`).

Custom Timeouts only work when the StopContext is wrapped with the appropriate timeouts helper (`timeouts.ForCreate`, `ForCreateUpdate`, `ForRead`, `ForUpdate` or `ForDelete`) and the resource configures `Timeouts` on its schema.

## Flagged Code

```go
func resourceExampleCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	ctx := meta.(*clients.Client).StopContext
	// ...
}
```

## Passing Code

```go
func resourceExampleCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()
	// ...
}
```

with the resource configuring Timeouts:

```go
Timeouts: &pluginsdk.ResourceTimeout{
	Create: pluginsdk.DefaultTimeout(30 * time.Minute),
	Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
	Update: pluginsdk.DefaultTimeout(30 * time.Minute),
	Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
},
```

## Ignoring Reports

When run via golangci-lint, reports can be ignored with a `//nolint:azproviderlint` Go code comment at the end of the offending line or on the line immediately preceding it:

```go
ctx := meta.(*clients.Client).StopContext //nolint:azproviderlint
```

To ignore only this check on a line — leaving any other azproviderlint checks active — use a `//azignore:AZR006 - <reason>` comment instead, in the same positions:

```go
ctx := meta.(*clients.Client).StopContext //azignore:AZR006 - <reason>
```
