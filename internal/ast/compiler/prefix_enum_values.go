package compiler

import (
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/tools"
)

var _ Pass = (*PrefixEnumValues)(nil)

// PrefixEnumValues prefixes enum members with the name of the enum object in
// which they are defined.
//
// Example:
//
//	```
//	VariableRefresh enum(Never: "never", Always: "always")
//	```
//
// Will become:
//
//	```
//	VariableRefresh enum(VariableRefreshNever: "never", VariableRefreshAlways: "always")
//	```
type PrefixEnumValues struct {
}

func (pass *PrefixEnumValues) Process(schemas []*ast.Schema) ([]*ast.Schema, error) {
	for i, schema := range schemas {
		schemas[i] = pass.processSchema(schema)
	}

	return schemas, nil
}

func (pass *PrefixEnumValues) processSchema(schema *ast.Schema) *ast.Schema {
	for i, object := range schema.Objects {
		schema.Objects[i].Type = pass.processType(object.Name, object.Type)
	}

	return schema
}

func (pass *PrefixEnumValues) processType(parentObjectName string, def ast.Type) ast.Type {
	if def.Kind != ast.KindEnum {
		return def
	}

	return pass.processEnum(parentObjectName, def)
}

func (pass *PrefixEnumValues) processEnum(parentName string, def ast.Type) ast.Type {
	values := make([]ast.EnumValue, 0, len(def.AsEnum().Values))
	for _, val := range def.AsEnum().Values {
		values = append(values, ast.EnumValue{
			Type:  val.Type,
			Name:  tools.UpperCamelCase(parentName) + tools.UpperCamelCase(val.Name),
			Value: val.Value,
		})
	}

	def.Enum.Values = values
	def.AddCompilerPassTrail("PrefixEnumValues")

	return def
}
