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

	"github.com/chuckha/silks/internal/usecases"
)

func TestSQLiteGenerator(t *testing.T) {
	// sql set up
	dir, err := ioutil.TempDir("", "testing")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)
	db, err := sql.Open("sqlite3", filepath.Join(dir, "temp"))
	if err != nil {
		t.Fatal(err)
	}

	// test set up
	s := &SQLiteGenerator{}

	tableName := "busses"
	defs := []*usecases.ColDef{
		{
			"number", "string",
		},
		{
			"route", "string",
		},
	}

	t.Run("create table", func(tt *testing.T) {
		createStmt, err := s.CreateTable(tableName, defs)
		if err != nil {
			tt.Fatal(err)
		}
		if strings.Contains(createStmt, "EventType") {
			tt.Fatal("needs to use colname to generate column names")
		}
		fmt.Println(createStmt)
		if _, err := db.Exec(createStmt); err != nil {
			tt.Fatal(err)
		}
	})

	t.Run("add field", func(tt *testing.T) {
		field := &usecases.ColDef{Name: "max_passengers", Type: "int"}
		addStmt := s.AddField(tableName, field)
		if _, err := db.Exec(addStmt); err != nil {
			tt.Fatal(err)
		}
	})
}
