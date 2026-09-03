package rust

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/jennies/common"
	"github.com/grafana/cog/internal/languages"
	"github.com/grafana/cog/internal/tools"
)

// Converter emits, for every builder, a Rust function that accepts an
// already-built object and returns the Rust source code (as a string) that
// reconstructs the object through that builder. This is the reverse direction
// of the Builder jenny: dashboard JSON in, builder calls out. The semantics
// mirror the Go and PHP converter jennies: the languages.ConverterGenerator IR
// drives what gets emitted, and the emitted function appends one builder-call
// string per applicable option, guarded so that absent or default-valued
// fields produce no call.
//
// All converters for builders sharing a package land in a single
// src/converters/<package>.rs module, matching how the Builder jenny groups
// its output.
type Converter struct {
	config          Config
	apiRefCollector *common.APIReferenceCollector
}

func (jenny Converter) JennyName() string {
	return "RustConverter"
}

// converterNullableConfig is the nullability model the converter generator
// must use for Rust. It intentionally differs from Language.NullableKinds():
// the generator uses it to decide where "is the value present?" guards are
// needed, so it must mirror what the emitted Rust types actually wrap in
// Option (isOptionWrapped: any type with Nullable set, except collections) and
// which types carry their own empty representation (bare Vec/HashMap, guarded
// by emptiness). Refs and structs are only Option-wrapped when nullable, so
// listing their kinds here (as NullableKinds does for the builder nil-check
// pass) would generate guards on expressions that are not Options.
func converterNullableConfig() languages.NullableConfig {
	return languages.NullableConfig{
		Kinds: []ast.Kind{ast.KindMap, ast.KindArray},
		// `any` fields mirror the Go target: assignment paths into an `any` field
		// (a serde_json::Value, Option-wrapped when the field is nullable) need a
		// presence guard before they are traversed. Path rendering resolves the
		// actual field type to decide whether an unwrap is required, and a guard
		// on a non-optional `any` renders as vacuous and is dropped.
		AnyIsNullable: true,
	}
}

func (jenny Converter) Generate(context languages.Context) (codejen.Files, error) {
	type packageGroup struct {
		pkg     string
		imports *importMap
		bodies  []string
	}
	groups := make(map[string]*packageGroup)
	order := make([]string, 0)

	for _, builder := range context.Builders {
		group, ok := groups[builder.Package]
		if !ok {
			group = &packageGroup{pkg: builder.Package, imports: newImportMap()}
			groups[builder.Package] = group
			order = append(order, builder.Package)
		}

		body, err := jenny.generateConverter(context, builder, group.imports)
		if err != nil {
			return nil, err
		}
		group.bodies = append(group.bodies, body)
	}

	files := make(codejen.Files, 0, len(order))
	for _, pkg := range order {
		group := groups[pkg]

		var out strings.Builder
		importStatements := group.imports.String()
		if importStatements != "" {
			out.WriteString(importStatements)
			out.WriteString("\n\n")
		}
		out.WriteString(strings.Join(group.bodies, "\n"))

		filename := filepath.Join("src", "converters", formatPackageName(pkg)+".rs")
		files = append(files, *codejen.NewFile(filename, []byte(out.String()), jenny))
	}

	return files, nil
}

func (jenny Converter) generateConverter(context languages.Context, builder ast.Builder, imports *importMap) (string, error) {
	converter := languages.NewConverterGenerator(converterNullableConfig(), context.ConverterConfig).FromBuilder(context, builder)

	formatter := newTypeFormatter(context, imports, builder.Package)
	// The converters module is distinct from the types module, so even
	// same-package type refs (enum constants in guards) must be imported.
	formatter.importSamePackageRefs = true

	inputPkg := formatPackageName(converter.Input.TypeRef.ReferredPkg)
	imports.Add("crate::types::" + inputPkg)

	emitter := &converterEmitter{
		context:    context,
		converter:  converter,
		imports:    imports,
		formatter:  formatter,
		currentPkg: builder.Package,
	}

	converterName := formatConverterName(converter.BuilderName)
	inputType := fmt.Sprintf("%s::%s", inputPkg, formatTypeName(converter.Input.TypeRef.ReferredType))

	jenny.apiRefCollector.RegisterFunction(builder.Package, common.FunctionReference{
		Name: converterName,
		Arguments: []common.ArgumentReference{
			{Name: "input", Type: "&" + inputType},
		},
		Comments: []string{
			fmt.Sprintf("%[1]s accepts a `%[2]s` object and generates the Rust code to build this object using builders.", converterName, formatTypeName(converter.BuilderName)),
		},
		Return: "String",
	})

	return emitter.emitFunction(converterName, inputType), nil
}

// formatConverterName derives the converter function name from a builder name,
// e.g. "SomeStruct" becomes "some_struct_converter".
func formatConverterName(builderName string) string {
	// The suffix is appended before keyword escaping: "struct" must become
	// "struct_converter", not "r#struct" + "_converter".
	return escapeRustKeyword(tools.SnakeCase(builderName) + "_converter")
}

// converterEmitter holds the state shared by every emission helper while a
// single converter function is rendered.
type converterEmitter struct {
	context    languages.Context
	converter  languages.Converter
	imports    *importMap
	formatter  *typeFormatter
	currentPkg string
}

// pathExpr is a rendered IR path: a Rust place expression plus the metadata
// needed to format guards and values against it.
type pathExpr struct {
	expr string
	// def is the type of the last path chunk.
	def ast.Type
	// json reports that the path traversed an `any` field with a type hint, so
	// expr is a serde_json::Value indexing chain rather than typed field access.
	json bool
	// isRef reports that expr already evaluates to a reference (an unwrapped
	// Option or a loop variable), so it must not be re-borrowed with `&`.
	isRef bool
}

// renderPath renders an IR path into a Rust place expression. Intermediate
// Option-wrapped chunks are unwrapped with `.as_ref().unwrap()`; the guards the
// converter generator emits ahead of every access make the unwrap safe, the
// same way the Go converter's pointer dereferences rely on its nil-check
// guards. When unwrapLast is set, a trailing Option chunk is unwrapped too.
//
// A chunk whose type is `any` and that carries a type hint switches the
// traversal to JSON mode: the serialized value is indexed by wire key from
// there on, mirroring the Go converter's `.(*Hint)` type assertion.
func (emitter *converterEmitter) renderPath(path ast.Path, unwrapLast bool) pathExpr {
	var out strings.Builder
	result := pathExpr{}

	// The path chunks carry the types the builder assignments were declared
	// with, which do not always reflect the emitted field's nullability (an
	// assignment into a nullable `any` field carries a non-nullable chunk, for
	// example). The parent type is tracked across the traversal so each chunk
	// can be resolved against the actual struct field, which is what decides
	// whether an Option unwrap is required.
	var parent ast.Type

	for i, chunk := range path {
		last := i == len(path)-1
		def := emitter.resolveChunkType(parent, chunk, i == 0)
		result.def = def

		switch {
		case result.json:
			fmt.Fprintf(&out, "[%q]", chunk.Identifier)
			result.isRef = true
		case i == 0:
			out.WriteString(formatArgName(chunk.Identifier))
			// Both the `input` parameter and the loop variables introduced by
			// repeated mappings are references.
			result.isRef = true
		default:
			out.WriteString(".")
			out.WriteString(formatFieldName(chunk.Identifier))
			result.isRef = false
		}

		if chunk.Index != nil {
			out.WriteString("[")
			if chunk.Index.Constant != nil {
				out.WriteString(formatValue(chunk.Index.Constant))
			} else {
				out.WriteString(formatArgName(chunk.Index.Argument.Name))
			}
			out.WriteString("]")
			result.isRef = false
			def = collectionValueType(emitter.context.ResolveRefs(def))
			result.def = def
		}

		if !result.json && isOptionWrapped(def) && (!last || unwrapLast) {
			out.WriteString(".as_ref().unwrap()")
			result.isRef = true
		}

		if !result.json && def.IsAny() && chunk.TypeHint != nil && !last {
			result.json = true
		}

		parent = def
	}

	result.expr = out.String()
	return result
}

// resolveChunkType returns the authoritative type of a path chunk. A non-root
// chunk names a field of its parent type; that field's declared type (which
// carries the nullability the emitted Rust struct actually has) wins over the
// type recorded on the chunk. Root chunks (the input parameter, loop
// variables) keep their recorded type.
func (emitter *converterEmitter) resolveChunkType(parent ast.Type, chunk ast.PathItem, isRoot bool) ast.Type {
	if isRoot {
		if chunk.Root && chunk.Identifier == emitter.converter.Input.ArgName {
			return emitter.converter.Input.TypeRef.AsType()
		}
		return chunk.Type
	}

	resolved := emitter.context.ResolveRefs(parent)
	if resolved.IsStruct() {
		if field, found := resolved.AsStruct().FieldByName(chunk.Identifier); found {
			return field.Type
		}
	}

	return chunk.Type
}

// collectionValueType returns the element type of an indexed collection.
func collectionValueType(def ast.Type) ast.Type {
	switch {
	case def.IsMap():
		return def.AsMap().ValueType
	case def.IsArray():
		return def.AsArray().ValueType
	default:
		return def
	}
}

// refExpr renders a path as a reference expression (&T), suitable as a
// converter-call or runtime-helper argument.
func (emitter *converterEmitter) refExpr(path ast.Path) string {
	rendered := emitter.renderPath(path, true)
	if rendered.isRef {
		return rendered.expr
	}
	return "&" + rendered.expr
}

// renderGuards renders mapping guards as a `&&`-joined condition, de-duplicated
// by rendered form (a bare collection produces the same emptiness check for
// both its not-nil and its minLength guard).
func (emitter *converterEmitter) renderGuards(guards []languages.MappingGuard) string {
	seen := make(map[string]struct{})
	rendered := make([]string, 0, len(guards))
	for _, guard := range guards {
		condition := emitter.renderGuard(guard)
		if condition == "" {
			// Vacuous guard (a presence check on a value that always exists).
			continue
		}
		if _, dup := seen[condition]; dup {
			continue
		}
		seen[condition] = struct{}{}
		rendered = append(rendered, condition)
	}
	return strings.Join(rendered, " && ")
}

func (emitter *converterEmitter) renderGuard(guard languages.MappingGuard) string {
	raw := emitter.renderPath(guard.Path, false)
	unwrapped := emitter.renderPath(guard.Path, true)

	if raw.json {
		return emitter.renderJSONGuard(guard, raw)
	}

	switch guard.Op {
	case ast.NotEqualOp, ast.EqualOp:
		if guard.Value == nil {
			return emitter.renderNilGuard(guard.Op, raw)
		}
		return emitter.renderComparisonGuard(guard, raw, unwrapped)
	case ast.MinLengthOp:
		if length, ok := guardLength(guard.Value); ok && length == 1 {
			return fmt.Sprintf("!%s.is_empty()", unwrapped.expr)
		}
		return fmt.Sprintf("%s.len() >= %s", unwrapped.expr, formatValue(guard.Value))
	case ast.MaxLengthOp:
		return fmt.Sprintf("%s.len() <= %s", unwrapped.expr, formatValue(guard.Value))
	default:
		return fmt.Sprintf("%s %s %s", unwrapped.expr, guard.Op, formatValue(guard.Value))
	}
}

// renderNilGuard renders a presence check. A bare collection models absence by
// emptiness; an Option-wrapped value by is_some()/is_none().
func (emitter *converterEmitter) renderNilGuard(op ast.Op, raw pathExpr) string {
	if isBareCollection(raw.def) {
		if op == ast.EqualOp {
			return fmt.Sprintf("%s.is_empty()", raw.expr)
		}
		return fmt.Sprintf("!%s.is_empty()", raw.expr)
	}

	if isOptionWrapped(raw.def) {
		if op == ast.EqualOp {
			return fmt.Sprintf("%s.is_none()", raw.expr)
		}
		return fmt.Sprintf("%s.is_some()", raw.expr)
	}

	// A presence guard on a type that is neither Option-wrapped nor a
	// collection is vacuous in Rust (the value always exists); the caller
	// drops it.
	return ""
}

func (emitter *converterEmitter) renderComparisonGuard(guard languages.MappingGuard, raw pathExpr, unwrapped pathExpr) string {
	optionWrapped := isOptionWrapped(raw.def)

	if boolValue, ok := guard.Value.(bool); ok {
		if optionWrapped {
			return fmt.Sprintf("%s == Some(%t)", raw.expr, (guard.Op == ast.EqualOp) == boolValue)
		}
		if (guard.Op == ast.EqualOp) == boolValue {
			return unwrapped.expr
		}
		return "!" + unwrapped.expr
	}

	if stringValue, ok := guard.Value.(string); ok {
		if stringValue == "" {
			if guard.Op == ast.EqualOp {
				return fmt.Sprintf("%s.is_empty()", unwrapped.expr)
			}
			return fmt.Sprintf("!%s.is_empty()", unwrapped.expr)
		}
		if optionWrapped {
			return fmt.Sprintf("%s.as_deref() %s Some(%q)", raw.expr, guard.Op, stringValue)
		}
		return fmt.Sprintf("%s %s %q", unwrapped.expr, guard.Op, stringValue)
	}

	literal := emitter.guardLiteral(raw.def, guard.Value)
	if optionWrapped {
		return fmt.Sprintf("%s %s Some(%s)", raw.expr, guard.Op, literal)
	}
	return fmt.Sprintf("%s %s %s", unwrapped.expr, guard.Op, literal)
}

// guardLiteral renders a guard constant against the guarded field's type, so a
// numeric default compares against a literal of the right Rust numeric type
// and an enum constant compares against the matching variant.
func (emitter *converterEmitter) guardLiteral(def ast.Type, value any) string {
	return formatConstantValue(emitter.formatter, emitter.context, def, value)
}

// renderJSONGuard renders a guard whose path traversed into a
// serde_json::Value (an `any` field with a type hint). The value is accessed
// through serde_json's typed accessors instead of struct fields.
func (emitter *converterEmitter) renderJSONGuard(guard languages.MappingGuard, raw pathExpr) string {
	if guard.Value == nil {
		if guard.Op == ast.EqualOp {
			return fmt.Sprintf("%s.is_null()", raw.expr)
		}
		return fmt.Sprintf("!%s.is_null()", raw.expr)
	}

	if stringValue, ok := guard.Value.(string); ok {
		accessor := fmt.Sprintf("%s.as_str().unwrap_or_default()", raw.expr)
		if stringValue == "" {
			if guard.Op == ast.EqualOp {
				return fmt.Sprintf("%s.is_empty()", accessor)
			}
			return fmt.Sprintf("!%s.is_empty()", accessor)
		}
		return fmt.Sprintf("%s %s %q", accessor, guard.Op, stringValue)
	}

	if boolValue, ok := guard.Value.(bool); ok {
		accessor := fmt.Sprintf("%s.as_bool().unwrap_or_default()", raw.expr)
		if (guard.Op == ast.EqualOp) == boolValue {
			return accessor
		}
		return "!" + accessor
	}

	return fmt.Sprintf("%s %s serde_json::json!(%s)", raw.expr, guard.Op, formatValue(guard.Value))
}

// valueExpr renders the Rust expression that produces the *source code string*
// for a value: a quoted literal for scalars (Debug formatting yields a valid
// Rust literal for strings, numbers and booleans), a serde_json::json! literal
// for `any` values, and a from_value round-trip for everything else (enums and
// structs without a builder), mirroring the Go converter's cog.Dump.
func (emitter *converterEmitter) valueExpr(def ast.Type, path ast.Path) string {
	rendered := emitter.renderPath(path, true)

	if rendered.json {
		if def.IsScalar() && def.AsScalar().ScalarKind == ast.KindString {
			return fmt.Sprintf("format!(\"{:?}.to_string()\", %s.as_str().unwrap_or_default())", rendered.expr)
		}
		return fmt.Sprintf("format!(\"{}\", %s)", rendered.expr)
	}

	if def.IsAny() {
		emitter.imports.Add("crate::cog")
		return fmt.Sprintf("cog::dump_json(%s)", emitter.refExpr(path))
	}

	if def.IsScalar() {
		if def.AsScalar().ScalarKind == ast.KindString {
			return fmt.Sprintf("format!(\"{:?}.to_string()\", %s)", rendered.expr)
		}
		return fmt.Sprintf("format!(\"{:?}\", %s)", rendered.expr)
	}

	emitter.imports.Add("crate::cog")
	return fmt.Sprintf("cog::dump(%s)", emitter.refExpr(path))
}

// converterCall renders a call to another generated converter, importing its
// module when it lives in a different package.
func (emitter *converterEmitter) converterCall(builderPkg string, builderName string, valuePath ast.Path) string {
	call := formatConverterName(builderName)
	if formatPackageName(builderPkg) != formatPackageName(emitter.currentPkg) {
		pkg := formatPackageName(builderPkg)
		emitter.imports.Add("crate::converters::" + pkg)
		call = pkg + "::" + call
	}
	return fmt.Sprintf("%s(%s)", call, emitter.refExpr(valuePath))
}

// emitFunction renders the complete converter function.
func (emitter *converterEmitter) emitFunction(converterName string, inputType string) string {
	var out strings.Builder

	converter := emitter.converter

	inputName := "input"
	if len(converter.Mappings) == 0 && len(converter.ConstructorArgs) == 0 {
		// Nothing reads the input; underscore the parameter to keep the
		// generated crate warning-free.
		inputName = "_input"
	}

	fmt.Fprintf(&out, "/// %s accepts a `%s` object and generates the Rust code to build this object using builders.\n", converterName, formatTypeName(converter.BuilderName))
	fmt.Fprintf(&out, "pub fn %s(%s: &%s) -> String {\n", converterName, inputName, inputType)

	callsBinding := "let calls: Vec<String>"
	if len(converter.Mappings) > 0 {
		callsBinding = "let mut calls: Vec<String>"
	}

	ctor := fmt.Sprintf("%s::%sBuilder::new(", formatPackageName(converter.Package), formatTypeName(converter.BuilderName))
	if len(converter.ConstructorArgs) == 0 {
		fmt.Fprintf(&out, "%s = vec![%q.to_string()];\n", callsBinding, ctor+")")
	} else {
		args := tools.Map(converter.ConstructorArgs, func(arg languages.DirectArgMapping) string {
			return emitter.valueExpr(arg.ValueType, arg.ValuePath)
		})
		placeholders := strings.TrimSuffix(strings.Repeat("{}, ", len(args)), ", ")
		fmt.Fprintf(&out, "%s = vec![format!(%q, %s)];\n", callsBinding, ctor+placeholders+")", strings.Join(args, ", "))
	}

	for _, mapping := range converter.Mappings {
		emitter.emitMapping(&out, mapping)
	}

	out.WriteString("\ncalls.join(\"\\n    .\")\n")
	out.WriteString("}\n")

	return out.String()
}

func (emitter *converterEmitter) emitMapping(out *strings.Builder, mapping languages.ConversionMapping) {
	firstGuards := ""
	if len(mapping.Options) > 0 {
		firstGuards = emitter.renderGuards(mapping.Options[0].Guards)
	}

	if firstGuards != "" {
		fmt.Fprintf(out, "if %s {\n", firstGuards)
	}

	if mapping.RepeatFor != nil {
		repeatExpr := emitter.renderPath(mapping.RepeatFor, true)
		iterated := repeatExpr.expr
		if !repeatExpr.isRef {
			iterated = "&" + iterated
		}
		if mapping.RepeatIndex != "" {
			fmt.Fprintf(out, "for (%s, %s) in %s {\n", formatArgName(mapping.RepeatIndex), formatArgName(mapping.RepeatAs), iterated)
		} else {
			fmt.Fprintf(out, "for %s in %s {\n", formatArgName(mapping.RepeatAs), iterated)
		}
	}

	for _, optMapping := range mapping.Options {
		emitter.emitOption(out, optMapping)
	}

	if mapping.RepeatFor != nil {
		out.WriteString("}\n")
	}
	if firstGuards != "" {
		out.WriteString("}\n")
	}
}

func (emitter *converterEmitter) emitOption(out *strings.Builder, optMapping languages.OptionMapping) {
	argumentGuards := emitter.renderGuards(optMapping.ArgumentGuards())

	// Mirror the Go template's block structure: argument guards wrap the option
	// in their own conditional; an entirely unguarded option still gets a bare
	// block so its buffer binding stays scoped.
	switch {
	case argumentGuards != "":
		fmt.Fprintf(out, "if %s {\n", argumentGuards)
	case len(optMapping.Guards) == 0:
		out.WriteString("{\n")
	}

	out.WriteString("let mut buffer = String::new();\n")
	fmt.Fprintf(out, "buffer.push_str(%q);\n", formatFieldName(optMapping.Option.Name)+"(")

	for i, arg := range optMapping.Args {
		intoVar := fmt.Sprintf("arg%d", i)
		emitter.emitPrepareArg(out, intoVar, arg)
		fmt.Fprintf(out, "buffer.push_str(&%s);\n", intoVar)
		if i != len(optMapping.Args)-1 {
			out.WriteString("buffer.push_str(\", \");\n")
		}
	}

	out.WriteString("buffer.push(')');\n")
	out.WriteString("calls.push(buffer);\n")

	if argumentGuards != "" || len(optMapping.Guards) == 0 {
		out.WriteString("}\n")
	}
}

// emitPrepareArg renders the statements computing the source-code string for a
// single option argument into the binding named intoVar, mirroring the Go
// template's prepare_arg block.
func (emitter *converterEmitter) emitPrepareArg(out *strings.Builder, intoVar string, arg languages.ArgumentMapping) {
	switch {
	case arg.Builder != nil:
		fmt.Fprintf(out, "let %s = %s;\n", intoVar, emitter.converterCall(arg.Builder.BuilderPkg, arg.Builder.BuilderName, arg.Builder.ValuePath))
	case arg.BuilderDisjunction != nil:
		fmt.Fprintf(out, "let mut %s = String::new();\n", intoVar)
		for _, choice := range arg.BuilderDisjunction {
			fmt.Fprintf(out, "if %s {\n", emitter.renderGuards(choice.Guards))
			fmt.Fprintf(out, "%s = %s;\n", intoVar, emitter.converterCall(choice.Builder.BuilderPkg, choice.Builder.BuilderName, choice.Builder.ValuePath))
			out.WriteString("}\n")
		}
	case arg.Array != nil:
		emitter.emitArrayArg(out, intoVar, arg.Array)
	case arg.Map != nil:
		emitter.emitMapArg(out, intoVar, arg.Map)
	case arg.Runtime != nil:
		emitter.emitRuntimeArg(out, intoVar, arg.Runtime)
	case arg.Direct != nil:
		fmt.Fprintf(out, "let %s = %s;\n", intoVar, emitter.valueExpr(arg.Direct.ValueType, arg.Direct.ValuePath))
	default:
		// A raw disjunction argument never survives the compiler passes the Rust
		// target runs (mirroring the Go template, which has no rendering for it
		// either). Emit an empty argument defensively rather than failing the
		// whole generation.
		fmt.Fprintf(out, "let %s = String::new();\n", intoVar)
	}
}

func (emitter *converterEmitter) emitArrayArg(out *strings.Builder, intoVar string, arg *languages.ArrayArgMapping) {
	tmpVar := "tmp_" + intoVar
	loopVar := formatArgName(arg.ValueAs.Last().Identifier)
	subIntoVar := fmt.Sprintf("tmp_%s_%s", formatArgName(arg.For.Last().Identifier), loopVar)

	forExpr := emitter.renderPath(arg.For, true)
	iterated := forExpr.expr
	if !forExpr.isRef {
		iterated = "&" + iterated
	}

	fmt.Fprintf(out, "let mut %s: Vec<String> = Vec::new();\n", tmpVar)
	fmt.Fprintf(out, "for %s in %s {\n", loopVar, iterated)
	emitter.emitPrepareArg(out, subIntoVar, *arg.ForArg)
	fmt.Fprintf(out, "%s.push(%s);\n", tmpVar, subIntoVar)
	out.WriteString("}\n")
	fmt.Fprintf(out, "let %s = format!(\"vec![{}]\", %s.join(\", \"));\n", intoVar, tmpVar)
}

func (emitter *converterEmitter) emitMapArg(out *strings.Builder, intoVar string, arg *languages.MapArgMapping) {
	loopVar := formatArgName(arg.ValueAs.Last().Identifier)
	subIntoVar := fmt.Sprintf("tmp_%s_%s", formatArgName(arg.For.Last().Identifier), loopVar)

	forExpr := emitter.renderPath(arg.For, true)
	iterated := forExpr.expr
	if !forExpr.isRef {
		iterated = "&" + iterated
	}

	keyCode := "format!(\"{:?}\", key)"
	if arg.IndexType.IsScalar() && arg.IndexType.AsScalar().ScalarKind == ast.KindString {
		keyCode = "format!(\"{:?}.to_string()\", key)"
	}

	fmt.Fprintf(out, "let mut %s = String::from(\"std::collections::HashMap::from([\");\n", intoVar)
	fmt.Fprintf(out, "for (key, %s) in %s {\n", loopVar, iterated)
	emitter.emitPrepareArg(out, subIntoVar, *arg.ForArg)
	// The key is rendered into its own binding: nesting the format! call as a
	// format! argument would trip clippy's format_in_format_args lint.
	fmt.Fprintf(out, "let %s_key = %s;\n", subIntoVar, keyCode)
	fmt.Fprintf(out, "%s.push_str(&format!(\"({}, {}), \", %s_key, %s));\n", intoVar, subIntoVar, subIntoVar)
	out.WriteString("}\n")
	fmt.Fprintf(out, "%s.push_str(\"])\");\n", intoVar)
}

func (emitter *converterEmitter) emitRuntimeArg(out *strings.Builder, intoVar string, arg *languages.RuntimeArgMapping) {
	emitter.imports.Add("crate::cog")

	args := tools.Map(arg.Args, func(runtimeArg *languages.DirectArgMapping) string {
		if runtimeArg.ValueType.IsComposableSlot() {
			// A composable slot is a Box<dyn Variant>; the runtime helpers take
			// the trait object by reference.
			return emitter.renderPath(runtimeArg.ValuePath, true).expr + ".as_ref()"
		}
		return emitter.refExpr(runtimeArg.ValuePath)
	})

	fmt.Fprintf(out, "let %s = cog::%s(%s);\n", intoVar, tools.SnakeCase(arg.FuncName), strings.Join(args, ", "))
}

// guardLength extracts an integer from a min/max length guard value, which the
// IR may carry as an int or (after a JSON round-trip) a float64.
func guardLength(value any) (int, bool) {
	switch v := value.(type) {
	case int:
		return v, true
	case int64:
		return int(v), true
	case float64:
		return int(v), true
	default:
		return 0, false
	}
}
