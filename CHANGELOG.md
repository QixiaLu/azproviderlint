## Unreleased

- add rule `AZS008`: detect `registration.go` `Registration` method map keys and slice elements that are not sorted alphabetically — entries assigned to a local variable and then returned are followed to their composite literal definition, and sections separated by blank lines or comment lines are validated independently so grouped registrations keep their sections; reports carry a suggested fix applied via `-fix` that reorders each unsorted section, moving every entry's trailing comment with it

## v0.6.0 (2026-09-02)

- add rule `AZG002`: single-use temporaries whose only use is `&v` should be inlined as `new(<expr>)` (go1.26; or `pointer.To` via `use`); new mode also rewrites existing `pointer.To(x)` calls unless `allow: pointer.To`; fixable with `-fix` ([#29](https://github.com/katbyte/azproviderlint/pull/29))
- **breaking**: rename `AZG002` to `AZV001` — it polices validation error messages, not general Go style, so it moves to the reserved AZV validation category; update any `//azignore:AZG002` comments and settings references ([#28](https://github.com/katbyte/azproviderlint/pull/28))
- plugin settings: `enable`/`disable` entries can name a whole category (`enable: [AZG]` runs every AZG rule); `disable` still applies after `enable`
- build and scan with Go 1.26.8

## v0.5.1 (2026-09-02)

- `AZG005`/`AZG006`: no longer report a temporary when a statement between the declaration and the consumer writes to, takes the address of, or shadows anything the initializer reads — the suggested fix inlined `oldKey := column[y]` past a `column[y] = ...` write, changing the value read

## v0.5.0 (2026-09-02)

- `AZG005` now also flags temporaries consumed by a later statement in the same block, up to `max-gap` source lines away (default 100) ([#25](https://github.com/katbyte/azproviderlint/pull/25))
- `AZR008` now also covers map results, naked returns with an unassigned named result, and provably-nil variables (`var out []T; return out`); a provably-nil error (`var noErr error`) no longer masks a finding as an error path ([#27](https://github.com/katbyte/azproviderlint/pull/27))
- add rule `AZG006`: single-use variables only used as an argument of a later call should be inlined (`x := flattenThing(...)` then `d.Set("key", x)`); siblings must be literals or plain identifiers and the initializer single-line; tuned by `max-gap` (default 100), `only-when-literals`, and `maximum-arguments`; fixable with `-fix` ([#26](https://github.com/katbyte/azproviderlint/pull/26))

## v0.4.0 (2026-09-01)

- add rule `AZR008`: `flatten*` functions should return an empty slice (`[]T{}`) instead of `nil`; error paths (`return nil, err`) are exempt, and naked returns are out of scope; fixable with `-fix` ([#22](https://github.com/katbyte/azproviderlint/pull/22))
- add rule `AZS007`: schema fields with both `Optional: true` and `Computed: true` must have a `// Note: O+C because ...` comment between the two fields explaining why; an `exclude-packages` setting skips listed package names (e.g. state-migration snapshot packages) ([#20](https://github.com/katbyte/azproviderlint/pull/20))

## v0.3.2 (2026-08-28)

- `AZS004`: track-1 advice now suggests `validation.StringInEnumSlice(cdn.PossibleTransformValues(), false)` when the call's validation package exports a generic `StringInEnumSlice` wrapper ([azurerm#33246](https://github.com/hashicorp/terraform-provider-azurerm/pull/33246)); `pointer.FromEnumSlice(pointer.To(...))` remains the fallback ([#21](https://github.com/katbyte/azproviderlint/pull/21))

## v0.3.1 (2026-08-28)

- plugin settings: rule names are matched case-insensitively (golangci's YAML decoding lowercases map keys) ([#19](https://github.com/katbyte/azproviderlint/pull/19))
- `AZS004`: advice for track-1 enums (`Possible<Enum>Values() []Enum`) now compiles, via `pointer.FromEnumSlice(pointer.To(...))` ([#19](https://github.com/katbyte/azproviderlint/pull/19))

## v0.3.0 (2026-08-27)

- `//azignore` directives now take an optional reason after the rule list (`//azignore:AZR001 - deliberate subset`) ([#18](https://github.com/katbyte/azproviderlint/pull/18))
- add rule `AZG000`: report `//azignore` directives without a reason ([#18](https://github.com/katbyte/azproviderlint/pull/18))
- `AZS004`: also report list values that are not part of the enum; new `allow-missing-values` and `allow-extra-values` flags ([#17](https://github.com/katbyte/azproviderlint/pull/17))
- `AZS006`: new `ignore-sensitive` flag, and `//azignore:AZS006` now works on individual resource properties ([#15](https://github.com/katbyte/azproviderlint/pull/15))
- add rule `AZR007`: `StateChangeConf` from the plugin SDK's `helper/retry` should be a custom poller implementing `pollers.PollerType` ([#6](https://github.com/katbyte/azproviderlint/pull/6))

## v0.2.0 (2026-08-18)

- `AZT002` now only checks `_test.go` files ([#13](https://github.com/katbyte/azproviderlint/pull/13))
- add rule `AZG003`: `pointer.To(sdk.SomeEnum(v))` should use the generic `pointer.ToEnum[sdk.SomeEnum](v)`; fixable with `-fix` ([#5](https://github.com/katbyte/azproviderlint/pull/5), [#10](https://github.com/katbyte/azproviderlint/pull/10))
- add rule `AZG004`: zero-value declaration plus nil-check dereference should use `pointer.From(x)`; fixable with `-fix` ([#5](https://github.com/katbyte/azproviderlint/pull/5), [#11](https://github.com/katbyte/azproviderlint/pull/11))
- add rule `AZG005`: single-use temporaries immediately consumed by the next statement should be inlined; fixable with `-fix` ([#12](https://github.com/katbyte/azproviderlint/pull/12))
- add rule `AZS002`: schema `Default` values must match the declared `Type` (ports tfproviderlint [#329](https://github.com/bflad/tfproviderlint/pull/329) S038) ([#4](https://github.com/katbyte/azproviderlint/pull/4))
- add rule `AZS003`: optional/required `TypeList` blocks must not allow empty blocks (ports tfproviderlint [#236](https://github.com/bflad/tfproviderlint/pull/236) XS003) ([#4](https://github.com/katbyte/azproviderlint/pull/4))
- add rule `AZS004`: enum validation should use the SDK's `PossibleValuesFor<Enum>()` helper instead of a hand-written list ([#7](https://github.com/katbyte/azproviderlint/pull/7))
- add rule `AZS005`: registered resources should have a data source of the same name ([#8](https://github.com/katbyte/azproviderlint/pull/8))
- add rule `AZS006`: data sources should not be missing schema properties that exist on the same-named resource ([#9](https://github.com/katbyte/azproviderlint/pull/9))
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
