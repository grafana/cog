package builder

import (
	"strings"

	"github.com/grafana/cog/internal/ast"
)

type Selector func(builder ast.Builder) bool

// EveryBuilder accepts any given builder.
func EveryBuilder() Selector {
	return func(_ ast.Builder) bool {
		return true
	}
}

// ByObjectName matches builders for the given the object (referred to by its
// package and name).
// Note: the comparison on object name is case-insensitive.
func ByObjectName(pkg string, objectName string) Selector {
	return func(builder ast.Builder) bool {
		return builder.For.SelfRef.ReferredPkg == pkg &&
			strings.EqualFold(builder.For.SelfRef.ReferredType, objectName)
	}
}

// ByName matches builders for the given name.
// Note: the comparison on builder name is case-insensitive.
func ByName(pkg string, builderName string) Selector {
	return func(builder ast.Builder) bool {
		return builder.For.SelfRef.ReferredPkg == pkg &&
			strings.EqualFold(builder.Name, builderName)
	}
}

// StructGeneratedFromDisjunction matches builders for structs that were
// generated from a disjunction (see the Disjunction compiler pass).
func StructGeneratedFromDisjunction() Selector {
	return func(builder ast.Builder) bool {
		resolved, found := builder.Schema.Resolve(builder.For.Type)
		if !found {
			return false
		}

		return resolved.IsStructGeneratedFromDisjunction()
	}
}

// ComposableDashboardPanel matches builders for Panel variants.
func ComposableDashboardPanel() Selector {
	return func(builder ast.Builder) bool {
		return builder.Schema.Metadata.Kind == ast.SchemaKindComposable &&
			builder.Schema.Metadata.Variant == ast.SchemaVariantPanel &&
			builder.Schema.Metadata.Identifier != ""
	}
}
