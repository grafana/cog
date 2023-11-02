package golang

import (
	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/jennies/context"
)

type BuilderInterface struct {
}

func (jenny BuilderInterface) JennyName() string {
	return "GolangOptionsBuilder"
}

func (jenny BuilderInterface) Generate(_ context.Builders) (*codejen.File, error) {
	output := jenny.generateFile()

	return codejen.NewFile("builder.go", []byte(output), jenny), nil
}

func (jenny BuilderInterface) generateFile() string {
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

type Builder[ResourceT any] interface {
  Build() (ResourceT, error)
}
`
}
