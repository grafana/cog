package compiler

import (
	"github.com/grafana/cog/internal/ast"
)

var _ Pass = (*FlattenDisjunctions)(nil)

// FlattenDisjunctions will traverse all the branches every given disjunctions
// and, for each disjunction it finds, flatten it into the top-level type.
//
// Example:
//
//	```
//	SomeStruct: {
//		foo: string
//	}
//	OtherStruct: {
//		bar: string
//	}
//	LastStruct: {
//		hello: string
//	}
//	SomeOrOther: SomeStruct | OtherStruct
//	AnyStruct: SomeOrOther | LastStruct
//	```
//
// Will become:
//
//	```
//	SomeStruct: {
//		foo: string
//	}
//	OtherStruct: {
//		bar: string
//	}
//	LastStruct: {
//		hello: string
//	}
//	SomeOrOther: SomeStruct | OtherStruct
//	AnyStruct: SomeStruct | OtherStruct | LastStruct # this disjunction has been flattened
//	```
type FlattenDisjunctions struct {
}

func (pass *FlattenDisjunctions) Process(schemas []*ast.Schema) ([]*ast.Schema, error) {
	for i, schema := range schemas {
		schemas[i] = pass.processSchema(schema)
	}

	return schemas, nil
}

func (pass *FlattenDisjunctions) processSchema(schema *ast.Schema) *ast.Schema {
	for i, object := range schema.Objects {
		schema.Objects[i] = pass.processObject(schema, object)
	}

	return schema
}

func (pass *FlattenDisjunctions) processObject(schema *ast.Schema, object ast.Object) ast.Object {
	object.Type = pass.processType(schema, object.Type)

	return object
}

func (pass *FlattenDisjunctions) processType(schema *ast.Schema, def ast.Type) ast.Type {
	if def.Kind == ast.KindArray {
		return pass.processArray(schema, def)
	}

	if def.Kind == ast.KindMap {
		return pass.processMap(schema, def)
	}

	if def.Kind == ast.KindStruct {
		return pass.processStruct(schema, def)
	}

	if def.Kind == ast.KindDisjunction {
		return pass.processDisjunction(schema, def)
	}

	return def
}

func (pass *FlattenDisjunctions) processArray(schema *ast.Schema, def ast.Type) ast.Type {
	def.Array.ValueType = pass.processType(schema, def.AsArray().ValueType)

	return def
}

func (pass *FlattenDisjunctions) processMap(schema *ast.Schema, def ast.Type) ast.Type {
	def.Map.ValueType = pass.processType(schema, def.AsMap().ValueType)

	return def
}

func (pass *FlattenDisjunctions) processStruct(schema *ast.Schema, def ast.Type) ast.Type {
	for i, field := range def.AsStruct().Fields {
		def.Struct.Fields[i].Type = pass.processType(schema, field.Type)
	}

	return def
}

func (pass *FlattenDisjunctions) processDisjunction(schema *ast.Schema, def ast.Type) ast.Type {
	def.Disjunction = pass.flattenDisjunction(schema, def.AsDisjunction())

	return def
}

func (pass *FlattenDisjunctions) flattenDisjunction(schema *ast.Schema, disjunction ast.DisjunctionType) *ast.DisjunctionType {
	newDisjunction := disjunction.DeepCopy()
	newDisjunction.Branches = nil

	refResolver := pass.makeReferenceResolver(schema)

	branchMap := make(map[string]struct{})
	addBranch := func(typeDef ast.Type) {
		typeName := ast.TypeName(typeDef)
		if _, exists := branchMap[typeName]; exists {
			return
		}

		branchMap[typeName] = struct{}{}
		newDisjunction.Branches = append(newDisjunction.Branches, typeDef)
	}

	for _, branch := range disjunction.Branches {
		if branch.Kind != ast.KindRef {
			addBranch(branch)
			continue
		}

		resolved := refResolver(branch)
		if resolved.Kind != ast.KindDisjunction {
			addBranch(branch)
			continue
		}

		for _, resolvedBranch := range resolved.AsDisjunction().Branches {
			addBranch(resolvedBranch)
		}
	}

	return &newDisjunction
}

func (pass *FlattenDisjunctions) makeReferenceResolver(schema *ast.Schema) func(typeDef ast.Type) ast.Type {
	return func(typeDef ast.Type) ast.Type {
		if typeDef.Kind != ast.KindRef {
			return typeDef
		}

		typeName := typeDef.AsRef().ReferredType

		// FIXME: what if the definition is itself a reference? Resolve recursively?
		// FIXME: we only try to resolve references within the same schema
		return schema.LocateObject(typeName).Type
	}
}
