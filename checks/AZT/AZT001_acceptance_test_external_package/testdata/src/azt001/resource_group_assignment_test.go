// Should NOT be flagged: a parse/validate style unit test whose name starts with
// "resource_" (an Azure resource group ID parser), not an acceptance test file.
package azt001

import "testing"

func TestResourceGroupAssignmentID(t *testing.T) {
	if Thing() != "thing" {
		t.Fatal("unexpected")
	}
}
