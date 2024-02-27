package compiler

import (
	"errors"
	"fmt"
	"strings"

	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/orderedmap"
)

var _ Pass = (*DisjunctionToType)(nil)

// DisjunctionToType transforms disjunction into a struct, mapping disjunction branches to
// an optional and nullable field in that struct.
//
// Example:
//
//		```
//		SomeType: {
//			type: "some-type"
//	 	}
//		SomeOtherType: {
//			type: "other-type"
//	 	}
//		SomeStruct: {
//			foo: string | bool
//		}
//		OtherStruct: {
//			bar: SomeType | SomeOtherType
//		}
//		```
//
// Will become:
//
//		```
//		SomeType: {
//			type: "some-type"
//	 	}
//		SomeOtherType: {
//			type: "other-type"
//	 	}
//		StringOrBool: {
//			string: *string
//			bool: *string
//		}
//		SomeStruct: {
//			foo: StringOrBool
//		}
//		SomeTypeOrSomeOtherType: {
//			SomeType: *SomeType
//			SomeOtherType: *SomeOtherType
//		}
//		OtherStruct: {
//			bar: SomeTypeOrSomeOtherType
//		}
//		```
type DisjunctionToType struct {
	newObjects *orderedmap.Map[string, ast.Object]
}

func (pass *DisjunctionToType) Process(schemas []*ast.Schema) ([]*ast.Schema, error) {
	newSchemas := make([]*ast.Schema, 0, len(schemas))

	for _, schema := range schemas {
		newSchema, err := pass.processSchema(schema)
		if err != nil {
			return nil, fmt.Errorf("[%s] %w", schema.Package, err)
		}

		newSchemas = append(newSchemas, newSchema)
	}

	return newSchemas, nil
}

func (pass *DisjunctionToType) processSchema(schema *ast.Schema) (*ast.Schema, error) {
	var err error
	pass.newObjects = orderedmap.New[string, ast.Object]()

	schema.Objects = schema.Objects.Map(func(_ string, object ast.Object) ast.Object {
		processedObject, innerErr := pass.processObject(schema, object)
		if innerErr != nil {
			err = innerErr
			return object
		}

		return processedObject
	})
	if err != nil {
		return nil, err
	}

	pass.newObjects.Iterate(func(_ string, object ast.Object) {
		schema.AddObject(object)
	})

	return schema, nil
}

func (pass *DisjunctionToType) processObject(schema *ast.Schema, object ast.Object) (ast.Object, error) {
	processedType, err := pass.processType(schema, object.Type)
	if err != nil {
		return object, errors.Join(
			fmt.Errorf("could not process object '%s'", object.Name),
			err,
		)
	}

	newObject := object
	newObject.Type = processedType

	return newObject, nil
}

func (pass *DisjunctionToType) processType(schema *ast.Schema, def ast.Type) (ast.Type, error) {
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

	return def, nil
}

func (pass *DisjunctionToType) processArray(schema *ast.Schema, def ast.Type) (ast.Type, error) {
	processedType, err := pass.processType(schema, def.AsArray().ValueType)
	if err != nil {
		return ast.Type{}, err
	}

	newArray := def
	newArray.Array.ValueType = processedType

	return newArray, nil
}

func (pass *DisjunctionToType) processMap(schema *ast.Schema, def ast.Type) (ast.Type, error) {
	processedValueType, err := pass.processType(schema, def.AsMap().ValueType)
	if err != nil {
		return ast.Type{}, err
	}

	newMap := def
	newMap.Map.ValueType = processedValueType

	return newMap, nil
}

func (pass *DisjunctionToType) processStruct(schema *ast.Schema, def ast.Type) (ast.Type, error) {
	processedFields := make([]ast.StructField, 0, len(def.AsStruct().Fields))
	for _, field := range def.AsStruct().Fields {
		processedType, err := pass.processType(schema, field.Type)
		if err != nil {
			return ast.Type{}, errors.Join(
				fmt.Errorf("could not process struct field '%s'", field.Name),
				err,
			)
		}

		newField := field
		newField.Type = processedType

		processedFields = append(processedFields, newField)
	}

	newStruct := def
	newStruct.Struct.Fields = processedFields

	return newStruct, nil
}

func (pass *DisjunctionToType) processDisjunction(schema *ast.Schema, def ast.Type) (ast.Type, error) {
	disjunction := def.AsDisjunction()

	// Ex: "some concrete value" | "some other value" | string
	if pass.hasOnlySingleTypeScalars(schema, disjunction) {
		resolvedType, _ := schema.Resolve(disjunction.Branches[0])
		scalarKind := resolvedType.AsScalar().ScalarKind

		return ast.NewScalar(scalarKind, ast.Default(def.Default)), nil
	}

	// type | otherType | something (| null)?
	// generate a type with a nullable field for every branch of the disjunction,
	// add it to preprocessor.types, and use it instead.
	newTypeName := pass.disjunctionTypeName(disjunction)

	// if we already generated a new object for this disjunction, let's return
	// a reference to it.
	if pass.newObjects.Has(newTypeName) {
		ref := ast.NewRef(schema.Package, newTypeName, ast.Hints(def.Hints))
		if def.Nullable || disjunction.Branches.HasNullType() {
			ref.Nullable = true
		}

		return ref, nil
	}

	/*
		TODO: return an error here. Some jennies won't be able to handle
		this type of disjunction.
		if !disjunction.Branches.HasOnlyScalarOrArray() || !disjunction.Branches.HasOnlyRefs() {
		}
	*/

	fields := make([]ast.StructField, 0, len(disjunction.Branches))
	for _, branch := range disjunction.Branches {
		// Handled below, by allowing the reference to the disjunction struct
		// to be null.
		if branch.IsNull() {
			continue
		}

		processedBranch := branch
		processedBranch.Nullable = true

		fields = append(fields, ast.StructField{
			Name:     ast.TypeName(processedBranch),
			Type:     processedBranch,
			Required: false,
		})
	}

	structType := ast.NewStruct(fields...)
	if disjunction.Branches.HasOnlyScalarOrArray() {
		structType.Hints[ast.HintDisjunctionOfScalars] = disjunction
	}
	if disjunction.Branches.HasOnlyRefs() {
		if len(disjunction.Discriminator) == 0 {
			return ast.Type{}, fmt.Errorf("discriminator not set")
		}
		if len(disjunction.DiscriminatorMapping) == 0 {
			return ast.Type{}, fmt.Errorf("discriminator mapping not set")
		}
		structType.Hints[ast.HintDiscriminatedDisjunctionOfRefs] = disjunction
	}

	pass.newObjects.Set(newTypeName, ast.Object{
		Name: newTypeName,
		Type: structType,
		SelfRef: ast.RefType{
			ReferredPkg:  schema.Package,
			ReferredType: newTypeName,
		},
	})

	ref := ast.NewRef(schema.Package, newTypeName, ast.Hints(def.Hints))
	if def.Nullable || disjunction.Branches.HasNullType() {
		ref.Nullable = true
	}

	return ref, nil
}

func (pass *DisjunctionToType) disjunctionTypeName(def ast.DisjunctionType) string {
	parts := make([]string, 0, len(def.Branches))

	for _, subType := range def.Branches {
		parts = append(parts, ast.TypeName(subType))
	}

	return strings.Join(parts, "Or")
}

func (pass *DisjunctionToType) hasOnlySingleTypeScalars(schema *ast.Schema, disjunction ast.DisjunctionType) bool {
	branches := disjunction.Branches

	if len(branches) == 0 {
		return false
	}

	firstBranchType, found := schema.Resolve(branches[0])
	if !found {
		return false
	}

	if firstBranchType.Kind != ast.KindScalar {
		return false
	}

	scalarKind := firstBranchType.AsScalar().ScalarKind
	for _, t := range branches {
		resolvedType, found := schema.Resolve(t)
		if !found {
			return false
		}

		if resolvedType.Kind != ast.KindScalar {
			return false
		}

		if resolvedType.AsScalar().ScalarKind != scalarKind {
			return false
		}
	}

	return true
}
