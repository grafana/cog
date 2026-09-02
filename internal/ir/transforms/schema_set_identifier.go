package transforms

import (
	"fmt"

	"github.com/grafana/cog/internal/ir"
)

var _ Transform = (*SchemaSetIdentifier)(nil)

// SchemaSetIdentifier overwrites the Metadata.Identifier field of a schema.
type SchemaSetIdentifier struct {
	Package     string // we don't have a "clear" identifier, so we use the package to identify a schema.
	Identifier  string
	schemaFound bool
}

func (pass *SchemaSetIdentifier) Process(schemas ir.Schemas) (ir.Schemas, error) {
	pass.schemaFound = false

	for _, schema := range schemas {
		if schema.Package != pass.Package {
			continue
		}

		schema.Metadata.Identifier = pass.Identifier
		pass.schemaFound = true
	}

	return schemas, nil
}

func (pass *SchemaSetIdentifier) Diagnostics() []string {
	if pass.schemaFound {
		return nil
	}

	return []string{
		fmt.Sprintf("schema '%s' not found", pass.Package),
	}
}
