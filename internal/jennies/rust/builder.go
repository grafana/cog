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

// Builder emits an idiomatic Rust builder for every builder in the IR. Each
// builder is a `<Name>Builder` struct holding the object under construction in
// an `internal` field, mirroring the Go target's `internal *Foo`. The
// constructor (`new`) applies the constructor assignments, every option is a
// chaining method consuming `self` and returning `Self`, and `build` implements
// the runtime `cog::Builder<T>` trait by cloning the assembled object.
//
// Phase 4a scope: direct assignment of constants and option arguments to
// single-level field paths. Phase 4b adds constraint validation (an `errors`
// field accumulating cog::BuildError, populated by constrained option setters
// and surfaced by build()), constant assignment, discriminator-without-option,
// and the `known_any` shape (a two-level path into an `any`-typed field with a
// concrete TypeHint). Array-append, index assignment, builder delegation,
// factories, envelopes and variants are emitted by later chunks; a builder
// requiring any of them is rejected with an error so out-of-scope fixtures are
// skipped explicitly rather than mis-generated.
type Builder struct {
	config          Config
	apiRefCollector *common.APIReferenceCollector
}

func (jenny Builder) JennyName() string {
	return "RustBuilder"
}

func (jenny Builder) Generate(context languages.Context) (codejen.Files, error) {
	// Several builders can share a package (for example the two discriminator
	// builders in the discriminator_without_option fixture). They all land in the
	// same `src/builders/<package>.rs` module, so builders are grouped by package
	// and their bodies concatenated under a single, merged import block. Package
	// order follows first appearance to keep the output deterministic.
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
			imports := newImportMap()
			imports.Add("crate::cog")
			group = &packageGroup{pkg: builder.Package, imports: imports}
			groups[builder.Package] = group
			order = append(order, builder.Package)
		}

		body, err := jenny.generateBuilder(context, builder, group.imports)
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

		filename := filepath.Join("src", "builders", formatPackageName(pkg)+".rs")
		files = append(files, *codejen.NewFile(filename, []byte(out.String()), jenny))
	}

	return files, nil
}

func (jenny Builder) generateBuilder(context languages.Context, builder ast.Builder, imports *importMap) (string, error) {
	if err := jenny.rejectUnsupported(context, builder); err != nil {
		return "", err
	}

	formatter := newTypeFormatter(context, imports, builder.Package)
	// A builder module is distinct from its types module, so even same-package type
	// refs it names must be imported from crate::types.
	formatter.importSamePackageRefs = true

	objectType := formatTypeName(builder.For.Name)
	builderType := objectType + "Builder"

	// The object lives in the builder's own module under crate::types so the
	// builder can name it without a same-package short name (the builder module
	// is distinct from the types module).
	objectPath := fmt.Sprintf("crate::types::%s::%s", formatter.packageName, objectType)
	imports.Add(objectPath)

	var body strings.Builder

	body.WriteString(formatComments(builder.For.Comments, ""))
	fmt.Fprintf(&body, "#[derive(Debug, Clone)]\n")
	fmt.Fprintf(&body, "pub struct %s {\n", builderType)
	fmt.Fprintf(&body, "    internal: %s,\n", objectType)
	fmt.Fprintf(&body, "    errors: Vec<cog::BuildError>,\n")
	body.WriteString("}\n\n")

	body.WriteString(jenny.formatConstructor(formatter, context, builder, builderType, objectType))

	for _, option := range builder.Options {
		method, err := jenny.formatOption(formatter, context, builder, option)
		if err != nil {
			return "", err
		}
		body.WriteString("\n\n")
		body.WriteString(method)
	}

	body.WriteString("\n\n")
	body.WriteString(jenny.formatBuild(builderType, objectType))
	body.WriteString("\n")

	return body.String(), nil
}

// rejectUnsupported reports an error for any builder feature outside Phase 4a so
// that out-of-scope fixtures fail loudly (and are skipped in the suite) rather
// than emitting incorrect code.
func (jenny Builder) rejectUnsupported(context languages.Context, builder ast.Builder) error {
	// Builder properties (extra builder-only state, e.g. the `someBuilderProperty`
	// in the `properties` fixture) are not referenced by any in-scope option's
	// assignments, so they are inert holders here and intentionally not emitted.
	// Later chunks that introduce options consuming a property will surface them.
	if len(builder.Factories) != 0 {
		return fmt.Errorf("rust builder %q: factories are not supported until a later phase", builder.Name)
	}

	allAssignments := append([]ast.Assignment{}, builder.Constructor.Assignments...)
	for _, option := range builder.Options {
		allAssignments = append(allAssignments, option.Assignments...)
	}

	for _, assignment := range allAssignments {
		if assignment.Method != ast.DirectAssignment {
			return fmt.Errorf("rust builder %q: assignment method %q is not supported until a later phase", builder.Name, assignment.Method)
		}
		// A single-level field assignment is the common case. The only supported
		// multi-level path is the `known_any` shape: a two-level path whose first
		// element is an `any`-typed field carrying a concrete ref TypeHint (the
		// builder constructs that concrete type, sets the nested field and stores
		// the result back into the `any` field). Any other multi-level or indexed
		// path is out of scope.
		if !isKnownAnyPath(assignment.Path) {
			if len(assignment.Path) != 1 || assignment.Path[0].Index != nil {
				return fmt.Errorf("rust builder %q: only single-level field assignments are supported in this phase", builder.Name)
			}
		}
		if assignment.Value.Envelope != nil || assignment.Value.ConstructorFor != nil {
			return fmt.Errorf("rust builder %q: envelope/constructor assignments are not supported until a later phase", builder.Name)
		}
		// The only nil-check tolerated here is the one the pipeline injects for the
		// known_any composition (the `if config is None` guard); that path
		// fresh-constructs the hinted type, subsuming the check. Any other nil-check
		// guard is out of scope.
		if len(assignment.NilChecks) != 0 && !isKnownAnyPath(assignment.Path) {
			return fmt.Errorf("rust builder %q: nil-check guards are not supported until a later phase", builder.Name)
		}
		// A direct assignment of an option argument whose type resolves to a struct
		// with its own builder is builder delegation, handled in a later chunk. A
		// ref to an enum (no builder) is a plain value assignment and is in scope.
		if assignment.Value.Argument != nil && context.ResolveToBuilder(assignment.Value.Argument.Type) {
			return fmt.Errorf("rust builder %q: ref-typed option arguments (builder delegation) are not supported until a later phase", builder.Name)
		}
	}

	return nil
}

// formatConstructor emits `new`. When the constructor takes no arguments the
// builder also derives a `Default` (via an argument-less `new`), matching the Go
// target's argument-less constructor. The object's own `Default` seeds the
// internal value, then each constructor assignment overrides a field.
func (jenny Builder) formatConstructor(formatter *typeFormatter, context languages.Context, builder ast.Builder, builderType, objectType string) string {
	var buffer strings.Builder

	fmt.Fprintf(&buffer, "impl %s {\n", builderType)

	args := jenny.formatArgs(formatter, builder.Constructor.Args)
	fmt.Fprintf(&buffer, "    pub fn new(%s) -> Self {\n", args)

	if len(builder.Constructor.Assignments) == 0 {
		// With no constructor assignments the builder is the bare struct literal;
		// binding it to a local and returning it would trip clippy's let_and_return.
		fmt.Fprintf(&buffer, "        Self {\n")
		fmt.Fprintf(&buffer, "            internal: %s::default(),\n", objectType)
		fmt.Fprintf(&buffer, "            errors: Vec::new(),\n")
		buffer.WriteString("        }\n")
		buffer.WriteString("    }\n")
		buffer.WriteString("}")

		if len(builder.Constructor.Args) == 0 {
			buffer.WriteString("\n\n")
			buffer.WriteString(jenny.formatDefaultImpl(builderType))
		}
		return buffer.String()
	}

	fmt.Fprintf(&buffer, "        let mut builder = Self {\n")
	fmt.Fprintf(&buffer, "            internal: %s::default(),\n", objectType)
	fmt.Fprintf(&buffer, "            errors: Vec::new(),\n")
	buffer.WriteString("        };\n")

	for _, assignment := range builder.Constructor.Assignments {
		buffer.WriteString("        ")
		buffer.WriteString(jenny.formatAssignment(formatter, context, assignment, "builder"))
		buffer.WriteString("\n")
	}

	buffer.WriteString("\n        builder\n")
	buffer.WriteString("    }\n")
	buffer.WriteString("}")

	if len(builder.Constructor.Args) == 0 {
		buffer.WriteString("\n\n")
		buffer.WriteString(jenny.formatDefaultImpl(builderType))
	}

	return buffer.String()
}

// formatDefaultImpl emits a Default impl delegating to the argument-less new(),
// matching the Go target's argument-less constructor.
func (jenny Builder) formatDefaultImpl(builderType string) string {
	var buffer strings.Builder
	fmt.Fprintf(&buffer, "impl Default for %s {\n", builderType)
	buffer.WriteString("    fn default() -> Self {\n")
	buffer.WriteString("        Self::new()\n")
	buffer.WriteString("    }\n")
	buffer.WriteString("}")
	return buffer.String()
}

// formatOption emits a single chaining option method. The method consumes `self`
// (`mut self`) and returns `Self`, applying each assignment to `self.internal`.
func (jenny Builder) formatOption(formatter *typeFormatter, context languages.Context, builder ast.Builder, option ast.Option) (string, error) {
	var buffer strings.Builder

	buffer.WriteString(formatComments(option.Comments, ""))
	fmt.Fprintf(&buffer, "impl %sBuilder {\n", formatTypeName(builder.For.Name))

	args := jenny.formatArgs(formatter, option.Args)
	methodName := formatFieldName(option.Name)
	if args == "" {
		fmt.Fprintf(&buffer, "    pub fn %s(mut self) -> Self {\n", methodName)
	} else {
		fmt.Fprintf(&buffer, "    pub fn %s(mut self, %s) -> Self {\n", methodName, args)
	}

	for _, assignment := range option.Assignments {
		for _, constraint := range assignment.Constraints {
			check, err := jenny.formatConstraint(constraint)
			if err != nil {
				return "", err
			}
			buffer.WriteString(check)
		}
	}

	for _, assignment := range option.Assignments {
		buffer.WriteString("        ")
		buffer.WriteString(jenny.formatAssignment(formatter, context, assignment, "self"))
		buffer.WriteString("\n")
	}

	buffer.WriteString("\n        self\n")
	buffer.WriteString("    }\n")
	buffer.WriteString("}")

	return buffer.String(), nil
}

// formatConstraint renders a single validation check inside an option setter. A
// violated constraint pushes a cog::BuildError onto self.errors (keyed by the
// argument name) and the setter still returns Self, mirroring the Go target's
// RecordError/errors-slice model where build() surfaces the accumulated errors.
//
// String-length ops (minLength/maxLength) compare the argument's char count;
// numeric ops compare the value directly; regex ops are not yet supported and
// are rejected so a regex dependency is introduced deliberately in a later phase.
func (jenny Builder) formatConstraint(constraint ast.AssignmentConstraint) (string, error) {
	argName := formatArgName(constraint.Argument.Name)

	// Render the comparison parameter against the argument's scalar kind so an
	// integer-typed argument compares against an integer literal (5) rather than a
	// float literal (5.0). Length ops always compare a usize count, so their
	// parameter is rendered as a plain integer.
	parameter := formatValue(constraint.Parameter)
	argType := constraint.Argument.Type
	switch constraint.Op {
	case ast.MinLengthOp, ast.MaxLengthOp:
		parameter = fmt.Sprintf("%v", constraint.Parameter)
	default:
		if argType.IsScalar() {
			parameter = formatScalarValue(constraint.Parameter, argType.AsScalar().ScalarKind)
		}
	}

	// The check guards against a violation, so the emitted condition is the
	// negation of the constraint. Negating the operator (rather than wrapping in
	// `!(...)`) keeps clippy's nonminimal_bool lint quiet.
	var left, op, message string
	switch constraint.Op {
	case ast.MinLengthOp:
		left = fmt.Sprintf("%s.chars().count()", argName)
		op = "<"
		message = fmt.Sprintf("%s length must be >= %s", argName, parameter)
	case ast.MaxLengthOp:
		left = fmt.Sprintf("%s.chars().count()", argName)
		op = ">"
		message = fmt.Sprintf("%s length must be <= %s", argName, parameter)
	case ast.GreaterThanOp, ast.GreaterThanEqualOp, ast.LessThanOp, ast.LessThanEqualOp, ast.EqualOp, ast.NotEqualOp:
		left = argName
		op = negateComparisonOp(constraint.Op)
		message = fmt.Sprintf("%s must be %s %s", argName, string(constraint.Op), parameter)
	case ast.RegexMatchOp, ast.NotRegexMatchOp:
		return "", fmt.Errorf("rust builder: regex constraint %q is not supported until a later phase (needs a regex dependency)", constraint.Op)
	default:
		return "", fmt.Errorf("rust builder: constraint op %q is not supported", constraint.Op)
	}

	var buffer strings.Builder
	fmt.Fprintf(&buffer, "        if %s %s %s {\n", left, op, parameter)
	buffer.WriteString(formatErrorPush(argName, message))
	buffer.WriteString("        }\n")
	return buffer.String(), nil
}

// formatErrorPush renders `self.errors.push(cog::BuildError::new(path, message))`
// at 12-space indentation, choosing the same line wrapping rustfmt would apply
// (max width 100): a single line, then the receiver/method split, then one
// argument per line. Emitting the rustfmt-canonical form keeps generated builders
// `cargo fmt --check` clean without a formatting post-process step.
func formatErrorPush(path, message string) string {
	newCall := fmt.Sprintf("cog::BuildError::new(%q, %q.to_string())", path, message)

	// rustfmt's chain_width default (60) means the `.push(...)` segment is always
	// too wide to stay inline with the `self.errors` receiver, so the receiver and
	// method split onto separate lines. The inner `cog::BuildError::new(...)` call
	// then stays on one line only while it fits rustfmt's fn_call_width (60);
	// beyond that rustfmt expands the call arguments one per line.
	if len(newCall) <= 60 {
		return fmt.Sprintf("            self.errors\n                .push(%s);\n", newCall)
	}

	return fmt.Sprintf("            self.errors.push(cog::BuildError::new(\n                %q,\n                %q.to_string(),\n            ));\n", path, message)
}

// negateComparisonOp returns the Rust operator that is true exactly when the
// given constraint comparison is violated.
func negateComparisonOp(op ast.Op) string {
	switch op {
	case ast.GreaterThanOp:
		return "<="
	case ast.GreaterThanEqualOp:
		return "<"
	case ast.LessThanOp:
		return ">="
	case ast.LessThanEqualOp:
		return ">"
	case ast.EqualOp:
		return "!="
	case ast.NotEqualOp:
		return "=="
	default:
		return string(op)
	}
}

// formatBuild implements the runtime cog::Builder<T> trait. Constrained option
// setters accumulate violations into self.errors; build() surfaces them (mirroring
// the Go target's RecordError/errors slice) and otherwise clones the assembled
// object. Required-field nil checks are layered in by later chunks.
func (jenny Builder) formatBuild(builderType, objectType string) string {
	var buffer strings.Builder
	fmt.Fprintf(&buffer, "impl cog::Builder<%s> for %s {\n", objectType, builderType)
	fmt.Fprintf(&buffer, "    fn build(&self) -> Result<%s, Vec<cog::BuildError>> {\n", objectType)
	buffer.WriteString("        if !self.errors.is_empty() {\n")
	buffer.WriteString("            return Err(self.errors.clone());\n")
	buffer.WriteString("        }\n\n")
	buffer.WriteString("        Ok(self.internal.clone())\n")
	buffer.WriteString("    }\n")
	buffer.WriteString("}")
	return buffer.String()
}

// formatArgs renders a method/constructor argument list. Argument types are
// rendered through the shared type formatter, so cross-package refs and
// collections match the rawtypes output exactly.
func (jenny Builder) formatArgs(formatter *typeFormatter, args []ast.Argument) string {
	if len(args) == 0 {
		return ""
	}
	parts := make([]string, 0, len(args))
	for _, arg := range args {
		parts = append(parts, fmt.Sprintf("%s: %s", formatArgName(arg.Name), formatter.formatType(arg.Type)))
	}
	return strings.Join(parts, ", ")
}

// formatAssignment renders a single direct assignment to a field of the object
// held in receiver.internal. The field type drives the value rendering: an
// Option-wrapped field is assigned Some(value); a bare field is assigned the
// value directly.
func (jenny Builder) formatAssignment(formatter *typeFormatter, context languages.Context, assignment ast.Assignment, receiver string) string {
	if isKnownAnyPath(assignment.Path) {
		return jenny.formatKnownAnyAssignment(formatter, context, assignment, receiver)
	}

	field := assignment.Path[0]
	fieldName := formatFieldName(field.Identifier)
	value := jenny.formatValue(formatter, context, field.Type, assignment.Value)

	if isOptionWrapped(field.Type) {
		return fmt.Sprintf("%s.internal.%s = Some(%s);", receiver, fieldName, value)
	}
	return fmt.Sprintf("%s.internal.%s = %s;", receiver, fieldName, value)
}

// isKnownAnyPath reports the `known_any` shape: a two-level path whose first
// element is an `any`-typed field carrying a concrete ref TypeHint. The builder
// composes a value of the hinted type into the otherwise-untyped field.
func isKnownAnyPath(path ast.Path) bool {
	if len(path) != 2 {
		return false
	}
	root := path[0]
	if root.Index != nil || root.TypeHint == nil || !root.TypeHint.IsRef() {
		return false
	}
	return root.Type.IsScalar() && root.Type.AsScalar().ScalarKind == ast.KindAny
}

// formatKnownAnyAssignment renders the `known_any` assignment. The `any` field is
// rendered as serde_json::Value; the builder constructs the concrete hinted type,
// sets the nested field on it, and stores it back serialized into the Value field.
// serde_json::to_value cannot fail for a generated serializable struct, so the
// Result is unwrapped via expect with a descriptive message.
func (jenny Builder) formatKnownAnyAssignment(formatter *typeFormatter, context languages.Context, assignment ast.Assignment, receiver string) string {
	root := assignment.Path[0]
	nested := assignment.Path[1]

	rootField := formatFieldName(root.Identifier)
	nestedField := formatFieldName(nested.Identifier)

	ref := root.TypeHint.AsRef()
	concreteType := formatter.formatRef(ref)
	formatter.imports.Add(fmt.Sprintf("crate::types::%s::%s", formatPackageName(ref.ReferredPkg), formatTypeName(ref.ReferredType)))

	value := jenny.formatValue(formatter, context, nested.Type, assignment.Value)
	if isOptionWrapped(nested.Type) {
		value = fmt.Sprintf("Some(%s)", value)
	}

	// A struct literal (rather than Default::default() followed by reassignment)
	// keeps clippy's field_reassign_with_default lint quiet. The remaining fields
	// come from Default so only the composed-in field is named.
	// Use field-init shorthand when the value is exactly the field name to avoid
	// clippy's redundant_field_names lint.
	fieldInit := fmt.Sprintf("%s: %s", nestedField, value)
	if value == nestedField {
		fieldInit = nestedField
	}

	// Emit the `..Default::default()` rest only when the hinted struct carries
	// fields beyond the one being set; otherwise it trips clippy's needless_update.
	needsRest := true
	if referred, found := context.LocateObjectByRef(ref); found && referred.Type.IsStruct() {
		needsRest = len(referred.Type.AsStruct().Fields) > 1
	}

	var buffer strings.Builder
	// rustfmt collapses a single-field struct literal with no rest onto one line;
	// a literal with a rest expression stays multi-line. Emit the form rustfmt
	// would settle on so the golden is `cargo fmt --check` clean.
	if needsRest {
		fmt.Fprintf(&buffer, "let %s = %s {\n", rootField, concreteType)
		fmt.Fprintf(&buffer, "            %s,\n", fieldInit)
		fmt.Fprintf(&buffer, "            ..%s::default()\n", concreteType)
		buffer.WriteString("        };\n")
	} else {
		fmt.Fprintf(&buffer, "let %s = %s { %s };\n", rootField, concreteType, fieldInit)
	}

	// rustfmt breaks the assignment after `=` when the value does not fit on the
	// statement line (max width 100).
	store := fmt.Sprintf("serde_json::to_value(%s).expect(\"%s should serialize to JSON\")", rootField, concreteType)
	stmt := fmt.Sprintf("        %s.internal.%s = Some(%s);", receiver, rootField, store)
	if len(stmt) <= 100 {
		buffer.WriteString(stmt)
	} else {
		fmt.Fprintf(&buffer, "        %s.internal.%s =\n            Some(%s);", receiver, rootField, store)
	}
	return buffer.String()
}

// formatValue renders the right-hand side of an assignment. A constant is
// rendered against the destination field type (a String field takes an owned
// String, an enum-ref field takes the matching variant); an option argument is
// rendered as the bare argument identifier (a String argument is converted to an
// owned String since the argument is taken by value as String already).
func (jenny Builder) formatValue(formatter *typeFormatter, context languages.Context, destType ast.Type, value ast.AssignmentValue) string {
	if value.Argument != nil {
		return formatArgName(value.Argument.Name)
	}

	// Constant assignment.
	return jenny.formatConstantValue(formatter, context, destType, value.Constant)
}

// formatConstantValue renders a constant assigned to a field of destType. A ref
// to an enum resolves to Enum::Variant; a String scalar yields an owned String;
// other scalars render as their literal.
func (jenny Builder) formatConstantValue(formatter *typeFormatter, context languages.Context, destType ast.Type, constant any) string {
	if destType.IsRef() {
		ref := destType.AsRef()
		if referred, found := context.LocateObject(ref.ReferredPkg, ref.ReferredType); found && referred.Type.IsEnum() {
			typeName := formatter.formatRef(ref)
			// The builder module is never the types module, so even a same-package
			// ref (which formatRef leaves un-imported) must be brought in explicitly
			// from crate::types so the enum variant resolves.
			formatter.imports.Add(fmt.Sprintf("crate::types::%s::%s", formatPackageName(ref.ReferredPkg), typeName))
			for _, member := range referred.Type.AsEnum().Values {
				if fmt.Sprintf("%v", member.Value) == fmt.Sprintf("%v", constant) {
					return fmt.Sprintf("%s::%s", typeName, formatTypeName(member.Name))
				}
			}
			return fmt.Sprintf("%s::default()", typeName)
		}
	}

	if destType.IsScalar() && destType.AsScalar().ScalarKind == ast.KindString {
		return fmt.Sprintf("%s.to_string()", formatValue(constant))
	}
	if destType.IsScalar() {
		return formatScalarValue(constant, destType.AsScalar().ScalarKind)
	}

	return formatValue(constant)
}

// formatArgName converts an IR argument name into an idiomatic Rust binding
// (snake_case), escaping Rust keywords. Argument names follow the same rules as
// field names.
func formatArgName(name string) string {
	return escapeRustKeyword(tools.SnakeCase(name))
}
