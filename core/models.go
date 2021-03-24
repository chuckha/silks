package core

import (
	"fmt"
	"go/ast"
	"go/token"

	"github.com/pkg/errors"
)

type Model struct {
	def       ast.Decl
	Name      string
	TableName string
	Fields    map[string]*field
}

func NewModel(name, tableName string) (*Model, error) {
	if name == "" {
		return nil, errors.New("name is required on a model")
	}
	if tableName == "" {
		tableName = name
	}
	model := &Model{
		def: &ast.GenDecl{
			Doc: &ast.CommentGroup{
				List: []*ast.Comment{
					{Text: fmt.Sprintf("// %s.tablename=%s", name, tableName)},
				},
			},
			Specs: []ast.Spec{
				&ast.TypeSpec{
					Name: &ast.Ident{Name: name},
					Type: &ast.StructType{
						Fields: &ast.FieldList{
							List: []*ast.Field{},
						},
					},
				},
			},
		},
		Name:      name,
		TableName: tableName,
		Fields:    map[string]*field{},
	}
	return model, nil
}

func (m *Model) AddField(field *field) {
	m.addASTField(field)
	m.addFieldRepr(field)
}

func (m *Model) addASTField(field *field) {
	gd := m.def.(*ast.GenDecl)
	ts := gd.Specs[0].(*ast.TypeSpec)
	st := ts.Type.(*ast.StructType)
	newField := &ast.Field{
		Names: []*ast.Ident{
			{
				Name: field.Name,
			},
		},
		Type: &ast.Ident{
			Name: field.Type,
		},
	}
	if field.ColName != field.Name {
		newField.Tag = &ast.BasicLit{
			Kind:  token.STRING,
			Value: fmt.Sprintf("`slk:\"%s\"`", field.ColName),
		}
	}

	// Modify the actual AST
	st.Fields.List = append(st.Fields.List, newField)
}

func (m *Model) addFieldRepr(field *field) {
	m.Fields[field.Name] = field
}

func (m *Model) SetAttribute(attr *Attributes) {
	m.addAttributesAST(attr)
	m.addAttributesRepr(attr)
}

func (m *Model) addAttributesRepr(attr *Attributes) {
	tn, ok := attr.Attributes["tablename"]
	if !ok {
		return
	}
	m.TableName = tn
}

func (m *Model) addAttributesAST(attr *Attributes) {
	// TODO: this isn't supported yet, but the idea is to define the types and build the ast in reverse
}

func ModelFromAST(decl *ast.GenDecl) (*Model, error) {
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

	fields := map[string]*field{}
	for _, field := range st.Fields.List {
		fld, err := FieldFromASTField(field)
		if err != nil {
			return nil, err
		}
		fields[fld.Name] = fld
	}

	return &Model{
		def:       decl,
		Name:      ts.Name.Name,
		TableName: ts.Name.Name, // Overwritten with attributes
		Fields:    fields,
	}, nil
}
