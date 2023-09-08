package compiler

import (
	"github.com/grafana/cog/internal/ast"
)

type Pass interface {
	Process(files []*ast.File) ([]*ast.File, error)
}
