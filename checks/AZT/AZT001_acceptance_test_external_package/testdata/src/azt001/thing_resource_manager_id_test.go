// Should NOT be flagged: a validator unit test containing "resource" mid-name
// (storage_queue_resource_manager_id_test.go style), not an acceptance test file.
package azt001

import "testing"

func TestThingResourceManagerID(t *testing.T) {
	if Thing() != "thing" {
		t.Fatal("unexpected")
	}
}
