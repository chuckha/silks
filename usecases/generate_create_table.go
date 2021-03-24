package usecases

import (
	"strings"

	"github.com/chuckha/silks/core"
)

type CreateTableSQLCreator interface {
	CreateTable(tableName string, colDefs []*ColDef) (string, error)
}

type CreateTableGenerator struct {
	Gen CreateTableSQLCreator
}

func (c *CreateTableGenerator) GenerateCreateTable(modelFile *core.ModelFile) (string, error) {
	createStmts := []string{}
	for _, model := range modelFile.Models {
		colDefs := []*ColDef{}
		for _, field := range model.Fields {
			colDefs = append(colDefs, NewColDef(field.ColName, field.Type))
		}
		stmt, err := c.Gen.CreateTable(model.TableName, colDefs)
		if err != nil {
			return "", err
		}
		createStmts = append(createStmts, stmt)
	}
	return strings.Join(createStmts, "\n"), nil
}

type ColDef struct {
	Name string
	Type string
}

func NewColDef(name, colType string) *ColDef {
	return &ColDef{Name: name, Type: colType}
}
