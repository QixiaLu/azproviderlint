package azg005

type input struct {
	Format *string
}

type output struct {
	Format string
	Names  []string
}

func from(p *string) string {
	if p != nil {
		return *p
	}
	return ""
}

// Should be flagged: the temporary is consumed by the next statement's assignment and used
// nowhere else.
func invalidAssignConsumer(in input, out *output) {
	format := from(in.Format) // want `"format" is only used by the following statement and should be inlined`
	out.Format = format
}

// Should be flagged: the temporary is returned by the next statement.
func invalidReturnConsumer(in input) string {
	format := from(in.Format) // want `"format" is only used by the following statement and should be inlined`
	return format
}

// Should be flagged: the pattern inside a case clause.
func invalidInCaseClause(in input, out *output, kind int) {
	switch kind {
	case 1:
		format := from(in.Format) // want `"format" is only used by the following statement and should be inlined`
		out.Format = format
	}
}

// Should NOT be flagged: the temporary is used twice.
func validUsedTwice(in input, out *output) {
	format := from(in.Format)
	out.Format = format
	out.Names = []string{format}
}

// Should be flagged: the consumer no longer needs to be adjacent, only within max-gap lines.
func invalidNotAdjacent(in input, out *output) {
	format := from(in.Format) // want `"format" is only used by the statement on line \d+ and should be inlined`
	out.Names = nil
	out.Format = format
}

// Should NOT be flagged: consumed as a call argument — naming an argument is usually
// intentional documentation.
func validCallArgument(in input) {
	format := from(in.Format)
	use(format)
}

// Should NOT be flagged: the next statement assigns a different value.
func validDifferentValue(in input, out *output) {
	format := from(in.Format)
	out.Format = "fixed"
	use(format)
}

// Should NOT be flagged: multi-value declarations are out of scope.
func validMultiValue(in input, out *output) {
	format, err := fromErr(in.Format)
	_ = err
	out.Format = format
}

// Should NOT be flagged: the consumer's left-hand side contains a call, so inlining would
// reorder it relative to the temporary's initializer.
func validCallOnLhs(in input) {
	format := from(in.Format)
	outputs()[0] = format
}

// Should NOT be flagged: assigning to blank is a discard, not a consumption.
func validBlankConsumer(in input) {
	format := from(in.Format)
	_ = format
}

// Should NOT be flagged: the temporary is also captured by a closure later on.
func validClosureCapture(in input, out *output) {
	format := from(in.Format)
	out.Format = format
	defer func() { use(format) }()
}

// Should NOT be flagged: the return has multiple results.
func validMultiReturn(in input) (string, bool) {
	format := from(in.Format)
	return format, true
}

// Should be flagged: the pattern inside a function literal body.
func invalidInFuncLit(in input, out *output) {
	fn := func() {
		format := from(in.Format) // want `"format" is only used by the following statement and should be inlined`
		out.Format = format
	}
	fn()
}

func use(string)                      {}
func outputs() []string               { return nil }
func fromErr(*string) (string, error) { return "", nil }

// Should be flagged: the temporary's only use is two statements later, within max-gap.
func invalidLaterConsumer(in input, out *output) {
	format := from(in.Format) // want `"format" is only used by the statement on line \d+ and should be inlined`
	out.Names = []string{"a"}
	out.Format = format
}
