# azproviderlint

[![GitHub release](https://img.shields.io/github/v/release/katbyte/azproviderlint?color=blueviolet)](https://github.com/katbyte/azproviderlint/releases/latest)
![build](https://github.com/katbyte/azproviderlint/actions/workflows/build.yaml/badge.svg)
![test](https://github.com/katbyte/azproviderlint/actions/workflows/test.yaml/badge.svg)
![lint](https://github.com/katbyte/azproviderlint/actions/workflows/lint.yaml/badge.svg)
![govulncheck](https://github.com/katbyte/azproviderlint/actions/workflows/govulncheck.yaml/badge.svg)
![CodeQL](https://github.com/katbyte/azproviderlint/actions/workflows/codeql-analysis.yml/badge.svg)
[![Go Version](https://img.shields.io/github/go-mod/go-version/katbyte/azproviderlint?color=00ADD8)](https://github.com/katbyte/azproviderlint/blob/main/go.mod)
[![License](https://img.shields.io/github/license/katbyte/azproviderlint?color=blue)](https://github.com/katbyte/azproviderlint/blob/main/LICENSE)

A custom [golangci-lint](https://golangci-lint.run/) module plugin providing Azure provider-specific linting rules built on Go's `analysis` framework.

## Installation

### Standalone

```bash
go install github.com/katbyte/azproviderlint@latest
```

Then run directly:
```bash
azproviderlint ./...
```

### As a golangci-lint Plugin

Add to your `.custom-gcl.yml`:
```yaml
version: v2.12.2
plugins:
  - module: "github.com/katbyte/azproviderlint"
    import: "github.com/katbyte/azproviderlint/plugin"
    version: v0.1.0
```

Build the custom binary:
```bash
golangci-lint custom
```

Then enable in `.golangci.yml`:
```yaml
linters:
  enable:
    - azproviderlint
  settings:
    custom:
      azproviderlint:
        type: module
```

Individual rules can be enabled/disabled via plugin settings (an empty `enable` list means all rules):
```yaml
linters:
  settings:
    custom:
      azproviderlint:
        type: module
        settings:
          disable: [AZR002]
```

## Rules

Rules are named `AZ<category letter><number>`, aligned with [tfproviderlint](https://github.com/bflad/tfproviderlint)'s category letters (`R`, `S`, `V`, `AT`) where they overlap.

### AZG — General Go Style / Readability

| Rule | Description |
|------|-------------|
| [AZG001](checks/AZG/AZG001_combine_err_assignment_and_check) | `_, err := SomeFunc()` followed by `if err != nil` should be combined into a single `if` init statement |
| [AZG002](checks/AZG/AZG002_error_should_describe_expected_format) | Error messages should describe the expected format instead of saying `invalid format of ...` |

### AZR — Resource Implementation

| Rule | Description |
|------|-------------|
| [AZR001](checks/AZR/AZR001_set_id_dereferenced_pointer) | `SetId` must not be passed a dereferenced pointer (`d.SetId(*read.ID)`) — use a generated Resource ID Formatter/Parser and `d.SetId(id.ID())` |
| [AZR002](checks/AZR/AZR002_combined_create_update_method) | Resources must register separate `Create` and `Update` methods instead of a combined `CreateUpdate` method |
| [AZR003](checks/AZR/AZR003_resource_data_get_in_delete) | `d.Get` / `metadata.ResourceData.Get` must not be used inside a resource's Delete function, where it does not work as expected |
| [AZR004](checks/AZR/AZR004_resource_id_equality_comparison) | Resource IDs must not be compared with `==`/`!=` — use `resourceids.Match` |
| [AZR005](checks/AZR/AZR005_case_insensitive_segments_feature_flag) | `features.TreatUserSpecifiedSegmentsAsCaseInsensitive` must not be set — the case-aware comparisons feature is not ready for use |
| [AZR006](checks/AZR/AZR006_stop_context_without_timeouts) | `ctx` must not be assigned directly from `meta.(*clients.Client).StopContext` — use `timeouts.ForCreate`/`ForRead`/`ForUpdate`/`ForDelete` so Custom Timeouts work |

### AZD — Data Sources

| Rule | Description |
|------|-------------|
| [AZD001](checks/AZD/AZD001_data_source_empty_set_id) | Data sources must return an error when a resource cannot be found, not call `d.SetId("")` |
| [AZD002](checks/AZD/AZD002_data_source_mark_as_gone) | Data sources must return an error when a resource cannot be found, not call `metadata.MarkAsGone` |

### AZS — Schema & Typed SDK Models

| Rule | Description |
|------|-------------|
| [AZS001](checks/AZS/AZS001_typed_sdk_model_64bit_types) | Typed SDK model fields (tagged `tfschema`) must use 64-bit numeric types — `int64` not `int`/`int16`/`int32`, `float64` not `float32` — including slices, maps, pointers, named types, and aliases of them |

### AZC — Clients & SDK Usage

| Rule | Description |
|------|-------------|
| [AZC001](checks/AZC/AZC001_client_missing_base_uri) | Azure SDK (track1 & kermit) clients must be created via `NewFoosClientWithBaseURI` with the resource manager endpoint explicitly specified, not `NewFoosClient(o.SubscriptionId)` |

### AZT — Acceptance Testing

| Rule | Description |
|------|-------------|
| [AZT001](checks/AZT/AZT001_acceptance_test_external_package) | Resource and data source acceptance test files must use an external `_test` package to prevent circular dependencies |
| [AZT002](checks/AZT/AZT002_credentials_from_environment) | Tests must not obtain credentials via `os.Getenv("ARM_CLIENT_ID"/"ARM_CLIENT_SECRET"/"ARM_CLIENT_SECRET_ALT")` — create an `azurerm_user_assigned_identity` with minimal permissions instead |

### AZN — Naming Conventions

_No rules yet — reserved for property naming convention rules (e.g. percentage properties using a `_percentage` suffix rather than `_in_percent`)._

### AZV — Validation

_No rules yet — reserved for missing/incorrect validation rules (e.g. string arguments without a `ValidateFunc`)._
