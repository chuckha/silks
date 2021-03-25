package core

import (
	"go/ast"
	"go/token"
	"testing"
)

func TestPrimaryKey(t *testing.T) {
	model, err := NewModel("Utensil")
	if err != nil {
		t.Fatal(err)
	}
	field, err := NewField("Shape", "string", "shape")
	if err != nil {
		t.Fatal(err)
	}
	if err := model.AddField(field); err != nil {
		t.Fatal(err)
	}
}

func TestNewModel(t *testing.T) {
	model, err := NewModel("Animal")
	if err != nil {
		t.Fatal(err)
	}
	if model.Name != "Animal" {
		t.Fatal("expected model name to be `Animal`, but it was not")
	}
	if model.GetTableName() != "Animal" {
		t.Fatal("expected table name to be `Animal` but it was not")
	}
}

func TestNewModelFromAST(t *testing.T) {
	mod, err := NewModelFromAST(typeSpec())
	if err != nil {
		t.Fatal(err)
	}
	field, err := NewField("Quantity", "int", "quant_ity")
	if err != nil {
		t.Fatal(err)
	}
	if err := mod.AddField(field); err != nil {
		t.Fatal(err)
	}
	//printer.Fprint(os.Stdout, fset, mod.def)
}

func typeSpec() *ast.GenDecl {
	return &ast.GenDecl{
		Tok: token.TYPE,
		Specs: []ast.Spec{
			&ast.TypeSpec{
				Name: &ast.Ident{Name: "Inventory"},
				Type: &ast.StructType{
					Fields: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{{Name: "Price"}},
								Type:  &ast.Ident{Name: "int"},
							},
						},
					},
				},
			},
		},
	}
}
