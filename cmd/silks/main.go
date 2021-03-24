package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/chuckha/silks"
	"github.com/chuckha/silks/infrastructure"
	"github.com/chuckha/silks/usecases"
)

type baseflags struct {
	sqlDialect string
	modelFile  string
	outputDir  string
}

type cfg struct {
	SQLDialect string
	ModelFile  string
	OutputDir  string
}

func main() {
	cfgfile := os.Getenv("SILKS_CONFIG")
	if cfgfile == "" {
		panic("a configuraiton is required (set SILKS_CONFIG env var)")
	}
	cfgdata, err := os.ReadFile(cfgfile)
	if err != nil {
		panic(err)
	}
	config := &cfg{}
	if err := json.Unmarshal(cfgdata, config); err != nil {
		panic(err)
	}
	addFlags := flag.NewFlagSet("add", flag.ExitOnError)
	model := addFlags.String("model", "", "the model to chnage")
	field := addFlags.String("field", "", "the field to add to the model")
	fieldType := addFlags.String("field-type", "", "the type of the field to add")
	colName := addFlags.String("col-name", "", "a rename of the field for sql (optional)")

	app, err := setup(config.SQLDialect)
	if err != nil {
		panic(err)
	}

	switch os.Args[1] {
	case "init":
		sqls, err := app.GenerateCreateTable(config.SQLDialect, config.ModelFile)
		if err != nil {
			panic(err)
		}
		fmt.Println(sqls)
	case "add":
		err = addFlags.Parse(os.Args[2:])
		if err != nil {
			panic(err)
		}
		addsql, updateModel, err := app.AddField(config.SQLDialect, config.ModelFile, *model, *field, *fieldType, *colName)
		if err != nil {
			panic(err)
		}
		fmt.Println(addsql)
		fmt.Println(updateModel)
	default:
		flag.PrintDefaults()
		os.Exit(1)
	}
}

func setup(dialect string) (*silks.App, error) {
	generator := &infrastructure.SQLGeneratorFactory{}
	sqlGen, err := generator.Get(dialect)
	if err != nil {
		return nil, err
	}
	adapter := &silks.AppAdapter{}
	usecaes := silks.AppUseCases{
		CreateTableGenerator: &usecases.CreateTableGenerator{sqlGen},
		FieldAdder:           &usecases.FieldAdder{sqlGen},
	}
	presenter := &silks.AppPresenter{}
	return silks.NewApp(adapter, usecaes, presenter), nil
}
