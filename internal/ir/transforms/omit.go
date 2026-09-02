package transforms

import (
	"fmt"

	"github.com/grafana/cog/internal/ir"
	"github.com/grafana/cog/internal/tools"
)

var _ Transform = (*Omit)(nil)

// Omit rewrites schemas to omit the configured objects.
type Omit struct {
	Objects      []ObjectReference
	objectsFound []string
}

func (pass *Omit) Process(schemas ir.Schemas) (ir.Schemas, error) {
	pass.objectsFound = nil

	for i, schema := range schemas {
		schemas[i] = pass.processSchema(schema)
	}

	return schemas, nil
}

func (pass *Omit) processSchema(schema *ir.Schema) *ir.Schema {
	schema.Objects = schema.Objects.Filter(func(_ string, object ir.Object) bool {
		// if any reference matches the current object, we filter it out
		for _, objectRef := range pass.Objects {
			if objectRef.Matches(object) {
				pass.objectsFound = append(pass.objectsFound, objectRef.String())
				return false
			}
		}

		return true
	})

	return schema
}

func (pass *Omit) Diagnostics() []string {
	if len(pass.objectsFound) == len(pass.Objects) {
		return nil
	}

	expected := tools.Map(pass.Objects, func(ref ObjectReference) string {
		return ref.String()
	})
	missing := tools.SliceFindMissing(pass.objectsFound, expected)

	return tools.Map(missing, func(ref string) string {
		return fmt.Sprintf("object not found '%s'", ref)
	})
}
