package silks

import (
	"go/parser"
	"go/token"
	"os"

	"github.com/pkg/errors"

	"github.com/chuckha/silks/internal/core"
)

type AppAdapter struct{}

func (s *AppAdapter) AdaptCreateTable(modelFile string) (*core.ModelFile, error) {
	return s.cfgInputToModelFile(modelFile)
}

func (s *AppAdapter) AddField(modelFile, model, field, fieldType, colName string) (*core.ModelFile, *core.AddFieldConfiguration, error) {
	mf, err := s.cfgInputToModelFile(modelFile)
	if err != nil {
		return nil, nil, err
	}
	afc, err := core.NewAddFieldConfiguration(model, field, fieldType, colName)
	if err != nil {
		return nil, nil, err
	}
	return mf, afc, nil
}

func (s *AppAdapter) RenameField(modelFile, model, from, to, toColName string) (*core.ModelFile, *core.RenameFieldConfiguration, error) {
	mf, err := s.cfgInputToModelFile(modelFile)
	if err != nil {
		return nil, nil, err
	}
	rfc, err := core.NewRenameFieldConfiguration(model, from, to, toColName)
	if err != nil {
		return nil, nil, err
	}
	return mf, rfc, nil
}

func (s *AppAdapter) cfgInputToModelFile(modelFile string) (*core.ModelFile, error) {
	// quickly ensure the file exists
	data, err := os.ReadFile(modelFile)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	cfg, err := core.NewConfiguration(data)
	if err != nil {
		return nil, err
	}
	// parse the model file into go (error if go syntax at this point)
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "", cfg.ModelFile, parser.ParseComments)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return core.NewModelFile(file)
}
