package AZS004

import (
	"path/filepath"
	"runtime"
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAZS004(t *testing.T) {
	t.Parallel()

	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Join(filepath.Dir(filename), "testdata")

	analysistest.Run(t, dir, Analyzer, "azs004")

	// the allow flags are package state read during run, so the flag-on fixture packages must
	// run sequentially within the same test rather than as parallel siblings
	allowMissingValues = true
	analysistest.Run(t, dir, Analyzer, "azs004allowmissing")
	allowMissingValues = false

	allowExtraValues = true
	analysistest.Run(t, dir, Analyzer, "azs004allowextra")
	allowExtraValues = false
}
