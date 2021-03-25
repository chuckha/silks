package core

import (
	"bytes"
	"fmt"
	"go/format"
	"go/parser"
	"go/token"
	"testing"
)

func TestNewModelFile(t *testing.T) {
	testInput := `package silks

import "time"

// htelkd
// sdfsdkf

// User defines the users model
// User.tablename=users
type User struct {
	ID string ` + "`json:\"id\" slk:\"id,pk\"`" + `
	Created time.Time ` + "`slk:\"created\"`" + `
	Updated time.Time
}


`

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "", testInput, parser.ParseComments)
	if err != nil {
		t.Fatal(fmt.Sprintf("%+v", err))
	}

	o, err := NewModelFile(file)
	if err != nil {
		t.Fatal(fmt.Sprintf("%+v", err))
	}
	if len(o.Models) != 1 {
		t.Fatal("expected 1 model but didn't get any")
	}
	if len(o.Models["User"].Fields) != 3 {
		t.Fatalf("expected 3 fields but got %d", len(o.Models["User"].Fields))
	}
	if o.Models["User"].GetTableName() != "users" {
		t.Fatalf("expected name to be users but it was %s", o.Models["User"].GetTableName())
	}
	if o.Models["User"].Fields["ID"].GetColName() != "id" {
		t.Fatalf("not reading the struct tag data correctly, should be %q but is %q", "id", o.Models["User"].Fields["ID"].GetColName())
	}
	if o.Models["User"].Fields["Updated"].GetColName() != "updated" {
		t.Fatal("failed to automatically downcase the field name")
	}
	var buf bytes.Buffer
	tree := o.ToAST()
	if err := format.Node(&buf, fset, tree); err != nil {
		t.Fatal(fmt.Sprintf("%+v", err))
	}
	//fmt.Println(buf.String())

}
