package azd002

type ThingId struct{}

type Metadata struct{}

func (m Metadata) MarkAsGone(id ThingId) error { return nil }

// Should NOT be flagged: MarkAsGone in a resource file is the correct behaviour
func exampleResourceRead(metadata Metadata, id ThingId) error {
	return metadata.MarkAsGone(id)
}
