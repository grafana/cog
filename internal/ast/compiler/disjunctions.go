package compiler

import (
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
		processedObjects = append(processedObjects, pass.processObject(object))
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

func (pass *DisjunctionToType) processObject(object ast.Object) ast.Object {
	newObject := object
	newObject.Type = pass.processType(object.Type)

	return newObject
}

func (pass *DisjunctionToType) processType(def ast.Type) ast.Type {
	if def.Kind == ast.KindArray {
		return pass.processArray(def.AsArray())
	}

	if def.Kind == ast.KindStruct {
		return pass.processStruct(def.AsStruct())
	}

	if def.Kind == ast.KindDisjunction {
		return pass.processDisjunction(def.AsDisjunction())
	}

	return def
}

func (pass *DisjunctionToType) processArray(def ast.ArrayType) ast.Type {
	return ast.NewArray(pass.processType(def.ValueType))
}

func (pass *DisjunctionToType) processStruct(def ast.StructType) ast.Type {
	processedFields := make([]ast.StructField, 0, len(def.Fields))
	for _, field := range def.Fields {
		processedFields = append(processedFields, ast.StructField{
			Name:     field.Name,
			Comments: field.Comments,
			Type:     pass.processType(field.Type),
			Required: field.Required,
			Default:  field.Default,
		})
	}

	return ast.NewStruct(processedFields)
}

func (pass *DisjunctionToType) processDisjunction(def ast.DisjunctionType) ast.Type {
	// Ex: type | null
	if len(def.Branches) == 2 && def.Branches.HasNullType() {
		finalType := def.Branches.NonNullTypes()[0]
		// FIXME: this should be propagated
		// finalType.Nullable = true

		return finalType
	}

	// type | otherType | something (| null)?
	// generate a type with a nullable field for every branch of the disjunction,
	// add it to preprocessor.types, and use it instead.
	newTypeName := pass.disjunctionTypeName(def)

	if _, ok := pass.newObjects[newTypeName]; !ok {
		fields := make([]ast.StructField, 0, len(def.Branches))
		for _, branch := range def.Branches {
			if branch.IsNull() {
				continue
			}

			fields = append(fields, ast.StructField{
				Name:     "Val" + tools.UpperCamelCase(pass.typeName(branch)),
				Type:     branch,
				Required: false,
			})
		}

		structType := ast.NewStruct(fields)
		if def.Branches.HasOnlyScalarOrArray() {
			structType.Struct.Hint[ast.HintDisjunctionOfScalars] = true
		}

		pass.newObjects[newTypeName] = ast.Object{
			Name: newTypeName,
			Type: structType,
		}
	}

	return ast.NewRef(newTypeName)
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
