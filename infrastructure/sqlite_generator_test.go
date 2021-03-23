package infrastructure

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	_ "github.com/mattn/go-sqlite3"

	"github.com/chuckha/silks/core"
)

func TestSQLiteGenerator_CreateTable(t *testing.T) {
	model := &core.Model{
		Name:      "Offset",
		TableName: "offsets",
		Fields: map[string]*core.Field{
			"EventType": {
				Name:    "EventType",
				ColName: "event_type",
				Type:    "string",
			},
			"LastKnownOffset": {
				Name:    "LastKnownOffset",
				ColName: "last_known_offset",
				Type:    "int",
			},
			"IsTrue": {
				Name:    "IsTrue",
				ColName: "is_true",
				Type:    "bool",
			},
		},
	}
	s := &SQLiteGenerator{}
	ct, err := s.CreateTable(model)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(ct, "EventType") {
		t.Fatal("needs to use colname to generate column names")
	}
	fmt.Println(ct)
	testSql(t, ct)
}

func testSql(t *testing.T, sqlStmt string) {
	dir, err := ioutil.TempDir("", "testing")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)
	db, err := sql.Open("sqlite3", filepath.Join(dir, "temp"))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(sqlStmt); err != nil {
		t.Fatal(err)
	}

}
