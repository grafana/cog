package golang

import (
	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/jennies/common"
)

type Tools struct {
}

func (jenny Tools) JennyName() string {
	return "GoTools"
}

func (jenny Tools) Generate(i common.Context) (*codejen.File, error) {
	return codejen.NewFile("cog/tools.go", []byte(jenny.generateToPtrFunc()), jenny), nil
}

func (jenny Tools) generateToPtrFunc() string {
	return `package cog

func ToPtr[T any](v T) *T {
  return &v
}

`
}
