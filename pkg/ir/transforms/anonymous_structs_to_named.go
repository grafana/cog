package transforms

import (
	"github.com/grafana/cog/internal/tools"
	"github.com/grafana/cog/pkg/ir"
)

var _ Transform = (*AnonymousStructsToNamed)(nil)

// AnonymousStructsToNamed turns "anonymous structs" into a named object.
//
// Example:
//
//	```
//	Panel struct {
//		Options struct {
//			Title string
//		}
//	}
//	```
//
// Will become:
//
//	```
//	Panel struct {
//		Options PanelOptions
//	}
//
//	PanelOptions struct {
//		Title string
//	}
//	```
type AnonymousStructsToNamed struct {
	newObjects []ir.Object
}

func (pass *AnonymousStructsToNamed) Process(schemas ir.Schemas) (ir.Schemas, error) {
	newSchemas := make(ir.Schemas, 0, len(schemas))
	for _, schema := range schemas {
		newSchemas = append(newSchemas, pass.processSchema(schema))
	}

	return newSchemas, nil
}

func (pass *AnonymousStructsToNamed) processSchema(schema *ir.Schema) *ir.Schema {
	pass.newObjects = nil

	schema.Objects = schema.Objects.Map(func(_ string, object ir.Object) ir.Object {
		return pass.processObject(object)
	})
	schema.AddObjects(pass.newObjects...)

	return schema
}

func (pass *AnonymousStructsToNamed) processObject(object ir.Object) ir.Object {
	newObject := object
	pkg := object.SelfRef.ReferredPkg
	parentName := tools.UpperCamelCase(pkg) + tools.UpperCamelCase(object.Name)

	if object.Type.IsAnyOf(ir.KindArray, ir.KindMap, ir.KindDisjunction) {
		newObject.Type = pass.processType(pkg, parentName, object.Type)
	}

	if object.Type.IsStruct() {
		for i, field := range object.Type.AsStruct().Fields {
			name := parentName + tools.UpperCamelCase(field.Name)
			object.Type.Struct.Fields[i].Type = pass.processType(pkg, name, field.Type)
		}
	}

	return newObject
}

func (pass *AnonymousStructsToNamed) processType(pkg string, parentName string, def ir.Type) ir.Type {
	if def.IsArray() {
		return pass.processArray(pkg, parentName, def)
	}

	if def.IsMap() {
		return pass.processMap(pkg, parentName, def)
	}

	if def.IsDisjunction() {
		return pass.processDisjunction(pkg, parentName, def)
	}

	if def.IsStruct() {
		return pass.processStruct(pkg, parentName, def)
	}

	return def
}

func (pass *AnonymousStructsToNamed) processArray(pkg string, parentName string, def ir.Type) ir.Type {
	def.Array.ValueType = pass.processType(pkg, parentName, def.Array.ValueType)

	return def
}

func (pass *AnonymousStructsToNamed) processMap(pkg string, parentName string, def ir.Type) ir.Type {
	def.Map.IndexType = pass.processType(pkg, parentName, def.Map.IndexType)
	def.Map.ValueType = pass.processType(pkg, parentName, def.Map.ValueType)

	return def
}

func (pass *AnonymousStructsToNamed) processDisjunction(pkg string, parentName string, def ir.Type) ir.Type {
	for i, branch := range def.Disjunction.Branches {
		def.Disjunction.Branches[i] = pass.processType(pkg, parentName, branch)
	}

	return def
}

func (pass *AnonymousStructsToNamed) processStruct(pkg string, parentName string, def ir.Type) ir.Type {
	objectDef := def.DeepCopy()
	objectDef.Nullable = false

	for i, field := range def.AsStruct().Fields {
		name := parentName + tools.UpperCamelCase(field.Name)
		objectDef.Struct.Fields[i].Type = pass.processType(pkg, name, field.Type)
	}

	newObject := ir.NewObject(pkg, parentName, objectDef)
	newObject.AddToPassesTrail("AnonymousStructsToNamed")

	pass.newObjects = append(pass.newObjects, newObject)

	ref := ir.NewRef(pkg, parentName)
	ref.Nullable = def.Nullable
	ref.Default = def.Default

	return ref
}
