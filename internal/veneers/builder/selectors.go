package builder

import (
	"fmt"
	"strings"

	"github.com/grafana/cog/internal/ast"
)

type Selector struct {
	description string
	matcher     func(schemas ast.Schemas, builder ast.Builder) bool
}

func (selector Selector) Matches(schemas ast.Schemas, builder ast.Builder) bool {
	return selector.matcher(schemas, builder)
}

func (selector Selector) String() string {
	return selector.description
}

// EveryBuilder accepts any given builder.
func EveryBuilder() *Selector {
	return &Selector{
		description: "every_builder",
		matcher: func(_ ast.Schemas, _ ast.Builder) bool {
			return true
		},
	}
}

// ByObjectName matches builders for the given the object (referred to by its
// package and name).
// Note: the comparison on object name is case-insensitive.
func ByObjectName(pkg string, objectName string) *Selector {
	return &Selector{
		description: fmt.Sprintf("by_object_name[builder.for.pkg='%s', builder.for.obj='%s']", pkg, objectName),
		matcher: func(_ ast.Schemas, builder ast.Builder) bool {
			return (strings.EqualFold(builder.For.SelfRef.ReferredPkg, pkg) || pkg == "*") &&
				strings.EqualFold(builder.For.SelfRef.ReferredType, objectName)
		},
	}
}

// ByName matches builders for the given name.
// Note: the comparison on builder name is case-insensitive.
func ByName(pkg string, builderName string) *Selector {
	return &Selector{
		description: fmt.Sprintf("by_name[builder.pkg='%s', builder.name='%s']", pkg, builderName),
		matcher: func(_ ast.Schemas, builder ast.Builder) bool {
			return (strings.EqualFold(builder.Package, pkg) || pkg == "*") &&
				strings.EqualFold(builder.Name, builderName)
		},
	}
}

// StructGeneratedFromDisjunction matches builders for structs that were
// generated from a disjunction (see the Disjunction compiler pass).
func StructGeneratedFromDisjunction() *Selector {
	return &Selector{
		description: "struct_generated_from_disjunction",
		matcher: func(schemas ast.Schemas, builder ast.Builder) bool {
			resolved := schemas.ResolveToType(builder.For.Type)

			return resolved.IsStructGeneratedFromDisjunction()
		},
	}
}

// ByVariant matches builders defined within a schema marked as "composable"
// and implementing the given variant.
func ByVariant(variant ast.SchemaVariant) *Selector {
	return &Selector{
		description: fmt.Sprintf("by_variant[kind='composable', identifier!='', variant='%s']", variant),
		matcher: func(schemas ast.Schemas, builder ast.Builder) bool {
			schema, found := schemas.Locate(builder.For.SelfRef.ReferredPkg)
			if !found {
				return false
			}

			return schema.Metadata.Kind == ast.SchemaKindComposable &&
				schema.Metadata.Variant == variant &&
				schema.Metadata.Identifier != ""
		},
	}
}
