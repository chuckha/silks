package core

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"

	"github.com/pkg/errors"
)

type Model struct {
	Name       string
	Fields     map[string]field
	Attributes Attributes
}

func NewModel(name string) (*Model, error) {
	if name == "" {
		return nil, errors.New("name is required on a model")
	}
	model := &Model{
		Name:       name,
		Fields:     map[string]field{},
		Attributes: NewAttributes(),
	}
	return model, nil
}

func NewModelFromAST(decl *ast.GenDecl) (*Model, error) {
	if len(decl.Specs) == 0 {
		return nil, errors.New("not sure why there are no specs in this generic declaration")
	}
	spec := decl.Specs[0]
	ts, ok := spec.(*ast.TypeSpec)
	if !ok {
		return nil, errors.New("can only make type specs in this file")
	}

	st, ok := ts.Type.(*ast.StructType)
	if !ok {
		return nil, errors.New("declaration must only be structs, no interfaces")
	}

	fields := map[string]field{}
	for _, field := range st.Fields.List {
		fld, err := FieldFromASTField(field)
		if err != nil {
			return nil, err
		}
		fields[fld.Name] = fld
	}

	attributes := NewAttributes()

	for _, line := range strings.Split(decl.Doc.Text(), "\n") {
		// an interesting line
		parts := strings.Split(line, "=")

		// ignore non-attribute comments preserving godocs
		if len(parts) <= 1 {
			continue
		}
		pieces := strings.Split(parts[0], ".")
		if len(pieces) <= 1 {
			continue
		}
		attr, err := NewAttribute(pieces[1], parts[1])
		if err != nil {
			return nil, err
		}
		attributes.AddAttribute(attr)
	}

	return &Model{
		Name:       ts.Name.Name,
		Fields:     fields,
		Attributes: attributes,
	}, nil
}

func (m *Model) GetTableName() string {
	if tn, ok := m.Attributes[tableNameAttribute]; ok {
		return tn
	}
	return m.Name
}

func (m *Model) AddField(field field) error {
	if _, ok := m.Fields[field.Name]; ok {
		return errors.Errorf("field %q already exists on model %s", field.Name, m.Name)
	}
	m.Fields[field.Name] = field
	return nil
}

func (m *Model) DeleteField(fieldName string) {
	delete(m.Fields, fieldName)
}

func (m *Model) RenameField(fromName, toName, toColName string) error {
	if _, ok := m.Fields[fromName]; !ok {
		return errors.Errorf("field %q does not exist on model %s", fromName, m.Name)
	}
	field := m.Fields[fromName]
	delete(m.Fields, fromName)
	field.Name = toName
	field.colName = toColName
	m.Fields[toName] = field
	return nil
}

func (m *Model) AttributeComments() *ast.CommentGroup {
	comments := []*ast.Comment{}
	for key, value := range m.Attributes {
		comments = append(comments, &ast.Comment{
			Text: fmt.Sprintf("// %s.%s=%s", m.Name, key, value),
		})
	}
	return &ast.CommentGroup{List: comments}
}

func (m *Model) FieldsAsFieldList() *ast.FieldList {
	fieldList := []*ast.Field{}
	for _, field := range m.Fields {
		fld := &ast.Field{
			Names: []*ast.Ident{{Name: field.Name}},
			Type:  &ast.Ident{Name: field.Type},
			Tag:   field.Tag(),
		}
		fieldList = append(fieldList, fld)
	}
	return &ast.FieldList{List: fieldList}
}

func (m *Model) ToAST() *ast.GenDecl {
	return &ast.GenDecl{
		Tok: token.TYPE,
		Doc: m.AttributeComments(),
		Specs: []ast.Spec{
			&ast.TypeSpec{
				Name: &ast.Ident{
					Name: m.Name,
				},
				Type: &ast.StructType{
					Fields: m.FieldsAsFieldList(),
				},
			},
		},
	}
}
