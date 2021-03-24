package core

import (
	"go/ast"
	"go/token"
	"testing"
)

func TestNewModel(t *testing.T) {
	model, err := NewModel("Animal", "")
	if err != nil {
		t.Fatal(err)
	}
	if model.Name != "Animal" {
		t.Fatal("expected model name to be `Animal`, but it was not")
	}
	if model.TableName != model.Name {
		t.Fatal("defaulting table name to model name is broken")
	}
}

func TestNewModelFromAST(t *testing.T) {
	mod, err := ModelFromAST(typeSpec())
	if err != nil {
		t.Fatal(err)
	}
	field, err := NewField("Quantity", "quant_ity", "int")
	if err != nil {
		t.Fatal(err)
	}
	mod.AddField(field)
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
