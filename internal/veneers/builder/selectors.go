package builder

import (
	"github.com/grafana/cog/internal/ast"
)

type Selector func(builder ast.Builder) bool

func ByName(objectName string) Selector {
	return func(builder ast.Builder) bool {
		return builder.For.Name == objectName
	}
}

func StructGeneratedFromDisjunction() Selector {
	return func(builder ast.Builder) bool {
		if builder.For.Type.Kind != ast.KindStruct {
			return false
		}

		return builder.For.Type.AsStruct().IsGeneratedFromDisjunction()
	}
}

func ComposableDashboardPanel() Selector {
	return func(builder ast.Builder) bool {
		return builder.Schema.Metadata.Kind == ast.SchemaKindComposable &&
			builder.Schema.Metadata.Variant == ast.SchemaVariantPanel &&
			builder.Schema.Metadata.Identifier != ""
	}
}

func EveryBuilder() Selector {
	return func(builder ast.Builder) bool {
		return true
	}
}
