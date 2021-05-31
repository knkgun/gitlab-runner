package runner

import (
	"fmt"
	"go/ast"
	"os"
	"strings"

	"golang.org/x/tools/go/packages"
)

type testcase struct {
	pkg  string
	name string
}

type visitor struct {
	tests []string
}

func (v *visitor) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.SelectorExpr:
		s, ok := n.X.(fmt.Stringer)

		switch {
		case !ok, n.Sel.Name == "TestMain", !strings.HasPrefix(n.Sel.Name, "Test"):
			return v

		case s.String() == "_test" || s.String() == "_xtest":
			v.tests = append(v.tests, n.Sel.Name)
		}
	}

	return v
}

// list is similar to `go test -list` but much faster.
func list(wd string, flags []string, pkgList []string) ([]testcase, error) {
	cfg := &packages.Config{
		Mode:       packages.NeedFiles | packages.NeedSyntax,
		Tests:      true,
		Env:        os.Environ(),
		Dir:        wd,
		BuildFlags: flags,
	}

	pkgs, err := packages.Load(cfg, pkgList...)
	var testcases []testcase
	for _, pkg := range pkgs {
		if !strings.HasSuffix(pkg.ID, ".test") {
			continue
		}

		v := &visitor{}
		for _, af := range pkg.Syntax {
			ast.Walk(v, af)
		}

		for _, test := range v.tests {
			testcases = append(testcases, testcase{
				pkg:  strings.TrimSuffix(pkg.ID, ".test"),
				name: test,
			})
		}
	}

	return testcases, err
}
