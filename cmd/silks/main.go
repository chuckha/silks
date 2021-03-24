package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/chuckha/silks"
	"github.com/chuckha/silks/internal/usecases"
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

type addcfg struct {
	model     string
	field     string
	fieldType string
	colName   string
}

type renamecfg struct {
	model     string
	from      string
	to        string
	toColName string
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

	addCfg := &addcfg{}
	addFlags := flag.NewFlagSet("add", flag.ExitOnError)
	addFlags.StringVar(&addCfg.model, "model", "", "the model to change")
	addFlags.StringVar(&addCfg.field, "field", "", "the field to add to the model")
	addFlags.StringVar(&addCfg.fieldType, "field-type", "", "the type of the field to add")
	addFlags.StringVar(&addCfg.colName, "col-name", "", "a rename of the field for sql (optional)")

	renameCfg := &renamecfg{}
	renameFlags := flag.NewFlagSet("rename", flag.ExitOnError)
	renameFlags.StringVar(&renameCfg.model, "model", "", "the model to change")
	renameFlags.StringVar(&renameCfg.from, "from", "", "the field that is changing")
	renameFlags.StringVar(&renameCfg.to, "to", "", "the name the field will change to")
	renameFlags.StringVar(&renameCfg.toColName, "to-col-name", "", "the name of the new field's column'")

	app, err := setup(config.SQLDialect)
	if err != nil {
		panic(err)
	}

	switch os.Args[1] {
	case "init":
		sqls, err := app.GenerateCreateTable(config.ModelFile)
		if err != nil {
			panic(err)
		}
		fmt.Println(sqls)
	case "add":
		err = addFlags.Parse(os.Args[2:])
		if err != nil {
			panic(fmt.Sprintf("%+v", err))
		}
		addsql, updateModel, err := app.AddField(config.ModelFile, addCfg.model, addCfg.field, addCfg.fieldType, addCfg.colName)
		if err != nil {
			panic(fmt.Sprintf("%+v", err))
		}
		fmt.Println(addsql)
		fmt.Println(updateModel)
	case "rename":
		err = renameFlags.Parse(os.Args[2:])
		if err != nil {
			panic(fmt.Sprintf("%+v", err))
		}
		renamesql, renameModel, err := app.RenameField(config.ModelFile, renameCfg.model, renameCfg.from, renameCfg.to, renameCfg.toColName)
		// rename a model field from x to y
		if err != nil {
			panic(fmt.Sprintf("%+v", err))
		}
		fmt.Println(renamesql)
		fmt.Println(renameModel)
	default:
		flag.PrintDefaults()
		os.Exit(1)
	}
}

func setup(dialect string) (*silks.App, error) {
	generator := &silks.SQLGeneratorFactory{}
	sqlGen, err := generator.Get(dialect)
	if err != nil {
		return nil, err
	}
	adapter := &silks.AppAdapter{}
	usecaes := silks.AppUseCases{
		CreateTableGenerator: &usecases.CreateTableGenerator{sqlGen},
		FieldAdder:           &usecases.FieldAdder{sqlGen},
		FieldRenamer:         &usecases.FieldRenamer{sqlGen},
	}
	presenter := &silks.AppPresenter{}
	return silks.NewApp(adapter, usecaes, presenter), nil
}
