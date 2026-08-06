// Should NOT be flagged: a plain unit test file whose name merely ends in list_test.go
package azt001

import "testing"

func TestList(t *testing.T) {
	if Thing() != "thing" {
		t.Fatal("unexpected")
	}
}
