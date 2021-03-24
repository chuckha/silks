package silks

import (
	"testing"
)

func TestSQLGeneratorFactory_Get(t *testing.T) {
	fact := &SQLGeneratorFactory{}
	_, err := fact.Get("")
	if err != nil {
		t.Fatal(err)
	}
	_, err = fact.Get("postgres")
	if err == nil {
		t.Fatal("there should an error but it was not received")
	}
}
