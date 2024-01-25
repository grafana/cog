package compiler

import (
	"fmt"
	"github.com/grafana/cog/internal/ast"
)

var _ Pass = (*GenericOptionToNamed)(nil)

type GenericOptionToNamed struct {
}

func (pass *GenericOptionToNamed) Process(schemas []*ast.Schema) ([]*ast.Schema, error) {
	newSchemas := make([]*ast.Schema, len(schemas))
	for i, schema := range schemas {
		newSchemas[i] = pass.processSchema(schema)
	}

	return newSchemas, nil
}

func (pass *GenericOptionToNamed) processSchema(schema *ast.Schema) *ast.Schema {
	objects := make([]ast.Object, len(schema.Objects))
	for i, obj := range schema.Objects {
		if !obj.Type.IsStruct() {
			objects[i] = obj
			continue
		}

		objects[i] = pass.processObject(schema.Package, obj)
	}

	schema.Objects = objects
	return schema
}

func (pass *GenericOptionToNamed) processObject(pkg string, object ast.Object) ast.Object {
	fmt.Printf("pkg: %s, name: %s, kind: %s\n", pkg, object.Name, object.Type.Kind)

	return object
}
