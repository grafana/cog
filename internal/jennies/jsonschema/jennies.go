package jsonschema

import (
	"bytes"
	"strings"

	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/ast/compiler"
	"github.com/grafana/cog/internal/jennies/common"
	schemaparser "github.com/santhosh-tekuri/jsonschema"
)

const LanguageRef = "jsonschema"

type Config struct {
	Debug bool
}

func (config Config) MergeWithGlobal(global common.Config) Config {
	newConfig := config
	newConfig.Debug = global.Debug

	return newConfig
}

type Language struct {
	config Config
}

func New() *Language {
	return &Language{
		config: Config{},
	}
}

func (language *Language) Jennies(globalConfig common.Config) *codejen.JennyList[common.Context] {
	config := language.config.MergeWithGlobal(globalConfig)
	jenny := codejen.JennyListWithNamer[common.Context](func(_ common.Context) string {
		return LanguageRef
	})

	jenny.AppendOneToMany(Schema{Config: config})

	if config.Debug {
		jenny.AddPostprocessors(ValidateSchemas)
	}

	return jenny
}

func (language *Language) CompilerPasses() compiler.Passes {
	return compiler.Passes{
		&compiler.DisjunctionWithNullToOptional{},
		&compiler.InferEntrypoint{},
	}
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
