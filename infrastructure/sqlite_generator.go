package infrastructure

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"

	"github.com/chuckha/silks/core"
)

type SQLiteGenerator struct{}

func (s *SQLiteGenerator) CreateTable(model *core.Model) (string, error) {
	cols := []string{}
	for _, field := range model.Fields {
		colDef, err := s.fieldToColDef(field)
		if err != nil {
			return "", err
		}
		cols = append(cols, colDef)
	}
	coldefs := strings.Join(cols, ", ")
	return fmt.Sprintf("CREATE TABLE %s ( %s );", model.TableName, coldefs), nil
}

func (s *SQLiteGenerator) AddField(model *core.Model, field *core.Field) error {
	return errors.New("implement me")
}

func (s *SQLiteGenerator) fieldToColDef(fld *core.Field) (string, error) {
	typ, err := s.goTypeToSQLiteType(fld.Type)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s %s", fld.ColName, typ), nil

}

func (s *SQLiteGenerator) goTypeToSQLiteType(goType string) (string, error) {
	switch goType {
	case "int", "bool":
		return "INTEGER", nil
	case "float64":
		return "REAL", nil
	case "string", "time.Time":
		return "TEXT", nil
	default:
		return "", errors.Errorf("SQLiteGenerator does not support go type %s", goType)
	}
}
