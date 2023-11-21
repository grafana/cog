package compiler

import (
	"strconv"

	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/tools"
)

var _ Pass = (*RenameNumericEnumValues)(nil)

// RenameNumericEnumValues turns any numeric enum member name to an alphanumeric name.
//
// Example:
//
//	```
//	Position enum(0: 0, 1: 1, 2: 2)
//	```
//
// Will become:
//
//	```
//	Position enum(N0: 0, N1: 1, N2: 2)
//	```
type RenameNumericEnumValues struct {
}

func (pass *RenameNumericEnumValues) Process(schemas []*ast.Schema) ([]*ast.Schema, error) {
	for i, schema := range schemas {
		schemas[i] = pass.processSchema(schema)
	}

	return schemas, nil
}

func (pass *RenameNumericEnumValues) processSchema(schema *ast.Schema) *ast.Schema {
	for i, object := range schema.Objects {
		schema.Objects[i].Type = pass.processType(object.Type)
	}

	return schema
}

func (pass *RenameNumericEnumValues) processType(def ast.Type) ast.Type {
	if def.Kind != ast.KindEnum {
		return def
	}

	return pass.processEnum(def)
}

func (pass *RenameNumericEnumValues) processEnum(def ast.Type) ast.Type {
	modified := false
	for i, val := range def.AsEnum().Values {
		if _, err := strconv.Atoi(val.Name); err != nil {
			continue
		}

		modified = true
		def.AsEnum().Values[i].Name = "N" + tools.UpperCamelCase(val.Name)
	}

	if modified {
		def.AddCompilerPassTrail("RenameNumericEnumValues")
	}

	return def
}
