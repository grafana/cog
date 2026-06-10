package compiler

import (
	"log/slog"
	"reflect"

	"github.com/grafana/cog/internal/ast"
)

type Passes []Pass

func (passes Passes) Concat(other Passes) Passes {
	concat := make([]Pass, 0, len(passes)+len(other))

	concat = append(concat, passes...)
	concat = append(concat, other...)

	return concat
}

func (passes Passes) Process(logger *slog.Logger, schemas ast.Schemas) (ast.Schemas, error) {
	var err error
	processedSchemas := schemas.DeepCopy()

	for _, compilerPass := range passes {
		processedSchemas, err = compilerPass.Process(processedSchemas)
		if err != nil {
			return nil, err
		}

		p, ok := compilerPass.(Diagnosable)
		if !ok {
			continue
		}

		passName := reflect.TypeOf(p).Elem().Name()
		for _, msg := range p.Diagnostics() {
			logger.Warn(msg, slog.String("pass", passName))
		}
	}

	return processedSchemas, nil
}

type Pass interface {
	Process(schemas []*ast.Schema) ([]*ast.Schema, error)
}

type Diagnosable interface {
	Diagnostics() []string
}
