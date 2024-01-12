package compiler

import (
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/tools"
)

var _ Pass = (*AnonymousStructsToNamed)(nil)

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
	newObjects []ast.Object
}

func (pass *AnonymousStructsToNamed) Process(schemas []*ast.Schema) ([]*ast.Schema, error) {
	newSchemas := make([]*ast.Schema, 0, len(schemas))
	for _, schema := range schemas {
		newSchemas = append(newSchemas, pass.processSchema(schema))
	}

	return newSchemas, nil
}

func (pass *AnonymousStructsToNamed) processSchema(schema *ast.Schema) *ast.Schema {
	pass.newObjects = make([]ast.Object, 0, len(schema.Objects))
	for _, object := range schema.Objects {
		newObj := pass.processObject(object)
		pass.newObjects = append(pass.newObjects, newObj)
	}

	schema.Objects = pass.newObjects

	return schema
}

func (pass *AnonymousStructsToNamed) processObject(object ast.Object) ast.Object {
	newObject := object
	pkg := object.SelfRef.ReferredPkg
	parentName := tools.UpperCamelCase(pkg) + tools.UpperCamelCase(object.Name)

	if object.Type.IsAnyOf(ast.KindArray, ast.KindMap, ast.KindDisjunction) {
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

func (pass *AnonymousStructsToNamed) processType(pkg string, parentName string, def ast.Type) ast.Type {
	if def.Kind == ast.KindArray {
		return pass.processArray(pkg, parentName, def)
	}

	if def.Kind == ast.KindMap {
		return pass.processMap(pkg, parentName, def)
	}

	if def.Kind == ast.KindDisjunction {
		return pass.processDisjunction(pkg, parentName, def)
	}

	if def.Kind == ast.KindStruct {
		return pass.processStruct(pkg, parentName, def)
	}

	return def
}

func (pass *AnonymousStructsToNamed) processArray(pkg string, parentName string, def ast.Type) ast.Type {
	def.Array.ValueType = pass.processType(pkg, parentName, def.Array.ValueType)

	return def
}

func (pass *AnonymousStructsToNamed) processMap(pkg string, parentName string, def ast.Type) ast.Type {
	def.Map.IndexType = pass.processType(pkg, parentName, def.Map.IndexType)
	def.Map.ValueType = pass.processType(pkg, parentName, def.Map.ValueType)

	return def
}

func (pass *AnonymousStructsToNamed) processDisjunction(pkg string, parentName string, def ast.Type) ast.Type {
	for i, branch := range def.Disjunction.Branches {
		def.Disjunction.Branches[i] = pass.processType(pkg, parentName, branch)
	}

	return def
}

func (pass *AnonymousStructsToNamed) processStruct(pkg string, parentName string, def ast.Type) ast.Type {
	for i, field := range def.AsStruct().Fields {
		name := parentName + tools.UpperCamelCase(field.Name)
		def.Struct.Fields[i].Type = pass.processType(pkg, name, field.Type)
	}

	pass.newObjects = append(pass.newObjects, ast.NewObject(pkg, parentName, def))

	ref := ast.NewRef(pkg, parentName)
	ref.Nullable = def.Nullable
	ref.Default = def.Default

	return ref
}
