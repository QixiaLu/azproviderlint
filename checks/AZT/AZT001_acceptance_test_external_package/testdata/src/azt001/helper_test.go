package azt001

import "testing"

// Should NOT be flagged: not a resource or data source test file
func TestHelper(t *testing.T) {
	if Thing() != "thing" {
		t.Fatal("unexpected")
	}
}
