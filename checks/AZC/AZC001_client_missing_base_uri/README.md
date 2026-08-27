# AZC001

The AZC001 analyzer reports Azure SDK (track1 & kermit) clients being created with `NewFoosClient(o.SubscriptionId)`.

Clients must be created with `NewFoosClientWithBaseURI(...)` so the resource manager endpoint is explicitly specified - without it the client silently defaults to Azure Public, breaking sovereign and non-public clouds (Azure China, US Government, etc.).

## Flagged Code

```go
client := servers.NewServersClient(o.SubscriptionId)
o.ConfigureClient(&client.Client, o.ResourceManagerAuthorizer)
```

## Passing Code

```go
client := servers.NewServersClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
o.ConfigureClient(&client.Client, o.ResourceManagerAuthorizer)
```

## Ignoring Reports

When run via golangci-lint, reports can be ignored with a `//nolint:azproviderlint` Go code comment at the end of the offending line or on the line immediately preceding it:

```go
client := servers.NewServersClient(o.SubscriptionId) //nolint:azproviderlint
```

To ignore only this check on a line — leaving any other azproviderlint checks active — use a `//azignore:AZC001 - <reason>` comment instead, in the same positions:

```go
client := servers.NewServersClient(o.SubscriptionId) //azignore:AZC001 - <reason>
```
