# AZR002

The AZR002 analyzer reports resources registering a combined `CreateUpdate` method as their `Create` function.

New resources should define separate `Create` and `Update` methods: combined methods hide which properties are actually updatable, make ignore-changes behaviour harder to reason about, and complicate the eventual migration to the typed SDK. Existing resources with combined `CreateUpdate` methods are being split gradually over time.

## Flagged Code

```go
return &pluginsdk.Resource{
	Create: resourceExampleCreateUpdate,
	Read:   resourceExampleRead,
	Update: resourceExampleCreateUpdate,
	Delete: resourceExampleDelete,
}
```

## Passing Code

```go
return &pluginsdk.Resource{
	Create: resourceExampleCreate,
	Read:   resourceExampleRead,
	Update: resourceExampleUpdate,
	Delete: resourceExampleDelete,
}
```

## Ignoring Reports

When run via golangci-lint, reports can be ignored with a `//nolint:azproviderlint` Go code comment at the end of the offending line or on the line immediately preceding it:

```go
Create: resourceExampleCreateUpdate, //nolint:azproviderlint
```
