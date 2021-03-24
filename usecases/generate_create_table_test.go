package usecases

import (
	"strings"
	"testing"

	"github.com/chuckha/silks/core"
)

type gen struct{}

func (g *gen) CreateTable(tableName string, coldefs []*ColDef) (string, error) {
	fields := []string{}
	for _, coldef := range coldefs {
		fields = append(fields, coldef.Name)
	}
	return strings.Join(fields, " "), nil
}

func TestCreateTableGenerator_GenerateCreateTable(t *testing.T) {
	createTableGen := &CreateTableGenerator{&gen{}}

	model, err := core.NewModel("User", "users")
	if err != nil {
		t.Fatal(err)
	}
	field1, err := core.NewField("Username", "username", "string")
	if err != nil {
		t.Fatal(err)
	}
	field2, err := core.NewField("Age", "age", "int")
	if err != nil {
		t.Fatal(err)
	}
	model.AddField(field1)
	model.AddField(field2)
	models := &core.ModelFile{
		Models:     map[string]*core.Model{model.Name: model},
		Attributes: nil,
	}
	o, err := createTableGen.GenerateCreateTable(models)
	if err != nil {
		t.Fatal(err)
	}
	if o != "username age" {
		t.Fatalf("expeced `usernameage` but got %s", o)
	}
}
