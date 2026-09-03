# AZG006 - inline single-use variable only used in a later function call

The AZG006 analyzer reports a single-line variable declaration whose variable is used exactly once, as an argument of a later call whose every other argument is a basic literal or a plain identifier — `x := flattenThing(...)` followed by `d.Set("key", x)` — and should be inlined into the call.

When the sibling arguments are literals or bare names like `ctx`/`id` the temporary's name documents nothing they do not already say (`client.ImportThenPoll(ctx, id, importParameters)` reads no worse as `client.ImportThenPoll(ctx, id, expandMsSqlServerImport(d))`). Calls with any more complex sibling — selector chains, calls, type assertions — are out of scope: there the name is documentation among expressions (`client.CreateOrUpdate(ctx, id, payload)`) and inlining could reorder their evaluation. Multi-line initializers are out of scope since splicing one into an argument list hurts readability. Nothing is reported when a statement between the declaration and the call writes to, takes the address of, or shadows anything the initializer reads — there the temporary preserves a value the intervening code changes, so inlining would change behavior.

Two settings tighten the rule further: `only-when-literals` restricts sibling arguments to basic literals (blocking the rare case where several temporaries feeding one call co-inline into a long line), and `maximum-arguments` skips calls carrying more than that many arguments (default 0 = unlimited).

The consuming statement may be a bare call or the call in an if statement's init (`if err := d.Set("key", x); err != nil`), in the same block as the declaration and at most `max-gap` source lines below it (default 100, settable via `-AZG006.max-gap` or the plugin settings).

The report carries a suggested fix, so `azproviderlint -AZG006 -fix` (or an editor applying the suggested fix) inlines the variable automatically.

## Flagged Code

```go
apns := flattenNotificationHubsAPNSCredentials(props.ApnsCredential)
if err := d.Set("apns_credential", apns); err != nil {
	return fmt.Errorf("setting `apns_credential`: %+v", err)
}
```

## Passing Code

```go
if err := d.Set("apns_credential", flattenNotificationHubsAPNSCredentials(props.ApnsCredential)); err != nil {
	return fmt.Errorf("setting `apns_credential`: %+v", err)
}
```

```go
// the name documents the payload among expression arguments — not flagged
payload := expandThing(d)
client.CreateOrUpdate(ctx, id, payload)
```

## Options

| Option | Default | Effect |
|---|---|---|
| `max-gap` | 100 | maximum source lines between the declaration and the consuming call |
| `only-when-literals` | false | require every sibling argument to be a basic literal (plain identifiers are otherwise also accepted) |
| `maximum-arguments` | 0 | skip calls with more than this many arguments (0 = unlimited) |

Set via `-AZG006.<option>` on the CLI or a rule-name key in the plugin's golangci settings.

## Ignoring Reports

When run via golangci-lint, reports can be ignored with a `//nolint:azproviderlint` Go code comment at the end of the offending line or on the line immediately preceding it. To ignore only this check on a line — leaving any other azproviderlint checks active — use `//azignore:AZG006 - <reason>` instead, in the same positions.
