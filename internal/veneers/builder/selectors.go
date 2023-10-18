package builder

import (
	"strings"

	"github.com/grafana/cog/internal/ast"
)

type Selector func(builder ast.Builder) bool

func ByObjectName(objectName string) Selector {
	return func(builder ast.Builder) bool {
		objectPkg, objectNameWithoutPkg, found := strings.Cut(objectName, ".")
		if !found {
			return builder.For.Name == objectName
		}

		return builder.For.SelfRef.ReferredPkg == objectPkg && builder.For.SelfRef.ReferredType == objectNameWithoutPkg
	}
}

func StructGeneratedFromDisjunction() Selector {
	return func(builder ast.Builder) bool {
		return builder.For.Type.IsStructGeneratedFromDisjunction()
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
