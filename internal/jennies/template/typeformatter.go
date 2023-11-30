package template

import (
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/jennies/common"
)

type TypeFormatter interface {
	PackageMapper(pkg string) string
	VariantInterface(variant string) string
	WithContext(ctx common.Context) TypeFormatter
	ForBuilder(builder ast.Builder) TypeFormatter

	FormatType(def ast.Type, resolveBuilders bool) string
	FormatStructBody(def ast.Type) string
	FormatField(def ast.StructField) string
	FormatScalar(def ast.Type) string
	FormatRef(def ast.Type, resolveBuilders bool) string
	FormatArray(def ast.ArrayType, resolveBuilders bool) string
	FormatMap(def ast.MapType) string
	FormatDisjunction(def ast.DisjunctionType) string
	FormatIntersection(def ast.IntersectionType) string
	FormatComposableSlot(def ast.ComposableSlotType, resolveBuilders bool) string
}
