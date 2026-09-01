package AZG005

import (
	"path/filepath"
	"runtime"
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAZG005(t *testing.T) {
	t.Parallel()

	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Join(filepath.Dir(filename), "testdata")

	analysistest.RunWithSuggestedFixes(t, dir, Analyzer, "azg005")

	// maxGap is package state read during run, so the flag fixture must run sequentially
	// within the same test rather than as a parallel sibling
	maxGap = 3
	analysistest.Run(t, dir, Analyzer, "azg005maxgap")
	maxGap = 100
}
