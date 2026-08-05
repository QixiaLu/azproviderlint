package azt001 // want `acceptance test files for resources and data sources must use a _test package to prevent circular dependencies`

import "testing"

func TestAccThingDataSource(t *testing.T) {
	if Thing() != "thing" {
		t.Fatal("unexpected")
	}
}
