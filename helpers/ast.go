package helpers

import (
	"go/ast"
	"go/constant"

	"golang.org/x/tools/go/analysis"
)

func IsTrueConstant(pass *analysis.Pass, e ast.Expr) bool {
	if e == nil {
		return false
	}
	value := pass.TypesInfo.Types[e].Value
	return value != nil && value.Kind() == constant.Bool && constant.BoolVal(value)
}
