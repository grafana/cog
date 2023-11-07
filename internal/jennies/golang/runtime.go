package golang

import (
	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/jennies/context"
)

type Runtime struct {
}

func (jenny Runtime) JennyName() string {
	return "GoRuntime"
}

func (jenny Runtime) Generate(context context.Builders) (codejen.Files, error) {
	files := make(codejen.Files, 0, len(context.Schemas))

	runtime, err := jenny.Runtime()
	if err != nil {
		return nil, err
	}

	files = append(files, *codejen.NewFile("cog/runtime.go", []byte(runtime), jenny))
	files = append(files, *codejen.NewFile("cog/builder.go", []byte(jenny.generateBuilderInterface()), jenny))
	files = append(files, *codejen.NewFile("cog/errors.go", []byte(jenny.generateErrorTools()), jenny))

	return files, nil
}

func (jenny Runtime) Runtime() (string, error) {
	imports := newImportMap()
	imports.Add("cogvariants", "github.com/grafana/cog/generated/cog/variants")

	return renderTemplate("runtime.tmpl", map[string]any{
		"imports": imports.Format(),
	})
}

func (jenny Runtime) generateBuilderInterface() string {
	return `package cog

type Builder[ResourceT any] interface {
  Build() (ResourceT, error)
}

`
}

func (jenny Runtime) generateErrorTools() string {
	return `package cog

import (
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
	if buildErrs, ok := err.(BuildErrors); ok {
		for _, buildErr := range buildErrs {
			buildErr.Path = rootPath + "." + buildErr.Path
		}

		return buildErrs
	}
	
	if buildErr, ok := err.(*BuildError); ok {
		return BuildErrors{buildErr}
	}

	return BuildErrors{&BuildError{
		Path:    rootPath,
		Message: err.Error(),
	}}
}

`
}
