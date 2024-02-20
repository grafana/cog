package compiler

import (
	"github.com/grafana/cog/internal/ast"
)

var _ Pass = (*DashboardTargetsRewrite)(nil)

// DashboardTargetsRewrite rewrites any reference to the dashboard.Target type
// as a ComposableSlot<SchemaVariantDataQuery>.
//
// Example:
//
//	```
//	Panel struct {
//		Targets array(ref(dashboard.Target))
//	}
//	```
//
// Will become:
//
//	```
//	Panel struct {
//		Targets array(composableSlot(SchemaVariantDataQuery))
//	}
//	```
//
// Note: this compiler pass will only rewrite schemas in the "dashboard" package.
type DashboardTargetsRewrite struct {
}

func (pass *DashboardTargetsRewrite) Process(schemas []*ast.Schema) ([]*ast.Schema, error) {
	newSchemas := make([]*ast.Schema, 0, len(schemas))

	for _, schema := range schemas {
		if schema.Package != dashboardPackage {
			newSchemas = append(newSchemas, schema)
			continue
		}

		newSchemas = append(newSchemas, pass.processSchema(schema))
	}

	return newSchemas, nil
}

func (pass *DashboardTargetsRewrite) processSchema(schema *ast.Schema) *ast.Schema {
	schema.Objects = schema.Objects.Map(func(_ string, object ast.Object) ast.Object {
		object.Type = pass.processType(object.Type)

		return object
	})

	return schema
}

func (pass *DashboardTargetsRewrite) processType(def ast.Type) ast.Type {
	if def.Kind == ast.KindArray {
		return pass.processArray(def)
	}

	if def.Kind == ast.KindMap {
		return pass.processMap(def)
	}

	if def.Kind == ast.KindStruct {
		return pass.processStruct(def)
	}

	if def.Kind == ast.KindDisjunction {
		return pass.processDisjunction(def)
	}

	if def.Kind == ast.KindRef && def.AsRef().ReferredType == dashboardTargetObject {
		newDef := def
		newDef.Kind = ast.KindComposableSlot
		newDef.Ref = nil
		newDef.ComposableSlot = &ast.ComposableSlotType{
			Variant: ast.SchemaVariantDataQuery,
		}

		return newDef
	}

	return def
}

func (pass *DashboardTargetsRewrite) processArray(def ast.Type) ast.Type {
	processedType := pass.processType(def.AsArray().ValueType)

	newArray := def
	newArray.Array.ValueType = processedType

	return newArray
}

func (pass *DashboardTargetsRewrite) processMap(def ast.Type) ast.Type {
	processedValueType := pass.processType(def.AsMap().ValueType)

	newMap := def
	newMap.Map.ValueType = processedValueType

	return newMap
}

func (pass *DashboardTargetsRewrite) processStruct(def ast.Type) ast.Type {
	processedFields := make([]ast.StructField, 0, len(def.AsStruct().Fields))
	for _, field := range def.AsStruct().Fields {
		processedType := pass.processType(field.Type)

		newField := field
		newField.Type = processedType

		processedFields = append(processedFields, newField)
	}

	newStruct := def
	newStruct.Struct.Fields = processedFields

	return newStruct
}

func (pass *DashboardTargetsRewrite) processDisjunction(def ast.Type) ast.Type {
	disjunction := def
	newBranches := make([]ast.Type, 0, len(def.AsDisjunction().Branches))

	for _, branch := range def.AsDisjunction().Branches {
		newBranches = append(newBranches, pass.processType(branch))
	}

	disjunction.Disjunction.Branches = newBranches

	return disjunction
}
