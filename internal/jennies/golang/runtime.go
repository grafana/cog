package golang

import (
	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/jennies/template"
	"github.com/grafana/cog/internal/languages"
)

type Runtime struct {
	Tmpl   *template.Template
	Config Config
}

func (jenny Runtime) JennyName() string {
	return "GoRuntime"
}

func (jenny Runtime) Generate(_ languages.Context) (codejen.Files, error) {
	runtime, err := jenny.Runtime()
	if err != nil {
		return nil, err
	}

	files := []codejen.File{
		*codejen.NewFile("cog/runtime.go", runtime, jenny),
	}

	if jenny.Config.generateBuilders {
		files = append(files,
			*codejen.NewFile("cog/builder.go", jenny.generateBuilderInterface(), jenny),
			*codejen.NewFile("cog/errors.go", jenny.generateErrorTools(), jenny),
			*codejen.NewFile("cog/tools.go", jenny.generateToPtrFunc(), jenny),
		)
	}

	return files, nil
}

func (jenny Runtime) generateBuilderInterface() []byte {
	return []byte(`package cog

type Builder[ResourceT any] interface {
  Build() (ResourceT, error)
}

`)
}

func (jenny Runtime) Runtime() ([]byte, error) {
	imports := NewImportMap()
	imports.Add("", jenny.Config.importPath("cog/variants"))

	return jenny.Tmpl.RenderAsBytes("runtime/runtime.tmpl", map[string]any{
		"imports": imports,
	})
}

func (jenny Runtime) generateErrorTools() []byte {
	return []byte(`package cog

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

func Unptr[T any](v *T) T {
	var val T
	if v == nil {
		return val
	}
	return *v
}

`)
}

func (jenny Runtime) generateToPtrFunc() []byte {
	return []byte(`package cog

func ToPtr[T any](v T) *T {
  return &v
}

`)
}
