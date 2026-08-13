## v0.2.0 (Unreleased)

- add rule `AZG003`: detect `pointer.To(sdk.SomeEnum(v))` explicit go-azure-sdk enum conversions that should use the generic `pointer.ToEnum[sdk.SomeEnum](v)` helper — enums are recognised via the generated `PossibleValuesFor<Name>()` helper (string-backed only, aliases resolved), and conversions of non-string values get a `string(...)` hint in the message; reports carry a suggested fix applied via `-fix`
- add rule `AZG004`: detect zero-value initialization followed by a nil check and pointer dereference (`y := <zero>; if x != nil { y = *x }`) that should use the generic `pointer.From(x)` helper — covers both `:=` and `var` declarations, and skips function-call expressions where the rewrite would change evaluation
- add rule `AZG005`: detect single-use temporaries immediately consumed by the next statement (`x := <expr>` then `y = x` / `return x` with no other use of `x`) that should be inlined — call arguments are deliberately out of scope, and consumers whose left-hand side contains a call are skipped to preserve evaluation order; reports carry a suggested fix applied via `-fix`
- add rule `AZR007`: detect `pluginsdk.StateChangeConf` usage that should use a custom poller implementing `pollers.PollerType`

## v0.1.0 (2026-08-07)

Initial release!

- add rule `AZG001`: detect `_, err := SomeFunc()` followed by `if err != nil` that should be combined into a single `if` init statement
- add rule `AZS001`: detect typed SDK model fields (tagged `tfschema`) using non-64-bit numeric types (`int`, `int16`, `int32`, `float32`) instead of `int64`/`float64` — resolves named types and aliases via the type checker
- add rule `AZS002`: detect schema `Default` values whose type does not match the declared `Type` — resolves named constants via the type checker (ports tfproviderlint [#329](https://github.com/bflad/tfproviderlint/pull/329) S038 with direct constant-kind comparison)
- add rule `AZS003`: detect optional/required `TypeList` blocks that allow empty blocks — every property optional with no default and no `AtLeastOneOf`/`ExactlyOneOf` constraint (ports tfproviderlint [#236](https://github.com/bflad/tfproviderlint/pull/236) XS003)
- port the grep/sed based checks from terraform-provider-azurerm's `scripts/checks/` (`gradually-deprecated.sh`, `timeouts-check.sh`, `test-package-check.sh`) to AST-based rules:
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
