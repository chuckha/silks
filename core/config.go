package core

import (
	"github.com/pkg/errors"
)

type SQLDialect string

const (
	SQLite SQLDialect = "sqlite"
	// todo....
)

func NewSQLDialect(s string) (SQLDialect, error) {
	switch s {
	case "sqlite":
		return SQLite, nil
	default:
		return "", errors.Errorf("unknown dialect %q", s)
	}
}

type File string

type Configuration struct {
	SQLDialect SQLDialect
	ModelFile  []byte
}

func NewConfiguration(dialect string, data []byte) (*Configuration, error) {
	sqlDialect, err := NewSQLDialect(dialect)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errors.New("model file cannot be empty")
	}
	return &Configuration{
		SQLDialect: sqlDialect,
		ModelFile:  data,
	}, nil
}

type AddFieldConfiguration struct {
	Model      string
	FieldToAdd string
	FieldType  string
	ColumnName string
}

func NewAddFieldConfiguration(model, fieldToAdd, fieldType, columnName string) (*AddFieldConfiguration, error) {
	if model == "" {
		return nil, errors.New("must specify model to add a field to")
	}
	if fieldToAdd == "" {
		return nil, errors.New("must specify the name of the field to add")
	}
	switch fieldType {
	case "string", "int", "[]byte", "time.Time", "float64", "int64":
		break
	default:
		return nil, errors.Errorf("%q is an unsupported go type", fieldType)
	}
	if columnName == "" {
		columnName = fieldToAdd
	}
	return &AddFieldConfiguration{
		Model:      model,
		FieldToAdd: fieldToAdd,
		FieldType:  fieldType,
		ColumnName: columnName,
	}, nil
}
