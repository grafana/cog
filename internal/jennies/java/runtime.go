package java

import (
	"fmt"
	"path/filepath"

	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/jennies/template"
	"github.com/grafana/cog/internal/languages"
)

type Runtime struct {
	config Config
	tmpl   *template.Template
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

	files := codejen.Files{
		*codejen.NewFile(filepath.Join(jenny.config.ProjectPath, "cog/variants/Dataquery.java"), []byte(variants), jenny),
		*codejen.NewFile(filepath.Join(jenny.config.ProjectPath, "cog/Builder.java"), []byte(builder), jenny),
	}

	if jenny.config.generateConverters {
		converter, err := jenny.renderConverterInterface()
		if err != nil {
			return nil, err
		}

		files = append(files, *codejen.NewFile(filepath.Join(jenny.config.ProjectPath, "cog/Converter.java"), []byte(converter), jenny))
	}

	return files, nil
}

func (jenny Runtime) renderDataQueryVariant(variant string) (string, error) {
	return jenny.tmpl.Render("runtime/variants.tmpl", map[string]any{
		"Package": jenny.formatPackage("cog.variants"),
		"Variant": variant,
	})
}

func (jenny Runtime) renderBuilderInterface() (string, error) {
	return jenny.tmpl.Render("runtime/builder.tmpl", map[string]any{
		"Package": jenny.formatPackage("cog"),
	})
}

func (jenny Runtime) renderConverterInterface() (string, error) {
	imports := NewImportMap(jenny.config.PackagePath)
	imports.Add("Dataquery", "cog.variants")

	return jenny.tmpl.Render("runtime/converter_interface.tmpl", map[string]any{
		"Package": jenny.formatPackage("cog"),
		"Imports": imports.String(),
	})
}

func (jenny Runtime) formatPackage(pkg string) string {
	if jenny.config.PackagePath != "" {
		return fmt.Sprintf("%s.%s", jenny.config.PackagePath, pkg)
	}

	return pkg
}
