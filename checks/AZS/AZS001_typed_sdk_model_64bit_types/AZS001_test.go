package AZS001

import (
	"path/filepath"
	"runtime"
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAZS001(t *testing.T) {
	t.Parallel()

	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Join(filepath.Dir(filename), "testdata")

	analysistest.Run(t, dir, Analyzer, "azs001")
}
