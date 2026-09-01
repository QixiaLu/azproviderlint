package tf

import (
	"go/ast"
	"path/filepath"
	"strings"

	"golang.org/x/tools/go/analysis"
)

// InDataSourceFile reports whether the node sits in a file whose name marks it as a data
// source file.
func InDataSourceFile(pass *analysis.Pass, n ast.Node) bool {
	filename := filepath.Base(pass.Fset.Position(n.Pos()).Filename)
	return strings.Contains(filename, "data_source")
}
