package compiler

import (
	"fmt"

	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/tools"
)

var _ Pass = (*DisjunctionOfAnonymousStructsToExplicit)(nil)

// DisjunctionOfAnonymousStructsToExplicit looks for anonymous structs used as
// branches of disjunctions and turns them into explicitly named types.
type DisjunctionOfAnonymousStructsToExplicit struct {
	newObjects []ast.Object
}

func (pass *DisjunctionOfAnonymousStructsToExplicit) Process(schemas []*ast.Schema) ([]*ast.Schema, error) {
	for i, schema := range schemas {
		schemas[i] = pass.processSchema(schema)
	}

	return schemas, nil
}

func (pass *DisjunctionOfAnonymousStructsToExplicit) processSchema(schema *ast.Schema) *ast.Schema {
	pass.newObjects = nil

	schema.Objects = schema.Objects.Map(func(_ string, object ast.Object) ast.Object {
		return pass.processObject(schema, object)
	})

	schema.AddObjects(pass.newObjects...)

	return schema
}

func (pass *DisjunctionOfAnonymousStructsToExplicit) processObject(schema *ast.Schema, object ast.Object) ast.Object {
	object.Type = pass.processType(schema, object.Type)

	return object
}

func (pass *DisjunctionOfAnonymousStructsToExplicit) processType(schema *ast.Schema, def ast.Type) ast.Type {
	if def.IsArray() {
		return pass.processArray(schema, def)
	}

	if def.IsMap() {
		return pass.processMap(schema, def)
	}

	if def.IsStruct() {
		return pass.processStruct(schema, def)
	}

	if def.IsDisjunction() {
		return pass.processDisjunction(schema, def)
	}

	if def.IsIntersection() {
		return pass.processIntersection(schema, def)
	}

	return def
}

func (pass *DisjunctionOfAnonymousStructsToExplicit) processArray(schema *ast.Schema, def ast.Type) ast.Type {
	def.Array.ValueType = pass.processType(schema, def.AsArray().ValueType)

	return def
}

func (pass *DisjunctionOfAnonymousStructsToExplicit) processMap(schema *ast.Schema, def ast.Type) ast.Type {
	def.Map.ValueType = pass.processType(schema, def.AsMap().ValueType)

	return def
}

func (pass *DisjunctionOfAnonymousStructsToExplicit) processStruct(schema *ast.Schema, def ast.Type) ast.Type {
	for i, field := range def.Struct.Fields {
		def.Struct.Fields[i].Type = pass.processType(schema, field.Type)
	}

	return def
}

func (pass *DisjunctionOfAnonymousStructsToExplicit) processIntersection(schema *ast.Schema, def ast.Type) ast.Type {
	for i, branch := range def.Intersection.Branches {
		def.Intersection.Branches[i] = pass.processType(schema, branch)
	}

	return def
}

func (pass *DisjunctionOfAnonymousStructsToExplicit) processDisjunction(schema *ast.Schema, def ast.Type) ast.Type {
	scalarCount := 0
	anonymousCount := 0
	for _, branch := range def.Disjunction.Branches {
		if branch.IsScalar() {
			scalarCount++
		} else if branch.IsStruct() {
			anonymousCount++
		}
	}

	if scalarCount == 1 && anonymousCount == 1 {
		return def
	}

	for i, branch := range def.Disjunction.Branches {
		if !branch.IsStruct() {
			continue
		}

		branchName := pass.generateBranchName(branch, i)

		newObject := ast.NewObject(schema.Package, branchName, pass.processType(schema, branch))
		pass.newObjects = append(pass.newObjects, newObject)

		def.Disjunction.Branches[i] = ast.NewRef(schema.Package, newObject.Name)
	}

	return def
}

func (pass *DisjunctionOfAnonymousStructsToExplicit) generateBranchName(branch ast.Type, index int) string {
	for _, field := range branch.Struct.Fields {
		if field.Type.IsConcreteScalar() {
			val := fmt.Sprintf("%v", field.Type.Scalar.Value)
			return fmt.Sprintf("%s%s", field.Name, tools.UpperCamelCase(val))
		}
	}

	return fmt.Sprintf("branch%d", index)
}
