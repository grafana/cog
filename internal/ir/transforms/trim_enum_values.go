package transforms

import (
	"strings"

	"github.com/grafana/cog/internal/ir"
)

var _ Transform = (*TrimEnumValues)(nil)

// TrimEnumValues removes leading and trailing spaces from string values.
// It could happen when they add them by mistake in jsonschema/openapi when they define the enums
type TrimEnumValues struct {
}

func (t TrimEnumValues) Process(schemas ir.Schemas) (ir.Schemas, error) {
	visitor := Visitor{
		OnEnum: t.processEnum,
	}

	return visitor.VisitSchemas(schemas)
}

func (t TrimEnumValues) processEnum(_ *Visitor, _ *ir.Schema, def ir.Type) (ir.Type, error) {
	for i, value := range def.AsEnum().Values {
		if stringType, ok := value.Value.(string); ok {
			def.AsEnum().Values[i].Value = strings.TrimSpace(stringType)
		}
	}

	return def, nil
}
