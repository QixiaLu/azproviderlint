package azd001

// Should be flagged: SetId("") in a data source file
func exampleDataSourceRead(d *ResourceData) error {
	d.SetId("") // want `data sources should return an error when a resource cannot be found instead of calling SetId with an empty string`
	return nil
}

// Should NOT be flagged: setting a real ID
func exampleDataSourceReadGood(d *ResourceData, id string) error {
	d.SetId(id)
	return nil
}
