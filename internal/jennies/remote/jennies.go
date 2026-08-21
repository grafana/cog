package remote

import (
	"github.com/grafana/codejen"
	"github.com/grafana/cog/pkg/ast/compiler"
	"github.com/grafana/cog/pkg/languages"
)

type Language struct {
	name   string
	config map[string]any
}

func New(name string, config map[string]any) *Language {
	return &Language{
		name:   name,
		config: config,
	}
}

func (language *Language) Name() string {
	return language.name
}

func (language *Language) Jennies(globalConfig languages.Config) *codejen.JennyList[languages.Context] {
	jenny := codejen.JennyListWithNamer(func(_ languages.Context) string {
		return language.Name()
	})
	jenny.AppendOneToMany(remote{
		globalConfig: globalConfig,
		config:       language.config,
	})

	return jenny
}

func (language *Language) CompilerPasses() compiler.Passes {
	return compiler.Passes{
		// TODO: get from plugin
	}
}

func (language *Language) NullableKinds() languages.NullableConfig {
	return languages.NullableConfig{
		// TODO: get from plugin
	}
}
