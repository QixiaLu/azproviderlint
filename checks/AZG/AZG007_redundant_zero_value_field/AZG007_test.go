package AZG007

import (
	"path/filepath"
	"runtime"
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAZG007(t *testing.T) {
	t.Parallel()

	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Join(filepath.Dir(filename), "testdata")

	analysistest.RunWithSuggestedFixes(t, dir, Analyzer, "azg007")

	checkTests = true
	analysistest.Run(t, dir, Analyzer, "azg007test")
	checkTests = false
}
