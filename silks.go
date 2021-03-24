package silks

import (
	"go/ast"

	"github.com/chuckha/silks/internal/core"
	"github.com/chuckha/silks/internal/usecases"
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

func (s *App) GenerateCreateTable(modelFile string) (string, error) {
	mf, err := s.adapter.AdaptCreateTable(modelFile)
	if err != nil {
		return "", err
	}
	return s.useCases.GenerateCreateTable(mf)
}

func (s *App) AddField(modelFile, model, field, fieldType, colName string) (string, string, error) {
	mf, afc, err := s.adapter.AddField(modelFile, model, field, fieldType, colName)
	if err != nil {
		return "", "", err
	}
	addStmt, err := s.useCases.AddField(mf, afc)
	if err != nil {
		return "", "", err
	}
	tree := mf.ToAST()
	updatedModel, err := s.presenter.RewriteModelFile(tree)
	return addStmt, updatedModel, err
}

func (s *App) RenameField(modelFile, model, field, toField, toColName string) (string, string, error) {
	mf, rfc, err := s.adapter.RenameField(modelFile, model, field, toField, toColName)
	if err != nil {
		return "", "", err
	}
	renameStmt, err := s.useCases.RenameField(mf, rfc)
	if err != nil {
		return "", "", err
	}
	tree := mf.ToAST()
	updatedModel, err := s.presenter.RewriteModelFile(tree)
	return renameStmt, updatedModel, err
}

// Adapter takes input from another system and converts it to something a usecase can understand
type Adapter interface {
	AdaptCreateTable(modelFile string) (*core.ModelFile, error)
	AddField(modelFile, model, field, fieldType, colName string) (*core.ModelFile, *core.AddFieldConfiguration, error)
	RenameField(modelFile, model, field, toField, toColName string) (*core.ModelFile, *core.RenameFieldConfiguration, error)
}

// UseCases define the behavior of the usecases
type UseCases interface {
	GenerateCreateTable(model *core.ModelFile) (string, error)
	AddField(modelFile *core.ModelFile, addcfg *core.AddFieldConfiguration) (string, error)
	RenameField(modelFile *core.ModelFile, renameCfg *core.RenameFieldConfiguration) (string, error)
}

// AppUseCases are the actual implementation
type AppUseCases struct {
	*usecases.CreateTableGenerator
	*usecases.FieldAdder
	*usecases.FieldRenamer
}

type GeneratorFactory interface {
	Get(dialect string) (SQLSyntaxGenerator, error)
}

// Presenter changes the result of a usecase to something user-friendly
type Presenter interface {
	RewriteModelFile(*ast.File) (string, error)
}
