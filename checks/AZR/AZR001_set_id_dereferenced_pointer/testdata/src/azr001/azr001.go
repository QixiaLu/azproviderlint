package azr001

type ResourceData struct{}

func (d *ResourceData) SetId(id string) {}

type ReadResponse struct {
	ID *string
}

type ThingId struct{}

func (ThingId) ID() string { return "" }

// Should be flagged: SetId with a dereferenced pointer from the API response
func badSetId(d *ResourceData, read ReadResponse) {
	d.SetId(*read.ID) // want `SetId should not be passed a dereferenced pointer, use a generated Resource ID Formatter/Parser and id\.ID\(\)`
}

// Should NOT be flagged: SetId with a Resource ID Formatter
func goodSetId(d *ResourceData, id ThingId) {
	d.SetId(id.ID())
}

// Should NOT be flagged: SetId with a plain string
func goodSetIdString(d *ResourceData, id string) {
	d.SetId(id)
}
