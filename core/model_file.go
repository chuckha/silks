package core

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"
)

type ModelFile struct {
	fset       *token.FileSet
	def        *ast.File
	Models     map[string]*Model
	Attributes map[string]*Attributes
}

func (m *ModelFile) AddModelField(model string, field *field) {
	m.Models[model].AddField(field)
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
		model, err := ModelFromAST(gd)
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
