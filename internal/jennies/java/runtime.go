package java

import (
	"fmt"
	"path/filepath"
	gotemplate "text/template"

	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/jennies/template"
	"github.com/grafana/cog/internal/languages"
)

type Runtime struct {
	config Config
	tmpl   *gotemplate.Template
}

func (jenny Runtime) JennyName() string {
	return "JavaRuntime"
}

func (jenny Runtime) Generate(_ languages.Context) (codejen.Files, error) {
	variants, err := jenny.renderDataQueryVariant("Dataquery")
	if err != nil {
		return nil, err
	}

	builder, err := jenny.renderBuilderInterface()
	if err != nil {
		return nil, err
	}

	return codejen.Files{
		*codejen.NewFile(filepath.Join(jenny.config.ProjectPath, "cog/variants/Dataquery.java"), []byte(variants), jenny),
		*codejen.NewFile(filepath.Join(jenny.config.ProjectPath, "cog/Builder.java"), []byte(builder), jenny),
	}, nil
}

func (jenny Runtime) renderDataQueryVariant(variant string) (string, error) {
	return template.Render(jenny.tmpl, "runtime/variants.tmpl", map[string]any{
		"Package": jenny.formatPackage("cog.variants"),
		"Variant": variant,
	})
}

func (jenny Runtime) renderBuilderInterface() (string, error) {
	return template.Render(jenny.tmpl, "runtime/builder.tmpl", map[string]any{
		"Package": jenny.formatPackage("cog"),
	})
}

func (jenny Runtime) formatPackage(pkg string) string {
	if jenny.config.PackagePath != "" {
		return fmt.Sprintf("%s.%s", jenny.config.PackagePath, pkg)
	}

	return pkg
}
