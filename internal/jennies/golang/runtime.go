package golang

import (
	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/jennies/context"
	"github.com/grafana/cog/internal/jennies/template"
)

type Runtime struct {
	Config Config
}

func (jenny Runtime) JennyName() string {
	return "GoRuntime"
}

func (jenny Runtime) Generate(_ context.Builders) (codejen.Files, error) {
	runtime, err := jenny.Runtime()
	if err != nil {
		return nil, err
	}

	return codejen.Files{
		*codejen.NewFile("cog/builder.go", []byte(jenny.generateBuilderInterface()), jenny),
		*codejen.NewFile("cog/errors.go", []byte(jenny.generateErrorTools()), jenny),
		*codejen.NewFile("cog/runtime.go", []byte(runtime), jenny),
	}, nil
}

func (jenny Runtime) generateBuilderInterface() string {
	return `package cog

type Builder[ResourceT any] interface {
  Build() (ResourceT, error)
}

`
}

func (jenny Runtime) Runtime() (string, error) {
	imports := template.NewImportMap()
	imports.Add("cogvariants", jenny.Config.importPath("cog/variants"))

	return renderTemplate("runtime/runtime.tmpl", map[string]any{
		"imports": formatImports(imports),
	})
}

func (jenny Runtime) generateErrorTools() string {
	return `package cog

import (
	"errors"
	"fmt"
)

type BuildErrors []*BuildError

func (errs BuildErrors) Error() string {
	var b []byte
	for i, err := range errs {
		if i > 0 {
			b = append(b, '\n')
		}
		b = append(b, err.Error()...)
	}
	return string(b)
}

type BuildError struct {
	Path string
	Message string
}

func (err *BuildError) Error() string {
	return fmt.Sprintf("%s: %s", err.Path, err.Message)
}

func MakeBuildErrors(rootPath string, err error) BuildErrors {
	var buildErrs BuildErrors
	if errors.As(err, &buildErrs) {
		for _, buildErr := range buildErrs {
			buildErr.Path = rootPath + "." + buildErr.Path
		}

		return buildErrs
	}
	
	var buildErr *BuildError
	if errors.As(err, &buildErr) {
		return BuildErrors{buildErr}
	}

	return BuildErrors{&BuildError{
		Path:    rootPath,
		Message: err.Error(),
	}}
}

`
}
