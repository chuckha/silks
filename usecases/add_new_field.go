package usecases

import (
	"github.com/chuckha/silks/core"
)

type SQLFieldAddGenerator interface {
	AddField(tableName string, colDef *ColDef) (string, error)
}

type FieldAdder struct {
	SQLFieldAddGenerator
}

func (fa *FieldAdder) AddField(modelFile *core.ModelFile, addcfg *core.AddFieldConfiguration) (string, error) {
	field, err := core.NewField(addcfg.FieldToAdd, addcfg.FieldType, addcfg.ColumnName)
	if err != nil {
		return "", err
	}
	// has a side effect of changing the AST
	modelFile.AddModelField(addcfg.Model, field)

	tableName := modelFile.Models[addcfg.Model].TableName
	// get the sql changes
	return fa.SQLFieldAddGenerator.AddField(tableName, NewColDef(addcfg.ColumnName, addcfg.FieldType))
}
