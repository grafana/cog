package transforms

import (
	"fmt"

	"github.com/grafana/cog/internal/ir"
)

var _ Transform = (*SchemaSetEntrypoint)(nil)

type SchemaSetEntrypoint struct {
	Package     string // we don't have a "clear" identifier, so we use the package to identify a schema.
	EntryPoint  string
	schemaFound bool
}

func (pass *SchemaSetEntrypoint) Process(schemas ir.Schemas) (ir.Schemas, error) {
	pass.schemaFound = false

	for _, schema := range schemas {
		if schema.Package != pass.Package {
			continue
		}

		if !schema.HasObject(pass.EntryPoint) {
			return nil, fmt.Errorf("can not use %s as entrypoint for schema %s: object not found", pass.EntryPoint, pass.Package)
		}

		schema.EntryPoint = pass.EntryPoint
		schema.EntryPointType = ir.NewRef(schema.Package, pass.EntryPoint)

		pass.schemaFound = true
	}

	return schemas, nil
}

func (pass *SchemaSetEntrypoint) Diagnostics() []string {
	if pass.schemaFound {
		return nil
	}

	return []string{
		fmt.Sprintf("schema '%s' not found", pass.Package),
	}
}
