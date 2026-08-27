# AZG000

The AZG000 analyzer reports `//azignore:` directives that do not carry a reason. A suppression without a reason tells the next reader nothing — whether the flagged code is a deliberate exception, a false positive, or debt someone meant to come back to — so every directive must say why the check does not apply:

```go
//azignore:<Rule>[,<Rule>...] - <reason>
```

The reason is free text following the rule list; the `-` separator (a `–` or `—` also works) is permitted but not required — `//azignore:AZR001 deliberate subset` parses the same. Rule names have a fixed shape (letters then digits), so the rule list ends at the first token that is not rule-shaped and the reason is never confused with it, dashes and commas included.

Bare directives still suppress their target checks — this check's report is the enforcement, so one problem produces one actionable message rather than the suppressed check's report reappearing alongside it.

## Flagged Code

```go
d.SetId(*read.ID) //azignore:AZR001

//azignore:AZG001,AZR003
err := client.Delete(ctx, id)

d.SetId(*read.ID) //azignore:AZR001 -
```

## Passing Code

```go
d.SetId(*read.ID) //azignore:AZR001 - legacy resource, ID formatter tracked in #1234

//azignore:AZG001,AZR003 combined form obscures the retry loop here
err := client.Delete(ctx, id)
```

## Ignoring Reports

AZG000 deliberately does not honour `//azignore:AZG000` — a bare directive could otherwise suppress the report about itself by adding `AZG000` to its rule list.

When run via golangci-lint, reports can still be ignored with a `//nolint:azproviderlint` Go code comment at the end of the offending line or on the line immediately preceding it, and providers that accept bare directives as policy can disable the check entirely via the plugin's `disable` setting:

```yaml
linters:
  settings:
    custom:
      azproviderlint:
        settings:
          disable: [AZG000]
```
