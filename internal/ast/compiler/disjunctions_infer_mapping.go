package compiler

import (
	"errors"
	"fmt"

	"github.com/grafana/cog/internal/ast"
)

var _ Pass = (*DisjunctionInferMapping)(nil)

var ErrCanNotInferDiscriminator = errors.New("can not infer discriminator mapping")

// DisjunctionInferMapping infers the discriminator field and mapping used to
// describe a disjunction of references.
// See https://swagger.io/docs/specification/data-models/inheritance-and-polymorphism/
type DisjunctionInferMapping struct {
}

func (pass *DisjunctionInferMapping) Process(schemas []*ast.Schema) ([]*ast.Schema, error) {
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

func (pass *DisjunctionInferMapping) processSchema(schema *ast.Schema) (*ast.Schema, error) {
	var err error
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

	return schema, nil
}

func (pass *DisjunctionInferMapping) processObject(schema *ast.Schema, object ast.Object) (ast.Object, error) {
	var err error

	object.Type, err = pass.processType(schema, object.Type)
	if err != nil {
		return object, errors.Join(
			fmt.Errorf("could not process object '%s'", object.Name),
			err,
		)
	}

	return object, nil
}

func (pass *DisjunctionInferMapping) processType(schema *ast.Schema, def ast.Type) (ast.Type, error) {
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

	return def, nil
}

func (pass *DisjunctionInferMapping) processArray(schema *ast.Schema, def ast.Type) (ast.Type, error) {
	var err error

	def.Array.ValueType, err = pass.processType(schema, def.AsArray().ValueType)
	if err != nil {
		return ast.Type{}, err
	}

	return def, nil
}

func (pass *DisjunctionInferMapping) processMap(schema *ast.Schema, def ast.Type) (ast.Type, error) {
	var err error

	def.Map.ValueType, err = pass.processType(schema, def.AsMap().ValueType)
	if err != nil {
		return ast.Type{}, err
	}

	return def, nil
}

func (pass *DisjunctionInferMapping) processStruct(schema *ast.Schema, def ast.Type) (ast.Type, error) {
	var err error

	for i, field := range def.Struct.Fields {
		def.Struct.Fields[i].Type, err = pass.processType(schema, field.Type)
		if err != nil {
			return ast.Type{}, errors.Join(
				fmt.Errorf("could not process struct field '%s'", field.Name),
				err,
			)
		}
	}

	return def, nil
}

func (pass *DisjunctionInferMapping) processDisjunction(schema *ast.Schema, def ast.Type) (ast.Type, error) {
	var err error

	if !def.Disjunction.Branches.HasOnlyRefs() {
		return def, nil
	}

	def.Disjunction, err = pass.ensureDiscriminator(schema, def)
	if err != nil {
		return ast.Type{}, err
	}

	return def, nil
}

func (pass *DisjunctionInferMapping) ensureDiscriminator(schema *ast.Schema, def ast.Type) (*ast.DisjunctionType, error) {
	disjunction := def.Disjunction

	// discriminator-related data was set during parsing: nothing to do.
	if disjunction.Discriminator != "" && len(disjunction.DiscriminatorMapping) != 0 {
		return disjunction, nil
	}

	if disjunction.Discriminator == "" {
		disjunction.Discriminator = pass.inferDiscriminatorField(schema, disjunction)
		def.AddToPassesTrail("DisjunctionInferMapping[discriminator inferred]")
	}

	if len(disjunction.DiscriminatorMapping) == 0 {
		mapping, err := pass.buildDiscriminatorMapping(schema, disjunction)
		if err != nil {
			return disjunction, err
		}

		disjunction.DiscriminatorMapping = mapping
		def.AddToPassesTrail("DisjunctionInferMapping[mapping inferred]")
	}

	return disjunction, nil
}

// inferDiscriminatorField tries to identify a field that might be used
// as a way to distinguish between the types in the disjunction branches.
// Such a field must:
//   - exist in all structs referred by the disjunction
//   - have a concrete, scalar value
//
// Note: this function assumes a disjunction of references to structs.
func (pass *DisjunctionInferMapping) inferDiscriminatorField(schema *ast.Schema, def *ast.DisjunctionType) string {
	fieldName := ""
	// map[typeName][fieldName]value
	candidates := make(map[string]map[string]any)

	// Identify candidates from each branch
	for _, branch := range def.Branches {
		referredType, found := schema.Resolve(branch)
		if !found {
			continue
		}

		if !referredType.IsStruct() {
			continue
		}

		typeName := branch.AsRef().ReferredType
		structType := referredType.AsStruct()
		candidates[typeName] = make(map[string]any)

		for _, field := range structType.Fields {
			if !field.Type.IsConcreteScalar() {
				continue
			}

			scalarField := field.Type.AsScalar()
			if scalarField.ScalarKind != ast.KindString {
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

func (pass *DisjunctionInferMapping) buildDiscriminatorMapping(schema *ast.Schema, def *ast.DisjunctionType) (map[string]string, error) {
	mapping := make(map[string]string, len(def.Branches))
	if def.Discriminator == "" {
		return nil, fmt.Errorf("discriminator field is empty: %w", ErrCanNotInferDiscriminator)
	}

	for _, branch := range def.Branches {
		typeName := branch.AsRef().ReferredType
		referredType, found := schema.Resolve(branch)
		if !found {
			return nil, fmt.Errorf("could not resolve reference '%s'", branch.AsRef().String())
		}

		structType := referredType.AsStruct()

		field, found := structType.FieldByName(def.Discriminator)
		if !found {
			return nil, fmt.Errorf("discriminator field '%s' not found in Object '%s': %w", def.Discriminator, typeName, ErrCanNotInferDiscriminator)
		}

		// trust, but verify: we need the field to be an actual scalar with a concrete value?
		if !field.Type.IsScalar() {
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
