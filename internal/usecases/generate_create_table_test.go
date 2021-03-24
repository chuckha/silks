package usecases

import (
	"strings"
	"testing"

	"github.com/chuckha/silks/internal/core"
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

	model, err := core.NewModel("User")
	if err != nil {
		t.Fatal(err)
	}
	for _, f := range []struct{ name, typ, col string }{
		{"Username", "string", "username"},
		{"Age", "int", "age"},
	} {
		field, err := core.NewField(f.name, f.typ, f.col)
		if err != nil {
			t.Fatal(err)
		}
		if err := model.AddField(field); err != nil {
			t.Fatal(err)
		}
	}

	models := &core.ModelFile{
		Models: map[string]*core.Model{model.Name: model},
	}
	o, err := createTableGen.GenerateCreateTable(models)
	if err != nil {
		t.Fatal(err)
	}
	if o != "username age" {
		t.Fatalf("expeced `usernameage` but got %s", o)
	}
}
