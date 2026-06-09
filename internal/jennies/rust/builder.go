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
		switch assignment.Method {
		case ast.DirectAssignment:
			// In scope since Phase 4a.
		case ast.AppendAssignment:
			// Phase 4c: append a single value to a Vec field. The path is a single
			// (non-indexed) array-typed field; the value is appended via push.
			if !isAppendPath(assignment.Path) {
				return fmt.Errorf("rust builder %q: unsupported append path shape", builder.Name)
			}
			// Phase 4e: an append whose value is an envelope constructs the nested
			// element inline (a struct literal of the envelope's ref type) and pushes
			// it. Reject an envelope whose own field paths are not the simple
			// single-level shape this phase renders.
			if assignment.Value.Envelope != nil {
				if err := checkEnvelopeSupported(*assignment.Value.Envelope); err != nil {
					return fmt.Errorf("rust builder %q: %w", builder.Name, err)
				}
			}
			continue
		case ast.IndexAssignment:
			// Phase 4c: insert a value at a key into a HashMap field. The path is a
			// two-level path: a map-typed field followed by an indexed element whose
			// Index carries the key (a constant or an option argument).
			if !isIndexPath(assignment.Path) {
				return fmt.Errorf("rust builder %q: unsupported index path shape", builder.Name)
			}
			continue
		default:
			return fmt.Errorf("rust builder %q: assignment method %q is not supported until a later phase", builder.Name, assignment.Method)
		}
		// A factory constructor value (ConstructorFor) is a Phase 5 feature.
		if assignment.Value.ConstructorFor != nil {
			return fmt.Errorf("rust builder %q: constructor-value assignments are not supported until a later phase", builder.Name)
		}
		// Supported direct-assignment path shapes:
		//   - single-level field (the common case),
		//   - the `known_any` two-level composition shape,
		//   - a multi-level ref/struct path (Phase 4e), where intermediate
		//     Option-wrapped refs are lazily initialised via get_or_insert_with using
		//     the NilChecks the pipeline injects, and a bare intermediate is accessed
		//     directly. Every non-leaf element must be a plain (non-indexed) ref so
		//     the get_or_insert_with chain has a concrete default type; an anonymous
		//     inline struct intermediate is unsupported (Rust models it as
		//     serde_json::Value, which cannot carry typed nested fields).
		if !isKnownAnyPath(assignment.Path) {
			if len(assignment.Path) == 1 {
				if assignment.Path[0].Index != nil {
					return fmt.Errorf("rust builder %q: indexed single-level paths are not supported in direct assignment", builder.Name)
				}
			} else if err := checkMultiLevelPathSupported(assignment.Path); err != nil {
				return fmt.Errorf("rust builder %q: %w", builder.Name, err)
			}
		}
		// A nil-check whose initialised path element is not a plain ref (for example
		// an anonymous inline struct) has no idiomatic default-construction in Rust.
		for _, check := range assignment.NilChecks {
			if isKnownAnyPath(assignment.Path) {
				continue
			}
			leaf := check.Path.Last()
			if !leaf.Type.IsArray() && !leaf.Type.IsMap() && !leaf.Type.IsRef() {
				return fmt.Errorf("rust builder %q: nil-check on non-ref intermediate %q is not supported", builder.Name, check.Path.String())
			}
		}
		// Builder delegation (an option argument whose type, or a collection's
		// element type, resolves to a builder) is in scope as of Phase 4d for the
		// ref / array-of-ref / map-of-ref shapes. A disjunction-typed delegated
		// argument (e.g. `Builder<T> | string`) has no clean idiomatic-Rust mapping
		// and is rejected (its fixture is skipped, matching the Go target which
		// eliminates disjunctions via a compiler pass the Rust target does not run).
		if assignment.Value.Argument != nil && context.ResolveToBuilder(assignment.Value.Argument.Type) {
			if !isDelegableType(assignment.Value.Argument.Type) {
				return fmt.Errorf("rust builder %q: disjunction-typed builder delegation is not supported", builder.Name)
			}
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
	fmt.Fprintf(&buffer, "impl %sBuilder {\n", formatTypeName(builder.Name))

	args := jenny.formatArgs(formatter, option.Args, context)
	methodName := formatFieldName(option.Name)
	if args == "" {
		fmt.Fprintf(&buffer, "    pub fn %s(mut self) -> Self {\n", methodName)
	} else {
		// rustfmt keeps the signature on one line only while it fits the 100-column
		// max width; beyond that it puts each parameter (including `mut self`) on its
		// own line. Emit the form rustfmt would settle on so the golden is
		// `cargo fmt --check` clean (builder-delegated args can be long, e.g. a
		// nested `Vec<Vec<impl cog::Builder<T>>>`).
		signature := fmt.Sprintf("    pub fn %s(mut self, %s) -> Self {\n", methodName, args)
		if len(signature)-1 <= 100 {
			buffer.WriteString(signature)
		} else {
			fmt.Fprintf(&buffer, "    pub fn %s(\n        mut self,\n        %s,\n    ) -> Self {\n", methodName, args)
		}
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
	value := jenny.formatValue(formatter, context, field.Type, assignment.Value)

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
	argName := formatArgName(arg.Name)

	const indent = "        "
	var buffer strings.Builder

	// Every statement is emitted fully indented; the caller prefixes the block's
	// first line with the same indentation, so the leading indent of the first line
	// is trimmed before returning to avoid a double indent.
	value := jenny.formatDelegatedValue(&buffer, field.Type, argName, receiver, indent, 0)

	if isOptionWrapped(field.Type) {
		fmt.Fprintf(&buffer, "%s%s.internal.%s = Some(%s);", indent, receiver, fieldName, value)
	} else {
		fmt.Fprintf(&buffer, "%s%s.internal.%s = %s;", indent, receiver, fieldName, value)
	}
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
		literal := jenny.formatEnvelopeLiteral(formatter, context, *assignment.Value.Envelope, "        ")
		if isOptionWrapped(elementType) {
			literal = fmt.Sprintf("Some(%s)", literal)
		}
		return fmt.Sprintf("%s.internal.%s.push(%s);", receiver, fieldName, literal)
	}

	value := jenny.formatValue(formatter, context, elementType, assignment.Value)
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

	value := jenny.formatValue(formatter, context, mapType.ValueType, assignment.Value)
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
	// Build the access chain as an ordered list of dotted segments so the assignment
	// can be rendered either inline or, when rustfmt would break the method chain,
	// one segment per line. `hasCall` records whether any segment is a method call
	// (get_or_insert_with): rustfmt only multi-lines a chain that contains a call.
	segments := []string{fmt.Sprintf("%s.internal", receiver)}
	hasCall := false

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
			// Two chain segments: the field access and the get_or_insert_with call.
			// rustfmt lays a broken method chain out one dot-segment per line, so the
			// field and the call must be separate segments to match its output.
			segments = append(segments, fieldName)
			segments = append(segments, fmt.Sprintf("get_or_insert_with(%s::default)", concreteType))
			hasCall = true
		} else {
			segments = append(segments, fieldName)
		}
	}

	leaf := assignment.Path.Last()
	segments = append(segments, formatFieldName(leaf.Identifier))

	value := jenny.formatValue(formatter, context, leaf.Type, assignment.Value)
	if isOptionWrapped(leaf.Type) {
		value = fmt.Sprintf("Some(%s)", value)
	}

	inline := fmt.Sprintf("%s = %s;", strings.Join(segments, "."), value)

	// rustfmt keeps the assignment inline only when the whole statement fits its
	// chain layout. A chain containing a method call (get_or_insert_with) is laid
	// out one segment per line once it no longer fits rustfmt's chain_width (60)
	// measured from the 8-space statement indent. Match that: a call-bearing chain
	// wider than 60 (statement body, excluding the leading indent) breaks; a
	// plain field chain (no call) always stays inline.
	const stmtIndent = "        "
	if !hasCall || len(inline) <= 60 {
		return inline
	}

	var buffer strings.Builder
	buffer.WriteString(segments[0])
	buffer.WriteString("\n")
	for i, seg := range segments[1:] {
		fmt.Fprintf(&buffer, "%s    .%s", stmtIndent, seg)
		if i < len(segments)-2 {
			buffer.WriteString("\n")
		}
	}
	fmt.Fprintf(&buffer, " = %s;", value)
	return buffer.String()
}

// formatEnvelopeLiteral renders an assignment envelope as a Rust struct literal
// of the envelope's ref type, setting one named field per envelope value and
// filling the rest from Default. The literal is emitted at `indent` so a
// multi-field literal nests cleanly inside the enclosing statement; rustfmt
// keeps a short single-field literal inline, so the form here matches what
// rustfmt would settle on.
func (jenny Builder) formatEnvelopeLiteral(formatter *typeFormatter, context languages.Context, envelope ast.AssignmentEnvelope, indent string) string {
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
		rendered := jenny.formatValue(formatter, context, field.Type, value.Value)
		if isOptionWrapped(field.Type) {
			rendered = fmt.Sprintf("Some(%s)", rendered)
		}
		if rendered == fieldName {
			inits = append(inits, fieldName)
		} else {
			inits = append(inits, fmt.Sprintf("%s: %s", fieldName, rendered))
		}
	}

	// rustfmt collapses a struct literal with no rest expression onto a single
	// line while it fits the 100-column max width; a literal carrying a
	// `..Default::default()` rest stays multi-line. Emit the form rustfmt would
	// settle on so the golden is `cargo fmt --check` clean. The inline width is
	// estimated from the enclosing statement indent plus the literal text; this is
	// conservative and only collapses when it comfortably fits.
	if !needsRest {
		singleLine := fmt.Sprintf("%s { %s }", typeName, strings.Join(inits, ", "))
		if len(indent)+len("self.internal..push();")+len(singleLine) <= 100 {
			return singleLine
		}
	}

	var buffer strings.Builder
	fmt.Fprintf(&buffer, "%s {\n", typeName)
	for _, init := range inits {
		fmt.Fprintf(&buffer, "%s    %s,\n", indent, init)
	}
	if needsRest {
		fmt.Fprintf(&buffer, "%s    ..Default::default()\n", indent)
	}
	fmt.Fprintf(&buffer, "%s}", indent)
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
