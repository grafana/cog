package compiler

import (
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/tools"
)

var _ Pass = (*AnonymousEnumToExplicitType)(nil)

// AnonymousEnumToExplicitType turns "anonymous enums" into a named
// object.
//
// Example:
//
//	```
//	Panel struct {
//		Type enum(Foo, Bar, Baz)
//	}
//	```
//
// Will become:
//
//	```
//	Panel struct {
//		Type PanelType
//	}
//
//	PanelType enum(Foo, Bar, Baz)
//	```
//
// Note: this compiler pass looks for anonymous enums in structs and arrays only.
type AnonymousEnumToExplicitType struct {
	newObjects     []ast.Object
	currentPackage string
}

func (pass *AnonymousEnumToExplicitType) Process(schemas []*ast.Schema) ([]*ast.Schema, error) {
	for i, schema := range schemas {
		schemas[i] = pass.processSchema(schema)
	}

	return schemas, nil
}

func (pass *AnonymousEnumToExplicitType) processSchema(schema *ast.Schema) *ast.Schema {
	pass.newObjects = nil
	pass.currentPackage = schema.Package

	for i, object := range schema.Objects {
		schema.Objects[i] = pass.processObject(object)
	}

	schema.Objects = append(schema.Objects, pass.newObjects...)

	return schema
}

func (pass *AnonymousEnumToExplicitType) processObject(object ast.Object) ast.Object {
	if object.Type.Kind == ast.KindEnum {
		return object
	}

	object.Type = pass.processType(object.SelfRef.ReferredPkg, object.Name, tools.UpperCamelCase(object.Name)+"Enum", object.Type)

	return object
}

func (pass *AnonymousEnumToExplicitType) processType(pkg string, currentObjectName string, suggestedEnumName string, def ast.Type) ast.Type {
	if def.Kind == ast.KindArray {
		return pass.processArray(pkg, currentObjectName, suggestedEnumName, def)
	}

	if def.Kind == ast.KindStruct {
		return pass.processStruct(pkg, currentObjectName, def)
	}

	if def.Kind == ast.KindEnum {
		return pass.processAnonymousEnum(pkg, suggestedEnumName, def.AsEnum())
	}

	// TODO: do the same for disjunctions?

	return def
}

func (pass *AnonymousEnumToExplicitType) processArray(pkg string, currentObjectName string, suggestedEnumName string, def ast.Type) ast.Type {
	def.Array.ValueType = pass.processType(pkg, currentObjectName, suggestedEnumName, def.Array.ValueType)

	return def
}

func (pass *AnonymousEnumToExplicitType) processStruct(pkg string, parentName string, def ast.Type) ast.Type {
	for i, field := range def.Struct.Fields {
		def.Struct.Fields[i].Type = pass.processType(pkg, parentName, tools.UpperCamelCase(parentName)+tools.UpperCamelCase(field.Name), field.Type)
	}

	return def
}

func (pass *AnonymousEnumToExplicitType) processAnonymousEnum(pkg string, parentName string, def ast.EnumType) ast.Type {
	enumTypeName := tools.UpperCamelCase(parentName)

	values := make([]ast.EnumValue, 0, len(def.Values))
	for _, val := range def.Values {
		values = append(values, ast.EnumValue{
			Type:  val.Type,
			Name:  tools.UpperCamelCase(val.Name),
			Value: val.Value,
		})
	}

	enumType := ast.NewEnum(values)
	enumType.AddCompilerPassTrail("AnonymousEnumToExplicitType")

	pass.newObjects = append(pass.newObjects, ast.NewObject(pkg, enumTypeName, enumType))

	return ast.NewRef(pass.currentPackage, enumTypeName)
}
