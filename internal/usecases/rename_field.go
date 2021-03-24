package usecases

import (
	"github.com/chuckha/silks/internal/core"
)

type SQLFieldRenamer interface {
	RenameField(tableName, fromColName, toColName string) string
}

type FieldRenamer struct {
	SQLFieldRenamer
}

func (f *FieldRenamer) RenameField(modelFile *core.ModelFile, renameCfg *core.RenameFieldConfiguration) (string, error) {
	tableName := modelFile.Models[renameCfg.Model].GetTableName()
	fromColName := modelFile.Models[renameCfg.Model].Fields[renameCfg.From].GetColName()
	renameSQL := f.SQLFieldRenamer.RenameField(tableName, fromColName, renameCfg.NewColumnName)
	if err := modelFile.RenameModelField(renameCfg.Model, renameCfg.From, renameCfg.To, renameCfg.NewColumnName); err != nil {
		return "", err
	}
	return renameSQL, nil
}
