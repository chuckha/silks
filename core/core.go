package core

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"

	"github.com/pkg/errors"
)

type ModelFile struct {
	fset       *token.FileSet
	def        *ast.File
	Models     map[string]*Model
	Attributes map[string]*Attributes
}

func (m *ModelFile) AddModelField(model, field, fieldType, colName string) {
	m.Models[model].AddField(field, fieldType, colName)
}

type Model struct {
	def       ast.Decl
	Name      string
	TableName string
	Fields    map[string]*Field
}

func (m *Model) AddField(field, fieldType, colName string) {
	gd := m.def.(*ast.GenDecl)
	ts := gd.Specs[0].(*ast.TypeSpec)
	st := ts.Type.(*ast.StructType)
	newField := &ast.Field{
		Names: []*ast.Ident{
			{
				Name: field,
			},
		},
		Type: &ast.Ident{
			Name: fieldType,
		},
		Tag:     nil,
		Comment: nil,
	}
	if colName != "" {
		newField.Tag = &ast.BasicLit{
			Kind:  token.STRING,
			Value: fmt.Sprintf("`slk:\"%s\"`", colName),
		}
	}

	// Modify the actual AST
	st.Fields.List = append(st.Fields.List, newField)

	// Save the code version of the new field
	m.Fields[field] = &Field{
		Name:    field,
		ColName: colName,
		Type:    fieldType,
	}
}

func (m *Model) SetAttribute(attr *Attributes) {
	tn, ok := attr.Attributes["tablename"]
	if !ok {
		return
	}
	m.TableName = tn
}

func NewModel(decl *ast.GenDecl) (*Model, error) {
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

	fields := map[string]*Field{}
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

// this is responsible for converting go source into our model file struct
// it will error if there are syntax errors in our model file
func NewModelFile(fset *token.FileSet, file *ast.File) (*ModelFile, error) {
	models := map[string]*Model{}
	for _, decl := range file.Decls {
		gd, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}
		if gd.Tok != token.TYPE {
			continue
		}
		model, err := NewModel(gd)
		if err != nil {
			return nil, err
		}
		models[model.Name] = model
	}

	attributes := map[string]*Attributes{}
	// parse through A.B=C comments and ignore anything that doesn't match that pattern
	for _, commentGroup := range file.Comments {
		for _, line := range strings.Split(commentGroup.Text(), "\n") {
			// an interesting line
			parts := strings.Split(line, "=")
			if len(parts) <= 1 {
				continue
			}
			pieces := strings.Split(parts[0], ".")
			if len(pieces) <= 1 {
				continue
			}
			model := pieces[0]
			if _, ok := attributes[model]; !ok {
				attributes[model] = NewAttributes(model)
			}
			attributes[model].AddAttribute(pieces[1], parts[1])
		}
	}

	mf := &ModelFile{
		fset:       fset,
		def:        file,
		Models:     models,
		Attributes: attributes,
	}

	// Delete non-model file comments
	cg := make([]*ast.CommentGroup, 0)

	// associate model with attribute comments found in file
	for modelName, model := range mf.Models {
		attr, ok := mf.Attributes[modelName]
		if !ok {
			continue
		}
		model.SetAttribute(attr)
		cg = append(cg, attr.AsComments(model.def.Pos()-1))
	}
	mf.def.Comments = cg

	return mf, nil
}

func (m *ModelFile) GetASTData() (*token.FileSet, *ast.File) {
	return m.fset, m.def
}

// Attribute are custom comments that affect the output of the system in the models file
type Attributes struct {

	// Model is the associated model
	Model string

	// Attribute is
	Attributes map[string]string
}

func (a *Attributes) AsComments(start token.Pos) *ast.CommentGroup {
	cg := &ast.CommentGroup{
		List: []*ast.Comment{},
	}
	for key, value := range a.Attributes {
		txt := fmt.Sprintf("// %s.%s=%s", a.Model, key, value)
		cg.List = append(cg.List, &ast.Comment{
			Slash: start,
			Text:  txt,
		})
	}
	return cg
}

func (a *Attributes) AddAttribute(key, value string) {
	if key != "tablename" {
		return
	}
	a.Attributes[key] = value
}

func NewAttributes(model string) *Attributes {
	return &Attributes{
		Model:      model,
		Attributes: make(map[string]string),
	}
}

// type comes from a sql conversion type
type Field struct {
	Name    string
	ColName string
	Type    string
}

func FieldFromASTField(field *ast.Field) (*Field, error) {
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

func NewField(name, colName, t string) (*Field, error) {
	if name == "" {
		return nil, errors.New("fields must have a name")
	}
	if t == "" {
		return nil, errors.New("fields must have a type")
	}
	if colName == "" {
		colName = strings.ToLower(name)
	}
	return &Field{name, colName, t}, nil
}

type SQLSyntaxGenerator interface {
	CreateTable(model *Model) (string, error)
	AddField(model *Model, field *Field) error
}
