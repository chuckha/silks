package core

import (
	"fmt"
	"go/ast"
	"strings"

	"github.com/pkg/errors"
)

func NewField(name, t, colName string) (field, error) {
	if name == "" {
		return field{}, errors.New("fields must have a name")
	}
	if t == "" {
		return field{}, errors.New("fields must have a type")
	}
	return field{
		Name:    name,
		Type:    t,
		colName: colName,
	}, nil
}

// type comes from a sql conversion type
type field struct {
	Name    string
	Type    string
	colName string
}

func (f field) GetColName() string {
	if f.colName == "" {
		return strings.ToLower(f.Name)
	}
	return f.colName
}

func (f field) Tag() *ast.BasicLit {
	if f.colName == "" {
		return nil
	}
	return &ast.BasicLit{
		Value: fmt.Sprintf("`slk:\"%s\"`", f.colName),
	}
}

func FieldFromASTField(fld *ast.Field) (field, error) {
	var typ string
	switch v := fld.Type.(type) {
	case *ast.Ident:
		typ = v.Name
	case *ast.SelectorExpr:
		pkg, ok := v.X.(*ast.Ident)
		if !ok {
			return field{}, errors.New("the package must be an identifier")
		}
		if pkg.Name != "time" {
			return field{}, errors.New("a field can only reference the time package")
		}
		typ = fmt.Sprintf("%s.%s", pkg.Name, v.Sel.String())
	default:
		return field{}, errors.Errorf("unknown field type on field `%s`", fld.Names[0].Name)
	}
	var col = ""
	if fld.Tag != nil {
		key, value, err := grabKeyValueFromTag(fld.Tag.Value)
		if err != nil {
			return field{}, err
		}
		if key == "slk" {
			col = value
		}
	}
	return NewField(fld.Names[0].Name, typ, col)
}

func grabKeyValueFromTag(tag string) (string, string, error) {
	tag = strings.Trim(tag, "`")
	parts := strings.Split(tag, ":")
	if parts[0] != "slk" {
		return "", "", nil
	}
	if len(parts) < 1 {
		return "", "", errors.New(`slk tag is malformed, must be key:"val"`)
	}
	key := parts[0]
	value := strings.Trim(parts[1], `"`)
	return key, value, nil
}
