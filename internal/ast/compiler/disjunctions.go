package compiler

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/grafana/cog/internal/ast"
)

var _ Pass = (*DisjunctionToType)(nil)

var ErrCanNotInferDiscriminator = errors.New("can not infer discriminator mapping")

// DisjunctionToType transforms disjunction into a struct, mapping disjunction branches to
// an optional and nullable field in that struct.
//
// This compiler pass also simplifies disjunction of two branches, where one is `null`. For those,
// it transforms `type | null` into `*type` (optional, nullable reference to `type`).
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
//	```
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
//	```
//
// Note: for disjunctions of `Ref`s, the pass attempts to infer a discriminator field and mapping. See https://swagger.io/docs/specification/data-models/inheritance-and-polymorphism/
type DisjunctionToType struct {
	newObjects map[string]ast.Object
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
	pass.newObjects = make(map[string]ast.Object)

	processedObjects := make([]ast.Object, 0, len(schema.Objects))
	for _, object := range schema.Objects {
		processedObject, err := pass.processObject(schema, object)
		if err != nil {
			return nil, err
		}

		processedObjects = append(processedObjects, processedObject)
	}

	newObjects := make([]ast.Object, 0, len(pass.newObjects))
	for _, obj := range pass.newObjects {
		newObjects = append(newObjects, obj)
	}

	// Since newly created objects are temporarily stored in a map, we need to
	// sort them to have a deterministic output.
	sort.SliceStable(newObjects, func(i, j int) bool {
		return newObjects[i].Name < newObjects[j].Name
	})

	newSchema := schema.DeepCopy()
	newSchema.Objects = processedObjects
	newSchema.Objects = append(newSchema.Objects, newObjects...)

	return &newSchema, nil
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

	// Ex: type | null
	if len(disjunction.Branches) == 2 && disjunction.Branches.HasNullType() {
		finalType := disjunction.Branches.NonNullTypes()[0]
		finalType.Nullable = true
		finalType.AddCompilerPassTrail("DisjunctionToType")

		return finalType, nil
	}

	// Ex: "some concrete value" | "some other value" | string
	if pass.hasOnlySingleTypeScalars(schema, disjunction) {
		refResolver := pass.makeReferenceResolver(schema)
		scalarKind := refResolver(disjunction.Branches[0]).AsScalar().ScalarKind
		scalarType := ast.NewScalar(scalarKind, ast.Default(def.Default))
		scalarType.AddCompilerPassTrail("DisjunctionToType")

		return scalarType, nil
	}

	// type | otherType | something (| null)?
	// generate a type with a nullable field for every branch of the disjunction,
	// add it to preprocessor.types, and use it instead.
	newTypeName := pass.disjunctionTypeName(disjunction)

	// if we already generated a new object for this disjunction, let's return
	// a reference to it.
	if _, ok := pass.newObjects[newTypeName]; ok {
		ref := ast.NewRef(schema.Package, newTypeName, ast.Hints(def.Hints))
		if def.Nullable || disjunction.Branches.HasNullType() {
			ref.Nullable = true
		}
		ref.AddCompilerPassTrail("DisjunctionToType")

		return ref, nil
	}

	/*
		TODO: return an error here. Some jennies won't be able to handle
		this type of disjunction.
		if !def.Branches.HasOnlyScalarOrArray() || !def.Branches.HasOnlyRefs() {
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
		structType.AddCompilerPassTrail("DisjunctionToType")
	}
	if disjunction.Branches.HasOnlyRefs() {
		newDisjunctionDef, err := pass.ensureDiscriminator(schema, disjunction)
		if err != nil {
			return ast.Type{}, err
		}

		structType.Hints[ast.HintDiscriminatedDisjunctionOfRefs] = newDisjunctionDef
		structType.AddCompilerPassTrail("DisjunctionToType")
	}

	pass.newObjects[newTypeName] = ast.Object{
		Name: newTypeName,
		Type: structType,
		SelfRef: ast.RefType{
			ReferredPkg:  schema.Package,
			ReferredType: newTypeName,
		},
	}

	ref := ast.NewRef(schema.Package, newTypeName, ast.Hints(def.Hints))
	if def.Nullable || disjunction.Branches.HasNullType() {
		ref.Nullable = true
	}
	ref.AddCompilerPassTrail("DisjunctionToType")

	return ref, nil
}

func (pass *DisjunctionToType) disjunctionTypeName(def ast.DisjunctionType) string {
	parts := make([]string, 0, len(def.Branches))

	for _, subType := range def.Branches {
		parts = append(parts, ast.TypeName(subType))
	}

	return strings.Join(parts, "Or")
}

func (pass *DisjunctionToType) ensureDiscriminator(schema *ast.Schema, def ast.DisjunctionType) (ast.DisjunctionType, error) {
	// discriminator-related data was set during parsing: nothing to do.
	if def.Discriminator != "" && len(def.DiscriminatorMapping) != 0 {
		return def, nil
	}

	newDef := def

	if def.Discriminator == "" {
		newDef.Discriminator = pass.inferDiscriminatorField(schema, newDef)
	}

	if len(def.DiscriminatorMapping) == 0 {
		mapping, err := pass.buildDiscriminatorMapping(schema, newDef)
		if err != nil {
			return newDef, err
		}

		newDef.DiscriminatorMapping = mapping
	}

	return newDef, nil
}

// inferDiscriminatorField tries to identify a field that might be used
// as a way to distinguish between the types in the disjunction branches.
// Such a field must:
//   - exist in all structs referred by the disjunction
//   - have a concrete, scalar value
//
// Note: this function assumes a disjunction of references to structs.
func (pass *DisjunctionToType) inferDiscriminatorField(schema *ast.Schema, def ast.DisjunctionType) string {
	fieldName := ""
	// map[typeName][fieldName]value
	candidates := make(map[string]map[string]any)

	refResolver := pass.makeReferenceResolver(schema)

	// Identify candidates from each branch
	for _, branch := range def.Branches {
		typeName := branch.AsRef().ReferredType
		structType := refResolver(branch).AsStruct()
		candidates[typeName] = make(map[string]any)

		for _, field := range structType.Fields {
			if field.Type.Kind != ast.KindScalar {
				continue
			}

			scalarField := field.Type.AsScalar()
			if !scalarField.IsConcrete() {
				continue
			}
			if field.Type.AsScalar().ScalarKind != ast.KindString {
				continue
			}

			candidates[typeName][field.Name] = scalarField.Value
		}
	}

	// At this point, if a discriminator exists, it will be listed under the candidates
	// of any type in our map.
	// We need to check if all other types have a similar field.
	someType := def.Branches[0].AsRef().ReferredType
	allTypes := make([]string, 0, len(candidates))

	for typeName := range candidates {
		allTypes = append(allTypes, typeName)
	}

	for candidateFieldName := range candidates[someType] {
		existsInAllBranches := true
		for _, branchTypeName := range allTypes {
			if _, ok := candidates[branchTypeName][candidateFieldName]; !ok {
				existsInAllBranches = false
				break
			}
		}

		if existsInAllBranches {
			fieldName = candidateFieldName
			break
		}
	}

	return fieldName
}

func (pass *DisjunctionToType) buildDiscriminatorMapping(schema *ast.Schema, def ast.DisjunctionType) (map[string]string, error) {
	mapping := make(map[string]string, len(def.Branches))
	if def.Discriminator == "" {
		return nil, fmt.Errorf("discriminator field is empty: %w", ErrCanNotInferDiscriminator)
	}

	refResolver := pass.makeReferenceResolver(schema)

	for _, branch := range def.Branches {
		typeName := branch.AsRef().ReferredType
		structType := refResolver(branch).AsStruct()

		field, found := structType.FieldByName(def.Discriminator)
		if !found {
			return nil, fmt.Errorf("discriminator field '%s' not found in Object '%s': %w", def.Discriminator, typeName, ErrCanNotInferDiscriminator)
		}

		// trust, but verify: we need the field to be an actual scalar with a concrete value?
		if field.Type.Kind != ast.KindScalar {
			return nil, fmt.Errorf("discriminator field is not a scalar: %w", ErrCanNotInferDiscriminator)
		}

		switch {
		case field.Type.AsScalar().IsConcrete():
			mapping[field.Type.AsScalar().Value.(string)] = typeName
		case field.Type.Default != nil:
			mapping[field.Type.Default.(string)] = typeName
		default:
			return nil, fmt.Errorf("discriminator field is not concrete: %w", ErrCanNotInferDiscriminator)
		}
	}

	return mapping, nil
}

func (pass *DisjunctionToType) hasOnlySingleTypeScalars(schema *ast.Schema, disjunction ast.DisjunctionType) bool {
	branches := disjunction.Branches

	if len(branches) == 0 {
		return false
	}

	refResolver := pass.makeReferenceResolver(schema)

	if refResolver(branches[0]).Kind != ast.KindScalar {
		return false
	}

	scalarKind := refResolver(branches[0]).AsScalar().ScalarKind
	for _, t := range branches {
		resolvedType := refResolver(t)

		if resolvedType.Kind != ast.KindScalar {
			return false
		}

		if resolvedType.AsScalar().ScalarKind != scalarKind {
			return false
		}
	}

	return true
}

func (pass *DisjunctionToType) makeReferenceResolver(schema *ast.Schema) func(typeDef ast.Type) ast.Type {
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
