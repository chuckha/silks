package usecases

import (
	"fmt"
	"strings"
	"testing"

	"github.com/chuckha/silks/internal/core"
)

func (g *gen) RenameField(tableName, fromColName, toColName string) string {
	return fmt.Sprintf("%s %s %s", tableName, fromColName, toColName)
}

func TestFieldRenamer_RenameField(t *testing.T) {
	renamer := &FieldRenamer{
		SQLFieldRenamer: &gen{},
	}
	mdl, err := core.NewModel("Tester")
	if err != nil {
		t.Fatal(err)
	}
	mf := &core.ModelFile{
		Name: "main",
		Models: map[string]*core.Model{
			"Tester": mdl,
		},
	}

	for _, f := range []struct{ name, typ, col string }{
		{"Username", "string", "username"},
		{"From", "string", "from"},
	} {
		field, err := core.NewField(f.name, f.typ, f.col)
		if err != nil {
			t.Fatal(err)
		}
		if err := mf.AddModelField("Tester", field); err != nil {
			t.Fatal(err)
		}
	}
	cfg, err := core.NewRenameFieldConfiguration("Tester", "From", "HomeTown", "home_town")
	if err != nil {
		t.Fatal(err)
	}

	sql, err := renamer.RenameField(mf, cfg)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(sql, "from") {
		t.Fatal("the sql does not have the original column name in it")
	}
}
