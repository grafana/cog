package transforms

import (
	"log/slog"
	"reflect"

	"github.com/grafana/cog/pkg/ir"
)

type Transforms []Transform

func (passes Transforms) Concat(other Transforms) Transforms {
	concat := make([]Transform, 0, len(passes)+len(other))

	concat = append(concat, passes...)
	concat = append(concat, other...)

	return concat
}

func (passes Transforms) Process(logger *slog.Logger, schemas ir.Schemas) (ir.Schemas, error) {
	var err error
	processedSchemas := schemas.DeepCopy()

	for _, compilerPass := range passes {
		processedSchemas, err = compilerPass.Process(processedSchemas)
		if err != nil {
			return nil, err
		}

		p, ok := compilerPass.(diagnosable)
		if !ok {
			continue
		}

		transformName := reflect.TypeOf(p).Elem().Name()
		for _, msg := range p.Diagnostics() {
			logger.Warn(msg, slog.String("transform", transformName))
		}
	}

	return processedSchemas, nil
}

type Transform interface {
	Process(schemas ir.Schemas) (ir.Schemas, error)
}

type diagnosable interface {
	Diagnostics() []string
}
