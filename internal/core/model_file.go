package core

import (
	"go/ast"
	"go/token"

	"github.com/pkg/errors"
)

type ModelFile struct {
	Name string
	// Models is a map from model name to model def
	Models map[string]*Model
}

func (m *ModelFile) AddModelField(model string, fld field) error {
	if _, ok := m.Models[model]; !ok {
		return errors.New("cannot change a model that does not exist")
	}
	return m.Models[model].AddField(fld)
}

func (m *ModelFile) RenameModelField(model, from, to, newColName string) error {
	if _, ok := m.Models[model]; !ok {
		return errors.New("cannot change a model that doesn't exist")
	}
	return m.Models[model].RenameField(from, to, newColName)
}

// this is responsible for converting go source into our model file struct
func NewModelFile(file *ast.File) (*ModelFile, error) {
	// models lists all models defined in the file
	models := map[string]*Model{}

	// Find all the models
	for _, decl := range file.Decls {
		gd, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}
		if gd.Tok != token.TYPE {
			continue
		}
		model, err := NewModelFromAST(gd)
		if err != nil {
			return nil, err
		}
		models[model.Name] = model
	}

	return &ModelFile{
		Name:   file.Name.Name,
		Models: models,
	}, nil
}

func (m *ModelFile) modelsAsDecls() []ast.Decl {
	out := make([]ast.Decl, 0)
	for _, model := range m.Models {
		out = append(out, model.ToAST())
	}
	return out
}

func (m *ModelFile) ToAST() *ast.File {
	return &ast.File{
		Name:    &ast.Ident{Name: m.Name},
		Decls:   m.modelsAsDecls(),
		Imports: nil, // TODO import time if necessary
	}
}
