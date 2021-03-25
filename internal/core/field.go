package core

import (
	"fmt"
	"go/ast"
	"strings"

	"github.com/pkg/errors"
)

const (
	// SilksTag is the struct tag used for metadata.
	// The format is a comma separated list.
	// The first item is the column name for a given field
	SilksTag = "slk"
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
	if fld.Tag != nil { // TODO: test a complex tag (`json:"hello" xml:"hi" slk:"greeting,pk"`)
		slkVals, err := slkValuesFromTag(fld.Tag.Value)
		if err != nil {
			return field{}, err
		}
		if len(slkVals) > 0 {
			if len(slkVals[0]) > 0 {
				col = slkVals[0]
			}
		}
	}
	return NewField(fld.Names[0].Name, typ, col)
}

func slkValuesFromTag(tag string) ([]string, error) {
	tag = strings.Trim(tag, "`")
	pkgTags := strings.Split(tag, " ")
	slksTag := ""
	for _, tag := range pkgTags {
		if strings.HasPrefix(tag, SilksTag) {
			slksTag = tag
			break
		}
	}

	// there is no silks tag
	if slksTag == "" {
		return nil, nil
	}

	parts := strings.Split(slksTag, ":")
	// slk tag exists but it's empty `slk`
	if len(parts) < 2 {
		return nil, errors.New(`slk tag is malformed, must be key:"val1,val2,val3"`)
	}
	value := strings.Trim(parts[1], `"`)
	return strings.Split(value, ","), nil
}
