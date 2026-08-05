package azd002

import "fmt"

// Should be flagged: MarkAsGone in a data source file
func exampleDataSourceRead(metadata Metadata, id ThingId) error {
	return metadata.MarkAsGone(id) // want `data sources should return an error when a resource cannot be found instead of calling MarkAsGone`
}

// Should NOT be flagged: returning an error instead
func exampleDataSourceReadGood(metadata Metadata, id ThingId) error {
	return fmt.Errorf("thing %v was not found", id)
}
