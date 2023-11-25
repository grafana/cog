package builder

import (
	"strings"

	"github.com/grafana/cog/internal/ast"
)

type Selector func(builder ast.Builder) bool

// Note: comparison on object name is case insensitive
func ByObjectName(pkg string, objectName string) Selector {
	return func(builder ast.Builder) bool {
		return builder.For.SelfRef.ReferredPkg == pkg &&
			strings.EqualFold(builder.For.SelfRef.ReferredType, objectName)
	}
}

func StructGeneratedFromDisjunction() Selector {
	return func(builder ast.Builder) bool {
		resolved, found := builder.Schema.Resolve(builder.For.Type)
		if !found {
			return false
		}

		return resolved.IsStructGeneratedFromDisjunction()
	}
}

func ComposableDashboardPanel() Selector {
	return func(builder ast.Builder) bool {
		return builder.Schema.Metadata.Kind == ast.SchemaKindComposable &&
			builder.Schema.Metadata.Variant == ast.SchemaVariantPanel &&
			builder.Schema.Metadata.Identifier != ""
	}
}

func ComposableDataQuery() Selector {
	return func(builder ast.Builder) bool {
		return builder.Schema.Metadata.Kind == ast.SchemaKindComposable &&
			builder.Schema.Metadata.Variant == ast.SchemaVariantDataQuery &&
			builder.Schema.Metadata.Identifier != ""
	}
}

func EveryBuilder() Selector {
	return func(builder ast.Builder) bool {
		return true
	}
}
