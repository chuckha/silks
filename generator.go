package silks

import (
	"github.com/pkg/errors"

	"github.com/chuckha/silks/internal/infrastructure"
	"github.com/chuckha/silks/internal/usecases"
)

type SQLSyntaxGenerator interface {
	CreateTable(tableName string, colDefs []*usecases.ColDef) (string, error)
	AddField(tableName string, colDef *usecases.ColDef) string
	RenameField(tableName, fromColName, toColName string) string
}

type SQLGeneratorFactory struct{}

func (*SQLGeneratorFactory) Get(dialect string) (SQLSyntaxGenerator, error) {
	switch dialect {
	case "sqlite", "": // empty string defaults to sqlite
		return &infrastructure.SQLiteGenerator{}, nil
	default:
		return nil, errors.Errorf("unknown sql dialect %q", dialect)
	}
}
