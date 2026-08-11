package azg005

import (
	"errors"
	"fmt"
)

// Should NOT be flagged: errors.New for fixed strings.
func validErrorsNew() {
	_ = errors.New("something went wrong")
	_ = errors.New("invalid input")
}

// Should NOT be flagged: fmt.Errorf with format placeholders.
func validErrorf() {
	value := "test"
	_ = fmt.Errorf("value %s is invalid", value)
	_ = fmt.Errorf("count: %d", 42)
	_ = fmt.Errorf("error: %v", errors.New("nested"))
	_ = fmt.Errorf("wrapped: %w", errors.New("cause"))
}

// Should be flagged: fmt.Errorf with a fixed string and no placeholders.
func invalidErrorf() {
	_ = fmt.Errorf("something went wrong") // want `errors\.New`
	_ = fmt.Errorf("invalid input")        // want `errors\.New`
	_ = fmt.Errorf("error occurred")       // want `errors\.New`
}
