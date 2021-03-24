package usecases

import (
	"github.com/chuckha/silks/internal/core"
)

type SQLFieldAddGenerator interface {
	AddField(tableName string, colDef *ColDef) string
}

type FieldAdder struct {
	SQLFieldAddGenerator
}

func (fa *FieldAdder) AddField(modelFile *core.ModelFile, addcfg *core.AddFieldConfiguration) (string, error) {
	field, err := core.NewField(addcfg.FieldToAdd, addcfg.FieldType, addcfg.ColumnName)
	if err != nil {
		return "", err
	}
	if err := modelFile.AddModelField(addcfg.Model, field); err != nil {
		return "", err
	}

	tableName := modelFile.Models[addcfg.Model].GetTableName()
	// get the sql changes
	return fa.SQLFieldAddGenerator.AddField(tableName, NewColDef(addcfg.ColumnName, addcfg.FieldType)), nil
}
