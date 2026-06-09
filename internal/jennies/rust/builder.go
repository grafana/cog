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
// single-level field paths. Constraints-as-guards, array-append, index
// assignment, builder delegation, factories, envelopes and variants are emitted
// by later chunks; a builder requiring any of them is rejected with an error so
// out-of-scope fixtures are skipped explicitly rather than mis-generated.
type Builder struct {
	config          Config
	apiRefCollector *common.APIReferenceCollector
}

func (jenny Builder) JennyName() string {
	return "RustBuilder"
}

func (jenny Builder) Generate(context languages.Context) (codejen.Files, error) {
	files := make(codejen.Files, 0, len(context.Builders))

	for _, builder := range context.Builders {
		output, err := jenny.generateBuilder(context, builder)
		if err != nil {
			return nil, err
		}

		filename := filepath.Join("src", "builders", formatPackageName(builder.Package)+".rs")
		files = append(files, *codejen.NewFile(filename, output, jenny))
	}

	return files, nil
}

func (jenny Builder) generateBuilder(context languages.Context, builder ast.Builder) ([]byte, error) {
	if err := jenny.rejectUnsupported(builder); err != nil {
		return nil, err
	}

	imports := newImportMap()
	imports.Add("crate::cog")
	formatter := newTypeFormatter(context, imports, builder.Package)

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
	body.WriteString("}\n\n")

	body.WriteString(jenny.formatConstructor(formatter, context, builder, builderType, objectType))

	for _, option := range builder.Options {
		method, err := jenny.formatOption(formatter, context, builder, option)
		if err != nil {
			return nil, err
		}
		body.WriteString("\n\n")
		body.WriteString(method)
	}

	body.WriteString("\n\n")
	body.WriteString(jenny.formatBuild(builderType, objectType))
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

// rejectUnsupported reports an error for any builder feature outside Phase 4a so
// that out-of-scope fixtures fail loudly (and are skipped in the suite) rather
// than emitting incorrect code.
func (jenny Builder) rejectUnsupported(builder ast.Builder) error {
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
		if len(assignment.Path) != 1 || assignment.Path[0].Index != nil {
			return fmt.Errorf("rust builder %q: only single-level field assignments are supported in this phase", builder.Name)
		}
		if assignment.Value.Envelope != nil || assignment.Value.ConstructorFor != nil {
			return fmt.Errorf("rust builder %q: envelope/constructor assignments are not supported until a later phase", builder.Name)
		}
		if len(assignment.NilChecks) != 0 {
			return fmt.Errorf("rust builder %q: nil-check guards are not supported until a later phase", builder.Name)
		}
		// A direct assignment of an option argument whose type is a ref to a struct
		// that has its own builder is builder delegation, handled in a later chunk.
		if assignment.Value.Argument != nil && assignment.Value.Argument.Type.IsRef() {
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
		buffer.WriteString("        ")
		buffer.WriteString(jenny.formatAssignment(formatter, context, assignment, "self"))
		buffer.WriteString("\n")
	}

	buffer.WriteString("\n        self\n")
	buffer.WriteString("    }\n")
	buffer.WriteString("}")

	return buffer.String(), nil
}

// formatBuild implements the runtime cog::Builder<T> trait. Phase 4a performs no
// validation, so it clones the assembled object. Constraint-as-guard validation
// and required-field nil checks are layered in by later chunks.
func (jenny Builder) formatBuild(builderType, objectType string) string {
	var buffer strings.Builder
	fmt.Fprintf(&buffer, "impl cog::Builder<%s> for %s {\n", objectType, builderType)
	fmt.Fprintf(&buffer, "    fn build(&self) -> Result<%s, Vec<cog::BuildError>> {\n", objectType)
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
	field := assignment.Path[0]
	fieldName := formatFieldName(field.Identifier)
	value := jenny.formatValue(formatter, context, field.Type, assignment.Value)

	if isOptionWrapped(field.Type) {
		return fmt.Sprintf("%s.internal.%s = Some(%s);", receiver, fieldName, value)
	}
	return fmt.Sprintf("%s.internal.%s = %s;", receiver, fieldName, value)
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
