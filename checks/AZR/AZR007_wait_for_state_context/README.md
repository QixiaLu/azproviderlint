# AZR007

The AZR007 analyzer reports use of `pluginsdk.StateChangeConf{...}`. Going forward the provider prefers custom pollers that implement the [go-azure-sdk](https://github.com/hashicorp/go-azure-sdk) `pollers.PollerType` interface and are driven via `pollers.NewPoller(...).PollUntilDone(ctx)`.

The check matches any `StateChangeConf` composite literal by type name (e.g. `pluginsdk.StateChangeConf{...}`), whether taken by value or by pointer.

Reference: [terraform-provider-azurerm#guide-new-resource](https://github.com/hashicorp/terraform-provider-azurerm/blob/main/contributing/topics/guide-new-resource.md)

## Flagged Code

```go
stateConf := &pluginsdk.StateChangeConf{
	Pending: []string{"Creating"},
	Target:  []string{"Created"},
	Refresh: refreshFunc,
	Timeout: 10 * time.Minute,
}
result, err := stateConf.WaitForStateContext(ctx)
```

## Passing Code

```go
pollerType := custompollers.NewMyCustomPoller(...)
poller := pollers.NewPoller(pollerType, 10*time.Second, pollers.DefaultNumberOfDroppedConnectionsToAllow)
if err := poller.PollUntilDone(ctx); err != nil {
	return err
}
```

## Ignoring Reports

When run via golangci-lint, reports can be ignored with a `//nolint:azproviderlint` Go code comment at the end of the offending line or on the line immediately preceding it:

```go
result, err := stateConf.WaitForStateContext(ctx) //nolint:azproviderlint
```
To ignore only this check on a line — leaving any other azproviderlint checks active — use a `//azignore:AZR007` comment instead, in the same positions:

```go
stateConf := &pluginsdk.StateChangeConf{ //azignore:AZR007
```
