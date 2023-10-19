package compiler

import (
	"strings"

	"github.com/grafana/cog/internal/ast"
)

var _ Pass = (*Unspec)(nil)

// Compiler pass to remove (for now) the Kubernetes-style envelope added by kindsys
type Unspec struct {
}

func (pass *Unspec) Process(schemas []*ast.Schema) ([]*ast.Schema, error) {
	newSchemas := make([]*ast.Schema, 0, len(schemas))

	for _, schema := range schemas {
		newSchema, err := pass.processSchema(schema)
		if err != nil {
			return nil, err
		}

		newSchemas = append(newSchemas, newSchema)
	}

	return newSchemas, nil
}

func (pass *Unspec) processSchema(schema *ast.Schema) (*ast.Schema, error) {
	newSchema := schema.DeepCopy()
	newSchema.Objects = nil

	for _, object := range schema.Objects {
		if strings.ToLower(object.Name) == "metadata" {
			continue
		}

		newObject := object.DeepCopy()

		if strings.ToLower(object.Name) == "spec" && object.Type.Kind == ast.KindStruct {
			newObject.Name = schema.Package
			if schema.Metadata.Identifier != "" {
				newObject.Name = schema.Metadata.Identifier
			}
		}

		newSchema.Objects = append(newSchema.Objects, newObject)
	}

	return &newSchema, nil
}
