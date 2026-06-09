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
	// Import the object's module and refer to the object module-qualified
	// (<pkg>::<Type>), consistent with formatRef. Bare-name imports of the object
	// would collide across packages that export the same type name (e.g. the
	// several dashboard package versions all defining MatcherConfig).
	imports.Add(fmt.Sprintf("crate::types::%s", objectPkg))
	objectType = objectPkg + "::" + objectType

	var body strings.Builder

	body.WriteString(formatComments(builder.For.Comments, ""))
	fmt.Fprintf(&body, "#[derive(Debug, Clone)]\n")
	fmt.Fprintf(&body, "pub struct %s {\n", builderType)
	fmt.Fprintf(&body, "    internal: %s,\n", objectType)
	// pub(crate) so cross-module factory functions (e.g. the dashboardv2 Manifest
	// factory, which builds a resource.Manifest builder and force-builds a spec
	// argument) can append accumulated build errors. Still crate-private, so the
	// public builder API is unchanged.
	fmt.Fprintf(&body, "    pub(crate) errors: Vec<cog::BuildError>,\n")
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
	body.WriteString(jenny.formatBuild(formatter, builder, builderType, objectType))
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
		if err := checkFactorySupported(context, factory); err != nil {
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
func checkBuilderDelegationSupported(_ languages.Context, _ ast.Assignment) error {
	// A delegated argument whose type resolves to a builder but is not itself a
	// delegable shape (most commonly a disjunction one of whose branches has a
	// builder, e.g. dashboardv2 Preferences.Layout's
	// AutoGridLayoutKindOrGridLayoutKind) has no `impl cog::Builder<T>` bound. The
	// Go target handles this by taking the disjunction value directly rather than
	// delegating; the Rust target does the same (see shouldDelegate and
	// formatArgs), so this is no longer rejected.
	return nil
}

// shouldDelegate reports whether an option argument of the given type is rendered
// with builder delegation (an `impl cog::Builder<T>` parameter whose build() the
// setter calls). This holds only when the type resolves to a builder AND is a
// delegable shape (a ref, or an array/map of delegable refs). A type that
// resolves to a builder but is not delegable (a disjunction with a builder-typed
// branch) is passed through as a plain value, mirroring the Go target.
func shouldDelegate(context languages.Context, def ast.Type) bool {
	return context.ResolveToBuilder(def) && isDelegableType(def)
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

// isComposableSlotType reports whether a type is a composable slot, or a
// collection (array/map) whose element type is (recursively) a composable slot.
// Such an option argument is filled by a builder of the slot's variant value.
func isComposableSlotType(def ast.Type) bool {
	switch {
	case def.IsArray():
		return isComposableSlotType(def.AsArray().ValueType)
	case def.IsMap():
		return isComposableSlotType(def.AsMap().ValueType)
	case def.IsComposableSlot():
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
	// A nullable argument is Option<T>; a constraint applies to the contained value,
	// so it is guarded by an `if let Some(...)` that unwraps the argument (a None
	// argument simply skips the check) and the check operates on the binding. A
	// bare collection (Vec/HashMap) is not Option-wrapped even when nullable, so it
	// needs no guard.
	guardNullable := isOptionWrapped(argType)
	operand := argName
	if guardNullable {
		operand = argName + "_val"
	}

	var left, op, message string
	switch constraint.Op {
	case ast.MinLengthOp:
		left = fmt.Sprintf("%s.chars().count()", operand)
		op = "<"
		message = fmt.Sprintf("%s length must be >= %s", argName, parameter)
	case ast.MaxLengthOp:
		left = fmt.Sprintf("%s.chars().count()", operand)
		op = ">"
		message = fmt.Sprintf("%s length must be <= %s", argName, parameter)
	case ast.GreaterThanOp, ast.GreaterThanEqualOp, ast.LessThanOp, ast.LessThanEqualOp, ast.EqualOp, ast.NotEqualOp:
		// A nullable operand is bound by reference (&Option contents); a numeric
		// comparison needs the value, so it is dereferenced.
		left = operand
		if guardNullable {
			left = "*" + operand
		}
		op = negateComparisonOp(constraint.Op)
		message = fmt.Sprintf("%s must be %s %s", argName, string(constraint.Op), parameter)
	case ast.RegexMatchOp, ast.NotRegexMatchOp:
		return "", fmt.Errorf("rust builder: regex constraint %q is not supported until a later phase (needs a regex dependency)", constraint.Op)
	default:
		return "", fmt.Errorf("rust builder: constraint op %q is not supported", constraint.Op)
	}

	var buffer strings.Builder
	if guardNullable {
		fmt.Fprintf(&buffer, "        if let Some(%s) = &%s {\n", operand, argName)
		fmt.Fprintf(&buffer, "            if %s %s %s {\n", left, op, parameter)
		buffer.WriteString("    " + formatErrorPush(argName, message))
		buffer.WriteString("            }\n")
		buffer.WriteString("        }\n")
		return buffer.String(), nil
	}
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
func (jenny Builder) formatBuild(formatter *typeFormatter, builder ast.Builder, builderType, objectType string) string {
	// A builder whose target type implements a variant (a dataquery-variant
	// builder, e.g. LokiBuilder building a Loki query) implements
	// cog::Builder<Box<dyn variants::Dataquery>> and boxes the assembled object,
	// mirroring the Go target whose Build() returns the variant interface. This
	// makes the variant builder usable wherever a composable dataquery slot is
	// filled (which takes impl cog::Builder<Box<dyn variants::Dataquery>>). A
	// non-variant builder implements cog::Builder<ConcreteType> and returns the
	// concrete object.
	buildType := objectType
	okExpr := "self.internal.clone()"
	if builder.For.Type.ImplementsVariant() && builder.For.Type.ImplementedVariant() == string(ast.SchemaVariantDataQuery) {
		formatter.imports.Add("crate::cog::variants")
		buildType = "Box<dyn variants::Dataquery>"
		okExpr = "Box::new(self.internal.clone())"
	}

	var buffer strings.Builder
	fmt.Fprintf(&buffer, "impl cog::Builder<%s> for %s {\n", buildType, builderType)
	fmt.Fprintf(&buffer, "    fn build(&self) -> Result<%s, Vec<cog::BuildError>> {\n", buildType)
	buffer.WriteString("        if !self.errors.is_empty() {\n")
	buffer.WriteString("            return Err(self.errors.clone());\n")
	buffer.WriteString("        }\n\n")
	fmt.Fprintf(&buffer, "        Ok(%s)\n", okExpr)
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
		if shouldDelegate(context, arg.Type) {
			argType = formatter.formatBuilderArgType(arg.Type)
		} else if isComposableSlotType(arg.Type) {
			// A composable-slot option argument (e.g. a dataquery slot, or a
			// collection of them) is filled by a builder of the slot's variant value,
			// mirroring the Go/Python target's `cog.Builder[variants.Dataquery]` arg.
			// The setter calls build() on it and stores the boxed trait object (see
			// formatComposableSlotAssignment).
			argType = formatter.formatComposableSlotBuilderArgType(arg.Type)
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
	if assignment.Value.Argument != nil && shouldDelegate(context, assignment.Value.Argument.Type) {
		return jenny.formatDelegatedAssignment(field, assignment.Value.Argument, receiver)
	}

	// Composable-slot delegation: the option argument is a builder of the slot's
	// variant value (e.g. cog::Builder<Box<dyn variants::Dataquery>>). Build it
	// and store the boxed trait object, reusing the builder-delegation machinery
	// (the destination field type drives single/Vec/HashMap shaping and the leaf
	// build() call). Errors accumulate exactly as for ordinary delegation.
	if assignment.Value.Argument != nil && isComposableSlotType(assignment.Value.Argument.Type) {
		return jenny.formatDelegatedAssignment(field, assignment.Value.Argument, receiver)
	}

	fieldName := formatFieldName(field.Identifier)
	value := jenny.formatAssignmentValue(formatter, context, field.Type, assignment.Value)

	if isOptionWrapped(field.Type) && !valueIsAlreadyOption(assignment.Value) {
		return fmt.Sprintf("%s.internal.%s = Some(%s);", receiver, fieldName, value)
	}
	return fmt.Sprintf("%s.internal.%s = %s;", receiver, fieldName, value)
}

// valueIsAlreadyOption reports whether an assignment value is an option argument
// whose own type is already Option-wrapped (a nullable setter parameter rendered
// as Option<T>). Such a value is assigned directly to an Option-typed field; it
// must not be wrapped again in Some(...), which would produce Some(Option<T>).
func valueIsAlreadyOption(value ast.AssignmentValue) bool {
	return value.Argument != nil && isOptionWrapped(value.Argument.Type)
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
func checkFactorySupported(context languages.Context, factory ast.BuilderFactory) error {
	// A factory's option calls apply to the builder named by its BuilderRef (which
	// may differ from the builder the factory is attached to: the dashboardv2
	// Manifest factory hangs off the Dashboard builder but constructs and configures
	// a resource.Manifest builder). Validate each call against that target builder.
	target, found := context.Builders.LocateByName(factory.BuilderRef.ReferredPkg, factory.BuilderRef.ReferredType)
	if !found {
		return fmt.Errorf("factory %q references unknown builder %s.%s", factory.Name, factory.BuilderRef.ReferredPkg, factory.BuilderRef.ReferredType)
	}
	for _, call := range factory.OptionCalls {
		if _, found := target.OptionByName(call.Name); !found {
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
	// A factory constructs the builder named by its BuilderRef. When that is the
	// builder the factory hangs off (the promql `abs` shape), it is emitted as an
	// associated function returning Self via a Self::new() chain. When it differs
	// (the dashboardv2 Manifest factory, which configures a resource.Manifest
	// builder), it is emitted as a free function returning the target builder type,
	// mirroring the Go target's package-level factory function.
	target, _ := context.Builders.LocateByName(factory.BuilderRef.ReferredPkg, factory.BuilderRef.ReferredType)
	sameBuilder := factory.BuilderRef.ReferredPkg == builder.Package && factory.BuilderRef.ReferredType == builder.Name

	// Whether any option-call parameter must be built (force_build): such a factory
	// cannot be a pure chain (it has to accumulate the built value's errors and
	// early-return), so it is emitted in statement form.
	hasForceBuild := false
	for _, call := range factory.OptionCalls {
		for _, param := range call.Parameters {
			if param.Argument != nil && param.ForceBuild {
				hasForceBuild = true
			}
		}
	}

	args := jenny.formatArgs(formatter, factory.Args, context)
	factoryName := formatFieldName(factory.Name)

	// A factory whose (snake-cased) name collides with an option setter on the same
	// builder cannot be an associated function (two `impl` items with the same name:
	// E0592, e.g. dashboardv2 RowsLayout has both a `rows` option and a `rows`
	// factory). Emit it as a free function instead, mirroring the Go target's
	// package-level factory functions.
	_, nameClashesWithOption := builder.OptionByName(factory.Name)

	if sameBuilder && !hasForceBuild && !nameClashesWithOption {
		return jenny.formatChainedFactory(formatter, context, target, builderType, factory, factoryName, args)
	}
	return jenny.formatFreeFactory(formatter, context, target, factory, factoryName, args)
}

// formatChainedFactory emits an associated function returning Self via a
// Self::new() chain of preset option calls (the promql `abs` shape).
func (jenny Builder) formatChainedFactory(formatter *typeFormatter, context languages.Context, target ast.Builder, builderType string, factory ast.BuilderFactory, factoryName, args string) string {
	var buffer strings.Builder
	buffer.WriteString(formatComments(factory.Comments, ""))
	fmt.Fprintf(&buffer, "impl %s {\n", builderType)

	chain := []string{"Self::new()"}
	for _, call := range factory.OptionCalls {
		option, _ := target.OptionByName(call.Name)
		chain = append(chain, jenny.formatFactoryOptionCall(formatter, context, option, call))
	}

	// Emit the signature and the chained call on single lines; the FormatRustFiles
	// postprocessor (rustfmt) lays them out across lines when they exceed the max
	// width.
	fmt.Fprintf(&buffer, "    pub fn %s(%s) -> Self {\n", factoryName, args)
	fmt.Fprintf(&buffer, "        %s\n", strings.Join(chain, "."))
	buffer.WriteString("    }\n")
	buffer.WriteString("}")
	return buffer.String()
}

// formatFreeFactory emits a free function that constructs and configures the
// target builder, returning it (the dashboardv2 Manifest shape: a factory hung
// off one builder that produces a differently-typed builder). Options consume
// and return the builder by value, so configuration is a reassigning chain of
// statements; a force_build parameter builds its argument first, extending the
// builder's errors and returning early on failure.
func (jenny Builder) formatFreeFactory(formatter *typeFormatter, context languages.Context, target ast.Builder, factory ast.BuilderFactory, factoryName, args string) string {
	targetType := formatTypeName(factory.BuilderRef.ReferredType) + "Builder"
	if formatPackageName(factory.BuilderRef.ReferredPkg) != formatter.packageName {
		formatter.imports.Add(fmt.Sprintf("crate::builders::%s::%s", formatPackageName(factory.BuilderRef.ReferredPkg), targetType))
	}

	var buffer strings.Builder
	buffer.WriteString(formatComments(factory.Comments, ""))
	fmt.Fprintf(&buffer, "pub fn %s(%s) -> %s {\n", factoryName, args, targetType)
	fmt.Fprintf(&buffer, "    let mut builder = %s::new();\n", targetType)

	for _, call := range factory.OptionCalls {
		option, _ := target.OptionByName(call.Name)
		methodName := formatFieldName(call.Name)
		params := make([]string, 0, len(call.Parameters))
		for i, param := range call.Parameters {
			var argType ast.Type
			if i < len(option.Args) {
				argType = option.Args[i].Type
			}
			if param.Argument != nil && param.ForceBuild {
				// Build the argument up front, extending the builder's errors and
				// returning it on failure (matching the Go target's RecordError + early
				// return). The built value is then passed to the option.
				built := formatArgName(param.Argument.Name) + "_built"
				fmt.Fprintf(&buffer, "    let %s = match %s.build() {\n", built, formatArgName(param.Argument.Name))
				buffer.WriteString("        Ok(val) => val,\n")
				buffer.WriteString("        Err(mut err) => {\n")
				buffer.WriteString("            builder.errors.append(&mut err);\n")
				buffer.WriteString("            return builder;\n")
				buffer.WriteString("        }\n")
				buffer.WriteString("    };\n")
				params = append(params, built)
				continue
			}
			params = append(params, jenny.formatFactoryParameter(formatter, context, argType, param))
		}
		fmt.Fprintf(&buffer, "    builder = builder.%s(%s);\n", methodName, strings.Join(params, ", "))
	}

	buffer.WriteString("    builder\n")
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
	if len(path) == 0 {
		return false
	}
	leaf := path.Last()
	if leaf.Index != nil || !leaf.Type.IsArray() {
		return false
	}
	// Intermediate path elements (a nested struct/ref the leaf array lives inside,
	// e.g. dashboard Panel's FieldConfig.Overrides) must be plain field steps
	// (no index): the append walks into them, lazily initialising Option-wrapped
	// refs, then pushes onto the leaf array.
	for _, item := range path[:len(path)-1] {
		if item.Index != nil {
			return false
		}
	}
	return true
}

// accessChainSegments builds the dotted access-chain segments for the
// intermediate elements of a builder path (everything before the leaf),
// initialising Option-wrapped ref intermediates with get_or_insert_with so a
// missing intermediate is created and an existing one is reused (no clobber).
// This is shared by multi-level direct assignment and multi-level append.
func (jenny Builder) accessChainSegments(formatter *typeFormatter, path ast.Path, receiver string) []string {
	return jenny.walkIntermediates(formatter, []string{fmt.Sprintf("%s.internal", receiver)}, path[:len(path)-1])
}

// walkIntermediates appends one access-chain segment per intermediate path
// element to the given seed, lazily initialising Option-wrapped ref
// intermediates with get_or_insert_with (creating a missing intermediate and
// reusing an existing one without clobbering). The seed is the root expression
// the chain is rooted at (e.g. `<receiver>.internal` or a local variable).
func (jenny Builder) walkIntermediates(formatter *typeFormatter, seed []string, intermediates ast.Path) []string {
	segments := seed
	for _, item := range intermediates {
		fieldName := formatFieldName(item.Identifier)
		if isOptionWrapped(item.Type) && item.Type.IsRef() {
			ref := item.Type.AsRef()
			// formatRef imports the (module-qualified) type and returns its name.
			concreteType := formatter.formatRef(ref)
			segments = append(segments, fieldName, fmt.Sprintf("get_or_insert_with(%s::default)", concreteType))
		} else {
			segments = append(segments, fieldName)
		}
	}
	return segments
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
	leaf := assignment.Path.Last()
	elementType := leaf.Type.AsArray().ValueType

	// The push target is the leaf array, reached through any intermediate elements
	// (a single-level path resolves to <receiver>.internal.<field>; a multi-level
	// path, e.g. dashboard Panel's FieldConfig.Overrides, walks into the
	// intermediates first, lazily initialising Option-wrapped refs).
	segments := jenny.accessChainSegments(formatter, assignment.Path, receiver)
	segments = append(segments, formatFieldName(leaf.Identifier))
	target := strings.Join(segments, ".")

	// An envelope value constructs the element inline as a struct literal (its ref
	// type) from the envelope's per-field values, then pushes it. This is the
	// envelope_assignment shape: each WithVariable call builds a fresh Variable.
	if assignment.Value.Envelope != nil {
		literal := jenny.formatEnvelopeLiteral(formatter, context, *assignment.Value.Envelope)
		if isOptionWrapped(elementType) {
			literal = fmt.Sprintf("Some(%s)", literal)
		}
		return fmt.Sprintf("%s.push(%s);", target, literal)
	}

	// Builder delegation: the appended element is itself produced by a nested
	// builder (the ArrayToAppend veneer shape, e.g. promql's `arg`). Build the
	// argument (accumulating errors, returning early on failure) then push the
	// built value, mirroring the Go target which appends `argResource` after
	// building it.
	// A composable-slot element (e.g. a panel's targets taking a Dataquery-slot
	// builder) is built and the boxed trait object pushed, reusing the same
	// delegation machinery as ordinary builder delegation.
	if assignment.Value.Argument != nil && (shouldDelegate(context, assignment.Value.Argument.Type) || isComposableSlotType(assignment.Value.Argument.Type)) {
		// The value is built against the argument's type but wrapped (Some) against
		// the array's element type, matching the original per-element shape.
		return jenny.formatDelegatedFieldStatement(assignment.Value.Argument.Type, isOptionWrapped(elementType), assignment.Value.Argument, receiver, func(builtValue string) string {
			return fmt.Sprintf("%s.push(%s);", target, builtValue)
		})
	}

	value := jenny.formatAssignmentValue(formatter, context, elementType, assignment.Value)
	if isOptionWrapped(elementType) && !valueIsAlreadyOption(assignment.Value) {
		value = fmt.Sprintf("Some(%s)", value)
	}
	return fmt.Sprintf("%s.push(%s);", target, value)
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
		if value.Value.ConstructorFor != nil {
			return fmt.Errorf("nested constructor in an envelope field is not supported")
		}
		// A nested envelope (the OverrideByName shape: the outer envelope's Matcher
		// field is itself a MatcherConfig envelope) is rendered as a nested struct
		// literal; validate it recursively.
		if value.Value.Envelope != nil {
			if err := checkEnvelopeSupported(*value.Value.Envelope); err != nil {
				return err
			}
		}
	}
	return nil
}

// isKnownAnyPath reports the `known_any` shape: a two-level path whose first
// element is an `any`-typed field carrying a concrete ref TypeHint. The builder
// composes a value of the hinted type into the otherwise-untyped field.
func isKnownAnyPath(path ast.Path) bool {
	idx := knownAnyIndex(path)
	if idx < 0 {
		return false
	}
	// At least one element follows the known-any field (the leaf set on the hinted
	// type, possibly nested through ref intermediates inside it, e.g. heatmap's
	// options.legend.show).
	if idx >= len(path)-1 {
		return false
	}
	// Every element before the known-any field is walked as an access chain on the
	// builder; every element after it is walked on the deserialized hinted value.
	// All must be plain (non-indexed) field steps.
	for _, item := range path {
		if item.Index != nil {
			return false
		}
	}
	return true
}

// knownAnyIndex returns the index of the known-any field in a path: an
// `any`-typed element carrying a concrete ref TypeHint (the panelcfg
// `FieldConfig.Defaults.Custom` shape), or -1 if absent.
func knownAnyIndex(path ast.Path) int {
	for i, item := range path {
		if item.Index != nil || item.TypeHint == nil || !item.TypeHint.IsRef() {
			continue
		}
		if item.Type.IsScalar() && item.Type.AsScalar().ScalarKind == ast.KindAny {
			return i
		}
	}
	return -1
}

// formatKnownAnyAssignment renders the `known_any` assignment. The `any` field is
// rendered as serde_json::Value; the builder constructs the concrete hinted type,
// sets the nested field on it, and stores it back serialized into the Value field.
// serde_json::to_value cannot fail for a generated serializable struct, so the
// Result is unwrapped via expect with a descriptive message.
func (jenny Builder) formatKnownAnyAssignment(formatter *typeFormatter, context languages.Context, assignment ast.Assignment, receiver string) string {
	anyIdx := knownAnyIndex(assignment.Path)
	anyItem := assignment.Path[anyIdx]
	leaf := assignment.Path.Last()

	anyField := formatFieldName(anyItem.Identifier)

	ref := anyItem.TypeHint.AsRef()
	concreteType := formatter.formatRef(ref)
	formatter.imports.Add(fmt.Sprintf("crate::types::%s::%s", formatPackageName(ref.ReferredPkg), formatTypeName(ref.ReferredType)))

	// Walk the intermediate elements (everything before the known-any field) to the
	// container that holds it, lazily initialising Option-wrapped refs. The
	// known-any field itself is then the location we read-modify-write.
	prefix := jenny.accessChainSegments(formatter, assignment.Path[:anyIdx+1], receiver)
	container := strings.Join(prefix, ".")
	anyLocation := fmt.Sprintf("%s.%s", container, anyField)

	// The path after the known-any field is an access chain rooted at the
	// deserialized hinted value (`custom`): plain field steps, with Option-wrapped
	// ref intermediates lazily initialised, ending at the leaf the option sets.
	// heatmap's options.legend.show walks into legend (a ref) before show.
	innerPath := assignment.Path[anyIdx+1:]
	innerSegments := jenny.walkIntermediates(formatter, []string{"custom"}, innerPath[:len(innerPath)-1])
	innerSegments = append(innerSegments, formatFieldName(leaf.Identifier))
	leafLocation := strings.Join(innerSegments, ".")

	// A leaf whose value is produced by a nested builder (a delegable builder ref or
	// a composable slot, e.g. timeseries' hide_from taking a HideSeriesConfig
	// builder) must be built before assignment. The build is emitted as a
	// preamble statement that accumulates errors onto the receiver and returns it
	// early on failure, exactly like an ordinary delegated assignment; the built
	// local is then assigned into the hinted value.
	var preamble string
	var value string
	delegated := false
	arg := assignment.Value.Argument
	if arg != nil && (shouldDelegate(context, arg.Type) || isComposableSlotType(arg.Type)) {
		delegated = true
		built := formatArgName(arg.Name) + "_built"
		preamble = jenny.formatDelegatedFieldStatement(arg.Type, false, arg, receiver, func(builtValue string) string {
			return fmt.Sprintf("let %s = %s;", built, builtValue)
		})
		value = built
	} else {
		value = jenny.formatAssignmentValue(formatter, context, leaf.Type, assignment.Value)
		// The option argument for a disjunction-typed leaf is rendered as
		// serde_json::Value (an undiscriminated disjunction has no single Rust type),
		// but the hinted struct field is the concrete (named) disjunction type. When
		// the argument is Value and the resolved field type is not `any`, deserialize
		// the Value into the field type (defaulting on failure), e.g. a panelcfg
		// custom option's insert_nulls (bool|float) into GraphFieldConfiginsertNulls.
		if arg != nil {
			// The option argument is rendered as serde_json::Value when its type is an
			// undiscriminated disjunction or `any` (no single Rust type). The hinted
			// struct field, however, is the concrete (named) type, so the Value is
			// deserialized into it. This covers a panelcfg custom option whose value is
			// a disjunction, e.g. insertNulls (bool|float) into GraphFieldConfiginsertNulls.
			argNonNull := arg.Type
			argNonNull.Nullable = false
			argRendersAsValue := formatter.formatType(argNonNull) == "serde_json::Value"
			if argRendersAsValue {
				if fieldType, ok := jenny.knownAnyLeafType(context, ref, innerPath); ok {
					nonNull := fieldType
					nonNull.Nullable = false
					concrete := formatter.formatType(nonNull)
					if concrete != "serde_json::Value" {
						if isOptionWrapped(arg.Type) {
							// A nullable Value argument maps each Some(v) through from_value,
							// yielding Option<T>; the outer Some-wrap below is then skipped.
							value = fmt.Sprintf("%s.and_then(|v| serde_json::from_value::<%s>(v).ok())", value, concrete)
						} else {
							value = fmt.Sprintf("serde_json::from_value::<%s>(%s).unwrap_or_default()", concrete, value)
						}
					}
				}
			}
		}
	}
	// The leaf is a field of the (possibly deeply nested) hinted struct; whether it
	// is Option-wrapped is determined by that struct field, not the IR path leaf
	// type (which can lag the NotRequiredFieldAsNullableType pass that made the
	// emitted struct field optional). Resolve the field optionality from the schema.
	// A built (delegated) value is always a bare T, so it is wrapped whenever the
	// field is optional; a plain argument that is itself already Option<T> is not
	// re-wrapped.
	leafOptional := jenny.knownAnyLeafOptional(context, ref, innerPath)
	if leafOptional && (delegated || !valueIsAlreadyOption(assignment.Value)) {
		value = fmt.Sprintf("Some(%s)", value)
	}

	// The known-any field is stored as Option<serde_json::Value>. Multiple option
	// setters compose into the same field (every panelcfg custom-option setter
	// targets FieldConfig.Defaults.Custom), so this is a read-modify-write: parse
	// any existing value into the hinted type (defaulting when absent or
	// unparseable), set the leaf field, then reserialize. serde_json::to_value
	// cannot fail for a generated serializable struct, so the Result is unwrapped
	// via expect with a descriptive message.
	expectMsg := fmt.Sprintf("%s should serialize to JSON", concreteType)
	var buffer strings.Builder
	if preamble != "" {
		fmt.Fprintf(&buffer, "%s\n        ", preamble)
	}
	fmt.Fprintf(&buffer, "let mut custom = %s.clone().and_then(|raw| serde_json::from_value::<%s>(raw).ok()).unwrap_or_default();\n", anyLocation, concreteType)
	fmt.Fprintf(&buffer, "        %s = %s;\n", leafLocation, value)
	fmt.Fprintf(&buffer, "        %s = Some(serde_json::to_value(custom).expect(%q));", anyLocation, expectMsg)
	return buffer.String()
}

// knownAnyLeafOptional reports whether the field a known-any assignment ultimately
// sets is Option-wrapped in the emitted hinted struct. It walks the hinted struct
// (ref) along innerPath: each non-leaf element descends into a nested ref struct,
// and the final element names the leaf field whose schema-declared nullability is
// returned. It falls back to the path leaf's own type nullability if the struct or
// field cannot be resolved.
func (jenny Builder) knownAnyLeafOptional(context languages.Context, ref ast.RefType, innerPath ast.Path) bool {
	currentPkg, currentType := ref.ReferredPkg, ref.ReferredType
	for i, item := range innerPath {
		obj, found := context.LocateObject(currentPkg, currentType)
		if !found || !obj.Type.IsStruct() {
			return isOptionWrapped(item.Type)
		}
		var field ast.StructField
		fieldFound := false
		for _, f := range obj.Type.AsStruct().Fields {
			if f.Name == item.Identifier {
				field, fieldFound = f, true
				break
			}
		}
		if !fieldFound {
			return isOptionWrapped(item.Type)
		}
		if i == len(innerPath)-1 {
			return isOptionWrapped(field.Type)
		}
		if !field.Type.IsRef() {
			return isOptionWrapped(item.Type)
		}
		currentPkg, currentType = field.Type.AsRef().ReferredPkg, field.Type.AsRef().ReferredType
	}
	return false
}

// knownAnyLeafType resolves the schema type of the field a known-any assignment
// ultimately sets, walking the hinted struct (ref) along innerPath the same way
// as knownAnyLeafOptional. It reports false when the struct or field cannot be
// resolved.
func (jenny Builder) knownAnyLeafType(context languages.Context, ref ast.RefType, innerPath ast.Path) (ast.Type, bool) {
	currentPkg, currentType := ref.ReferredPkg, ref.ReferredType
	for i, item := range innerPath {
		obj, found := context.LocateObject(currentPkg, currentType)
		if !found || !obj.Type.IsStruct() {
			return ast.Type{}, false
		}
		var field ast.StructField
		fieldFound := false
		for _, f := range obj.Type.AsStruct().Fields {
			if f.Name == item.Identifier {
				field, fieldFound = f, true
				break
			}
		}
		if !fieldFound {
			return ast.Type{}, false
		}
		if i == len(innerPath)-1 {
			return field.Type, true
		}
		if !field.Type.IsRef() {
			return ast.Type{}, false
		}
		currentPkg, currentType = field.Type.AsRef().ReferredPkg, field.Type.AsRef().ReferredType
	}
	return ast.Type{}, false
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
	segments := jenny.accessChainSegments(formatter, assignment.Path, receiver)

	leaf := assignment.Path.Last()
	segments = append(segments, formatFieldName(leaf.Identifier))
	target := strings.Join(segments, ".")

	// A leaf produced by a nested builder (a delegable builder ref or composable
	// slot, e.g. panel's field_config.defaults.thresholds taking a ThresholdsConfig
	// builder) is built then assigned, accumulating errors and returning early on
	// failure - the same delegation machinery used for single-level assignments,
	// applied at the end of the deep access chain.
	// Whether the leaf field is Option-wrapped is taken from the owning struct (the
	// penultimate path element's ref), not the IR path leaf type, which can lag the
	// NotRequiredFieldAsNullableType pass that made the emitted field optional.
	leafOptional := jenny.pathLeafOptional(context, assignment.Path)

	arg := assignment.Value.Argument
	if arg != nil && (shouldDelegate(context, arg.Type) || isComposableSlotType(arg.Type)) {
		return jenny.formatDelegatedFieldStatement(leaf.Type, leafOptional, arg, receiver, func(builtValue string) string {
			return fmt.Sprintf("%s = %s;", target, builtValue)
		})
	}

	value := jenny.formatAssignmentValue(formatter, context, leaf.Type, assignment.Value)
	if leafOptional && !valueIsAlreadyOption(assignment.Value) {
		value = fmt.Sprintf("Some(%s)", value)
	}

	return fmt.Sprintf("%s = %s;", target, value)
}

// pathLeafOptional reports whether the leaf field of a multi-level path is
// Option-wrapped, taking the optionality from the field declared on the owning
// struct (the penultimate path element's ref type) rather than the path leaf's
// own type. The two can disagree because the rawtypes pass that nullifies a
// not-required field updates the struct definition the builder writes into. Falls
// back to the path leaf type's nullability when the owning struct or field cannot
// be resolved.
func (jenny Builder) pathLeafOptional(context languages.Context, path ast.Path) bool {
	leaf := path.Last()
	if len(path) < 2 {
		return isOptionWrapped(leaf.Type)
	}
	parent := path[len(path)-2]
	if !parent.Type.IsRef() {
		return isOptionWrapped(leaf.Type)
	}
	ref := parent.Type.AsRef()
	obj, found := context.LocateObject(ref.ReferredPkg, ref.ReferredType)
	if !found || !obj.Type.IsStruct() {
		return isOptionWrapped(leaf.Type)
	}
	for _, f := range obj.Type.AsStruct().Fields {
		if f.Name == leaf.Identifier {
			return isOptionWrapped(f.Type)
		}
	}
	return isOptionWrapped(leaf.Type)
}

// formatEnvelopeLiteral renders an assignment envelope as a single-line Rust
// struct literal of the envelope's ref type, setting one named field per
// envelope value and filling the rest from Default. The FormatRustFiles
// postprocessor (rustfmt) expands the literal across lines when it exceeds the
// max width.
func (jenny Builder) formatEnvelopeLiteral(formatter *typeFormatter, context languages.Context, envelope ast.AssignmentEnvelope) string {
	ref := envelope.Type.AsRef()
	// formatRef imports the (module-qualified) type and returns its name.
	typeName := formatter.formatRef(ref)

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
		var rendered string
		// An envelope field whose value is itself an envelope (the
		// OverrideByName/ByRegexp shape: the outer DashboardFieldConfigSourceOverrides
		// envelope's Matcher field is a nested MatcherConfig{ id, options } literal)
		// recurses into a nested struct literal.
		if value.Value.Envelope != nil {
			rendered = jenny.formatEnvelopeLiteral(formatter, context, *value.Value.Envelope)
		} else {
			rendered = jenny.formatAssignmentValue(formatter, context, field.Type, value.Value)
		}
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
		rendered := formatArgName(value.Argument.Name)
		// An argument that is already serde_json::Value (its own type is `any`) is
		// assigned to an `any` field as-is; only a typed argument is serialized
		// through coerceToAny. (Both an Option<Value> arg into an Option<Value> field
		// and a bare Value into a Value field need no conversion.)
		if value.Argument.Type.IsScalar() && value.Argument.Type.AsScalar().ScalarKind == ast.KindAny {
			return rendered
		}
		return jenny.coerceToAny(destType, rendered)
	}

	// A constructor-value assignment default-constructs a value of the named ref
	// type (the Go target emits New<Type>(), which returns that type's default
	// value; the Rust equivalent is <Type>::default(), backed by the
	// rawtypes-generated Default impl). This is how a builder seeds a nested
	// object before deep-assigning into it, e.g. dashboard Panel's FieldConfig.
	if value.ConstructorFor != nil {
		ref := *value.ConstructorFor
		// formatRef imports the (module-qualified) type and returns its name.
		typeName := formatter.formatRef(ref)
		return jenny.coerceToAny(destType, fmt.Sprintf("%s::default()", typeName))
	}

	// Constant assignment.
	return jenny.coerceToAny(destType, jenny.formatConstantValue(formatter, context, destType, value.Constant))
}

// coerceToAny converts a rendered value into a serde_json::Value when the
// destination field is `any` (the Rust representation of an untyped field is
// serde_json::Value). A typed value (a String matcher option, a default-built
// Options struct) assigned to such a field is serialized via serde_json::to_value,
// e.g. the OverrideByName veneer assigns a byName name into MatcherConfig.options
// (any), and a panelcfg constructor seeds Panel.options with an Options struct
// (any). serde_json::to_value cannot fail for a generated serializable value, so
// the Result is unwrapped via expect. A non-any destination is unchanged.
func (jenny Builder) coerceToAny(destType ast.Type, rendered string) string {
	if destType.IsScalar() && destType.AsScalar().ScalarKind == ast.KindAny {
		return fmt.Sprintf("serde_json::to_value(%s).expect(\"value should serialize to JSON\")", rendered)
	}
	return rendered
}

// formatConstantValue renders a constant assigned to a field of destType. A ref
// to an enum resolves to Enum::Variant; a String scalar yields an owned String;
// other scalars render as their literal.
func (jenny Builder) formatConstantValue(formatter *typeFormatter, context languages.Context, destType ast.Type, constant any) string {
	if destType.IsRef() {
		ref := destType.AsRef()
		if referred, found := context.LocateObject(ref.ReferredPkg, ref.ReferredType); found && referred.Type.IsEnum() {
			// formatRef imports the (module-qualified) enum type and returns its name;
			// the builder formatter sets importSamePackageRefs so even a same-package
			// enum is module-qualified and resolves from the builder module.
			typeName := formatter.formatRef(ref)
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
