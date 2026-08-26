package transforms

import (
	"strings"

	"github.com/grafana/cog/internal/ir"
)

var _ Transform = (*InferEntrypoint)(nil)

type InferEntrypoint struct {
}

func (pass *InferEntrypoint) Process(schemas ir.Schemas) (ir.Schemas, error) {
	for _, schema := range schemas {
		if schema.EntryPoint != "" {
			continue
		}

		schema.EntryPoint = pass.inferEntrypoint(schema)
		if schema.EntryPoint != "" {
			schema.EntryPointType = schema.Objects.Get(schema.EntryPoint).SelfRef.AsType()
		}
	}

	return schemas, nil
}

func (pass *InferEntrypoint) inferEntrypoint(schema *ir.Schema) string {
	entrypoint := ""

	schema.Objects.Iterate(func(_ string, object ir.Object) {
		if strings.EqualFold(schema.Package, object.Name) {
			entrypoint = object.Name
		}
	})

	return entrypoint
}
