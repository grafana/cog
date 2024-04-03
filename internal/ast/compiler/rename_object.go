package compiler

import (
	"fmt"

	"github.com/grafana/cog/internal/ast"
)

var _ Pass = (*RenameObject)(nil)

type RenameObject struct {
	From ObjectReference
	To   string
}

func (pass *RenameObject) Process(schemas []*ast.Schema) ([]*ast.Schema, error) {
	for i, schema := range schemas {
		schemas[i] = pass.processSchema(schema)
	}

	return schemas, nil
}

func (pass *RenameObject) processSchema(schema *ast.Schema) *ast.Schema {
	var renamedObject ast.Object

	schema.Objects = schema.Objects.Map(func(_ string, object ast.Object) ast.Object {
		if pass.From.Matches(object) {
			object.AddToPassesTrail(fmt.Sprintf("RenameObject[%s → %s]", object.Name, pass.To))

			object.Name = pass.To
			object.SelfRef.ReferredType = pass.To

			renamedObject = object
		}

		object.Type = pass.processType(object.Type)

		return object
	})

	if renamedObject.Name != "" {
		schema.Objects.Remove(pass.From.Object)
		schema.AddObject(renamedObject)
	}

	return schema
}

func (pass *RenameObject) processType(def ast.Type) ast.Type {
	if def.IsArray() {
		return pass.processArray(def)
	}

	if def.IsMap() {
		return pass.processMap(def)
	}

	if def.IsStruct() {
		return pass.processStruct(def)
	}

	if def.IsDisjunction() {
		return pass.processDisjunction(def)
	}

	if def.IsIntersection() {
		return pass.processIntersection(def)
	}

	if def.IsRef() {
		return pass.processRef(def)
	}

	return def
}

func (pass *RenameObject) processArray(def ast.Type) ast.Type {
	def.Array.ValueType = pass.processType(def.Array.ValueType)

	return def
}

func (pass *RenameObject) processMap(def ast.Type) ast.Type {
	def.Map.IndexType = pass.processType(def.Map.IndexType)
	def.Map.ValueType = pass.processType(def.Map.ValueType)

	return def
}

func (pass *RenameObject) processDisjunction(def ast.Type) ast.Type {
	for i, branch := range def.Disjunction.Branches {
		def.Disjunction.Branches[i] = pass.processType(branch)
	}

	return def
}

func (pass *RenameObject) processIntersection(def ast.Type) ast.Type {
	for i, branch := range def.Intersection.Branches {
		def.Intersection.Branches[i] = pass.processType(branch)
	}

	return def
}

func (pass *RenameObject) processRef(def ast.Type) ast.Type {
	if def.Ref.ReferredType == pass.From.Object {
		def.Ref.ReferredType = pass.To
	}

	return def
}

func (pass *RenameObject) processStruct(def ast.Type) ast.Type {
	for i, field := range def.Struct.Fields {
		def.Struct.Fields[i].Type = pass.processType(field.Type)
	}

	return def
}
