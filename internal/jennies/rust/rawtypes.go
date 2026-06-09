package rust

import (
	"fmt"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/jennies/common"
	"github.com/grafana/cog/internal/languages"
	"github.com/grafana/cog/internal/tools"
)

// RawTypes emits idiomatic Rust type definitions (structs, type aliases and
// constants) for every schema, one Rust module per schema package.
type RawTypes struct {
	config          Config
	apiRefCollector *common.APIReferenceCollector
}

func (jenny RawTypes) JennyName() string {
	return "RustRawTypes"
}

func (jenny RawTypes) Generate(context languages.Context) (codejen.Files, error) {
	files := make(codejen.Files, 0, len(context.Schemas))

	for _, schema := range context.Schemas {
		output, err := jenny.generateSchema(context, schema)
		if err != nil {
			return nil, err
		}

		filename := filepath.Join("src", "types", formatPackageName(schema.Package)+".rs")
		files = append(files, *codejen.NewFile(filename, output, jenny))
	}

	return files, nil
}

func (jenny RawTypes) generateSchema(context languages.Context, schema *ast.Schema) ([]byte, error) {
	imports := newImportMap()
	formatter := newTypeFormatter(context, imports, schema.Package)

	var body strings.Builder
	blocks := make([]string, 0, schema.Objects.Len())
	var genErr error

	schema.Objects.Iterate(func(_ string, object ast.Object) {
		if genErr != nil {
			return
		}

		block, err := jenny.formatObject(formatter, imports, object)
		if err != nil {
			genErr = err
			return
		}
		if block != "" {
			blocks = append(blocks, block)
		}
	})
	if genErr != nil {
		return nil, genErr
	}

	body.WriteString(strings.Join(blocks, "\n\n"))
	body.WriteString("\n")

	var out strings.Builder
	importStatements := imports.String()
	if importStatements != "" {
		out.WriteString(importStatements)
		out.WriteString("\n\n")
	}
	out.WriteString(body.String())

	return []byte(out.String()), nil
}

func (jenny RawTypes) formatObject(formatter *typeFormatter, imports *importMap, object ast.Object) (string, error) {
	switch {
	case object.Type.IsConcreteScalar():
		return jenny.formatConstant(object), nil
	case object.Type.IsStruct():
		return jenny.formatStruct(formatter, imports, object), nil
	case object.Type.IsEnum():
		return jenny.formatEnum(imports, object), nil
	case object.Type.IsDisjunction():
		// Every disjunction (scalar, ref, or mixed; discriminated or not) is emitted as
		// a single untagged enum. The branch structs already serialize their own
		// constant discriminator field, so an internally tagged enum (#[serde(tag =
		// "...")]) would emit that key twice. Untagged serialization yields the key
		// exactly once, byte-identical to the Go/Python/TS targets, and round-trips
		// because the branches have distinct shapes. Variant identifiers are
		// disambiguated so branches that render to the same Rust type still compile.
		def := object.Type.AsDisjunction()
		return jenny.formatUntaggedEnum(formatter, imports, formatTypeName(object.Name), object.Comments, def.Branches), nil
	case object.Type.IsScalar(), object.Type.IsArray(), object.Type.IsMap(), object.Type.IsRef():
		return jenny.formatTypeAlias(formatter, object), nil
	default:
		return "", fmt.Errorf("rust rawtypes: unsupported top-level object kind %q for %q (Phase 3b+)", object.Type.Kind, object.Name)
	}
}

// formatEnum emits a top-level enum object as an idiomatic Rust enum.
//
// String-valued enums serialize to their string value via #[serde(rename)] when
// the variant identifier differs from the wire value, and derive Eq + Hash so
// they can be used as map keys. Integer-valued enums serialize to their numeric
// value using serde_repr (#[repr(...)] + Serialize_repr/Deserialize_repr) with
// explicit discriminants, and likewise derive Eq + Hash.
func (jenny RawTypes) formatEnum(imports *importMap, object ast.Object) string {
	enum := object.Type.AsEnum()
	numeric := enumIsNumeric(enum)

	var buffer strings.Builder
	buffer.WriteString(formatComments(object.Comments, ""))

	if numeric {
		imports.Add("serde_repr::Serialize_repr")
		imports.Add("serde_repr::Deserialize_repr")
		fmt.Fprintf(&buffer, "#[derive(Serialize_repr, Deserialize_repr, Debug, Clone, Copy, PartialEq, Eq, Hash, Default)]\n")
		fmt.Fprintf(&buffer, "#[repr(%s)]\n", enumReprType(enum))
	} else {
		imports.Add("serde::Serialize")
		imports.Add("serde::Deserialize")
		buffer.WriteString("#[derive(Serialize, Deserialize, Debug, Clone, PartialEq, Eq, Hash, Default)]\n")
	}

	fmt.Fprintf(&buffer, "pub enum %s {\n", formatTypeName(object.Name))
	for i, value := range enum.Values {
		variant := formatTypeName(value.Name)
		// The first variant is the enum's Default. A constant reference to an enum
		// pins a specific variant in the owning struct's Default impl, but a plain
		// enum-typed field still needs a derivable Default; the first declared value
		// is the conventional choice and matches the other targets.
		if i == 0 {
			buffer.WriteString("    #[default]\n")
		}
		if numeric {
			fmt.Fprintf(&buffer, "    %s = %s,\n", variant, formatScalarValue(value.Value, value.Type.AsScalar().ScalarKind))
			continue
		}

		wire := fmt.Sprintf("%v", value.Value)
		if variant != wire {
			fmt.Fprintf(&buffer, "    #[serde(rename = %q)]\n", wire)
		}
		fmt.Fprintf(&buffer, "    %s,\n", variant)
	}
	buffer.WriteString("}")

	return buffer.String()
}

// enumIsNumeric reports whether an enum's values are integer-valued (as opposed
// to string-valued). Enum values share a single scalar kind across all members.
func enumIsNumeric(enum ast.EnumType) bool {
	if len(enum.Values) == 0 {
		return false
	}
	return enum.Values[0].Type.AsScalar().ScalarKind != ast.KindString
}

// enumReprType returns the Rust integer type used for a numeric enum's #[repr].
func enumReprType(enum ast.EnumType) string {
	return formatScalarKind(enum.Values[0].Type.AsScalar().ScalarKind)
}

// formatUntaggedEnum emits a disjunction (scalar, ref, or a mix of scalars,
// arrays, maps and refs; discriminated or not) as a single untagged Rust enum.
// serde tries the variants top-down on deserialization, so branch order is
// preserved, and untagged (de)serialization makes each variant round-trip as the
// bare value, matching the JSON the other targets produce. Variant identifiers
// are disambiguated so two branches rendering to the same Rust type (e.g. two
// distinct refs that resolve to the same scalar) still emit unique, compiling
// identifiers. Eq and Hash are omitted because a branch may carry a float
// payload (directly or transitively through a ref), which is not Eq/Hash in Rust.
func (jenny RawTypes) formatUntaggedEnum(formatter *typeFormatter, imports *importMap, name string, comments []string, branches ast.Types) string {
	imports.Add("serde::Serialize")
	imports.Add("serde::Deserialize")

	var buffer strings.Builder
	buffer.WriteString(formatComments(comments, ""))
	buffer.WriteString("#[derive(Serialize, Deserialize, Debug, Clone, PartialEq)]\n")
	buffer.WriteString("#[serde(untagged)]\n")
	fmt.Fprintf(&buffer, "pub enum %s {\n", name)

	used := make(map[string]struct{}, len(branches))
	for _, branch := range branches {
		inner := formatter.formatType(branch)
		variant := disambiguateVariant(disjunctionVariantName(branch, inner), used)
		fmt.Fprintf(&buffer, "    %s(%s),\n", variant, inner)
	}
	buffer.WriteString("}")

	return buffer.String()
}

// disjunctionVariantName derives a PascalCase variant identifier for a branch of
// an untagged disjunction: a ref uses its referred type name, anything else
// (scalar, array, map) uses the PascalCase form of its rendered Rust type.
func disjunctionVariantName(branch ast.Type, rendered string) string {
	if branch.IsRef() {
		return formatTypeName(branch.AsRef().ReferredType)
	}
	return formatTypeName(rendered)
}

// disambiguateVariant guarantees a unique variant identifier within an enum by
// appending an incrementing numeric suffix on collision.
func disambiguateVariant(name string, used map[string]struct{}) string {
	candidate := name
	for i := 2; ; i++ {
		if _, taken := used[candidate]; !taken {
			used[candidate] = struct{}{}
			return candidate
		}
		candidate = fmt.Sprintf("%s%d", name, i)
	}
}

// formatConstant emits a top-level constant scalar as a Rust `const`.
func (jenny RawTypes) formatConstant(object ast.Object) string {
	var buffer strings.Builder
	buffer.WriteString(formatComments(object.Comments, ""))

	scalar := object.Type.AsScalar()
	constType := constReferenceType(scalar.ScalarKind)

	fmt.Fprintf(
		&buffer,
		"pub const %s: %s = %s;",
		formatConstName(object.Name),
		constType,
		formatScalarValue(scalar.Value, scalar.ScalarKind),
	)

	return buffer.String()
}

// formatTypeAlias emits a top-level scalar object as a Rust type alias.
func (jenny RawTypes) formatTypeAlias(formatter *typeFormatter, object ast.Object) string {
	var buffer strings.Builder
	buffer.WriteString(formatComments(object.Comments, ""))
	fmt.Fprintf(
		&buffer,
		"pub type %s = %s;",
		formatTypeName(object.Name),
		formatter.formatType(object.Type),
	)

	return buffer.String()
}

func (jenny RawTypes) formatStruct(formatter *typeFormatter, imports *importMap, object ast.Object) string {
	imports.Add("serde::Serialize")
	imports.Add("serde::Deserialize")

	fields := object.Type.AsStruct().Fields

	// An inline (anonymous) multi-arm non-null disjunction field has no named type
	// to reference, so it is hoisted here into a named untagged enum emitted in this
	// same module, and the field is rewritten to reference it. This mirrors the
	// Python target (typing.Union[...]) and the Go target (a named union type)
	// rather than degrading to the type-erased serde_json::Value. The synthesized
	// enums are emitted after the struct (and its Default impl / default fns) so the
	// struct stays first, matching how every other top-level object is ordered.
	fields, hoistedEnums := jenny.hoistInlineDisjunctions(formatter, imports, object.Name, fields)

	var buffer strings.Builder
	buffer.WriteString(formatComments(object.Comments, ""))

	buffer.WriteString(structDerives(fields))
	buffer.WriteString("\n")

	if len(fields) == 0 {
		// rustfmt renders a field-less struct as `pub struct X {}` on one line.
		fmt.Fprintf(&buffer, "pub struct %s {}", formatTypeName(object.Name))
	} else {
		fmt.Fprintf(&buffer, "pub struct %s {\n", formatTypeName(object.Name))
		for i, field := range fields {
			if i != 0 {
				buffer.WriteString("\n")
			}
			buffer.WriteString(jenny.formatStructField(formatter, object.Name, field))
		}
		buffer.WriteString("}")
	}

	if defaultImpl := jenny.formatDefaultImpl(formatter, object); defaultImpl != "" {
		buffer.WriteString("\n\n")
		buffer.WriteString(defaultImpl)
	}

	for _, field := range fields {
		if fn := jenny.formatFieldDefaultFn(formatter, object.Name, field); fn != "" {
			buffer.WriteString("\n\n")
			buffer.WriteString(fn)
		}
	}

	for _, enum := range hoistedEnums {
		buffer.WriteString("\n\n")
		buffer.WriteString(enum)
	}

	return buffer.String()
}

// hoistInlineDisjunctions rewrites every struct field whose type is an inline
// (anonymous) multi-arm non-null disjunction into a reference to a freshly named
// untagged enum, and returns the rewritten field list alongside the rendered enum
// blocks. The synthesized enum is named <StructName><FieldName> in PascalCase,
// which is deterministic and unique within the schema (a struct cannot have two
// fields of the same name). A null branch never survives to this point for a
// two-arm disjunction (DisjunctionWithNullToOptional rewrites string|null to
// Option<String>); should a null branch reach here as part of a larger
// disjunction it is dropped from the enum because untagged serde already models
// JSON null via the surrounding field, mirroring the other targets.
func (jenny RawTypes) hoistInlineDisjunctions(formatter *typeFormatter, imports *importMap, structName string, fields []ast.StructField) ([]ast.StructField, []string) {
	var enums []string

	rewritten := make([]ast.StructField, len(fields))
	copy(rewritten, fields)

	for i, field := range fields {
		if !field.Type.IsDisjunction() {
			continue
		}

		branches := nonNullBranches(field.Type.AsDisjunction().Branches)
		if len(branches) < 2 {
			// A single surviving branch is not a union; leave it for the generic
			// formatter (it renders as that branch's type or serde_json::Value).
			continue
		}

		enumName := formatTypeName(structName + field.Name)
		enums = append(enums, jenny.formatUntaggedEnum(formatter, imports, enumName, nil, branches))
		enums = append(enums, hoistedEnumDefaultImpl(formatter, enumName, branches[0]))

		rewritten[i].Type = ast.NewRef(formatter.rawPackage, enumName)
		rewritten[i].Type.Nullable = field.Type.Nullable
	}

	return rewritten, enums
}

// hoistedEnumDefaultImpl emits a manual Default impl for a synthesized untagged
// enum. derive(Default) cannot mark a tuple variant as #[default], so the impl is
// written by hand, defaulting to the first branch wrapped in its variant. The
// first branch is the conventional Default and matches what the other targets do
// (Python defaults a Union field to the first arm's zero value).
func hoistedEnumDefaultImpl(formatter *typeFormatter, enumName string, firstBranch ast.Type) string {
	inner := formatter.formatType(firstBranch)
	variant := disjunctionVariantName(firstBranch, inner)

	var buffer strings.Builder
	fmt.Fprintf(&buffer, "impl Default for %s {\n", enumName)
	buffer.WriteString("    fn default() -> Self {\n")
	fmt.Fprintf(&buffer, "        Self::%s(Default::default())\n", variant)
	buffer.WriteString("    }\n")
	buffer.WriteString("}")

	return buffer.String()
}

// nonNullBranches returns the disjunction branches that are not the bare null
// type. Untagged serde models JSON null through the surrounding optional field,
// so a null branch needs no enum variant.
func nonNullBranches(branches ast.Types) ast.Types {
	filtered := make(ast.Types, 0, len(branches))
	for _, branch := range branches {
		if branch.IsNull() {
			continue
		}
		filtered = append(filtered, branch)
	}
	return filtered
}

// fieldNeedsConstantDefaultFn reports whether a field is a required constant or
// constant reference whose pinned value must be restored by serde when the key
// is absent from the input (most notably when the struct is a branch of an
// internally tagged disjunction, where serde consumes the discriminator key for
// the tag). serde's bare #[serde(default)] would fall back to the field type's
// Default (e.g. an empty String) rather than the pinned constant, so a dedicated
// default function is generated instead.
func fieldNeedsConstantDefaultFn(field ast.StructField) bool {
	if !field.Required {
		return false
	}
	// A nullable constant (rendered as Option<T>) takes the Option skip/default
	// serde attribute instead: an absent key deserializes to None and the pinned
	// constant lives in the struct's Default impl. Generating a const default
	// function here would both shadow that attribute and return the bare T rather
	// than Option<T>, so Option-wrapped fields are excluded.
	if isOptionWrapped(field.Type) {
		return false
	}
	return field.Type.IsConcreteScalar() || field.Type.IsConstantRef()
}

// fieldDefaultFnName returns the deterministic, unique name of the serde default
// function generated for a constant field.
func fieldDefaultFnName(structName string, field ast.StructField) string {
	return fmt.Sprintf("default_%s_%s", tools.SnakeCase(structName), tools.SnakeCase(field.Name))
}

// formatFieldDefaultFn renders the serde default function for a constant field,
// returning the field's pinned constant value. It returns the empty string for
// fields that need no such function.
func (jenny RawTypes) formatFieldDefaultFn(formatter *typeFormatter, structName string, field ast.StructField) string {
	if !fieldNeedsConstantDefaultFn(field) {
		return ""
	}

	var value string
	if field.Type.IsConstantRef() {
		value = constantRefDefault(formatter, field.Type.AsConstantRef())
	} else {
		value = constScalarLiteral(field.Type.AsScalar())
	}

	var buffer strings.Builder
	fmt.Fprintf(&buffer, "fn %s() -> %s {\n", fieldDefaultFnName(structName, field), formatter.formatType(field.Type))
	fmt.Fprintf(&buffer, "    %s\n", value)
	buffer.WriteString("}")
	return buffer.String()
}

func (jenny RawTypes) formatStructField(formatter *typeFormatter, structName string, field ast.StructField) string {
	var buffer strings.Builder

	buffer.WriteString(formatComments(field.Comments, "    "))

	for _, attr := range fieldSerdeAttributes(structName, field) {
		fmt.Fprintf(&buffer, "    %s\n", attr)
	}

	fmt.Fprintf(&buffer, "    pub %s: %s,\n", formatFieldName(field.Name), formatter.formatType(field.Type))

	return buffer.String()
}

// fieldSerdeAttributes returns the serde attribute lines for a struct field:
// rename for non-snake-case keys and skip/default handling for optionals.
func fieldSerdeAttributes(structName string, field ast.StructField) []string {
	var attrs []string

	if rustFieldNeedsRename(field.Name) {
		attrs = append(attrs, fmt.Sprintf("#[serde(rename = %q)]", field.Name))
	}

	// Collections (arrays and maps) are rendered as bare Vec/HashMap rather than
	// Option<T> (see typeFormatter.formatType). The Go target marks every
	// not-required field ",omitempty", which for an empty collection omits the key
	// from the JSON entirely. To stay wire-compatible with the Go/TS/Python SDKs and
	// Grafana, a not-required collection field is given default + skip-when-empty so
	// an empty Vec/HashMap is omitted and an absent key deserializes to empty.
	switch {
	case !field.Required && field.Type.IsArray():
		attrs = append(attrs, `#[serde(default, skip_serializing_if = "Vec::is_empty")]`)
	case !field.Required && field.Type.IsMap():
		attrs = append(attrs, `#[serde(default, skip_serializing_if = "HashMap::is_empty")]`)
	case isOptionWrapped(field.Type):
		// A nullable scalar/struct/ref field is rendered as Option<T> and so gets the
		// Option skip/default attribute.
		attrs = append(attrs, `#[serde(default, skip_serializing_if = "Option::is_none")]`)
	case fieldNeedsConstantDefaultFn(field):
		// A required constant (or constant reference) field carries a value pinned by
		// the struct's Default impl. A dedicated default function restores the pinned
		// constant on deserialization (bare #[serde(default)] would yield the field
		// type's zero value instead) while the value is still serialized. This keeps
		// the discriminator present on both standalone use and as an untagged
		// disjunction branch, matching the Go/Python/TS wire format.
		attrs = append(attrs, fmt.Sprintf("#[serde(default = %q)]", fieldDefaultFnName(structName, field)))
	}

	return attrs
}

// structDerives returns the #[derive(...)] line for a struct. Default is derived
// only when every field's default can be produced by #[derive(Default)] (i.e. no
// field carries an explicit non-zero default and no constant field). Otherwise a
// manual Default impl is emitted separately.
func structDerives(fields []ast.StructField) string {
	derives := []string{"Serialize", "Deserialize", "Debug", "Clone", "PartialEq"}

	if !needsManualDefault(fields) {
		derives = append(derives, "Default")
	}

	return fmt.Sprintf("#[derive(%s)]", strings.Join(derives, ", "))
}

func needsManualDefault(fields []ast.StructField) bool {
	for _, field := range fields {
		if field.Type.IsConcreteScalar() {
			return true
		}
		if field.Type.IsConstantRef() {
			return true
		}
		if field.Type.Default != nil {
			return true
		}
	}
	return false
}

func (jenny RawTypes) formatDefaultImpl(formatter *typeFormatter, object ast.Object) string {
	fields := object.Type.AsStruct().Fields
	if !needsManualDefault(fields) {
		return ""
	}

	var buffer strings.Builder
	fmt.Fprintf(&buffer, "impl Default for %s {\n", formatTypeName(object.Name))
	buffer.WriteString("    fn default() -> Self {\n")
	buffer.WriteString("        Self {\n")

	for _, field := range fields {
		fmt.Fprintf(&buffer, "            %s: %s,\n", formatFieldName(field.Name), defaultExpression(formatter, field.Type, "            "))
	}

	buffer.WriteString("        }\n")
	buffer.WriteString("    }\n")
	buffer.WriteString("}")

	return buffer.String()
}

// defaultExpression returns the Rust expression used to initialise a field in a
// manual Default impl. indent is the leading whitespace of the line the
// expression begins on, so any multi-line struct literal it produces nests its
// inner lines correctly and matches rustfmt's layout.
func defaultExpression(formatter *typeFormatter, typeDef ast.Type, indent string) string {
	if typeDef.IsConstantRef() {
		pinned := constantRefDefault(formatter, typeDef.AsConstantRef())
		// A nullable constant reference is stored as Option<T>; the pinned constant
		// is wrapped in Some so the field still round-trips to its constant value
		// rather than to None.
		if typeDef.Nullable {
			return fmt.Sprintf("Some(%s)", pinned)
		}
		return pinned
	}

	if typeDef.IsConcreteScalar() {
		return constScalarLiteral(typeDef.AsScalar())
	}

	// A ref-typed field carrying a struct-literal default (a map of overridden
	// fields) delegates to the referred struct's Default impl, overriding only the
	// named fields. Anything not named falls back via `..Default::default()`. The
	// referred object is resolved so each override is rendered with the correct
	// field type (snake-cased name, scalar/struct literal value).
	if typeDef.IsRef() && typeDef.Default != nil {
		if literal, ok := refStructLiteralDefault(formatter, typeDef.AsRef(), typeDef.Default, indent); ok {
			return literal
		}
	}

	if typeDef.Default != nil {
		return defaultLiteral(typeDef)
	}

	return "Default::default()"
}

// refStructLiteralDefault renders a struct-literal Default expression for a
// ref-typed field whose IR default is a map of overridden field values. It
// resolves the referred struct so each override is rendered with the field's
// real Rust type and snake-cased identifier, and closes the literal with
// `..Default::default()` so unspecified fields inherit the struct's own Default.
// It reports false when the default is not a map or the referred type cannot be
// resolved to a struct, letting the caller fall back to the generic path. The
// literal is rendered multi-line (one field per line) to match rustfmt, with
// inner lines indented one level deeper than indent.
func refStructLiteralDefault(formatter *typeFormatter, ref ast.RefType, def any, indent string) (string, bool) {
	overrides, ok := def.(map[string]any)
	if !ok {
		return "", false
	}

	referredObject, found := formatter.context.LocateObject(ref.ReferredPkg, ref.ReferredType)
	if !found || !referredObject.Type.IsStruct() {
		return "", false
	}

	typeName := formatter.formatRef(ref)
	if len(overrides) == 0 {
		return fmt.Sprintf("%s::default()", typeName), true
	}

	// Walk the struct's declared fields in order so the emitted literal is
	// deterministic and only includes fields the default actually overrides.
	structFields := referredObject.Type.AsStruct().Fields
	inner := indent + "    "
	var parts []string
	for _, field := range structFields {
		value, present := overrides[field.Name]
		if !present {
			continue
		}
		parts = append(parts, fmt.Sprintf(
			"%s%s: %s,",
			inner,
			formatFieldName(field.Name),
			fieldOverrideLiteral(formatter, field.Type, value, inner),
		))
	}

	if len(parts) == 0 {
		return fmt.Sprintf("%s::default()", typeName), true
	}

	// A `..Default::default()` rest-update is only emitted when the default leaves
	// at least one field unspecified; clippy's needless_update lint rejects it when
	// every field is already named.
	if len(parts) != len(structFields) {
		parts = append(parts, inner+"..Default::default()")
	}

	return fmt.Sprintf("%s {\n%s\n%s}", typeName, strings.Join(parts, "\n"), indent), true
}

// fieldOverrideLiteral renders a single overridden field value inside a struct
// literal default. A nested ref-with-map value recurses to produce a nested
// struct literal (indented relative to indent); scalars, arrays and maps use the
// existing default-literal renderers; a String scalar yields an owned String.
func fieldOverrideLiteral(formatter *typeFormatter, typeDef ast.Type, value any, indent string) string {
	if typeDef.IsRef() {
		if literal, ok := refStructLiteralDefault(formatter, typeDef.AsRef(), value, indent); ok {
			return literal
		}
	}

	switch {
	case typeDef.IsArray():
		return arrayDefaultLiteral(typeDef.AsArray(), value)
	case typeDef.IsMap():
		return mapDefaultLiteral(typeDef.AsMap(), value)
	case typeDef.IsScalar():
		return scalarDefaultLiteral(typeDef.AsScalar().ScalarKind, value)
	default:
		return formatValue(value)
	}
}

// constantRefDefault renders the pinned constant value of a constant reference.
// A reference to an enum resolves to the matching enum variant
// (Enum::Variant); a reference to a scalar constant resolves to that scalar's
// literal (a string becomes an owned String).
func constantRefDefault(formatter *typeFormatter, ref ast.ConstantReferenceType) string {
	referredObject, found := formatter.context.LocateObject(ref.ReferredPkg, ref.ReferredType)
	if found && referredObject.Type.IsEnum() {
		typeName := formatter.formatRef(ast.RefType{ReferredPkg: ref.ReferredPkg, ReferredType: ref.ReferredType})
		for _, value := range referredObject.Type.AsEnum().Values {
			if fmt.Sprintf("%v", value.Value) == fmt.Sprintf("%v", ref.ReferenceValue) {
				return fmt.Sprintf("%s::%s", typeName, formatTypeName(value.Name))
			}
		}
		return fmt.Sprintf("%s::default()", typeName)
	}

	if found && referredObject.Type.IsScalar() {
		scalar := referredObject.Type.AsScalar()
		if scalar.ScalarKind == ast.KindString {
			return fmt.Sprintf("%s.to_string()", formatValue(ref.ReferenceValue))
		}
		return formatScalarValue(ref.ReferenceValue, scalar.ScalarKind)
	}

	return formatValue(ref.ReferenceValue)
}

// defaultLiteral renders an explicit IR default into a Rust expression, using
// the field's type to pick the correct literal form.
func defaultLiteral(typeDef ast.Type) string {
	switch {
	case typeDef.IsArray():
		return arrayDefaultLiteral(typeDef.AsArray(), typeDef.Default)
	case typeDef.IsMap():
		return mapDefaultLiteral(typeDef.AsMap(), typeDef.Default)
	case typeDef.IsScalar():
		return scalarDefaultLiteral(typeDef.AsScalar().ScalarKind, typeDef.Default)
	default:
		return formatValue(typeDef.Default)
	}
}

// scalarDefaultLiteral renders a scalar default. The `any` kind maps to
// serde_json::Value, whose JSON-object default (the only shape seen so far, an
// empty `{}`) is an empty JSON object.
func scalarDefaultLiteral(kind ast.ScalarKind, value any) string {
	switch kind {
	case ast.KindString:
		return fmt.Sprintf("%s.to_string()", formatValue(value))
	case ast.KindAny:
		return "serde_json::Value::Object(serde_json::Map::new())"
	default:
		return formatScalarValue(value, kind)
	}
}

// arrayDefaultLiteral renders an array default. An empty default produces
// Vec::new(); a non-empty default produces a vec![...] literal whose elements
// are rendered from the array's value type.
func arrayDefaultLiteral(def ast.ArrayType, value any) string {
	elements, ok := value.([]any)
	if !ok || len(elements) == 0 {
		return "Vec::new()"
	}

	rendered := make([]string, 0, len(elements))
	for _, element := range elements {
		rendered = append(rendered, elementDefaultLiteral(def.ValueType, element))
	}

	return fmt.Sprintf("vec![%s]", strings.Join(rendered, ", "))
}

// mapDefaultLiteral renders a map default. An empty default produces
// HashMap::new(); a non-empty default produces a HashMap::from([...]) literal.
func mapDefaultLiteral(def ast.MapType, value any) string {
	entries, ok := value.(map[string]any)
	if !ok || len(entries) == 0 {
		return "HashMap::new()"
	}

	// Sort keys so the emitted literal is deterministic.
	keys := make([]string, 0, len(entries))
	for key := range entries {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	pairs := make([]string, 0, len(entries))
	for _, key := range keys {
		pairs = append(pairs, fmt.Sprintf(
			"(%s.to_string(), %s)",
			strconv.Quote(key),
			elementDefaultLiteral(def.ValueType, entries[key]),
		))
	}

	return fmt.Sprintf("HashMap::from([%s])", strings.Join(pairs, ", "))
}

// elementDefaultLiteral renders a single element of an array or map default,
// dispatching on the element's type.
func elementDefaultLiteral(typeDef ast.Type, value any) string {
	switch {
	case typeDef.IsArray():
		return arrayDefaultLiteral(typeDef.AsArray(), value)
	case typeDef.IsMap():
		return mapDefaultLiteral(typeDef.AsMap(), value)
	case typeDef.IsScalar():
		return scalarDefaultLiteral(typeDef.AsScalar().ScalarKind, value)
	default:
		return formatValue(value)
	}
}

// formatScalarValue renders a value according to the target scalar kind. JSON
// unmarshalling decodes every number as float64, so integer kinds must be
// rendered without a fractional part to produce valid Rust integer literals.
func formatScalarValue(value any, kind ast.ScalarKind) string {
	if isIntegerScalarKind(kind) {
		if f, ok := value.(float64); ok {
			return strconv.FormatInt(int64(f), 10)
		}
	}
	return formatValue(value)
}

func isIntegerScalarKind(kind ast.ScalarKind) bool {
	switch kind {
	case ast.KindUint8, ast.KindUint16, ast.KindUint32, ast.KindUint64,
		ast.KindInt8, ast.KindInt16, ast.KindInt32, ast.KindInt64:
		return true
	default:
		return false
	}
}

func formatComments(comments []string, indent string) string {
	if len(comments) == 0 {
		return ""
	}

	var buffer strings.Builder
	for _, line := range comments {
		buffer.WriteString(strings.TrimRight(fmt.Sprintf("%s/// %s", indent, line), " "))
		buffer.WriteString("\n")
	}

	return buffer.String()
}

// constReferenceType returns the Rust type used for a top-level constant of the
// given scalar kind. String constants use &str rather than String so the
// constant can live in a `const` item; in a `const NAME: &str = ...` the
// reference desugars to &'static str.
func constReferenceType(kind ast.ScalarKind) string {
	if kind == ast.KindString {
		return "&str"
	}
	return formatScalarKind(kind)
}

// constScalarLiteral renders a constant scalar field's value for use in a Default
// impl, where String fields need an owned String.
func constScalarLiteral(scalar ast.ScalarType) string {
	if scalar.ScalarKind == ast.KindString {
		return fmt.Sprintf("%s.to_string()", formatValue(scalar.Value))
	}
	return formatScalarValue(scalar.Value, scalar.ScalarKind)
}
