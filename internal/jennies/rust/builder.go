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
	// The builder struct is named after the builder (which can differ from the
	// object it builds: a foreign builder lives in one package and builds a type in
	// another, e.g. SomeNiceBuilder builds some_pkg::SomeStruct).
	builderType := formatTypeName(builder.Name) + "Builder"

	// The object lives in its own (possibly foreign) package under crate::types, so
	// resolve the path from the builder's For.SelfRef rather than assuming the
	// builder's own package. The builder module is distinct from the types module,
	// so the object is always imported explicitly.
	objectPkg := formatter.packageName
	if builder.For.SelfRef.ReferredPkg != "" {
		objectPkg = formatPackageName(builder.For.SelfRef.ReferredPkg)
	}
	objectPath := fmt.Sprintf("crate::types::%s::%s", objectPkg, objectType)
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

	for _, factory := range builder.Factories {
		body.WriteString("\n\n")
		body.WriteString(jenny.formatFactory(formatter, context, builder, builderType, factory))
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
	//
	// Factories (named alternate constructors that preset option calls) are in
	// scope as of Phase 4f. Validate every option call references a real option
	// and every call parameter is a constant, an argument, or a nested factory
	// call (the shapes the in-scope fixtures use); reject anything else loudly.
	for _, factory := range builder.Factories {
		if err := checkFactorySupported(builder, factory); err != nil {
			return fmt.Errorf("rust builder %q: %w", builder.Name, err)
		}
	}

	allAssignments := append([]ast.Assignment{}, builder.Constructor.Assignments...)
	for _, option := range builder.Options {
		allAssignments = append(allAssignments, option.Assignments...)
	}

	for _, assignment := range allAssignments {
		if err := checkAssignmentSupported(context, assignment); err != nil {
			return fmt.Errorf("rust builder %q: %w", builder.Name, err)
		}
	}

	return nil
}

// checkAssignmentSupported validates a single builder assignment against the
// shapes the Rust target renders, dispatching on the assignment method. It
// returns a plain (unwrapped) error so the caller can attach the builder name.
func checkAssignmentSupported(context languages.Context, assignment ast.Assignment) error {
	switch assignment.Method {
	case ast.DirectAssignment:
		return checkDirectAssignmentSupported(context, assignment)
	case ast.AppendAssignment:
		return checkAppendAssignmentSupported(assignment)
	case ast.IndexAssignment:
		return checkIndexAssignmentSupported(assignment)
	default:
		return fmt.Errorf("assignment method %q is not supported until a later phase", assignment.Method)
	}
}

// checkAppendAssignmentSupported validates a Phase 4c append: a single
// (non-indexed) array-typed field path the value is pushed onto. An envelope
// value (Phase 4e) additionally must construct a ref via simple single-level
// field paths.
func checkAppendAssignmentSupported(assignment ast.Assignment) error {
	if !isAppendPath(assignment.Path) {
		return fmt.Errorf("unsupported append path shape")
	}
	if assignment.Value.Envelope != nil {
		if err := checkEnvelopeSupported(*assignment.Value.Envelope); err != nil {
			return err
		}
	}
	return nil
}

// checkIndexAssignmentSupported validates a Phase 4c index assignment: a
// two-level path whose first element is a map-typed field and whose second
// element carries the key (a constant or an option argument).
func checkIndexAssignmentSupported(assignment ast.Assignment) error {
	if !isIndexPath(assignment.Path) {
		return fmt.Errorf("unsupported index path shape")
	}
	return nil
}

// checkDirectAssignmentSupported validates a direct assignment: its value shape,
// its path shape, the nil-checks the pipeline injected, and any builder
// delegation.
func checkDirectAssignmentSupported(context languages.Context, assignment ast.Assignment) error {
	// A factory constructor value (ConstructorFor) is a Phase 5 feature.
	if assignment.Value.ConstructorFor != nil {
		return fmt.Errorf("constructor-value assignments are not supported until a later phase")
	}
	if err := checkDirectAssignmentPathSupported(assignment.Path); err != nil {
		return err
	}
	if err := checkNilChecksSupported(assignment); err != nil {
		return err
	}
	return checkBuilderDelegationSupported(context, assignment)
}

// checkDirectAssignmentPathSupported validates the path of a direct assignment.
// Supported shapes are a single-level field, the `known_any` two-level
// composition shape, and a multi-level ref/struct path (Phase 4e) whose
// non-leaf elements are plain (non-indexed) refs.
func checkDirectAssignmentPathSupported(path ast.Path) error {
	if isKnownAnyPath(path) {
		return nil
	}
	if len(path) == 1 {
		if path[0].Index != nil {
			return fmt.Errorf("indexed single-level paths are not supported in direct assignment")
		}
		return nil
	}
	return checkMultiLevelPathSupported(path)
}

// checkNilChecksSupported rejects a nil-check whose initialised path element is
// not a plain ref/array/map (for example an anonymous inline struct), which has
// no idiomatic default-construction in Rust. The `known_any` shape carries its
// own composition and is exempt.
func checkNilChecksSupported(assignment ast.Assignment) error {
	if isKnownAnyPath(assignment.Path) {
		return nil
	}
	for _, check := range assignment.NilChecks {
		leaf := check.Path.Last()
		if !leaf.Type.IsArray() && !leaf.Type.IsMap() && !leaf.Type.IsRef() {
			return fmt.Errorf("nil-check on non-ref intermediate %q is not supported", check.Path.String())
		}
	}
	return nil
}

// checkBuilderDelegationSupported validates builder delegation: an option
// argument whose type (or a collection's element type) resolves to a builder is
// in scope for the ref / array-of-ref / map-of-ref shapes. A disjunction-typed
// delegated argument has no clean idiomatic-Rust mapping and is rejected,
// matching the Go target.
func checkBuilderDelegationSupported(context languages.Context, assignment ast.Assignment) error {
	if assignment.Value.Argument != nil && context.ResolveToBuilder(assignment.Value.Argument.Type) {
		if !isDelegableType(assignment.Value.Argument.Type) {
			return fmt.Errorf("disjunction-typed builder delegation is not supported")
		}
	}
	return nil
}

// isDelegableType reports whether a builder-valued type can be rendered with the
// `impl cog::Builder<T>` delegation shape: a bare ref, or an array/map whose
// (recursive) element type is a delegable ref. A disjunction branch is not
// delegable (no single static builder bound exists for it).
func isDelegableType(def ast.Type) bool {
	switch {
	case def.IsArray():
		return isDelegableType(def.AsArray().ValueType)
	case def.IsMap():
		return isDelegableType(def.AsMap().ValueType)
	case def.IsRef():
		return true
	default:
		return false
	}
}

// formatConstructor emits `new`. When the constructor takes no arguments the
// builder also derives a `Default` (via an argument-less `new`), matching the Go
// target's argument-less constructor. The object's own `Default` seeds the
// internal value, then each constructor assignment overrides a field.
func (jenny Builder) formatConstructor(formatter *typeFormatter, context languages.Context, builder ast.Builder, builderType, objectType string) string {
	var buffer strings.Builder

	fmt.Fprintf(&buffer, "impl %s {\n", builderType)

	args := jenny.formatArgs(formatter, builder.Constructor.Args, context)
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
	} else {
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
	}

	// An argument-less constructor also derives a Default impl that delegates to
	// new(), matching the Go target.
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
	fmt.Fprintf(&buffer, "impl %sBuilder {\n", formatTypeName(builder.Name))

	args := jenny.formatArgs(formatter, option.Args, context)
	methodName := formatFieldName(option.Name)
	// Emit the signature on a single line and let the FormatRustFiles
	// postprocessor (rustfmt) wrap it across lines if it exceeds the max width.
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
// at 12-space indentation. It emits the single-line form and lets the
// FormatRustFiles postprocessor (rustfmt) apply any chain/call wrapping.
func formatErrorPush(path, message string) string {
	return fmt.Sprintf("            self.errors.push(cog::BuildError::new(%q, %q.to_string()));\n", path, message)
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
func (jenny Builder) formatArgs(formatter *typeFormatter, args []ast.Argument, context languages.Context) string {
	if len(args) == 0 {
		return ""
	}
	parts := make([]string, 0, len(args))
	for _, arg := range args {
		argType := formatter.formatType(arg.Type)
		// An argument whose value is produced by a nested builder (its type, or a
		// collection's element type, resolves to a builder) is rendered with the
		// generic `impl cog::Builder<T>` bound, mirroring the Go target's
		// `cog.Builder[T]` arg. The setter calls build() on it (see
		// formatDelegatedAssignment).
		if context.ResolveToBuilder(arg.Type) {
			argType = formatter.formatBuilderArgType(arg.Type)
		}
		parts = append(parts, fmt.Sprintf("%s: %s", formatArgName(arg.Name), argType))
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

	switch assignment.Method {
	case ast.AppendAssignment:
		return jenny.formatAppendAssignment(formatter, context, assignment, receiver)
	case ast.IndexAssignment:
		return jenny.formatIndexAssignment(formatter, context, assignment, receiver)
	}

	// Multi-level ref path (Phase 4e): walk into nested objects, lazily
	// initialising Option-wrapped intermediates via get_or_insert_with, and assign
	// the leaf.
	if len(assignment.Path) > 1 {
		return jenny.formatMultiLevelAssignment(formatter, context, assignment, receiver)
	}

	field := assignment.Path[0]

	// Builder delegation: the option argument is itself produced by one or more
	// nested builders. The setter calls build() on the argument (or each element of
	// a collection of builders), accumulating any errors and assigning the built
	// value(s).
	if assignment.Value.Argument != nil && context.ResolveToBuilder(assignment.Value.Argument.Type) {
		return jenny.formatDelegatedAssignment(field, assignment.Value.Argument, receiver)
	}

	fieldName := formatFieldName(field.Identifier)
	value := jenny.formatAssignmentValue(formatter, context, field.Type, assignment.Value)

	if isOptionWrapped(field.Type) {
		return fmt.Sprintf("%s.internal.%s = Some(%s);", receiver, fieldName, value)
	}
	return fmt.Sprintf("%s.internal.%s = %s;", receiver, fieldName, value)
}

// formatDelegatedAssignment renders builder delegation: the option argument is a
// builder (or a collection of builders) whose build() produces the value(s)
// assigned to the field. Errors are accumulated into <receiver>.errors and, on
// the first failure, the setter returns early - mirroring the Go target, which
// appends the built builder's errors and `return`s the builder rather than
// continuing. (Early return keeps the partially-built object untouched by a
// failed delegation and avoids inserting a zero value; build() will surface the
// accumulated errors.)
//
// The destination field type drives the shape: a single ref builds one value; an
// array builds a Vec (nested arrays recurse); a map builds a HashMap. An
// Option-wrapped destination is assigned Some(value).
func (jenny Builder) formatDelegatedAssignment(field ast.PathItem, arg *ast.Argument, receiver string) string {
	fieldName := formatFieldName(field.Identifier)
	return jenny.formatDelegatedFieldStatement(field.Type, isOptionWrapped(field.Type), arg, receiver, func(builtValue string) string {
		return fmt.Sprintf("%s.internal.%s = %s;", receiver, fieldName, builtValue)
	})
}

// formatDelegatedFieldStatement renders the shared builder-delegation sequence:
// build the delegated value of buildType from the argument (emitting the
// error-accumulating build() calls), wrap it in Some(...) when wrapInSome is
// set, then emit the final statement produced by `finalStmt` (an assignment or
// a push). Every line is emitted at the standard 8-space body indent; the caller
// prefixes the block's first line with the same indentation, so the leading
// indent of the first line is trimmed to avoid a double indent.
func (jenny Builder) formatDelegatedFieldStatement(buildType ast.Type, wrapInSome bool, arg *ast.Argument, receiver string, finalStmt func(builtValue string) string) string {
	argName := formatArgName(arg.Name)

	const indent = "        "
	var buffer strings.Builder

	value := jenny.formatDelegatedValue(&buffer, buildType, argName, receiver, indent, 0)
	if wrapInSome {
		value = fmt.Sprintf("Some(%s)", value)
	}

	fmt.Fprintf(&buffer, "%s%s", indent, finalStmt(value))
	return strings.TrimPrefix(buffer.String(), indent)
}

// formatDelegatedValue writes the statements needed to build the delegated value
// of destType from the source `expr` (a builder, or a collection of builders),
// emitting an error-accumulating match for each build() call, and returns the
// Rust expression holding the built result. Every emitted line carries `indent`;
// the returned expression is spliced into the caller's assignment line. Nested
// loops carry `depth` so their bindings never collide.
func (jenny Builder) formatDelegatedValue(buffer *strings.Builder, destType ast.Type, expr, receiver, indent string, depth int) string {
	switch {
	case destType.IsArray():
		elem := destType.AsArray().ValueType
		// Build each element with a for-loop, pushing onto a fresh Vec. A nested
		// array recurses (the inner build runs per innermost builder).
		built := fmt.Sprintf("built%d", depth)
		item := fmt.Sprintf("item%d", depth)
		fmt.Fprintf(buffer, "%slet mut %s = Vec::new();\n", indent, built)
		fmt.Fprintf(buffer, "%sfor %s in %s {\n", indent, item, expr)
		inner := jenny.formatDelegatedValue(buffer, elem, item, receiver, indent+"    ", depth+1)
		fmt.Fprintf(buffer, "%s    %s.push(%s);\n", indent, built, inner)
		fmt.Fprintf(buffer, "%s}\n", indent)
		return built
	case destType.IsMap():
		valueType := destType.AsMap().ValueType
		built := fmt.Sprintf("built%d", depth)
		key := fmt.Sprintf("key%d", depth)
		item := fmt.Sprintf("item%d", depth)
		fmt.Fprintf(buffer, "%slet mut %s = std::collections::HashMap::new();\n", indent, built)
		fmt.Fprintf(buffer, "%sfor (%s, %s) in %s {\n", indent, key, item, expr)
		inner := jenny.formatDelegatedValue(buffer, valueType, item, receiver, indent+"    ", depth+1)
		fmt.Fprintf(buffer, "%s    %s.insert(%s, %s);\n", indent, built, key, inner)
		fmt.Fprintf(buffer, "%s}\n", indent)
		return built
	default:
		// A single builder: build() it, accumulating errors and returning early on
		// failure, then bind the built value.
		binding := fmt.Sprintf("built%d", depth)
		fmt.Fprintf(buffer, "%slet %s = match %s.build() {\n", indent, binding, expr)
		fmt.Fprintf(buffer, "%s    Ok(val) => val,\n", indent)
		fmt.Fprintf(buffer, "%s    Err(mut err) => {\n", indent)
		fmt.Fprintf(buffer, "%s        %s.errors.append(&mut err);\n", indent, receiver)
		fmt.Fprintf(buffer, "%s        return %s;\n", indent, receiver)
		fmt.Fprintf(buffer, "%s    }\n", indent)
		fmt.Fprintf(buffer, "%s};\n", indent)
		return binding
	}
}

// checkFactorySupported validates a builder factory against the shapes the Rust
// target renders: every option call must name a real option on the builder, and
// every call parameter must be a constant, an argument, or a nested factory call
// (recursively validated). ForceBuild is tolerated: an argument that is itself a
// builder is delegated through the option's own build() call, so no special
// pre-build is needed here.
func checkFactorySupported(builder ast.Builder, factory ast.BuilderFactory) error {
	for _, call := range factory.OptionCalls {
		if _, found := builder.OptionByName(call.Name); !found {
			return fmt.Errorf("factory %q calls unknown option %q", factory.Name, call.Name)
		}
		for _, param := range call.Parameters {
			if err := checkFactoryParameterSupported(factory, param); err != nil {
				return err
			}
		}
	}
	return nil
}

// checkFactoryParameterSupported validates a single option-call parameter: it is
// a constant, an argument, or a nested factory call whose own parameters are also
// supported.
func checkFactoryParameterSupported(factory ast.BuilderFactory, param ast.OptionCallParameter) error {
	switch {
	case param.Constant != nil, param.Argument != nil:
		return nil
	case param.Factory != nil:
		for _, nested := range param.Factory.Parameters {
			if err := checkFactoryParameterSupported(factory, nested); err != nil {
				return err
			}
		}
		return nil
	default:
		return fmt.Errorf("factory %q has an unsupported option-call parameter", factory.Name)
	}
}

// formatFactory emits a named alternate constructor as an associated function on
// the builder. It mirrors the Go/Python factories (a function that creates a
// fresh builder and applies a fixed sequence of option calls), rendered the
// idiomatic-Rust way: `Self::new(...)` chained through the preset options,
// returning the builder for further chaining. The factory name is snake_cased to
// match the Rust option/function convention.
func (jenny Builder) formatFactory(formatter *typeFormatter, context languages.Context, builder ast.Builder, builderType string, factory ast.BuilderFactory) string {
	var buffer strings.Builder

	buffer.WriteString(formatComments(factory.Comments, ""))
	fmt.Fprintf(&buffer, "impl %s {\n", builderType)

	args := jenny.formatArgs(formatter, factory.Args, context)
	factoryName := formatFieldName(factory.Name)

	// Build the chained call expression: Self::new() followed by one option call
	// per preset, each consuming and returning Self.
	chain := []string{"Self::new()"}
	for _, call := range factory.OptionCalls {
		option, _ := builder.OptionByName(call.Name)
		chain = append(chain, jenny.formatFactoryOptionCall(formatter, context, option, call))
	}

	// Emit the signature and the chained call on single lines; the FormatRustFiles
	// postprocessor (rustfmt) lays the signature and method chain out across lines
	// when they exceed the max width.
	fmt.Fprintf(&buffer, "    pub fn %s(%s) -> Self {\n", factoryName, args)
	fmt.Fprintf(&buffer, "        %s\n", strings.Join(chain, "."))

	buffer.WriteString("    }\n")
	buffer.WriteString("}")

	return buffer.String()
}

// formatFactoryOptionCall renders one preset option call inside a factory, e.g.
// `function("abs".to_string())` or `arg(v)`. Each parameter is rendered against
// the matching option argument's type so constants, arguments, and nested factory
// calls all produce the value the option setter expects.
func (jenny Builder) formatFactoryOptionCall(formatter *typeFormatter, context languages.Context, option ast.Option, call ast.OptionCall) string {
	methodName := formatFieldName(call.Name)

	params := make([]string, 0, len(call.Parameters))
	for i, param := range call.Parameters {
		var argType ast.Type
		if i < len(option.Args) {
			argType = option.Args[i].Type
		}
		params = append(params, jenny.formatFactoryParameter(formatter, context, argType, param))
	}

	return fmt.Sprintf("%s(%s)", methodName, strings.Join(params, ", "))
}

// formatFactoryParameter renders a single option-call parameter. A constant is
// rendered against the option argument's type (so a string constant becomes an
// owned String, an enum-ref constant becomes the matching variant); an argument
// is the bare factory-argument identifier; a nested factory call becomes a call
// to the referenced builder's factory associated function (which yields a builder
// the option delegates through build()).
func (jenny Builder) formatFactoryParameter(formatter *typeFormatter, context languages.Context, argType ast.Type, param ast.OptionCallParameter) string {
	switch {
	case param.Constant != nil:
		return jenny.formatConstantValue(formatter, context, argType, param.Constant.Value)
	case param.Argument != nil:
		return formatArgName(param.Argument.Name)
	case param.Factory != nil:
		ref := param.Factory.Ref
		builderType := formatTypeName(ref.Builder) + "Builder"
		// The referenced builder lives in some package's builders module; bring it
		// in so the associated-function call resolves. A same-package builder is in
		// this very module, so only import a foreign one.
		if formatPackageName(ref.Package) != formatter.packageName {
			formatter.imports.Add(fmt.Sprintf("crate::builders::%s::%s", formatPackageName(ref.Package), builderType))
		}
		nestedParams := make([]string, 0, len(param.Factory.Parameters))
		for _, nested := range param.Factory.Parameters {
			nestedParams = append(nestedParams, jenny.formatFactoryParameter(formatter, context, ast.Type{}, nested))
		}
		return fmt.Sprintf("%s::%s(%s)", builderType, formatFieldName(ref.Factory), strings.Join(nestedParams, ", "))
	default:
		return ""
	}
}

// isAppendPath reports the append shape: a single, non-indexed path element
// whose type is an array (Vec) field. The option argument is pushed onto it.
func isAppendPath(path ast.Path) bool {
	if len(path) != 1 {
		return false
	}
	root := path[0]
	return root.Index == nil && root.Type.IsArray()
}

// isIndexPath reports the index shape: a two-level path whose first element is a
// map (HashMap) field and whose second element carries an Index (the key, either
// a constant or an option argument). The option value is inserted at that key.
func isIndexPath(path ast.Path) bool {
	if len(path) != 2 {
		return false
	}
	root := path[0]
	leaf := path[1]
	return root.Index == nil && root.Type.IsMap() && leaf.Index != nil
}

// formatAppendAssignment renders an append to a Vec field. The IR models a
// single-element append (the option argument is a single value, the field a
// Vec<value>), so this emits the idiomatic `push`. The value is rendered against
// the array's element type so an Option-valued element is wrapped in Some(...).
func (jenny Builder) formatAppendAssignment(formatter *typeFormatter, context languages.Context, assignment ast.Assignment, receiver string) string {
	field := assignment.Path[0]
	fieldName := formatFieldName(field.Identifier)
	elementType := field.Type.AsArray().ValueType

	// An envelope value constructs the element inline as a struct literal (its ref
	// type) from the envelope's per-field values, then pushes it. This is the
	// envelope_assignment shape: each WithVariable call builds a fresh Variable.
	if assignment.Value.Envelope != nil {
		literal := jenny.formatEnvelopeLiteral(formatter, context, *assignment.Value.Envelope)
		if isOptionWrapped(elementType) {
			literal = fmt.Sprintf("Some(%s)", literal)
		}
		return fmt.Sprintf("%s.internal.%s.push(%s);", receiver, fieldName, literal)
	}

	// Builder delegation: the appended element is itself produced by a nested
	// builder (the ArrayToAppend veneer shape, e.g. promql's `arg`). Build the
	// argument (accumulating errors, returning early on failure) then push the
	// built value, mirroring the Go target which appends `argResource` after
	// building it.
	if assignment.Value.Argument != nil && context.ResolveToBuilder(assignment.Value.Argument.Type) {
		// The value is built against the argument's type but wrapped (Some) against
		// the array's element type, matching the original per-element shape.
		return jenny.formatDelegatedFieldStatement(assignment.Value.Argument.Type, isOptionWrapped(elementType), assignment.Value.Argument, receiver, func(builtValue string) string {
			return fmt.Sprintf("%s.internal.%s.push(%s);", receiver, fieldName, builtValue)
		})
	}

	value := jenny.formatAssignmentValue(formatter, context, elementType, assignment.Value)
	if isOptionWrapped(elementType) {
		value = fmt.Sprintf("Some(%s)", value)
	}
	return fmt.Sprintf("%s.internal.%s.push(%s);", receiver, fieldName, value)
}

// formatIndexAssignment renders an insert into a HashMap field. The key comes
// from the leaf path element's Index (a constant or an option argument); the
// value is the assignment value rendered against the map's value type (an
// Option-valued map value is wrapped in Some(...)). Per the rawtypes convention
// maps are bare HashMap even when nullable, so no initialize-if-absent guard is
// needed: the field is always a constructed HashMap.
func (jenny Builder) formatIndexAssignment(formatter *typeFormatter, context languages.Context, assignment ast.Assignment, receiver string) string {
	field := assignment.Path[0]
	fieldName := formatFieldName(field.Identifier)
	mapType := field.Type.AsMap()

	leaf := assignment.Path[1]
	key := jenny.formatIndexKey(formatter, context, mapType.IndexType, *leaf.Index)

	value := jenny.formatAssignmentValue(formatter, context, mapType.ValueType, assignment.Value)
	if isOptionWrapped(mapType.ValueType) {
		value = fmt.Sprintf("Some(%s)", value)
	}

	return fmt.Sprintf("%s.internal.%s.insert(%s, %s);", receiver, fieldName, key, value)
}

// formatIndexKey renders a map key from a PathIndex. An argument key is the bare
// argument identifier (a String argument is already owned); a constant key is
// rendered against the map's index type (a String index yields an owned String).
func (jenny Builder) formatIndexKey(formatter *typeFormatter, context languages.Context, indexType ast.Type, index ast.PathIndex) string {
	if index.Argument != nil {
		return formatArgName(index.Argument.Name)
	}
	return jenny.formatConstantValue(formatter, context, indexType, index.Constant)
}

// checkMultiLevelPathSupported validates a multi-level direct-assignment path.
// Every non-leaf element must be a plain (non-indexed) ref so its intermediate
// object can be default-constructed when absent; the leaf may be any field.
func checkMultiLevelPathSupported(path ast.Path) error {
	for i := 0; i < len(path)-1; i++ {
		item := path[i]
		if item.Index != nil {
			return fmt.Errorf("indexed multi-level paths are not supported")
		}
		if !item.Type.IsRef() {
			return fmt.Errorf("multi-level path through a non-ref field %q is not supported", item.Identifier)
		}
	}
	if path.Last().Index != nil {
		return fmt.Errorf("indexed leaf in multi-level path is not supported")
	}
	return nil
}

// checkEnvelopeSupported validates an assignment envelope: it must construct a
// ref type, and each of its field values must target a single-level field path
// (the nested struct literal sets one named field per envelope value).
func checkEnvelopeSupported(envelope ast.AssignmentEnvelope) error {
	if !envelope.Type.IsRef() {
		return fmt.Errorf("envelope of a non-ref type is not supported")
	}
	for _, value := range envelope.Values {
		if len(value.Path) != 1 || value.Path[0].Index != nil {
			return fmt.Errorf("envelope field with a multi-level or indexed path is not supported")
		}
		if value.Value.Envelope != nil || value.Value.ConstructorFor != nil {
			return fmt.Errorf("nested envelope/constructor in an envelope field is not supported")
		}
	}
	return nil
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

	value := jenny.formatAssignmentValue(formatter, context, nested.Type, assignment.Value)
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
	// Emit the struct literal on a single line; rustfmt expands it across lines
	// when it carries a rest expression or exceeds the max width. The rest
	// expression is included only when the hinted struct has more fields than the
	// one being set (otherwise it trips clippy's needless_update).
	if needsRest {
		fmt.Fprintf(&buffer, "let %s = %s { %s, ..%s::default() };\n", rootField, concreteType, fieldInit, concreteType)
	} else {
		fmt.Fprintf(&buffer, "let %s = %s { %s };\n", rootField, concreteType, fieldInit)
	}

	// Emit the store assignment on a single line; rustfmt breaks it after `=` when
	// the value does not fit.
	store := fmt.Sprintf("serde_json::to_value(%s).expect(\"%s should serialize to JSON\")", rootField, concreteType)
	fmt.Fprintf(&buffer, "        %s.internal.%s = Some(%s);", receiver, rootField, store)
	return buffer.String()
}

// formatMultiLevelAssignment renders a direct assignment whose path descends
// through one or more nested ref objects (the initialization_safeguards shape:
// `options.legend.show = show`). Each intermediate that is Option-wrapped is
// lazily initialised with `get_or_insert_with(T::default)` so the leaf can be
// reached; this both initialises a missing intermediate and reuses an existing
// one (no clobber on a repeat call), matching the Go target's
// `if x == nil { x = NewT() }` then deep-assign. A bare (non-Option)
// intermediate is accessed directly. The leaf is assigned the value, wrapped in
// Some(...) when the leaf field is itself Option-wrapped.
//
// The EmptyValueType of the corresponding NilCheck names the type to
// default-construct; the path element types provide the same information, so the
// access chain is built from the path itself and validated against the
// NilChecks (every Option-wrapped intermediate must have a NilCheck).
func (jenny Builder) formatMultiLevelAssignment(formatter *typeFormatter, context languages.Context, assignment ast.Assignment, receiver string) string {
	// Build the access chain as a list of dotted segments and emit it on a single
	// line; the FormatRustFiles postprocessor (rustfmt) breaks a method chain out
	// one segment per line when it exceeds the chain width.
	segments := []string{fmt.Sprintf("%s.internal", receiver)}

	for i := 0; i < len(assignment.Path)-1; i++ {
		item := assignment.Path[i]
		fieldName := formatFieldName(item.Identifier)
		if isOptionWrapped(item.Type) {
			// The intermediate ref is Option-wrapped: initialise-if-absent (preserving
			// an already-present value: get_or_insert_with does not overwrite Some),
			// then descend into the contained value. The default type is the ref's
			// concrete type (imported so the path resolves).
			ref := item.Type.AsRef()
			concreteType := formatter.formatRef(ref)
			formatter.imports.Add(fmt.Sprintf("crate::types::%s::%s", formatPackageName(ref.ReferredPkg), formatTypeName(ref.ReferredType)))
			segments = append(segments, fieldName)
			segments = append(segments, fmt.Sprintf("get_or_insert_with(%s::default)", concreteType))
		} else {
			segments = append(segments, fieldName)
		}
	}

	leaf := assignment.Path.Last()
	segments = append(segments, formatFieldName(leaf.Identifier))

	value := jenny.formatAssignmentValue(formatter, context, leaf.Type, assignment.Value)
	if isOptionWrapped(leaf.Type) {
		value = fmt.Sprintf("Some(%s)", value)
	}

	return fmt.Sprintf("%s = %s;", strings.Join(segments, "."), value)
}

// formatEnvelopeLiteral renders an assignment envelope as a single-line Rust
// struct literal of the envelope's ref type, setting one named field per
// envelope value and filling the rest from Default. The FormatRustFiles
// postprocessor (rustfmt) expands the literal across lines when it exceeds the
// max width.
func (jenny Builder) formatEnvelopeLiteral(formatter *typeFormatter, context languages.Context, envelope ast.AssignmentEnvelope) string {
	ref := envelope.Type.AsRef()
	typeName := formatter.formatRef(ref)
	formatter.imports.Add(fmt.Sprintf("crate::types::%s::%s", formatPackageName(ref.ReferredPkg), typeName))

	// Determine whether the envelope sets every field of the struct; if not, a
	// `..Default::default()` rest is required (and otherwise omitted to keep
	// clippy's needless_update quiet).
	needsRest := true
	if referred, found := context.LocateObjectByRef(ref); found && referred.Type.IsStruct() {
		needsRest = len(referred.Type.AsStruct().Fields) > len(envelope.Values)
	}

	inits := make([]string, 0, len(envelope.Values))
	for _, value := range envelope.Values {
		field := value.Path[0]
		fieldName := formatFieldName(field.Identifier)
		rendered := jenny.formatAssignmentValue(formatter, context, field.Type, value.Value)
		if isOptionWrapped(field.Type) {
			rendered = fmt.Sprintf("Some(%s)", rendered)
		}
		if rendered == fieldName {
			inits = append(inits, fieldName)
		} else {
			inits = append(inits, fmt.Sprintf("%s: %s", fieldName, rendered))
		}
	}

	// Emit the struct literal on a single line; the FormatRustFiles postprocessor
	// (rustfmt) expands it across lines when it exceeds the max width. The
	// `..Default::default()` rest is included only when the envelope does not set
	// every field (otherwise it trips clippy's needless_update).
	fields := strings.Join(inits, ", ")
	if needsRest {
		return fmt.Sprintf("%s { %s, ..Default::default() }", typeName, fields)
	}
	return fmt.Sprintf("%s { %s }", typeName, fields)
}

// formatValue renders the right-hand side of an assignment. A constant is
// rendered against the destination field type (a String field takes an owned
// String, an enum-ref field takes the matching variant); an option argument is
// rendered as the bare argument identifier (a String argument is converted to an
// owned String since the argument is taken by value as String already).
func (jenny Builder) formatAssignmentValue(formatter *typeFormatter, context languages.Context, destType ast.Type, value ast.AssignmentValue) string {
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
