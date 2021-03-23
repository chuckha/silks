package usecases

import (
	"strings"
	"testing"

	"github.com/chuckha/silks/core"
)

type gen struct{}

func (g *gen) CreateTable(model *core.Model) (string, error) {
	fields := []string{}
	for _, field := range model.Fields {
		fields = append(fields, field.ColName)
	}
	return strings.Join(fields, " "), nil
}

func TestCreateTableGenerator_GenerateCreateTable(t *testing.T) {
	createTableGen := &CreateTableGenerator{&gen{}}
	models := &core.ModelFile{
		Models: map[string]*core.Model{
			"User": {
				Name:      "User",
				TableName: "users",
				Fields: map[string]*core.Field{
					"Username": {
						Name:    "Username",
						ColName: "username",
						Type:    "string",
					},
					"Age": {
						Name:    "Age",
						ColName: "age",
						Type:    "int",
					},
				},
			},
		},
	}
	o, err := createTableGen.GenerateCreateTable(models)
	if err != nil {
		t.Fatal(err)
	}
	if o != "username age" {
		t.Fatalf("expeced `usernameage` but got %s", o)
	}
}
