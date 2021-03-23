package usecases

import (
	"github.com/chuckha/silks/core"
)

type FieldAdder struct {
	core.SQLSyntaxGenerator
}

func (fa *FieldAdder) AddField(modelFile *core.ModelFile, addcfg *core.AddFieldConfiguration) error {
	modelFile.AddModelField(addcfg.Model, addcfg.FieldToAdd, addcfg.FieldType, addcfg.ColumnName)
	return nil
}
