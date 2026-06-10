package compiler

import (
	"fmt"

	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/tools"
)

var _ Pass = (*Omit)(nil)

// Omit rewrites schemas to omit the configured objects.
type Omit struct {
	Objects      []ObjectReference
	objectsFound []string
}

func (pass *Omit) Process(schemas []*ast.Schema) ([]*ast.Schema, error) {
	pass.objectsFound = nil

	for i, schema := range schemas {
		schemas[i] = pass.processSchema(schema)
	}

	return schemas, nil
}

func (pass *Omit) processSchema(schema *ast.Schema) *ast.Schema {
	schema.Objects = schema.Objects.Filter(func(_ string, object ast.Object) bool {
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
