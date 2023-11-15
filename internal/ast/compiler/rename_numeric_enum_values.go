package compiler

import (
	"strconv"

	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/tools"
)

var _ Pass = (*RenameNumericEnumValues)(nil)

type RenameNumericEnumValues struct {
}

func (pass *RenameNumericEnumValues) Process(schemas []*ast.Schema) ([]*ast.Schema, error) {
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

func (pass *RenameNumericEnumValues) processSchema(schema *ast.Schema) (*ast.Schema, error) {
	processedObjects := make([]ast.Object, 0, len(schema.Objects))
	for _, object := range schema.Objects {
		newObject := object
		newObject.Type = pass.processType(object.Type)

		processedObjects = append(processedObjects, newObject)
	}

	newSchema := schema.DeepCopy()
	newSchema.Objects = processedObjects

	return &newSchema, nil
}

func (pass *RenameNumericEnumValues) processType(def ast.Type) ast.Type {
	if def.Kind != ast.KindEnum {
		return def
	}

	return pass.processEnum(def)
}

func (pass *RenameNumericEnumValues) processEnum(def ast.Type) ast.Type {
	newType := def

	values := make([]ast.EnumValue, 0, len(def.AsEnum().Values))
	for _, val := range def.AsEnum().Values {
		if _, err := strconv.Atoi(val.Name); err != nil {
			values = append(values, val)
			continue
		}

		values = append(values, ast.EnumValue{
			Type:  val.Type,
			Name:  "N" + tools.UpperCamelCase(val.Name),
			Value: val.Value,
		})
	}

	newType.Enum.Values = values

	return newType
}
