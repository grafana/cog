package golang

import (
	"fmt"
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/jennies/common"
	"github.com/grafana/cog/internal/jennies/template"
	"github.com/grafana/cog/internal/tools"
	"strings"
)

type Formatter struct {
	context common.Context
	im      *common.DirectImportMap

	pkgRoot string
	builder ast.Builder
}

func NewFormatter(pkgRoot string, im *common.DirectImportMap) *Formatter {
	f := &Formatter{
		pkgRoot: pkgRoot,
		im:      im,
	}

	f.PackageMapper("cog")
	return f
}

func (f *Formatter) WithContext(ctx common.Context) template.TypeFormatter {
	f.context = ctx
	return f
}

func (f *Formatter) ForBuilder(builder ast.Builder) template.TypeFormatter {
	f.builder = builder
	return f
}

func (f *Formatter) VariantInterface(variant string) string {
	referredPkg := f.PackageMapper("cog/variants")

	return fmt.Sprintf("%s.%s", referredPkg, tools.UpperCamelCase(variant))
}

func (f *Formatter) PackageMapper(pkg string) string {
	if pkg == f.builder.Package {
		return ""
	}

	return f.im.Add(pkg, f.importPath(pkg))
}

func (f *Formatter) FormatType(def ast.Type, resolveBuilders bool) string {
	if def.IsAny() {
		return "any"
	}

	switch def.Kind {
	case ast.KindComposableSlot:
		return f.FormatComposableSlot(def.AsComposableSlot(), resolveBuilders)
	case ast.KindArray:
		return f.FormatArray(def.AsArray(), resolveBuilders)
	case ast.KindMap:
		return f.FormatMap(def.AsMap())
	case ast.KindScalar:
		return f.FormatScalar(def)
	case ast.KindRef:
		return f.FormatRef(def, resolveBuilders)
	case ast.KindStruct:
		// anonymous struct or struct body
		return f.FormatStructBody(def)
	case ast.KindIntersection:
		return f.FormatIntersection(def.AsIntersection())
	default:
		return "unknown"
	}
}

func (f *Formatter) FormatStructBody(def ast.Type) string {
	var buffer strings.Builder

	buffer.WriteString("struct {\n")

	for _, fieldDef := range def.AsStruct().Fields {
		buffer.WriteString("\t" + f.FormatField(fieldDef))
	}

	buffer.WriteString("}")
	output := buffer.String()

	if def.Nullable {
		output = "*" + output
	}

	return output
}

func (f *Formatter) FormatField(def ast.StructField) string {
	var buffer strings.Builder

	for _, commentLine := range def.Comments {
		buffer.WriteString(fmt.Sprintf("// %s\n", commentLine))
	}

	jsonOmitEmpty := ""
	if !def.Required {
		jsonOmitEmpty = ",omitempty"
	}

	buffer.WriteString(fmt.Sprintf(
		"%s %s `json:\"%s%s\"`\n",
		tools.UpperCamelCase(def.Name),
		f.FormatType(def.Type, false),
		def.Name,
		jsonOmitEmpty,
	))

	return buffer.String()
}

func (f *Formatter) FormatScalar(def ast.Type) string {
	typeName := def.AsScalar().ScalarKind
	if def.Nullable {
		typeName = "*" + typeName
	}

	return string(typeName)
}

func (f *Formatter) FormatRef(def ast.Type, resolveBuilders bool) string {
	referredPkg := f.PackageMapper(def.AsRef().ReferredPkg)
	typeName := tools.UpperCamelCase(def.AsRef().ReferredType)

	if referredPkg != "" {
		typeName = referredPkg + "." + typeName
	}

	if resolveBuilders && f.context.ResolveToBuilder(def) {
		cogAlias := f.PackageMapper("cog")

		return fmt.Sprintf("%s.Builder[%s]", cogAlias, typeName)
	}

	if def.Nullable {
		typeName = "*" + typeName
	}

	return typeName
}

func (f *Formatter) FormatArray(def ast.ArrayType, resolveBuilders bool) string {
	subTypeString := f.FormatType(def.ValueType, resolveBuilders)

	return fmt.Sprintf("[]%s", subTypeString)
}

func (f *Formatter) FormatMap(def ast.MapType) string {
	keyTypeString := f.FormatType(def.IndexType, false)
	valueTypeString := f.FormatType(def.ValueType, false)

	return fmt.Sprintf("map[%s]%s", keyTypeString, valueTypeString)
}

func (f *Formatter) FormatDisjunction(_ ast.DisjunctionType) string {
	return ""
}

func (f *Formatter) FormatIntersection(def ast.IntersectionType) string {
	var buffer strings.Builder

	buffer.WriteString("struct {\n")

	refs := make([]ast.Type, 0)
	rest := make([]ast.Type, 0)
	for _, b := range def.Branches {
		if b.Ref != nil {
			refs = append(refs, b)
			continue
		}

		rest = append(rest, b)
	}

	for _, ref := range refs {
		buffer.WriteString("\t" + f.FormatRef(ref, false) + "\n")
	}

	if len(refs) > 0 {
		buffer.WriteString("\n")
	}

	for _, r := range rest {
		if r.Struct != nil {
			for _, fieldDef := range r.AsStruct().Fields {
				buffer.WriteString("\t" + f.FormatField(fieldDef))
			}
			continue
		}
		buffer.WriteString("\t" + f.FormatType(r, false) + "\n")
	}

	buffer.WriteString("}")

	return buffer.String()
}

func (f *Formatter) FormatComposableSlot(def ast.ComposableSlotType, resolveBuilders bool) string {
	formatted := f.VariantInterface(string(def.Variant))

	if !resolveBuilders {
		return formatted
	}

	cogAlias := f.PackageMapper("cog")

	return fmt.Sprintf("%s.Builder[%s]", cogAlias, formatted)
}

func (f *Formatter) importPath(suffix string) string {
	root := strings.TrimSuffix(f.pkgRoot, "/")
	return fmt.Sprintf("%s/%s", root, suffix)
}
