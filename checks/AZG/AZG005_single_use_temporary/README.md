# AZG005 - inline single-use variable only used in a later assignment or return

The AZG005 analyzer reports single-use temporaries immediately consumed by the next statement: `x := <expr>` followed by `y = x` or `return x`, where `x` has no other use in the function. Such a temporary adds a name without adding information — the inlined form reads just as well in one statement.

The consuming statement must be a plain single assignment (`y = x`) or a single-value `return x`, in the same block as the declaration and at most `max-gap` source lines below it (default 100, settable via `-AZG005.max-gap` or the plugin settings) — inlining across a short gap reads fine, while teleporting an initializer hundreds of lines down its function is a readability loss. Call arguments are deliberately out of scope, since naming an argument is usually intentional documentation. Assignments whose left-hand side contains a function call are skipped, because inlining would move the temporary's initializer after the left-hand side's operands in evaluation order. Multi-value declarations, non-adjacent consumers, blank-identifier discards and temporaries captured by closures are all ignored.

The report carries a suggested fix, so `azproviderlint -AZG005 -fix` (or an editor applying the suggested fix) inlines the temporary automatically — the initializer is spliced as raw source text, so multi-line initializers keep their exact original formatting.

## Flagged Code

```go
format := pointer.From(input.Format)
output.Format = format
```

```go
name := buildName(input)
return name
```

## Passing Code

```go
output.Format = pointer.From(input.Format)
```

```go
return buildName(input)
```

## Ignoring Reports

A named temporary can be deliberate documentation of an otherwise opaque expression, so suppressions are expected where the name genuinely helps. When run via golangci-lint, reports can be ignored with a `//nolint:azproviderlint` Go code comment at the end of the offending line or on the line immediately preceding it:

```go
format := pointer.From(input.Format) //nolint:azproviderlint
```

To ignore only this check on a line — leaving any other azproviderlint checks active — use a `//azignore:AZG005 - <reason>` comment instead, in the same positions:

```go
format := pointer.From(input.Format) //azignore:AZG005 - <reason>
```
