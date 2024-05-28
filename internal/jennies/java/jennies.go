package java

import (
	"fmt"
	"strings"

	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/ast/compiler"
	"github.com/grafana/cog/internal/jennies/common"
	"github.com/grafana/cog/internal/languages"
)

const LanguageRef = "java"

type Config struct {
	GeneratePOM  bool   `yaml:"gen_pom"`
	MavenVersion string `yaml:"maven_version"`
	ProjectPath  string `yaml:"-"`
	PackagePath  string `yaml:"package_path"`
}

func (config *Config) InterpolateParameters(interpolator func(input string) string) {
	config.PackagePath = interpolator(config.PackagePath)
	config.ProjectPath = fmt.Sprintf("src/main/java/%s", strings.ReplaceAll(config.PackagePath, ".", "/"))
}

type Language struct {
	config Config
}

func New(config Config) *Language {
	return &Language{config}
}

func (language *Language) Name() string {
	return LanguageRef
}

func (language *Language) Jennies(globalConfig common.Config) *codejen.JennyList[common.Context] {

	jenny := codejen.JennyListWithNamer[common.Context](func(_ common.Context) string {
		return LanguageRef
	})
	jenny.AppendOneToMany(
		Runtime{config: language.config},
		RawTypes{config: language.config},
		common.If[common.Context](language.config.GeneratePOM, Pom{language.config}),
	)
	jenny.AddPostprocessors(common.GeneratedCommentHeader(globalConfig))

	return jenny
}

func (language *Language) CompilerPasses() compiler.Passes {
	return compiler.Passes{
		&compiler.AnonymousEnumToExplicitType{},
		&compiler.AnonymousStructsToNamed{},
		&compiler.NotRequiredFieldAsNullableType{},
		&compiler.FlattenDisjunctions{},
		&compiler.DisjunctionWithNullToOptional{},
		&compiler.DisjunctionInferMapping{},
		&compiler.DisjunctionToType{},
	}
}

func (language *Language) NullableKinds() languages.NullableConfig {
	return languages.NullableConfig{
		Kinds:              []ast.Kind{ast.KindMap, ast.KindArray, ast.KindRef, ast.KindStruct},
		ProtectArrayAppend: true,
		AnyIsNullable:      true,
	}
}
