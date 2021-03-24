package core

import (
	"fmt"
	"go/ast"
	"strings"

	"github.com/pkg/errors"
)

func NewField(name, colName, t string) (*field, error) {
	if name == "" {
		return nil, errors.New("fields must have a name")
	}
	if t == "" {
		return nil, errors.New("fields must have a type")
	}
	if colName == "" {
		colName = strings.ToLower(name)
	}
	return &field{name, colName, t}, nil
}

// type comes from a sql conversion type
type field struct {
	Name    string
	ColName string
	Type    string
}

func FieldFromASTField(field *ast.Field) (*field, error) {
	var typ string
	switch v := field.Type.(type) {
	case *ast.Ident:
		typ = v.Name
	case *ast.SelectorExpr:
		pkg, ok := v.X.(*ast.Ident)
		if !ok {
			return nil, errors.New("the package must be an identifier")
		}
		if pkg.Name != "time" {
			return nil, errors.New("a field can only reference the time package")
		}
		typ = fmt.Sprintf("%s.%s", pkg.Name, v.Sel.String())
	default:
		return nil, errors.Errorf("unknown field type on field `%s`", field.Names[0].Name)
	}
	var col = ""
	if field.Tag != nil {
		key, value, err := grabKeyValueFromTag(field.Tag.Value)
		if err != nil {
			return nil, err
		}
		if key == "slk" {
			col = value
		}
	}
	return NewField(field.Names[0].Name, col, typ)
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
