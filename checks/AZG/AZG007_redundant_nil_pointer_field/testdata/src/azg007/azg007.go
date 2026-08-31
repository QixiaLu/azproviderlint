package azg007

type Config struct {
	Name    *string
	Options *Options
	Items   []string          // slice - should NOT be flagged
	Data    map[string]string // map - should NOT be flagged
}

type Options struct {
	Enabled bool
}

// Should be flagged: redundant nil on pointer fields
func invalidCases() *Config {
	return &Config{
		Name:    nil, // want `redundant nil assignment to pointer field "Name" - omit the field`
		Options: nil, // want `redundant nil assignment to pointer field "Options" - omit the field`
	}
}

// Should be flagged but must preserve the comment leading the following field.
func invalidPreservesFollowingComment() *Config {
	return &Config{
		Name: nil, // want `redundant nil assignment to pointer field "Name" - omit the field`
		// keep this comment for Options
		Options: &Options{},
	}
}

// Should NOT be flagged: slices and maps default to nil meaningfully
func validSliceAndMap() *Config {
	return &Config{
		Items: nil,
		Data:  nil,
	}
}

// Should NOT be flagged: fields omitted entirely
func validOmitted() *Config {
	return &Config{}
}
