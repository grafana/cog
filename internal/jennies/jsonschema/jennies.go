package jsonschema

import (
	"bytes"
	"log/slog"
	"strings"

	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/ir"
	"github.com/grafana/cog/internal/ir/transforms"
	"github.com/grafana/cog/internal/languages"
	schemaparser "github.com/santhosh-tekuri/jsonschema/v6"
)

const LanguageRef = "jsonschema"

type Config struct {
	Debug bool `yaml:"-"`

	// Compact controls whether the generated JSON should be pretty printed or
	// not.
	Compact bool `yaml:"compact"`
}

func (config Config) MergeWithGlobal(global languages.Config) Config {
	newConfig := config
	newConfig.Debug = global.Debug

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
	jenny := codejen.JennyListWithNamer[languages.Context](func(_ languages.Context) string {
		return LanguageRef
	})

	jenny.AppendOneToMany(Schema{Config: config})

	if config.Debug {
		jenny.AddPostprocessors(ValidateSchemas)
	}

	return jenny
}

func (language *Language) Transform(schemas ir.Schemas) (ir.Schemas, error) {
	passes := transforms.Transforms{
		&transforms.DisjunctionWithNullToOptional{},
		&transforms.InferEntrypoint{},
	}

	return passes.Process(language.logger, schemas)
}

func ValidateSchemas(file codejen.File) (codejen.File, error) {
	if !strings.HasSuffix(file.RelativePath, ".json") {
		return file, nil
	}

	schemaReader := bytes.NewReader(file.Data)
	schemaCompiler := schemaparser.NewCompiler()

	if err := schemaCompiler.AddResource("schema", schemaReader); err != nil {
		return file, err
	}

	return file, nil
}
