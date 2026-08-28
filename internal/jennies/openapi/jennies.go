package openapi

import (
	"log/slog"

	"github.com/grafana/codejen"
	"github.com/grafana/cog/pkg/ir"
	"github.com/grafana/cog/pkg/ir/transforms"
	"github.com/grafana/cog/pkg/languages"
)

const LanguageRef = "openapi"

type Config struct {
	debug bool

	// Compact controls whether the generated JSON should be pretty printed or
	// not.
	Compact bool `yaml:"compact"`
}

func (config Config) MergeWithGlobal(global languages.Config) Config {
	newConfig := config
	newConfig.debug = global.Debug

	return newConfig
}

type Language struct {
	logger *slog.Logger
	config Config
}

func New(logger *slog.Logger, config Config) *Language {
	return &Language{
		logger: logger,
		config: config,
	}
}

func (language *Language) Name() string {
	return LanguageRef
}

func (language *Language) Jennies(globalConfig languages.Config) *codejen.JennyList[languages.Context] {
	config := language.config.MergeWithGlobal(globalConfig)
	jenny := codejen.JennyListWithNamer(func(_ languages.Context) string {
		return LanguageRef
	})

	jenny.AppendOneToMany(Schema{Config: config})

	return jenny
}

func (language *Language) Transform(schemas ir.Schemas) (ir.Schemas, error) {
	passes := transforms.Transforms{
		// should be a superset of the compiler passes defined for jsonschema jennies
		&transforms.DisjunctionWithNullToOptional{},
		&transforms.InferEntrypoint{},
	}

	return passes.Process(language.logger, schemas)
}
