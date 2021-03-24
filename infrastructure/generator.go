package infrastructure

import (
	"github.com/pkg/errors"

	"github.com/chuckha/silks/usecases"
)

type SQLSyntaxGenerator interface {
	CreateTable(tableName string, colDefs []*usecases.ColDef) (string, error)
	AddField(tableName string, colDef *usecases.ColDef) (string, error)
}

type SQLGeneratorFactory struct{}

func (*SQLGeneratorFactory) Get(dialect string) (SQLSyntaxGenerator, error) {
	switch dialect {
	case "sqlite", "": // empty string defaults to sqlite
		return &SQLiteGenerator{}, nil
	default:
		return nil, errors.Errorf("unknown sql dialect %q", dialect)
	}
}
