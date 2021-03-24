package core

import (
	"bytes"
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

type User struct {
	ID string ` + "`slk:\"id\"`" + `
	Created time.Time ` + "`slk:\"created\"`" + `
	Updated time.Time
}

// User.tablename=users

`

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "", testInput, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}

	o, err := NewModelFile(fset, file)
	if err != nil {
		t.Fatal(err)
	}
	if len(o.Models) != 1 {
		t.Fatal("expected 1 model but didn't get any")
	}
	if len(o.Models["User"].Fields) != 3 {
		t.Fatalf("expected 3 fields but got %d", len(o.Models["User"].Fields))
	}
	if o.Models["User"].TableName != "users" {
		t.Fatalf("expected name to be users but it was %s", o.Models["User"].TableName)
	}
	if o.Models["User"].Fields["ID"].ColName != "id" {
		t.Fatalf("not reading the struct tag data correctly, should be %q but is %q", "id", o.Models["User"].Fields["ID"].ColName)
	}
	if o.Models["User"].Fields["Updated"].ColName != "updated" {
		t.Fatal("failed to automatically downcase the field name")
	}
	var buf bytes.Buffer
	fset, tree := o.GetASTData()
	if err := format.Node(&buf, fset, tree); err != nil {
		t.Fatal(err)
	}
	//fmt.Println(buf.String())

}
