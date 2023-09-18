package compiler

import (
	"fmt"
	"sort"
	"strings"

	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/tools"
)

var _ Pass = (*DisjunctionToType)(nil)

type DisjunctionToType struct {
	newObjects map[string]ast.Object
}

func (pass *DisjunctionToType) Process(files []*ast.File) ([]*ast.File, error) {
	newFiles := make([]*ast.File, 0, len(files))

	for _, file := range files {
		newFile, err := pass.processFile(file)
		if err != nil {
			return nil, err
		}

		newFiles = append(newFiles, newFile)
	}

	return newFiles, nil
}

func (pass *DisjunctionToType) processFile(file *ast.File) (*ast.File, error) {
	pass.newObjects = make(map[string]ast.Object)

	processedObjects := make([]ast.Object, 0, len(file.Definitions))
	for _, object := range file.Definitions {
		processedObject, err := pass.processObject(file, object)
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

	return &ast.File{
		Package:     file.Package,
		Definitions: append(processedObjects, newObjects...),
	}, nil
}

func (pass *DisjunctionToType) processObject(file *ast.File, object ast.Object) (ast.Object, error) {
	processedType, err := pass.processType(file, object.Type)
	if err != nil {
		return object, err
	}

	newObject := object
	newObject.Type = processedType

	return newObject, nil
}

func (pass *DisjunctionToType) processType(file *ast.File, def ast.Type) (ast.Type, error) {
	if def.Kind == ast.KindArray {
		return pass.processArray(file, def.AsArray())
	}

	if def.Kind == ast.KindStruct {
		return pass.processStruct(file, def.AsStruct())
	}

	if def.Kind == ast.KindDisjunction {
		return pass.processDisjunction(file, def.AsDisjunction())
	}

	return def, nil
}

func (pass *DisjunctionToType) processArray(file *ast.File, def ast.ArrayType) (ast.Type, error) {
	processedType, err := pass.processType(file, def.ValueType)
	if err != nil {
		return ast.Type{}, err
	}

	return ast.NewArray(processedType), nil
}

func (pass *DisjunctionToType) processStruct(file *ast.File, def ast.StructType) (ast.Type, error) {
	processedFields := make([]ast.StructField, 0, len(def.Fields))
	for _, field := range def.Fields {
		processedType, err := pass.processType(file, field.Type)
		if err != nil {
			return ast.Type{}, err
		}

		processedFields = append(processedFields, ast.StructField{
			Name:     field.Name,
			Comments: field.Comments,
			Type:     processedType,
			Required: field.Required,
			Default:  field.Default,
		})
	}

	return ast.NewStruct(processedFields...), nil
}

func (pass *DisjunctionToType) processDisjunction(file *ast.File, def ast.DisjunctionType) (ast.Type, error) {
	// Ex: type | null
	if len(def.Branches) == 2 && def.Branches.HasNullType() {
		finalType := def.Branches.NonNullTypes()[0]
		// FIXME: this should be propagated
		// finalType.Nullable = true

		return finalType, nil
	}

	// type | otherType | something (| null)?
	// generate a type with a nullable field for every branch of the disjunction,
	// add it to preprocessor.types, and use it instead.
	newTypeName := pass.disjunctionTypeName(def)

	// if we already generated a new object for this disjunction, let's return
	// a reference to it.
	if _, ok := pass.newObjects[newTypeName]; ok {
		return ast.NewRef(newTypeName), nil
	}

	/*
		TODO: return an error here. Some jennies won't be able to handle
		this type of disjunction.
		if !def.Branches.HasOnlyScalarOrArray() || !def.Branches.HasOnlyRefs() {
		}
	*/

	fields := make([]ast.StructField, 0, len(def.Branches))
	for _, branch := range def.Branches {
		// FIXME: should ignore this completely.
		// ie: if there was a nullable branch, the whole resulting type should be nullable.
		if branch.IsNull() {
			continue
		}

		fields = append(fields, ast.StructField{
			Name:     "Val" + tools.UpperCamelCase(pass.typeName(branch)),
			Type:     branch,
			Required: false,
		})
	}

	structType := ast.NewStruct(fields...)
	if def.Branches.HasOnlyScalarOrArray() {
		structType.Struct.Hint[ast.HintDisjunctionOfScalars] = def
	}
	if def.Branches.HasOnlyRefs() {
		newDisjunctionDef, err := pass.ensureDiscriminator(file, def)
		if err != nil {
			return ast.Type{}, err
		}

		structType.Struct.Hint[ast.HintDiscriminatedDisjunctionOfStructs] = newDisjunctionDef
	}

	pass.newObjects[newTypeName] = ast.Object{
		Name: newTypeName,
		Type: structType,
	}

	return ast.NewRef(newTypeName), nil
}

func (pass *DisjunctionToType) disjunctionTypeName(def ast.DisjunctionType) string {
	parts := make([]string, 0, len(def.Branches))

	for _, subType := range def.Branches {
		parts = append(parts, tools.UpperCamelCase(pass.typeName(subType)))
	}

	return strings.Join(parts, "Or")
}

func (pass *DisjunctionToType) typeName(typeDef ast.Type) string {
	if typeDef.Kind == ast.KindRef {
		return typeDef.AsRef().ReferredType
	}
	if typeDef.Kind == ast.KindScalar {
		return string(typeDef.AsScalar().ScalarKind)
	}

	return string(typeDef.Kind)
}

func (pass *DisjunctionToType) ensureDiscriminator(file *ast.File, def ast.DisjunctionType) (ast.DisjunctionType, error) {
	// discriminator-related data was set during parsing: nothing to do.
	if def.Discriminator != "" && len(def.DiscriminatorMapping) != 0 {
		return def, nil
	}

	newDef := def

	if def.Discriminator == "" {
		newDef.Discriminator = pass.inferDiscriminatorField(file, newDef)
	}

	if len(def.DiscriminatorMapping) == 0 {
		mapping, err := pass.buildDiscriminatorMapping(file, newDef)
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
func (pass *DisjunctionToType) inferDiscriminatorField(file *ast.File, def ast.DisjunctionType) string {
	fieldName := ""
	// map[typeName][fieldName]value
	candidates := make(map[string]map[string]any)

	// Identify candidates from each branch
	for _, branch := range def.Branches {
		// FIXME: what if the definition is itself a reference? Resolve recursively?
		typeName := branch.AsRef().ReferredType
		structType := file.LocateDefinition(typeName).Type.AsStruct()
		candidates[typeName] = make(map[string]any)

		for _, field := range structType.Fields {
			if field.Type.Kind != ast.KindScalar {
				continue
			}

			scalarField := field.Type.AsScalar()
			if !scalarField.IsConcrete() {
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

func (pass *DisjunctionToType) buildDiscriminatorMapping(file *ast.File, def ast.DisjunctionType) (map[string]any, error) {
	mapping := make(map[string]any, len(def.Branches))
	if def.Discriminator == "" {
		return nil, fmt.Errorf("can not build discriminator mapping: discriminator field is empty")
	}

	for _, branch := range def.Branches {
		// FIXME: what if the definition is itself a reference? Resolve recursively?
		typeName := branch.AsRef().ReferredType
		structType := file.LocateDefinition(typeName).Type.AsStruct()

		field, found := structType.FieldByName(def.Discriminator)
		if !found {
			return nil, fmt.Errorf("can not build discriminator mapping: could not locate the definition of Ref<%s>", typeName)
		}

		// trust, but verify: we need the field to be an actual scalar with a concrete value?
		if field.Type.Kind != ast.KindScalar {
			return nil, fmt.Errorf("can not build discriminator mapping: discriminator field is not a scalar")
		}
		if !field.Type.AsScalar().IsConcrete() {
			return nil, fmt.Errorf("can not build discriminator mapping: discriminator field is not concrete")
		}

		mapping[typeName] = field.Type.AsScalar().Value
	}

	return mapping, nil
}
