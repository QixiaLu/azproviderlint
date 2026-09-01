package AZG006

import (
	"path/filepath"
	"runtime"
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAZG006(t *testing.T) {
	t.Parallel()

	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Join(filepath.Dir(filename), "testdata")

	analysistest.RunWithSuggestedFixes(t, dir, Analyzer, "azg006")

	// the flags are package state read during run, so flag fixtures must run sequentially
	// within the same test rather than as parallel siblings
	onlyWhenLiterals = true
	analysistest.Run(t, dir, Analyzer, "azg006literals")
	onlyWhenLiterals = false

	maximumArguments = 2
	analysistest.Run(t, dir, Analyzer, "azg006maxargs")
	maximumArguments = 0

	maxGap = 3
	analysistest.Run(t, dir, Analyzer, "azg006maxgap")
	maxGap = 100
}
