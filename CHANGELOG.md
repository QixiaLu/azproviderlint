## Unreleased

- add rule `AZG006`: detect `flatten*` functions that return `nil` for a slice result instead of an empty slice (`[]T{}`/`make([]T, 0)`) — covers the slice position of multi-value returns, named slice types (`type ACLs []T`), and multi-name result fields matched by returned-value position; skips error paths (`return nil, err`) and returns inside nested closures so only the empty/nil-input branch is flagged, and leaves non-slice results and `expand*` functions alone; reports carry a suggested fix applied via `-fix`

## v0.3.2 (2026-08-28)

- `AZS004`: track-1 advice now suggests `validation.StringInEnumSlice(cdn.PossibleTransformValues(), false)` when the call's validation package exports a generic `StringInEnumSlice` wrapper ([azurerm#33246](https://github.com/hashicorp/terraform-provider-azurerm/pull/33246)); `pointer.FromEnumSlice(pointer.To(...))` remains the fallback

## v0.3.1 (2026-08-28)

- plugin settings: rule names are matched case-insensitively (golangci's YAML decoding lowercases map keys)
- `AZS004`: advice for track-1 enums (`Possible<Enum>Values() []Enum`) now compiles, via `pointer.FromEnumSlice(pointer.To(...))`

## v0.3.0 (2026-08-27)

- `//azignore` directives now take an optional reason after the rule list (`//azignore:AZR001 - deliberate subset`)
- add rule `AZG000`: report `//azignore` directives without a reason
- `AZS004`: also report list values that are not part of the enum; new `allow-missing-values` and `allow-extra-values` flags
- `AZS006`: new `ignore-sensitive` flag, and `//azignore:AZS006` now works on individual resource properties
- add rule `AZR007`: `StateChangeConf` from the plugin SDK's `helper/retry` should be a custom poller implementing `pollers.PollerType`

## v0.2.0 (2026-08-18)

- `AZT002` now only checks `_test.go` files
- add rule `AZG003`: `pointer.To(sdk.SomeEnum(v))` should use the generic `pointer.ToEnum[sdk.SomeEnum](v)`; fixable with `-fix`
- add rule `AZG004`: zero-value declaration plus nil-check dereference should use `pointer.From(x)`; fixable with `-fix`
- add rule `AZG005`: single-use temporaries immediately consumed by the next statement should be inlined; fixable with `-fix`
- add rule `AZS002`: schema `Default` values must match the declared `Type` (ports tfproviderlint [#329](https://github.com/bflad/tfproviderlint/pull/329) S038)
- add rule `AZS003`: optional/required `TypeList` blocks must not allow empty blocks (ports tfproviderlint [#236](https://github.com/bflad/tfproviderlint/pull/236) XS003)
- add rule `AZS004`: enum validation should use the SDK's `PossibleValuesFor<Enum>()` helper instead of a hand-written list
- add rule `AZS005`: registered resources should have a data source of the same name
- add rule `AZS006`: data sources should not be missing schema properties that exist on the same-named resource
- build and scan with Go 1.25.13 (fixes GO-2026-6218); govulncheck honours `.go-version`

## v0.1.0 (2026-08-07)

Initial release!

- add rule `AZG001`: `_, err := SomeFunc()` followed by `if err != nil` should be a single `if` init statement
- add rule `AZS001`: typed SDK model numeric fields (tagged `tfschema`) must be `int64`/`float64`
- port the grep/sed based checks from terraform-provider-azurerm's `scripts/checks/` to AST-based rules:
  - `AZG002`: unclear `invalid format of ...` error messages
  - `AZR001`: `d.SetId(*ptr)` instead of a Resource ID Formatter/Parser's `id.ID()`
  - `AZR002`: combined `CreateUpdate` methods instead of separate Create and Update
  - `AZR003`: `d.Get`/`metadata.ResourceData.Get` inside Delete functions
  - `AZC001`: Azure SDK clients created without an explicit resource manager endpoint
  - `AZR004`: Resource IDs compared with `==`/`!=` instead of `resourceids.Match`
  - `AZR005`: assignments to the unreleased `TreatUserSpecifiedSegmentsAsCaseInsensitive` feature flag
  - `AZD001`: data sources calling `d.SetId("")` instead of returning an error
  - `AZD002`: data sources calling `metadata.MarkAsGone` instead of returning an error
  - `AZR006`: `ctx` assigned from `meta.(*clients.Client).StopContext` without a timeouts wrapper
  - `AZT001`: resource/data source acceptance test files not using a `_test` package
  - `AZT002`: tests reading `ARM_CLIENT_ID`/`ARM_CLIENT_SECRET` credentials from the environment
- release binaries with goreleaser (linux/darwin/windows/freebsd/openbsd/solaris) on tagged releases
- add a `version` subcommand printing the version and git commit
- support per-rule `enable`/`disable` lists via golangci-lint plugin settings
