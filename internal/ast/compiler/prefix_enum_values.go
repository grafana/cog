package compiler

import (
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/tools"
)

var _ Pass = (*PrefixEnumValues)(nil)

type PrefixEnumValues struct {
}

func (pass *PrefixEnumValues) Process(files []*ast.File) ([]*ast.File, error) {
	newFiles := make([]*ast.File, 0, len(files))

	for _, file := range files {
		newFile, err := pass.processFile(file)
		if err != nil {
			return nil, err
		}

		newFiles = append(newFiles, newFile)
	}

	return newFiles, nil
}

func (pass *PrefixEnumValues) processFile(file *ast.File) (*ast.File, error) {
	processedObjects := make([]ast.Object, 0, len(file.Definitions))
	for _, object := range file.Definitions {
		newObject := object
		newObject.Type = pass.processType(object.Name, object.Type)

		processedObjects = append(processedObjects, newObject)
	}

	return &ast.File{
		Package:     file.Package,
		Definitions: processedObjects,
	}, nil
}

func (pass *PrefixEnumValues) processType(parentObjectName string, def ast.Type) ast.Type {
	if def.Kind != ast.KindEnum {
		return def
	}

	return pass.processEnum(parentObjectName, def)
}

func (pass *PrefixEnumValues) processEnum(parentName string, def ast.Type) ast.Type {
	newType := def

	values := make([]ast.EnumValue, 0, len(def.AsEnum().Values))
	for _, val := range def.AsEnum().Values {
		values = append(values, ast.EnumValue{
			Type:  val.Type,
			Name:  tools.UpperCamelCase(parentName) + tools.UpperCamelCase(val.Name),
			Value: val.Value,
		})
	}

	newType.Enum.Values = values

	return newType
}
