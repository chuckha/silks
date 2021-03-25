package core

import (
	"testing"
)

func TestFieldFromASTField(t *testing.T) {
	vls, err := slkValuesFromTag("`json:\"id\" slk:\"id,pk\" xml:\"hello\"`")
	if err != nil {
		t.Fatal(err)
	}
	if len(vls) != 2 {
		t.Fatalf("did not read tag correctly, wanted %v but got %v", []string{"id", "pk"}, vls)
	}
}
