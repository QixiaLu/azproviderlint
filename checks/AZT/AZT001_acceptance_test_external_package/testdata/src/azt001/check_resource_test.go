// Should NOT be flagged: a tool's unit test file whose name ends _resource_test.go but
// contains no acceptance test functions (document-lint check_resource_test.go style).
package azt001

import "testing"

func TestSliceDiff(t *testing.T) {
	if Thing() != "thing" {
		t.Fatal("unexpected")
	}
}
