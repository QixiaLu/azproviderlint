// Package nilguard is the shared guard engine behind AZG008 and AZG009: it decides whether a
// pointer-typed variable or selector chain is provably non-nil at a given point in a function
// body.
package nilguard

import (
	"go/ast"
	"go/token"
	"go/types"
	"strings"

	"github.com/katbyte/azproviderlint/lib/astx"
	"github.com/katbyte/azproviderlint/lib/pointerpkg"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/ast/inspector"
)

// ForEachFunc visits every function body in the pass exactly once — package-level function
// literals (var f = func() {...}) are reached through their own visit, nested literals through
// their enclosing FuncDecl, keeping guards outside a closure visible to dereferences inside
// it. Bodies in _test.go files are skipped when tests is false. visit receives the body, the
// function's (and receiver's) parameter objects, and a child-to-parent map covering the body.
func ForEachFunc(pass *analysis.Pass, insp *inspector.Inspector, tests bool, visit func(body *ast.BlockStmt, params map[types.Object]bool, parents map[ast.Node]ast.Node)) {
	var declRanges [][2]token.Pos
	insp.Preorder([]ast.Node{(*ast.FuncDecl)(nil)}, func(n ast.Node) {
		declRanges = append(declRanges, [2]token.Pos{n.Pos(), n.End()})
	})
	inDecl := func(pos token.Pos) bool {
		for _, r := range declRanges {
			if r[0] <= pos && pos < r[1] {
				return true
			}
		}
		return false
	}

	insp.Preorder([]ast.Node{(*ast.FuncDecl)(nil), (*ast.FuncLit)(nil)}, func(n ast.Node) {
		var body *ast.BlockStmt
		params := map[types.Object]bool{}
		collectParams := func(fl *ast.FieldList) {
			if fl == nil {
				return
			}
			for _, field := range fl.List {
				for _, name := range field.Names {
					if obj := pass.TypesInfo.Defs[name]; obj != nil {
						params[obj] = true
					}
				}
			}
		}
		switch fn := n.(type) {
		case *ast.FuncDecl:
			body = fn.Body
			collectParams(fn.Recv)
			collectParams(fn.Type.Params)
		case *ast.FuncLit:
			if inDecl(fn.Pos()) {
				return // visited via its enclosing FuncDecl
			}
			body = fn.Body
			collectParams(fn.Type.Params)
		}
		if body == nil {
			return
		}

		if !tests && strings.HasSuffix(pass.Fset.Position(body.Pos()).Filename, "_test.go") {
			return
		}

		parents := map[ast.Node]ast.Node{}
		var stack []ast.Node
		ast.Inspect(body, func(x ast.Node) bool {
			if x == nil {
				stack = stack[:len(stack)-1]
				return true
			}
			if len(stack) > 0 {
				parents[x] = stack[len(stack)-1]
			}
			stack = append(stack, x)
			return true
		})

		visit(body, params, parents)
	})
}

// PathKey canonicalizes a variable or field chain (`props.Status`) so guard conditions and
// dereference operands can be compared; the root must resolve to a variable.
func PathKey(pass *analysis.Pass, e ast.Expr) (string, bool) {
	switch x := ast.Unparen(e).(type) {
	case *ast.Ident:
		obj := pass.TypesInfo.Uses[x]
		if obj == nil {
			obj = pass.TypesInfo.Defs[x] // the ident being defined by `x := ...`
		}
		if v, ok := obj.(*types.Var); ok {
			return v.Id() + "@" + pass.Fset.Position(v.Pos()).String(), true
		}
	case *ast.SelectorExpr:
		if base, ok := PathKey(pass, x.X); ok {
			return base + "." + x.Sel.Name, true
		}
	}
	return "", false
}

// DerefNeedsPointer reports whether star's context requires the pointer itself rather than a
// copy of the pointee — an assignment target (*x = v), the operand of & (&*x keeps pointer
// identity), or an inc/dec statement — so no pointer.From rewrite can apply.
func DerefNeedsPointer(parents map[ast.Node]ast.Node, star *ast.StarExpr) bool {
	var outer ast.Node = star
	for {
		p, ok := parents[outer].(*ast.ParenExpr)
		if !ok {
			break
		}
		outer = p
	}
	switch p := parents[outer].(type) {
	case *ast.AssignStmt:
		for _, lhs := range p.Lhs {
			if lhs == outer {
				return true
			}
		}
	case *ast.UnaryExpr:
		return p.Op == token.AND
	case *ast.IncDecStmt:
		return true
	}
	return false
}

// Guarded walks at's ancestors looking for a condition, early-exit, or non-nil assignment
// that proves key is not nil where at sits. When the walk finds that key aliases another chain
// (`payload := existing.Model`), it restarts with the source chain so guards on either side
// of the assignment count; the depth cap bounds alias-of-alias chains.
func Guarded(pass *analysis.Pass, parents map[ast.Node]ast.Node, at ast.Node, key string) bool {
	return guardedKey(pass, parents, at, key, 0)
}

func guardedKey(pass *analysis.Pass, parents map[ast.Node]ast.Node, at ast.Node, key string, depth int) bool {
	if depth > 4 {
		return false
	}
	child := at
	for node := parents[at]; node != nil; child, node = node, parents[node] {
		var stmts []ast.Stmt
		switch p := node.(type) {
		case *ast.BinaryExpr:
			// short-circuit: in `x != nil && *x...` / `x == nil || *x...` the right operand
			// only evaluates once the left proved x non-nil
			if p.Y == child && ((p.Op == token.LAND && impliesNonNil(pass, p.X, key)) ||
				(p.Op == token.LOR && impliedByNil(pass, p.X, key))) {
				return true
			}
		case *ast.IfStmt:
			if p.Body == child && impliesNonNil(pass, p.Cond, key) {
				return true
			}
			// the else branch runs when the condition is false; a pure-|| condition with an
			// `x == nil` disjunct being false proves x non-nil
			if p.Else == child && impliedByNil(pass, p.Cond, key) {
				return true
			}
		case *ast.ForStmt:
			if p.Body == child && p.Cond != nil && impliesNonNil(pass, p.Cond, key) {
				return true
			}
		case *ast.CaseClause:
			// the clause fires when any listed expression is true, so all must prove it
			if len(p.List) > 0 && containsStmt(p.Body, child) {
				all := true
				for _, ce := range p.List {
					all = all && impliesNonNil(pass, ce, key)
				}
				if all {
					return true
				}
			}
			stmts = p.Body
		case *ast.BlockStmt:
			stmts = p.List
		case *ast.CommClause:
			stmts = p.Body
		}
		if stmts == nil {
			continue
		}
		switch ok, alias, settled := precededByGuard(pass, stmts, child, key); {
		case ok:
			return true
		case alias != "":
			return guardedKey(pass, parents, at, alias, depth+1)
		case settled:
			return false // an assignment made the value unknown; guards further out are stale
		}
	}
	return false
}

func containsStmt(stmts []ast.Stmt, n ast.Node) bool {
	for _, s := range stmts {
		if s == n {
			return true
		}
	}
	return false
}

// precededByGuard reports whether a statement before child in stmts proves key non-nil: an
// `if key == nil { <terminating> }` early exit, an assignment of a provably non-nil value, or
// a multi-result call whose companion error/ok result was checked before child. When key was
// assigned from another chain (`payload := existing.Model`), the source chain is returned as
// alias so the caller can restart the guard search with it. Any other assignment settles the
// value as unknown: settled tells the caller to stop — guards in enclosing scopes predate the
// assignment and no longer hold.
func precededByGuard(pass *analysis.Pass, stmts []ast.Stmt, child ast.Node, key string) (guarded bool, alias string, settled bool) {
	idx := -1
	for i, s := range stmts {
		if s == child {
			idx = i
			break
		}
	}
	for i := idx - 1; i >= 0; i-- {
		switch s := stmts[i].(type) {
		case *ast.IfStmt:
			if s.Else == nil && impliedByNil(pass, s.Cond, key) && terminates(s.Body) {
				return true, "", false
			}
		case *ast.AssignStmt:
			for j, lhs := range s.Lhs {
				lk, ok := PathKey(pass, lhs)
				if !ok {
					continue
				}
				if lk != key && !strings.HasPrefix(key, lk+".") {
					continue
				}
				// the assignment settles the value: a provable source proves the guard, an
				// alias of another chain redirects the search, anything else is unknown
				if len(s.Lhs) == len(s.Rhs) {
					if lk == key && isNonNilSource(pass, s.Rhs[j]) {
						return true, "", false
					}
					if rk, ok := PathKey(pass, s.Rhs[j]); ok {
						return false, rk + strings.TrimPrefix(key, lk), false
					}
					return false, "", true
				}
				// `x, err := f()` / `x, ok := f()` — the Go contract makes x valid once a
				// later `if err != nil { return }` / `if !ok { return }` exited
				if lk == key && len(s.Rhs) == 1 {
					if _, isCall := ast.Unparen(s.Rhs[0]).(*ast.CallExpr); isCall {
						return companionCheckedBetween(pass, s, stmts[i+1:idx]), "", true
					}
				}
				return false, "", true
			}
		}
	}
	return false, "", false
}

// companionCheckedBetween reports whether assign returns an error or ok-bool alongside the
// pointer and one of the following statements exits on it (`if err != nil { return ... }`,
// `if !ok { return ... }`) — the conventional Go contracts under which the other results are
// valid.
func companionCheckedBetween(pass *analysis.Pass, assign *ast.AssignStmt, between []ast.Stmt) bool {
	errType := types.Universe.Lookup("error").Type()
	var errKey, okKey string
	for _, lhs := range assign.Lhs {
		id, ok := ast.Unparen(lhs).(*ast.Ident)
		if !ok {
			continue
		}
		obj := pass.TypesInfo.Uses[id]
		if obj == nil {
			obj = pass.TypesInfo.Defs[id]
		}
		if obj == nil {
			continue
		}
		if types.Identical(obj.Type(), errType) {
			errKey, _ = PathKey(pass, lhs)
		} else if basic, isBasic := obj.Type().Underlying().(*types.Basic); isBasic && basic.Kind() == types.Bool {
			okKey, _ = PathKey(pass, lhs)
		}
	}
	if errKey == "" && okKey == "" {
		return false
	}
	for _, s := range between {
		ifs, ok := s.(*ast.IfStmt)
		if !ok || ifs.Else != nil || !terminates(ifs.Body) {
			continue
		}
		// after `if err != nil { return }` err is nil; after `if !ok { return }` ok is true
		if errKey != "" && orContainsNeqNil(pass, ifs.Cond, errKey) {
			return true
		}
		if okKey != "" && orContainsNot(pass, ifs.Cond, okKey) {
			return true
		}
	}
	return false
}

// orContainsNot reports whether cond is `!ok` or a pure-|| chain containing it.
func orContainsNot(pass *analysis.Pass, cond ast.Expr, key string) bool {
	switch x := ast.Unparen(cond).(type) {
	case *ast.BinaryExpr:
		if x.Op == token.LOR {
			return orContainsNot(pass, x.X, key) || orContainsNot(pass, x.Y, key)
		}
	case *ast.UnaryExpr:
		if x.Op == token.NOT {
			k, ok := PathKey(pass, x.X)
			return ok && k == key
		}
	}
	return false
}

// orContainsNeqNil reports whether cond is `key != nil` or a pure-|| chain containing it, so
// cond being false proves key == nil.
func orContainsNeqNil(pass *analysis.Pass, cond ast.Expr, key string) bool {
	x, ok := ast.Unparen(cond).(*ast.BinaryExpr)
	if !ok {
		return false
	}
	if x.Op == token.LOR {
		return orContainsNeqNil(pass, x.X, key) || orContainsNeqNil(pass, x.Y, key)
	}
	return x.Op == token.NEQ && nilComparison(pass, x, key)
}

// terminates reports whether a block's final statement leaves the enclosing flow: a return,
// branch, panic, or fatal call.
func terminates(block *ast.BlockStmt) bool {
	if len(block.List) == 0 {
		return false
	}
	switch last := block.List[len(block.List)-1].(type) {
	case *ast.ReturnStmt, *ast.BranchStmt:
		return true
	case *ast.ExprStmt:
		call, ok := last.X.(*ast.CallExpr)
		if !ok {
			return false
		}
		switch fun := call.Fun.(type) {
		case *ast.Ident:
			return fun.Name == "panic"
		case *ast.SelectorExpr:
			name := fun.Sel.Name
			return name == "Exit" || strings.HasPrefix(name, "Fatal")
		}
	}
	return false
}

// isNonNilSource reports whether expr can never be nil: an address-of, new(...), or
// pointer.To(...) call.
func isNonNilSource(pass *analysis.Pass, expr ast.Expr) bool {
	switch x := ast.Unparen(expr).(type) {
	case *ast.UnaryExpr:
		return x.Op == token.AND
	case *ast.CallExpr:
		switch fun := x.Fun.(type) {
		case *ast.Ident:
			_, isBuiltin := pass.TypesInfo.Uses[fun].(*types.Builtin)
			return isBuiltin && fun.Name == "new"
		case *ast.SelectorExpr:
			fn, ok := pass.TypesInfo.Uses[fun.Sel].(*types.Func)
			if !ok || fn.Pkg() == nil {
				return false
			}
			// pointer.To* always allocates; the flag package's constructors return
			// pointers into the flag set, never nil
			return (fn.Pkg().Path() == pointerpkg.PkgPath && strings.HasPrefix(fn.Name(), "To")) ||
				fn.Pkg().Path() == "flag"
		}
	}
	return false
}

// impliesNonNil reports whether cond being true proves key != nil.
func impliesNonNil(pass *analysis.Pass, cond ast.Expr, key string) bool {
	x, ok := ast.Unparen(cond).(*ast.BinaryExpr)
	if !ok {
		return false
	}
	if x.Op == token.LAND {
		return impliesNonNil(pass, x.X, key) || impliesNonNil(pass, x.Y, key)
	}
	return x.Op == token.NEQ && nilComparison(pass, x, key)
}

// impliedByNil reports whether key == nil forces cond to be true — so cond being false (an
// else branch, a short-circuited ||, a not-taken early exit) proves key != nil.
func impliedByNil(pass *analysis.Pass, cond ast.Expr, key string) bool {
	x, ok := ast.Unparen(cond).(*ast.BinaryExpr)
	if !ok {
		return false
	}
	if x.Op == token.LOR {
		return impliedByNil(pass, x.X, key) || impliedByNil(pass, x.Y, key)
	}
	return x.Op == token.EQL && nilComparison(pass, x, key)
}

// nilComparison reports whether cmp compares key's path against nil (either operand order).
func nilComparison(pass *analysis.Pass, cmp *ast.BinaryExpr, key string) bool {
	matches := func(e, other ast.Expr) bool {
		k, ok := PathKey(pass, e)
		return ok && k == key && astx.IsNilValue(pass, other)
	}
	return matches(cmp.X, cmp.Y) || matches(cmp.Y, cmp.X)
}
