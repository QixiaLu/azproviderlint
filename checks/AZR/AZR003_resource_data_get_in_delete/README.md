# AZR003

The AZR003 analyzer reports `d.Get(...)` (untyped resources, in the function registered as `Delete:`) and `metadata.ResourceData.Get(...)` (typed resources, in the `Delete() sdk.ResourceFunc` method) being used inside a resource's Delete function.

During deletion the state may be partial or the config unavailable, so schema reads do not work as expected in Delete. Everything a Delete function needs should come from parsing the Resource ID.

## Flagged Code

```go
func resourceExampleDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	name := d.Get("name").(string)
	// ...
}
```

```go
func (r ExampleResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			name := metadata.ResourceData.Get("name").(string)
			// ...
		},
	}
}
```

## Passing Code

```go
func resourceExampleDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	id, err := parse.ExampleID(d.Id())
	if err != nil {
		return err
	}
	// use id.Name, id.ResourceGroup, ...
}
```

## Ignoring Reports

When run via golangci-lint, reports can be ignored with a `//nolint:azproviderlint` Go code comment at the end of the offending line or on the line immediately preceding it:

```go
name := d.Get("name").(string) //nolint:azproviderlint
```

To ignore only this check on a line — leaving any other azproviderlint checks active — use a `//azignore:AZR003` comment instead, in the same positions:

```go
name := d.Get("name").(string) //azignore:AZR003
```
