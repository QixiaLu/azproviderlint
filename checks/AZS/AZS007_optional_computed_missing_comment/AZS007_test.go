package AZS007

import (
	"path/filepath"
	"runtime"
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAZS007(t *testing.T) {
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Join(filepath.Dir(filename), "testdata")

	analysistest.Run(t, dir, Analyzer, "azs007")
}

func TestAZS007_MigrationPackageSkipped(t *testing.T) {
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Join(filepath.Dir(filename), "testdata")

	excludePackages = "migration"

	analysistest.Run(t, dir, Analyzer, "migration")
}
