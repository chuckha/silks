package infrastructure

import (
	"github.com/pkg/errors"

	"github.com/chuckha/silks/core"
)

type SQLGeneratorFactory struct{}

func (*SQLGeneratorFactory) Get(dialect string) (core.SQLSyntaxGenerator, error) {
	dlct, err := core.NewSQLDialect(dialect)
	if err != nil {
		return nil, err
	}
	switch dlct {
	case core.SQLite:
		return &SQLiteGenerator{}, nil
	default:
		return nil, errors.Errorf("unknown sql dialect %q", dialect)
	}
}
