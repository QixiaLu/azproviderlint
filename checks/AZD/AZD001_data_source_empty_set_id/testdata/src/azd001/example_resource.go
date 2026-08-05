package azd001

type ResourceData struct{}

func (d *ResourceData) SetId(id string) {}

// Should NOT be flagged: SetId("") in a resource file (used to remove from state)
func exampleResourceRead(d *ResourceData) error {
	d.SetId("")
	return nil
}
