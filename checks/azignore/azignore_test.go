package azignore

import (
	"path/filepath"
	"runtime"
	"testing"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/analysistest"

	AZG001 "github.com/katbyte/azproviderlint/checks/AZG/AZG001_combine_err_assignment_and_check"
)

func TestLintIgnore(t *testing.T) {
	t.Parallel()

	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Join(filepath.Dir(filename), "testdata")

	analyzers := Wrap([]*analysis.Analyzer{AZG001.Analyzer})

	analysistest.Run(t, dir, analyzers[0], "azignore")
}
