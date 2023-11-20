package compiler

import "github.com/grafana/cog/internal/ast"

var _ Pass = (*CandlestickColors)(nil)

const (
	CandleStickPackage       = "candlestick"
	CandleStickOptionsObject = "Options"
	CandleStickColorsObject  = "CandlestickColors"
	CandleStickColorsField   = "colors"
)

type CandlestickColors struct {
}

func (c *CandlestickColors) Process(schemas []*ast.Schema) ([]*ast.Schema, error) {
	newSchemas := make([]*ast.Schema, len(schemas))

	for i, schema := range schemas {
		if schema.Package != CandleStickPackage {
			newSchemas[i] = schema
			continue
		}

		newSchemas[i] = c.processSchema(schema)
	}

	return newSchemas, nil
}

func (c *CandlestickColors) processSchema(schema *ast.Schema) *ast.Schema {
	newSchema := schema.DeepCopy()
	newSchema.Objects = nil

	for _, object := range schema.Objects {
		if object.Name != CandleStickOptionsObject {
			newSchema.Objects = append(newSchema.Objects, object)
			continue
		}

		newSchema.Objects = append(newSchema.Objects, c.processObject(object))
	}

	return &newSchema
}

func (c *CandlestickColors) processObject(object ast.Object) ast.Object {
	if object.Type.Kind != ast.KindStruct {
		return object
	}

	for i, field := range object.Type.AsStruct().Fields {
		if field.Name != CandleStickColorsField {
			continue
		}

		if field.Type.Kind != ast.KindStruct {
			continue
		}

		object.Type.AsStruct().Fields[i].Type = ast.NewRef(object.SelfRef.ReferredPkg, CandleStickColorsObject)
	}

	return object
}
