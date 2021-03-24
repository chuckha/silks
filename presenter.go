package silks

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/token"
)

type AppPresenter struct{}

func (a *AppPresenter) RewriteModelFile(file *ast.File) (string, error) {
	var buf bytes.Buffer
	fset := token.NewFileSet()
	//err := printer.Fprint(&buf, fset, file)
	err := format.Node(&buf, fset, file)
	return buf.String(), err
}
