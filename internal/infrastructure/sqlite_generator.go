package infrastructure

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"

	"github.com/chuckha/silks/internal/usecases"
)

type SQLiteGenerator struct{}

func (s *SQLiteGenerator) CreateTable(tableName string, colDefs []*usecases.ColDef) (string, error) {
	cols := []string{}
	for _, colDef := range colDefs {
		cols = append(cols, fmt.Sprintf("%s %s", colDef.Name, colDef.Type))
	}
	coldefs := strings.Join(cols, ", ")
	return fmt.Sprintf("CREATE TABLE %s ( %s );", tableName, coldefs), nil
}

func (s *SQLiteGenerator) AddField(tableName string, colDef *usecases.ColDef) string {
	return fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s %s;", tableName, colDef.Name, colDef.Type)
}

func (s *SQLiteGenerator) RenameField(tableName, fromColName, toColName string) string {
	return fmt.Sprintf("ALTER TABLE %s RENAME COLUMN %s TO %s;", tableName, fromColName, toColName)
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
