package azr003

type ResourceData struct{}

func (d *ResourceData) Get(key string) interface{}           { return nil }
func (d *ResourceData) GetOk(key string) (interface{}, bool) { return nil, false }
func (d *ResourceData) Id() string                           { return "" }

type Resource struct {
	Create func(d *ResourceData, meta interface{}) error
	Read   func(d *ResourceData, meta interface{}) error
	Delete func(d *ResourceData, meta interface{}) error
}

// Should be flagged: d.Get inside the registered Delete function
func resourceBadThing() *Resource {
	return &Resource{
		Read:   resourceBadThingRead,
		Delete: resourceBadThingDelete,
	}
}

func resourceBadThingDelete(d *ResourceData, meta interface{}) error {
	name := d.Get("name") // want `d\.Get should not be used within a Delete function as it does not work as expected during deletion`
	_ = name
	return nil
}

// Should NOT be flagged: d.Get in a Read function
func resourceBadThingRead(d *ResourceData, meta interface{}) error {
	name := d.Get("name")
	_ = name
	return nil
}

// Should NOT be flagged: d.GetOk and d.Id in the Delete function
func resourceOkThing() *Resource {
	return &Resource{
		Delete: resourceOkThingDelete,
	}
}

func resourceOkThingDelete(d *ResourceData, meta interface{}) error {
	if _, ok := d.GetOk("name"); ok {
		return nil
	}
	_ = d.Id()
	return nil
}

// typed SDK style resources

type ResourceMetaData struct {
	ResourceData *ResourceData
}

type ResourceFunc struct {
	Func func(metadata ResourceMetaData) error
}

type BadTypedResource struct{}

// Should be flagged: metadata.ResourceData.Get inside the typed Delete method
func (r BadTypedResource) Delete() ResourceFunc {
	return ResourceFunc{
		Func: func(metadata ResourceMetaData) error {
			name := metadata.ResourceData.Get("name") // want `ResourceData\.Get should not be used within a Delete function as it does not work as expected during deletion`
			_ = name
			return nil
		},
	}
}

// Should NOT be flagged: metadata.ResourceData.Get inside a typed Read method
func (r BadTypedResource) Read() ResourceFunc {
	return ResourceFunc{
		Func: func(metadata ResourceMetaData) error {
			name := metadata.ResourceData.Get("name")
			_ = name
			return nil
		},
	}
}
