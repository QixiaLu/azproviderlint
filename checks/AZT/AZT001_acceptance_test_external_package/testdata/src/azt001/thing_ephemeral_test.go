package azt001 // want `acceptance test files must use a _test package to prevent circular dependencies`

import "testing"

func TestAccThingEphemeral(t *testing.T) {
	if Thing() != "thing" {
		t.Fatal("unexpected")
	}
}
