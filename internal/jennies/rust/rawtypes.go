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

		block, err := jenny.formatObject(formatter, imports, schema, object)
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

func (jenny RawTypes) formatObject(formatter *typeFormatter, imports *importMap, schema *ast.Schema, object ast.Object) (string, error) {
	switch {
	case object.Type.IsConcreteScalar():
		return jenny.formatConstant(object), nil
	case object.Type.IsStruct():
		return jenny.formatStruct(formatter, imports, schema, object), nil
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
		enum := jenny.formatUntaggedEnum(formatter, imports, formatTypeName(object.Name), object.Comments, def.Branches)
		// A disjunction object can itself be a dataquery variant (a schema whose root
		// query type is a union of query kinds, e.g. cloudwatch's Request or expr's
		// Expr). It is boxed as Box<dyn Dataquery> by the plugin registry, so it needs
		// the Dataquery trait impl just like a struct query type.
		if impl := jenny.formatVariantImpl(imports, schema, object); impl != "" {
			enum += "\n\n" + impl
		}
		return enum, nil
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

	// A numeric enum with duplicate discriminant values cannot be a Rust C-like
	// enum (each discriminant must be unique: E0081). Grafana's common.ScaleDirection
	// is one such enum (Up=1, Right=1, Down=-1, Left=-1: distinct directions sharing
	// numeric codes). Emit it the way the Go target models all numeric enums: a type
	// alias over the integer plus one associated const per value. Duplicate values
	// are then legal (several consts simply share a number) and serde treats the
	// field as the plain integer.
	if numeric && enumHasDuplicateValues(enum) {
		return jenny.formatNumericEnumAsConsts(object, enum)
	}

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

// enumHasDuplicateValues reports whether two or more of a numeric enum's members
// share the same discriminant value. Such an enum cannot be rendered as a Rust
// C-like enum (E0081: each discriminant must be unique).
func enumHasDuplicateValues(enum ast.EnumType) bool {
	seen := map[string]struct{}{}
	for _, value := range enum.Values {
		key := fmt.Sprintf("%v", value.Value)
		if _, dup := seen[key]; dup {
			return true
		}
		seen[key] = struct{}{}
	}
	return false
}

// formatNumericEnumAsConsts emits a numeric enum as a type alias over its
// underlying integer plus one associated const per member. This is the fallback
// for enums whose discriminants are not unique (which a Rust enum forbids); it
// mirrors the Go target's representation of numeric enums (type + const block)
// and serializes as the plain integer with no serde glue.
func (jenny RawTypes) formatNumericEnumAsConsts(object ast.Object, enum ast.EnumType) string {
	var buffer strings.Builder
	buffer.WriteString(formatComments(object.Comments, ""))

	typeName := formatTypeName(object.Name)
	repr := enumReprType(enum)
	fmt.Fprintf(&buffer, "pub type %s = %s;\n", typeName, repr)

	for _, value := range enum.Values {
		constName := formatConstName(value.Name)
		fmt.Fprintf(&buffer, "\npub const %s: %s = %s;", constName, typeName, formatScalarValue(value.Value, value.Type.AsScalar().ScalarKind))
	}

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

	// An untagged enum cannot derive Default (no variant is implicitly the default),
	// yet a disjunction-typed field of a struct that derives Default needs one (e.g.
	// a v2 variable's `value: VariableValueSingle`). Emit a manual Default selecting
	// the first branch's zero value, the conventional default matching the other
	// targets and consistent with hoisted inline disjunctions.
	if nonNull := nonNullBranches(branches); len(nonNull) > 0 {
		buffer.WriteString("\n\n")
		buffer.WriteString(hoistedEnumDefaultImpl(formatter, name, nonNull[0]))
	}

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

func (jenny RawTypes) formatStruct(formatter *typeFormatter, imports *importMap, schema *ast.Schema, object ast.Object) string {
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

	// Anchor recursive-reference detection on this object while its fields render,
	// so a ref field pointing (transitively) back here is wrapped in Box<T>.
	formatter.currentObjectPkg = schema.Package
	formatter.currentObjectName = object.Name
	defer func() {
		formatter.currentObjectPkg = ""
		formatter.currentObjectName = ""
	}()

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

	if impl := jenny.formatVariantImpl(imports, schema, object); impl != "" {
		buffer.WriteString("\n\n")
		buffer.WriteString(impl)
	}

	return buffer.String()
}

// formatVariantImpl emits the runtime variant trait implementation for an object
// that implements a composable variant. A dataquery variant object gets a
// `cog::variants::Dataquery` impl whose `dataquery_type` returns the schema's
// datasource type identifier (the wire discriminator) and whose
// `dataquery_equals` compares the two queries structurally via their JSON form.
// This mirrors the Go target's `Implements<Variant>Variant()` marker plus the
// `variants.Dataquery` interface (DataqueryType/Equals), giving a composable
// query slot a `Box<dyn Dataquery>` it can resolve and round-trip by
// discriminator. Panelcfg variant objects are plain structs and need no impl
// here: their options/fieldConfig payloads are dispatched by the registry,
// matching the Go panelcfg variant (a bare marker interface).
func (jenny RawTypes) formatVariantImpl(imports *importMap, schema *ast.Schema, object ast.Object) string {
	if !object.Type.ImplementsVariant() || object.Type.IsRef() {
		return ""
	}
	// HintSkipVariantPluginRegistration only suppresses registration in the variant
	// plugin registry (handled in the plugins jenny); the type still needs its
	// Dataquery trait impl so a builder can box it as Box<dyn Dataquery> (e.g.
	// cloudwatch's MetricsQuery/LogsQuery/AnnotationQuery, which are dataquery
	// variants flagged skip-registration but still built through the slot).
	if !object.Type.IsDataqueryVariant() {
		// Panelcfg (and any other) variants carry no rawtypes-level impl.
		return ""
	}

	imports.Add("crate::cog::variants")

	typeName := formatTypeName(object.Name)
	identifier := schema.Metadata.Identifier

	var buffer strings.Builder
	fmt.Fprintf(&buffer, "impl variants::Dataquery for %s {\n", typeName)
	buffer.WriteString("    fn dataquery_type(&self) -> String {\n")
	fmt.Fprintf(&buffer, "        %q.to_string()\n", identifier)
	buffer.WriteString("    }\n\n")
	buffer.WriteString("    fn dataquery_equals(&self, other: &dyn variants::Dataquery) -> bool {\n")
	buffer.WriteString("        match (\n")
	buffer.WriteString("            variants::DataquerySerialize::to_json_value(self),\n")
	buffer.WriteString("            other.to_json_value(),\n")
	buffer.WriteString("        ) {\n")
	buffer.WriteString("            (Ok(a), Ok(b)) => a == b,\n")
	buffer.WriteString("            _ => false,\n")
	buffer.WriteString("        }\n")
	buffer.WriteString("    }\n")
	buffer.WriteString("}")

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
		// formatUntaggedEnum now appends the manual Default impl itself (selecting the
		// first non-null branch), so it is not emitted again here.
		enums = append(enums, jenny.formatUntaggedEnum(formatter, imports, enumName, nil, branches))

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

	// A composable dataquery slot (a `Box<dyn variants::Dataquery>`, directly or in
	// a Vec) cannot derive Deserialize: the concrete query type is only known at
	// runtime from the datasource type discriminator. A serde `deserialize_with`
	// helper reads the raw JSON and dispatches through the variant registry, while
	// serialization goes through the trait object's own Serialize impl. An array
	// slot additionally keeps the empty-collection omit behaviour below.
	if attr, ok := dataquerySlotSerdeAttribute(field); ok {
		attrs = append(attrs, attr)
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

// dataquerySlotSerdeAttribute returns the serde `deserialize_with` attribute for
// a composable dataquery slot field, and whether the field is such a slot. A
// scalar slot (`Box<dyn Dataquery>`) uses the single-value helper; an array slot
// (`Vec<Box<dyn Dataquery>>`) uses the vector helper. Panelcfg slots render as
// serde_json::Value and need no custom (de)serialization, so they are excluded.
func dataquerySlotSerdeAttribute(field ast.StructField) (string, bool) {
	switch {
	case field.Type.IsComposableSlot() && field.Type.AsComposableSlot().Variant == ast.SchemaVariantDataQuery:
		// A nullable scalar slot is Option<Box<dyn Dataquery>> and needs the
		// Option-aware deserializer (the non-optional helper returns a bare Box, which
		// the `?`-based generated visitor cannot coerce into the Option field).
		if field.Type.Nullable {
			return `#[serde(deserialize_with = "crate::cog::variants::deserialize_opt_dataquery")]`, true
		}
		return `#[serde(deserialize_with = "crate::cog::variants::deserialize_dataquery")]`, true
	case field.Type.IsArray() && field.Type.AsArray().ValueType.IsComposableSlot() &&
		field.Type.AsArray().ValueType.AsComposableSlot().Variant == ast.SchemaVariantDataQuery:
		return `#[serde(deserialize_with = "crate::cog::variants::deserialize_dataquery_vec")]`, true
	default:
		return "", false
	}
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
			// A nullable ref field with a struct-literal default is stored as
			// Option<T>; wrap the literal in Some so it type-checks (e.g. heatmap's
			// filter_values: Option<FilterValueRange>).
			if isOptionWrapped(typeDef) {
				return fmt.Sprintf("Some(%s)", literal)
			}
			return literal
		}
		// A ref-to-enum field default resolves to the matching enum variant (or the
		// matching const for a type-aliased enum), wrapped in Some when the field is
		// nullable.
		if literal, ok := enumRefDefaultVariant(formatter, typeDef.AsRef(), typeDef.Default); ok {
			if typeDef.Nullable {
				return fmt.Sprintf("Some(%s)", literal)
			}
			return literal
		}
		// A disjunction-typed ref with a scalar default falls back to the
		// disjunction's own Default (its first branch's zero value); a bare scalar
		// literal would not type-check against the untagged enum.
		if isDisjunctionRef(formatter, typeDef.AsRef()) {
			literal := fmt.Sprintf("%s::default()", formatter.formatRef(typeDef.AsRef()))
			if typeDef.Nullable {
				return fmt.Sprintf("Some(%s)", literal)
			}
			return literal
		}
	}

	if typeDef.Default != nil {
		literal := defaultLiteral(typeDef)
		// A nullable scalar/array/map field with a default is stored as Option<T>; the
		// default literal is the bare T, so wrap it in Some so it type-checks against
		// the Option field (e.g. a nullable bool defaulting to false -> Some(false)).
		// Bare collections (Vec/HashMap) are not Option-wrapped even when nullable, so
		// they are excluded.
		if isOptionWrapped(typeDef) {
			return fmt.Sprintf("Some(%s)", literal)
		}
		return literal
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
	literal := fieldOverrideLiteralInner(formatter, typeDef, value, indent)
	// A nullable scalar/ref/enum field is stored as Option<T>; a bare override
	// literal must be wrapped in Some to type-check (e.g. heatmap's nested
	// FilterValueRange.le: Option<f32> defaulting to a float). Bare collections
	// (Vec/HashMap) are not Option-wrapped even when nullable, so they are excluded
	// via isOptionWrapped. A struct-literal default already returns Some when the
	// ref field is nullable, so it is not double-wrapped (it is not a bare literal).
	if isOptionWrapped(typeDef) && !strings.HasPrefix(literal, "Some(") {
		return fmt.Sprintf("Some(%s)", literal)
	}
	return literal
}

func fieldOverrideLiteralInner(formatter *typeFormatter, typeDef ast.Type, value any, indent string) string {
	if typeDef.IsRef() {
		if literal, ok := refStructLiteralDefault(formatter, typeDef.AsRef(), value, indent); ok {
			return literal
		}
		// A ref to an enum with a scalar default resolves to the matching variant
		// (e.g. a `hide` field defaulting to "dontHide" becomes VariableHide::DontHide),
		// not the bare scalar literal. The enum type is module-qualified and imported
		// via formatRef.
		if literal, ok := enumRefDefaultVariant(formatter, typeDef.AsRef(), value); ok {
			return literal
		}
		// A disjunction-typed ref (emitted as an untagged enum, e.g. a variable
		// option's text/value of type String|Vec<String>) cannot take a bare scalar
		// literal. Fall back to the disjunction's own Default, which selects its first
		// branch's zero value - the conventional default and a valid value.
		if isDisjunctionRef(formatter, typeDef.AsRef()) {
			return fmt.Sprintf("%s::default()", formatter.formatRef(typeDef.AsRef()))
		}
	}

	switch {
	case typeDef.IsArray():
		return arrayDefaultLiteral(typeDef.AsArray(), value)
	case typeDef.IsMap():
		return mapDefaultLiteral(typeDef.AsMap(), value)
	case typeDef.IsScalar():
		return scalarDefaultLiteral(typeDef.AsScalar().ScalarKind, value)
	case typeDef.IsDisjunction():
		// An inline disjunction field (a hoisted untagged enum, e.g. a variable
		// option's text/value of type String|Vec<String>) cannot take a bare scalar
		// default literal. Defer to Default::default(): Rust infers the field's
		// concrete (hoisted) type from the struct literal, and the disjunction's
		// derived Default is its first branch's zero value.
		return "Default::default()"
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
// enumRefDefaultVariant renders a default value for a ref-to-enum field as the
// matching enum variant (Enum::Variant), or the matching associated const for an
// enum emitted as a type alias (duplicate-discriminant enums, see
// formatNumericEnumAsConsts). It reports false when the ref does not resolve to
// an enum or no member matches the default value.
func enumRefDefaultVariant(formatter *typeFormatter, ref ast.RefType, value any) (string, bool) {
	referred, found := formatter.context.LocateObject(ref.ReferredPkg, ref.ReferredType)
	if !found || !referred.Type.IsEnum() {
		return "", false
	}
	enum := referred.Type.AsEnum()
	typeName := formatter.formatRef(ref)
	for _, member := range enum.Values {
		if fmt.Sprintf("%v", member.Value) != fmt.Sprintf("%v", value) {
			continue
		}
		// A numeric enum with duplicate discriminants is emitted as a type alias plus
		// associated consts; the default is the matching const name. A normal enum
		// uses the Enum::Variant form.
		if enumIsNumeric(enum) && enumHasDuplicateValues(enum) {
			return formatConstName(member.Name), true
		}
		return fmt.Sprintf("%s::%s", typeName, formatTypeName(member.Name)), true
	}
	return "", false
}

// isDisjunctionRef reports whether a ref resolves to a disjunction-typed object
// (emitted as an untagged enum). Such a field cannot take a bare scalar default
// literal.
func isDisjunctionRef(formatter *typeFormatter, ref ast.RefType) bool {
	referred, found := formatter.context.LocateObject(ref.ReferredPkg, ref.ReferredType)
	return found && referred.Type.IsDisjunction()
}

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
	// A float-typed field whose default parsed as a whole number (JSON has no
	// integer/float distinction, so 0 arrives as 0.0 -> rendered "0") needs an
	// explicit float literal; a bare integer literal does not coerce to f32/f64
	// (E0308). strconv renders a whole float as e.g. "0" so append ".0" when the
	// rendered form has no decimal point or exponent.
	if kind == ast.KindFloat32 || kind == ast.KindFloat64 {
		rendered := formatValue(value)
		if _, err := strconv.Atoi(rendered); err == nil {
			// A whole-number rendering (e.g. "0") is an integer literal that does not
			// coerce to f32/f64; append ".0" so it is a float literal.
			rendered += ".0"
		}
		return rendered
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
