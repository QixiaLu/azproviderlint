package azg004

import (
	"bytes" // nolint: staticcheck
)

type payload struct {
	Body *string
}

// Should be flagged: the file lacks the pointer import and the preceding import carries a
// trailing comment — the inserted import must go after the comment, not between the path and
// its comment (which would re-attach the comment to the new import).
func invalidAfterCommentedImport(p *payload, buf *bytes.Buffer) {
	body := "" // want `pointer\.From`
	if p.Body != nil {
		body = *p.Body
	}
	buf.WriteString(body)
}
