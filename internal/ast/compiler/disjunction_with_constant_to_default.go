package compiler

import (
	"github.com/grafana/cog/internal/ast"
)

var _ Pass = (*DisjunctionWithConstantToDefault)(nil)

type DisjunctionWithConstantToDefault struct {
}

func (pass *DisjunctionWithConstantToDefault) Process(schemas []*ast.Schema) ([]*ast.Schema, error) {
	for i, schema := range schemas {
		schemas[i] = pass.processSchema(schema)
	}

	return schemas, nil
}

func (pass *DisjunctionWithConstantToDefault) processSchema(schema *ast.Schema) *ast.Schema {
	schema.Objects = schema.Objects.Map(func(_ string, object ast.Object) ast.Object {
		return pass.processObject(object)
	})

	return schema
}

func (pass *DisjunctionWithConstantToDefault) processObject(object ast.Object) ast.Object {
	object.Type = pass.processType(object.Type)

	return object
}

func (pass *DisjunctionWithConstantToDefault) processType(def ast.Type) ast.Type {
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

	return def
}

func (pass *DisjunctionWithConstantToDefault) processArray(def ast.Type) ast.Type {
	def.Array.ValueType = pass.processType(def.AsArray().ValueType)

	return def
}

func (pass *DisjunctionWithConstantToDefault) processMap(def ast.Type) ast.Type {
	def.Map.ValueType = pass.processType(def.AsMap().ValueType)

	return def
}

func (pass *DisjunctionWithConstantToDefault) processStruct(def ast.Type) ast.Type {
	for i, field := range def.Struct.Fields {
		def.Struct.Fields[i].Type = pass.processType(field.Type)
	}

	return def
}

func (pass *DisjunctionWithConstantToDefault) processDisjunction(def ast.Type) ast.Type {
	branches := def.Disjunction.Branches

	if len(branches) != 2 {
		return def
	}

	if branches[0].Kind != branches[1].Kind {
		return def
	}

	if !branches[0].IsScalar() {
		return def
	}

	if branches[0].Scalar.ScalarKind != branches[1].Scalar.ScalarKind {
		return def
	}

	if branches[0].Scalar.IsConcrete() == branches[1].Scalar.IsConcrete() {
		return def
	}

	if branches[0].Scalar.IsConcrete() {
		def = branches[1]
		def.Default = branches[0].Scalar.Value
	} else {
		def = branches[0]
		def.Default = branches[1].Scalar.Value
	}

	def.AddToPassesTrail("DisjunctionWithConstantToDefault")

	return def
}

func (pass *DisjunctionWithConstantToDefault) processIntersection(def ast.Type) ast.Type {
	for i, branch := range def.Intersection.Branches {
		def.Intersection.Branches[i] = pass.processType(branch)
	}

	return def
}
