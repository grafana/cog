package transforms

import (
	"fmt"

	"github.com/grafana/cog/internal/ir"
)

var _ Transform = (*NameAnonymousStruct)(nil)

// NameAnonymousStruct rewrites the definition of a struct field typed as an
// anonymous struct to instead refer to a named type.
//
// Example:
//
//	```
//	Field = Panel.DataSource
//	As = DataSourceRef
//	```
//
//	```
//	Panel: {
//		DataSource: {
//			Type: string
//			name: string
//		}
//	}
//	```
//
// Will become:
//
//	```
//	Panel: {
//		DataSource: DataSourceRef
//	}
//
//	DataSourceRef: {
//		Type: string
//		name: string
//	}
//	```
type NameAnonymousStruct struct {
	Field      FieldReference
	As         string
	fieldFound bool
}

func (pass *NameAnonymousStruct) Process(schemas ir.Schemas) (ir.Schemas, error) {
	pass.fieldFound = false

	for i, schema := range schemas {
		schemas[i] = pass.processSchema(schema)
	}

	return schemas, nil
}

func (pass *NameAnonymousStruct) processSchema(schema *ir.Schema) *ir.Schema {
	var newObject ir.Object

	schema.Objects = schema.Objects.Map(func(_ string, object ir.Object) ir.Object {
		currentObject, newObjectCandidate := pass.processObject(object)
		if newObjectCandidate.Name != "" {
			newObject = newObjectCandidate
		}

		return currentObject
	})

	// did we actually define a new object?
	if newObject.Name != "" {
		schema.AddObject(newObject)
	}

	return schema
}

func (pass *NameAnonymousStruct) processObject(object ir.Object) (ir.Object, ir.Object) {
	var newObject ir.Object

	if !object.Type.IsStruct() {
		return object, newObject
	}

	pkg := object.SelfRef.ReferredPkg

	for i, field := range object.Type.AsStruct().Fields {
		if !pass.Field.Matches(object, field) {
			continue
		}

		// we expect the target field to be defined as an inline struct
		if !field.Type.IsStruct() {
			continue
		}

		pass.fieldFound = true

		newObject = ir.NewObject(pkg, pass.As, field.Type)
		newObject.AddToPassesTrail("NameAnonymousStruct")

		object.Type.AsStruct().Fields[i].Type = ir.NewRef(pkg, pass.As)
	}

	return object, newObject
}

func (pass *NameAnonymousStruct) Diagnostics() []string {
	if pass.fieldFound {
		return nil
	}

	return []string{
		fmt.Sprintf("field '%s' not found", pass.Field),
	}
}
