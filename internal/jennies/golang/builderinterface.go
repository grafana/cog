package golang

import (
	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/ast"
)

type BuilderInterface struct {
}

func (jenny BuilderInterface) JennyName() string {
	return "GolangOptionsBuilder"
}

func (jenny BuilderInterface) Generate(_ []*ast.Schema) (*codejen.File, error) {
	output := jenny.generateFile()

	return codejen.NewFile("types/builder.go", []byte(output), jenny), nil
}

func (jenny BuilderInterface) generateFile() string {
	return `package types

type Builder[T any] interface {
  Build() *T
}
`
}
