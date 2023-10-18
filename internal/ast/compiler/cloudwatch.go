package compiler

import (
	"github.com/grafana/cog/internal/ast"
)

var _ Pass = (*Cloudwatch)(nil)

// Cloudwatch defines some kind of "types IR veneer",
// where we use a compiler pass to rewrite a part of the cloudwatch schema.
//
// In that schema, the `QueryEditorExpression` type is defined as a disjunction
// for which the discriminator and mapping can not be inferred.
// This compiler pass is here to define that mapping.
//
// The `QueryEditorArrayExpression` struct type is also modified to simplify the
// definition of its `expression` field from `[...#QueryEditorExpression] | [...#QueryEditorArrayExpression]` to
// `[...#QueryEditorExpression]`.
// This should be semantically equivalent since `#QueryEditorExpression` is a
// union type that includes `#QueryEditorArrayExpression`.
type Cloudwatch struct {
}

func (pass *Cloudwatch) Process(schemas []*ast.Schema) ([]*ast.Schema, error) {
	newSchemas := make([]*ast.Schema, 0, len(schemas))

	for _, schema := range schemas {
		if schema.Package != "cloudwatch" {
			newSchemas = append(newSchemas, schema)
			continue
		}

		newSchema, err := pass.processSchema(schema)
		if err != nil {
			return nil, err
		}

		newSchemas = append(newSchemas, newSchema)
	}

	return newSchemas, nil
}

func (pass *Cloudwatch) processSchema(schema *ast.Schema) (*ast.Schema, error) {
	newSchema := schema.DeepCopy()
	newSchema.Objects = nil

	for _, object := range schema.Objects {
		if object.Name == "QueryEditorExpression" {
			newSchema.Objects = append(newSchema.Objects, pass.processQueryEditorExpression(object))
			continue
		}
		if object.Name == "QueryEditorArrayExpression" {
			newSchema.Objects = append(newSchema.Objects, pass.processQueryEditorArrayExpression(object))
			continue
		}

		newSchema.Objects = append(newSchema.Objects, object)
	}

	return &newSchema, nil
}

func (pass *Cloudwatch) processQueryEditorExpression(object ast.Object) ast.Object {
	if object.Type.Kind != ast.KindDisjunction {
		return object
	}

	object.Type.Disjunction.Discriminator = "type"
	object.Type.Disjunction.DiscriminatorMapping = map[string]string{
		"and":               "QueryEditorArrayExpression",
		"or":                "QueryEditorArrayExpression",
		"property":          "QueryEditorPropertyExpression",
		"groupBy":           "QueryEditorGroupByExpression",
		"function":          "QueryEditorFunctionExpression",
		"functionParameter": "QueryEditorFunctionParameterExpression",
		"operator":          "QueryEditorOperatorExpression",
	}

	return object
}

func (pass *Cloudwatch) processQueryEditorArrayExpression(object ast.Object) ast.Object {
	if object.Type.Kind != ast.KindStruct {
		return object
	}

	structDef := object.Type.AsStruct()
	fields := make([]ast.StructField, 0, len(structDef.Fields)-1)

	for _, field := range structDef.Fields {
		if field.Name != "expressions" {
			fields = append(fields, field)
			continue
		}

		newField := field.DeepCopy()
		newField.Type = newField.Type.Disjunction.Branches[0]

		fields = append(fields, newField)
	}

	object.Type.Struct.Fields = fields

	return object
}
