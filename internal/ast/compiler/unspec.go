package compiler

import (
	"strings"

	"github.com/grafana/cog/internal/ast"
)

var _ Pass = (*Unspec)(nil)

// Unspec removes the Kubernetes-style envelope added by kindsys.
//
// Objects named "spec" will be renamed, using the package as new name.
type Unspec struct {
}

func (pass *Unspec) Process(schemas []*ast.Schema) ([]*ast.Schema, error) {
	for i, schema := range schemas {
		schemas[i] = pass.processSchema(schema)
	}

	return schemas, nil
}

func (pass *Unspec) processSchema(schema *ast.Schema) *ast.Schema {
	newSchema := schema.DeepCopy()
	newSchema.Objects = nil

	for _, object := range schema.Objects {
		if strings.EqualFold(object.Name, "metadata") {
			continue
		}

		newObject := object.DeepCopy()

		if strings.EqualFold(object.Name, "spec") && object.Type.Kind == ast.KindStruct {
			newObject.Name = schema.Package
			if schema.Metadata.Identifier != "" {
				newObject.Name = schema.Metadata.Identifier
			}

			newObject.SelfRef.ReferredType = newObject.Name
		}

		newSchema.Objects = append(newSchema.Objects, newObject)
	}

	return &newSchema
}
