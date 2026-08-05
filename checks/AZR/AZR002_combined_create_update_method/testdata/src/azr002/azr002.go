package azr002

type Resource struct {
	Create func() error
	Read   func() error
	Update func() error
	Delete func() error
}

// Should be flagged: combined CreateUpdate registered as Create
func resourceBadThing() *Resource {
	return &Resource{
		Create: resourceBadThingCreateUpdate, // want `new resources should use separate Create and Update methods instead of a combined CreateUpdate method`
		Read:   resourceBadThingRead,
		Update: resourceBadThingCreateUpdate,
		Delete: resourceBadThingDelete,
	}
}

// Should NOT be flagged: separate Create and Update methods
func resourceGoodThing() *Resource {
	return &Resource{
		Create: resourceGoodThingCreate,
		Read:   resourceGoodThingRead,
		Update: resourceGoodThingUpdate,
		Delete: resourceGoodThingDelete,
	}
}

func resourceBadThingCreateUpdate() error { return nil }
func resourceBadThingRead() error         { return nil }
func resourceBadThingDelete() error       { return nil }

func resourceGoodThingCreate() error { return nil }
func resourceGoodThingRead() error   { return nil }
func resourceGoodThingUpdate() error { return nil }
func resourceGoodThingDelete() error { return nil }
