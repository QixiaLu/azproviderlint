# AZT002

The AZT002 analyzer reports tests reading the provider's credentials from the environment via `os.Getenv("ARM_CLIENT_ID")`, `os.Getenv("ARM_CLIENT_SECRET")` or `os.Getenv("ARM_CLIENT_SECRET_ALT")`.

Test configurations should not reuse the credentials the test framework runs with. Instead, create a User Assigned Identity as part of the test configuration - with as minimal permissions as possible - which is then cleaned up with the rest of the test resources.

## Flagged Code

```go
clientId := os.Getenv("ARM_CLIENT_ID")
clientSecret := os.Getenv("ARM_CLIENT_SECRET")
```

## Passing Code

```hcl
resource "azurerm_user_assigned_identity" "test" {
  name                = "acctest-uai-${var.random_integer}"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_role_assignment" "test" {
  scope                = azurerm_resource_group.test.id
  role_definition_name = "Reader"
  principal_id         = azurerm_user_assigned_identity.test.principal_id
}
```

## Ignoring Reports

When run via golangci-lint, reports can be ignored with a `//nolint:azproviderlint` Go code comment at the end of the offending line or on the line immediately preceding it:

```go
clientId := os.Getenv("ARM_CLIENT_ID") //nolint:azproviderlint
```

To ignore only this check on a line — leaving any other azproviderlint checks active — use a `//azignore:AZT002` comment instead, in the same positions:

```go
clientId := os.Getenv("ARM_CLIENT_ID") //azignore:AZT002
```
