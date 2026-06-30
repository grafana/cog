package terraform

import (
	"fmt"
	"io"
	"strings"

	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/languages"
)

type attributes struct {
	cfg           Config
	context       languages.Context
	packageMapper func(pkg string) string
	validators    *validators
	types         *typeFormatter

	refsInDisjunction map[string]struct{}
}

func newAttributesGenerator(context languages.Context, cfg Config, types *typeFormatter, packageMapper func(pkg string) string) *attributes {
	return &attributes{
		cfg:           cfg,
		context:       context,
		packageMapper: packageMapper,
		types:         types,
		validators:    newValidators(context, packageMapper),
	}
}

// Disjunction branches (that are struct) need to be generated with all their fields set to "Optional: true"
// For that to happen, we need to identify them before the object codegen starts
func (a *attributes) identifyDisjunctionBranches(schema *ast.Schema) {
	a.refsInDisjunction = map[string]struct{}{}

	schema.Objects.Iterate(func(_ string, obj ast.Object) {
		if !obj.Type.IsStructGeneratedFromDisjunction() {
			return
		}

		for _, field := range obj.Type.Struct.Fields {
			if !field.Type.IsRef() {
				continue
			}

			if a.context.ResolveRefs(field.Type).IsStruct() {
				a.refsInDisjunction[field.Type.Ref.String()] = struct{}{}
			}
		}
	})
}

func (a *attributes) generateForObject(obj ast.Object) string {
	if !obj.Type.IsStruct() {
		return ""
	}

	var buffer strings.Builder

	a.packageMapper("github.com/hashicorp/terraform-plugin-framework/resource/schema")

	_, usedInDisjunction := a.refsInDisjunction[obj.SelfRef.String()]

	fmt.Fprintf(&buffer, "var %sSchema = schema.SingleNestedBlock{\n", formatObjectName(obj.Name))
	if obj.DeprecationMessage != "" {
		fmt.Fprintf(&buffer, "\tDeprecationMessage: %#v,\n", obj.DeprecationMessage)
	}
	a.writeDescriptionAttrs(&buffer, obj.Comments)
	fmt.Fprintf(&buffer, "\tAttributes: %s,\n", a.formatStructAttributes(obj.Type, usedInDisjunction))
	fmt.Fprintf(&buffer, "\tBlocks: %s,\n", a.formatStructBlocks(obj.Type))

	if obj.Type.IsStructGeneratedFromDisjunction() {
		a.packageMapper("github.com/hashicorp/terraform-plugin-framework/schema/validator")

		branchNames := make([]string, 0, len(obj.Type.Struct.Fields))
		for _, field := range obj.Type.Struct.Fields {
			branchNames = append(branchNames, formatScalar(formatTfSDKAttrName(field.Name)))
		}

		buffer.WriteString("\tValidators: []validator.Object{\n")
		// TODO: error/warning if `a.cfg.Validators.AttributeCountExactly` isn't set?
		// Or should we provide an implementation?
		fmt.Fprintf(&buffer, "\t\t%s(1, %s),\n", a.cfg.Validators.AttributeCountExactly, strings.Join(branchNames, ", "))
		buffer.WriteString("\t},\n")
	}

	fmt.Fprintf(&buffer, "}")

	return buffer.String()
}

// Attributes generation

func (a *attributes) formatStructAttributes(typeDef ast.Type, forceFieldsAsOptional bool) string {
	var buffer strings.Builder

	a.packageMapper("github.com/hashicorp/terraform-plugin-framework/resource/schema")

	buffer.WriteString("map[string]schema.Attribute{\n")
	for _, field := range typeDef.Struct.Fields {
		if field.Type.IsIntersection() {
			continue
		}

		// fields pointing to a struct shouldn't be part of the attributes
		// instead, they should be referenced as a block
		if a.context.ResolveToStruct(field.Type) {
			continue
		}

		// constant refs shouldn't be exposed to users as their value can be set directly by the provider
		if field.Type.IsConstantRef() {
			continue
		}

		fmt.Fprintf(&buffer, "\t\"%s\": %s,\n", formatTfSDKAttrName(field.Name), a.formatTypeAttribute(field.Type, forceFieldsAsOptional, field.Comments))
	}

	buffer.WriteString("}")
	return buffer.String()
}

func (a *attributes) formatTypeAttribute(typeDef ast.Type, forceOptional bool, comments []string) string {
	switch typeDef.Kind {
	case ast.KindScalar:
		return a.formatScalarAttribute(typeDef, forceOptional, comments)
	case ast.KindStruct:
		return a.formatStructAttributes(typeDef, forceOptional)
	case ast.KindArray:
		return a.formatArrayAttributes(typeDef, comments)
	case ast.KindMap:
		return a.formatMapAttributes(typeDef, comments)
	case ast.KindRef:
		return a.formatReferenceAttribute(typeDef, comments)
	case ast.KindEnum:
		return a.formatEnumAttribute(typeDef, comments)
	default:
		return ""
	}
}

func (a *attributes) formatArrayAttributes(def ast.Type, comments []string) string {
	var buffer strings.Builder

	buffer.WriteString("schema.ListAttribute{\n")
	fmt.Fprintf(&buffer, "\tElementType: %s,\n", a.types.formatType(def.AsArray().ValueType))

	if def.Default != nil {
		a.packageMapper("github.com/hashicorp/terraform-plugin-framework/resource/schema/listdefault")
		fmt.Fprintf(&buffer, "\tDefault: listdefault.StaticValue(%s),\n", a.parseArrayOrMapDefaults(def.AsArray().ValueType, def.Default, ListDefault))
		buffer.WriteString("\tComputed: true,\n")
	}

	arrayValidator := a.validators.arrayConstraintValidator(def.AsArray().Constraints)
	validator := a.validators.validateList(def.AsArray().ValueType)

	if arrayValidator != "" {
		fmt.Fprintf(&buffer, "\tValidators: %s,\n", arrayValidator)
	} else if validator != "" {
		fmt.Fprintf(&buffer, "\tValidators: %s,\n", validator)
	}

	if def.Nullable || def.Default != nil {
		buffer.WriteString("\tOptional: true,\n")
	} else {
		buffer.WriteString("\tRequired: true,\n")
	}

	a.writeDescriptionAttrs(&buffer, comments)

	fmt.Fprintf(&buffer, "}")

	return buffer.String()
}

func (a *attributes) formatMapAttributes(def ast.Type, comments []string) string {
	var buffer strings.Builder

	buffer.WriteString("schema.MapAttribute{\n")
	fmt.Fprintf(&buffer, "\tElementType: %s,\n", a.types.formatType(def.AsMap().ValueType))

	if def.Default != nil {
		a.packageMapper("github.com/hashicorp/terraform-plugin-framework/resource/schema/mapdefault")
		fmt.Fprintf(&buffer, "\tDefault: mapdefault.StaticValue(%s),\n", a.parseArrayOrMapDefaults(def.AsMap().ValueType, def.Default, MapDefault))
		buffer.WriteString("\tComputed: true,\n")
	}

	if def.Nullable || def.Default != nil {
		buffer.WriteString("\tOptional: true,\n")
	} else {
		buffer.WriteString("\tRequired: true,\n")
	}

	a.writeDescriptionAttrs(&buffer, comments)

	fmt.Fprintf(&buffer, "}")

	return buffer.String()
}

func (a *attributes) formatReferenceAttribute(def ast.Type, comments []string) string {
	obj, ok := a.context.LocateObject(def.AsRef().ReferredPkg, def.AsRef().ReferredType)
	if !ok {
		return "unknown"
	}

	obj.Type.Default = def.Default
	obj.Type.Nullable = def.Nullable

	return a.formatTypeAttribute(obj.Type, false, comments)
}

func (a *attributes) formatEnumAttribute(def ast.Type, comments []string) string {
	scalarDef := def.AsEnum().Values[0].Type
	scalarDef.Scalar.Constraints = formatEnumValuesAsConstraints(def.AsEnum().Values)
	scalarDef.Default = def.Default
	return a.formatScalarAttribute(scalarDef, false, comments)
}

func (a *attributes) formatScalarAttribute(def ast.Type, forceOptional bool, comments []string) string {
	var attrs = map[ast.ScalarKind]attributeFormatter{
		ast.KindString:  {name: "StringAttribute", defImport: "stringdefault", defVal: "StaticString"},
		ast.KindBytes:   {name: "StringAttribute", defImport: "stringdefault", defVal: "StaticString"},
		ast.KindNull:    {name: "StringAttribute", defImport: "stringdefault", defVal: "StaticString"},
		ast.KindBool:    {name: "BoolAttribute", defImport: "booldefault", defVal: "StaticBool"},
		ast.KindInt32:   {name: "Int32Attribute", defImport: "int32default", defVal: "StaticInt32"},
		ast.KindUint32:  {name: "Int32Attribute", defImport: "int32default", defVal: "StaticInt32"},
		ast.KindInt64:   {name: "Int64Attribute", defImport: "int64default", defVal: "StaticInt64"},
		ast.KindUint64:  {name: "Int64Attribute", defImport: "int64default", defVal: "StaticInt64"},
		ast.KindFloat32: {name: "Float32Attribute", defImport: "float32default", defVal: "StaticFloat32"},
		ast.KindFloat64: {name: "Float64Attribute", defImport: "float64default", defVal: "StaticFloat64"},
		ast.KindAny:     {name: "StringAttribute", defImport: "stringdefault", defVal: "StaticString"}, // `any` should be represented as a string holding a JSON payload
		ast.KindInt8:    {name: "NumberAttribute", defImport: "numberdefault", defVal: "StaticBigFloat"},
		ast.KindUint8:   {name: "NumberAttribute", defImport: "numberdefault", defVal: "StaticBigFloat"},
		ast.KindInt16:   {name: "NumberAttribute", defImport: "numberdefault", defVal: "StaticBigFloat"},
		ast.KindUint16:  {name: "NumberAttribute", defImport: "numberdefault", defVal: "StaticBigFloat"},
	}

	attr, ok := attrs[def.AsScalar().ScalarKind]
	if !ok {
		return "unknown"
	}

	var buffer strings.Builder

	fmt.Fprintf(&buffer, "schema.%s{\n", attr.name)

	if forceOptional || def.Nullable || def.Default != nil {
		buffer.WriteString("\tOptional: true,\n")
	} else {
		buffer.WriteString("\tRequired: true,\n")
	}

	if def.Default != nil {
		a.packageMapper(fmt.Sprintf("github.com/hashicorp/terraform-plugin-framework/resource/schema/%s", attr.defImport))
		fmt.Fprintf(&buffer, "\tDefault: %s.%s(%s),\n", attr.defImport, attr.defVal, formatScalar(def.Default))
		buffer.WriteString("\tComputed: true,\n")
	} else if def.AsScalar().Value != nil {
		a.packageMapper(fmt.Sprintf("github.com/hashicorp/terraform-plugin-framework/resource/schema/%s", attr.defImport))
		fmt.Fprintf(&buffer, "\tDefault: %s.%s(%s),\n", attr.defImport, attr.defVal, formatScalar(def.AsScalar().Value))
		buffer.WriteString("\tComputed: true,\n")
	}

	if def.HasHint(ast.HintStringFormatDateTime) {
		a.packageMapper("github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes")
		buffer.WriteString("\tCustomType: timetypes.RFC3339Type{},\n")
	}

	if def.HasHint(ast.HintStringFormatDuration) {
		a.packageMapper("github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes")
		buffer.WriteString("\tCustomType: timetypes.GoDurationType{},\n")
	}

	validator := a.validators.scalarValidator(def.AsScalar().ScalarKind, def.AsScalar().Constraints)
	if validator != "" {
		fmt.Fprintf(&buffer, "\tValidators: %s,\n", validator)
	}

	a.writeDescriptionAttrs(&buffer, comments)

	buffer.WriteString("}")

	return buffer.String()
}

func (a *attributes) writeDescriptionAttrs(out io.Writer, comments []string) {
	if len(comments) == 0 {
		return
	}

	joined := strings.Join(comments, "")

	fmt.Fprintf(out, "\tDescription: %#v,\n", joined)
	fmt.Fprintf(out, "\tMarkdownDescription: %#v,\n", joined)
}

type attributeFormatter struct {
	name      string
	defImport string
	defVal    string
}

type defaultFormatter struct {
	name         string
	defFunc      string
	extraDefFunc string
}

type defaultType string

const MapDefault = "Map"
const ListDefault = "List"

func (a *attributes) parseArrayOrMapDefaults(def ast.Type, defVal any, t defaultType) string {
	var scalarDef = map[ast.ScalarKind]defaultFormatter{
		ast.KindString:  {name: "StringType", defFunc: "StringValue"},
		ast.KindBytes:   {name: "StringType", defFunc: "StringValue"},
		ast.KindNull:    {name: "StringType", defFunc: "StringValue"},
		ast.KindBool:    {name: "BoolType", defFunc: "BoolValue"},
		ast.KindInt32:   {name: "Int32Type", defFunc: "Int32Value"},
		ast.KindUint32:  {name: "Int32Type", defFunc: "Int32Value"},
		ast.KindInt64:   {name: "Int64Type", defFunc: "Int64Value"},
		ast.KindUint64:  {name: "Int64Type", defFunc: "Int64Value"},
		ast.KindFloat32: {name: "Float32Type", defFunc: "Float32Value"},
		ast.KindFloat64: {name: "Float64Type", defFunc: "Float64Value"},
		ast.KindInt8:    {name: "NumberType", defFunc: "NumberValue", extraDefFunc: "big.NewFloat"},
		ast.KindUint8:   {name: "NumberType", defFunc: "NumberValue", extraDefFunc: "big.NewFloat"},
		ast.KindInt16:   {name: "NumberType", defFunc: "NumberValue", extraDefFunc: "big.NewFloat"},
		ast.KindUint16:  {name: "NumberType", defFunc: "NumberValue", extraDefFunc: "big.NewFloat"},
	}

	var enumDef = map[ast.ScalarKind]defaultFormatter{
		ast.KindString: {name: "StringType", defFunc: "StringValue"},
		ast.KindInt64:  {name: "StringType", defFunc: "StringValue"},
	}

	v := defVal.([]any)

	a.packageMapper("github.com/hashicorp/terraform-plugin-framework/attr")

	attrValue := "[]attr.Value"
	mustValue := "ListValueMust"
	if t == MapDefault {
		attrValue = "map[string]attr.Value"
		mustValue = "MapValueMust"
	}

	var buffer strings.Builder
	fmt.Fprintf(&buffer, "types.%s(\n", mustValue)

	switch def.Kind {
	case ast.KindScalar:
		if scalar, ok := scalarDef[def.AsScalar().ScalarKind]; ok {
			fmt.Fprintf(&buffer, "types.%s, %s{\n", scalar.name, attrValue)
			for _, val := range v {
				defaultValue := formatScalar(val)
				if scalar.extraDefFunc != "" {
					defaultValue = fmt.Sprintf("%s(%s)", scalar.extraDefFunc, defaultValue)
				}
				fmt.Fprintf(&buffer, "types.%s(%s),\n", scalar.defFunc, defaultValue)
			}
			buffer.WriteString("},\n")
		}

	case ast.KindRef:
		resolved := a.context.ResolveRefs(def.Ref.AsType())
		if !resolved.IsEnum() {
			return "unknown"
		}

		if enum, ok := enumDef[resolved.AsEnum().Values[0].Type.AsScalar().ScalarKind]; ok {
			fmt.Fprintf(&buffer, "types.%s, %s{\n", enum.name, attrValue)
			for _, val := range v {
				fmt.Fprintf(&buffer, "types.%s(%s),\n", enum.defFunc, formatScalar(val))
			}
			buffer.WriteString("},\n")
		}
	default:
		fmt.Println("Unknown type", def.Kind)
	}

	buffer.WriteString(")")
	return buffer.String()
}

// Blocks generation

func (a *attributes) formatStructBlocks(typeDef ast.Type) string {
	var buffer strings.Builder

	a.packageMapper("github.com/hashicorp/terraform-plugin-framework/resource/schema")

	buffer.WriteString("map[string]schema.Block{\n")
	for _, field := range typeDef.Struct.Fields {
		if field.Type.IsIntersection() {
			continue
		}

		// fields pointing to a struct shouldn't be part of the attributes
		// instead, they should be referenced as a block
		if !field.Type.IsRef() || !a.context.ResolveToStruct(field.Type) {
			continue
		}

		ref := a.context.ResolveRefsChain(field.Type)
		referredSchema := formatObjectName(ref.Ref.ReferredType) + "Schema"

		a.packageMapper("github.com/hashicorp/terraform-plugin-framework/resource/schema")
		fmt.Fprintf(&buffer, "\t\"%s\": schema.SingleNestedBlock{\n", formatTfSDKAttrName(field.Name))
		fmt.Fprintf(&buffer, "\t\tAttributes: %s,\n", referredSchema+".Attributes")
		fmt.Fprintf(&buffer, "\t\tBlocks: %s,\n", referredSchema+".Blocks")

		if typeDef.IsStructGeneratedFromDisjunction() {
			a.packageMapper("github.com/hashicorp/terraform-plugin-framework/schema/validator")

			branchStruct := a.context.ResolveRefs(ref).Struct
			requiredFields := make([]string, 0, len(branchStruct.Fields))
			for _, branchField := range branchStruct.Fields {
				if !branchField.Required || branchField.Type.IsConstantRef() {
					continue
				}

				requiredFields = append(requiredFields, formatScalar(formatTfSDKAttrName(branchField.Name)))
			}

			// TODO: error/warning if `a.cfg.Validators.RequireAttrsWhenPresent` isn't set?
			// Or should we provide an implementation?
			buffer.WriteString("\t\tValidators: []validator.Object{\n")
			fmt.Fprintf(&buffer, "\t\t\t%s(%s),\n", a.cfg.Validators.RequireAttrsWhenPresent, strings.Join(requiredFields, ", "))
			buffer.WriteString("\t\t},\n")
		} else {
			fmt.Fprintf(&buffer, "\t\tValidators: %s,\n", referredSchema+".Validators")
		}

		buffer.WriteString("\t},\n")
	}
	buffer.WriteString("}")

	return buffer.String()
}
