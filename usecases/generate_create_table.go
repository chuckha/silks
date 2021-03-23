package usecases

import (
	"strings"

	"github.com/chuckha/silks/core"
)

type CreateTableSQLCreator interface {
	CreateTable(model *core.Model) (string, error)
}

type CreateTableGenerator struct {
	Gen CreateTableSQLCreator
}

func (c *CreateTableGenerator) GenerateCreateTable(modelFile *core.ModelFile) (string, error) {
	createStmts := []string{}
	for _, model := range modelFile.Models {
		stmt, err := c.Gen.CreateTable(model)
		if err != nil {
			return "", err
		}
		createStmts = append(createStmts, stmt)
	}
	return strings.Join(createStmts, "\n"), nil
}
