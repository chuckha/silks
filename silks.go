package silks

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"

	"github.com/pkg/errors"

	"github.com/chuckha/silks/core"
	"github.com/chuckha/silks/infrastructure"
	"github.com/chuckha/silks/usecases"
)

type App struct {
	adapter   Adapter
	useCases  UseCases
	presenter Presenter
}

func NewApp(a Adapter, u UseCases, p Presenter) *App {
	return &App{
		adapter:   a,
		useCases:  u,
		presenter: p,
	}
}

func (s *App) GenerateCreateTable(sqldialect, modelFile string) (string, error) {
	mf, err := s.adapter.AdaptCreateTable(sqldialect, modelFile)
	if err != nil {
		return "", err
	}
	return s.useCases.GenerateCreateTable(mf)
}

func (s *App) AddField(sqldialect, modelFile, model, field, fieldType, colName string) (string, string, error) {
	mf, afc, err := s.adapter.AddField(sqldialect, modelFile, model, field, fieldType, colName)
	if err != nil {
		return "", "", err
	}
	addStmt, err := s.useCases.AddField(mf, afc)
	if err != nil {
		return "", "", err
	}
	fset, tree := mf.GetASTData()
	updatedModel, err := s.presenter.RewriteModelFile(fset, tree)
	return addStmt, updatedModel, err
}

type Adapter interface {
	AdaptCreateTable(sqldialect, modelFile string) (*core.ModelFile, error)
	AddField(sqldialect, modelFile, model, field, fieldType, colName string) (*core.ModelFile, *core.AddFieldConfiguration, error)
}

type UseCases interface {
	GenerateCreateTable(model *core.ModelFile) (string, error)
	AddField(modelFile *core.ModelFile, addcfg *core.AddFieldConfiguration) (string, error)
}

type AppUseCases struct {
	*usecases.CreateTableGenerator
	*usecases.FieldAdder
}

type GeneratorFactory interface {
	Get(dialect string) (infrastructure.SQLSyntaxGenerator, error)
}

type AppAdapter struct{}

func (s *AppAdapter) AdaptCreateTable(sqldialect, modelFile string) (*core.ModelFile, error) {
	return s.cfgInputToModelFile(sqldialect, modelFile)
}

func (s *AppAdapter) AddField(sqldialect, modelFile, model, field, fieldType, colName string) (*core.ModelFile, *core.AddFieldConfiguration, error) {
	mf, err := s.cfgInputToModelFile(sqldialect, modelFile)
	if err != nil {
		return nil, nil, err
	}
	afc, err := core.NewAddFieldConfiguration(model, field, fieldType, colName)
	if err != nil {
		return nil, nil, err
	}
	return mf, afc, nil
}

func (s *AppAdapter) cfgInputToModelFile(sqlDialect, modelFile string) (*core.ModelFile, error) {
	// quickly ensure the file exists
	data, err := os.ReadFile(modelFile)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	cfg, err := core.NewConfiguration(sqlDialect, data)
	if err != nil {
		return nil, err
	}
	// parse the model file into go (error if go syntax at this point)
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "", cfg.ModelFile, parser.ParseComments)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return core.NewModelFile(fset, file)
}

type Presenter interface {
	RewriteModelFile(*token.FileSet, *ast.File) (string, error)
}

type AppPresenter struct{}

func (a *AppPresenter) RewriteModelFile(fset *token.FileSet, file *ast.File) (string, error) {
	var buf bytes.Buffer
	//err := printer.Fprint(&buf, fset, file)
	err := format.Node(&buf, fset, file)
	return buf.String(), err
}
