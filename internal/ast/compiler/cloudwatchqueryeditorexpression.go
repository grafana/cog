package compiler

import (
	"github.com/grafana/cog/internal/ast"
)

var _ Pass = (*CloudwatchQueryEditorExpression)(nil)

type CloudwatchQueryEditorExpression struct {
}

func (pass *CloudwatchQueryEditorExpression) Process(schemas []*ast.Schema) ([]*ast.Schema, error) {
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

func (pass *CloudwatchQueryEditorExpression) processSchema(schema *ast.Schema) (*ast.Schema, error) {
	newSchema := schema.DeepCopy()
	newSchema.Objects = nil

	for _, object := range schema.Objects {
		if object.Name != "QueryEditorExpression" {
			newSchema.Objects = append(newSchema.Objects, object)
			continue
		}

		newSchema.Objects = append(newSchema.Objects, pass.processQueryEditorExpression(object))
	}

	return &newSchema, nil
}

func (pass *CloudwatchQueryEditorExpression) processQueryEditorExpression(object ast.Object) ast.Object {
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
