package rust

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/jennies/common"
	"github.com/grafana/cog/internal/languages"
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

		filename := filepath.Join("src", "types", schema.Package+".rs")
		files = append(files, *codejen.NewFile(filename, output, jenny))
	}

	return files, nil
}

func (jenny RawTypes) generateSchema(context languages.Context, schema *ast.Schema) ([]byte, error) {
	imports := newImportMap()
	formatter := newTypeFormatter(context, imports)

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
	case object.Type.IsScalar():
		return jenny.formatTypeAlias(formatter, object), nil
	default:
		return "", fmt.Errorf("rust rawtypes: unsupported top-level object kind %q for %q (Phase 3b+)", object.Type.Kind, object.Name)
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
		formatConstValue(scalar.Value, scalar.ScalarKind),
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
			buffer.WriteString(jenny.formatStructField(formatter, field))
		}
		buffer.WriteString("}")
	}

	if defaultImpl := jenny.formatDefaultImpl(object); defaultImpl != "" {
		buffer.WriteString("\n\n")
		buffer.WriteString(defaultImpl)
	}

	return buffer.String()
}

func (jenny RawTypes) formatStructField(formatter *typeFormatter, field ast.StructField) string {
	var buffer strings.Builder

	buffer.WriteString(formatComments(field.Comments, "    "))

	for _, attr := range fieldSerdeAttributes(field) {
		fmt.Fprintf(&buffer, "    %s\n", attr)
	}

	fmt.Fprintf(&buffer, "    pub %s: %s,\n", formatFieldName(field.Name), formatter.formatType(field.Type))

	return buffer.String()
}

// fieldSerdeAttributes returns the serde attribute lines for a struct field:
// rename for non-snake-case keys and skip/default handling for optionals.
func fieldSerdeAttributes(field ast.StructField) []string {
	var attrs []string

	if rustFieldNeedsRename(field.Name) {
		attrs = append(attrs, fmt.Sprintf("#[serde(rename = %q)]", field.Name))
	}

	if field.Type.Nullable {
		attrs = append(attrs, `#[serde(default, skip_serializing_if = "Option::is_none")]`)
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
		if field.Type.Default != nil {
			return true
		}
	}
	return false
}

func (jenny RawTypes) formatDefaultImpl(object ast.Object) string {
	fields := object.Type.AsStruct().Fields
	if !needsManualDefault(fields) {
		return ""
	}

	var buffer strings.Builder
	fmt.Fprintf(&buffer, "impl Default for %s {\n", formatTypeName(object.Name))
	buffer.WriteString("    fn default() -> Self {\n")
	buffer.WriteString("        Self {\n")

	for _, field := range fields {
		fmt.Fprintf(&buffer, "            %s: %s,\n", formatFieldName(field.Name), defaultExpression(field.Type))
	}

	buffer.WriteString("        }\n")
	buffer.WriteString("    }\n")
	buffer.WriteString("}")

	return buffer.String()
}

// defaultExpression returns the Rust expression used to initialise a field in a
// manual Default impl.
func defaultExpression(typeDef ast.Type) string {
	if typeDef.IsConcreteScalar() {
		return constScalarLiteral(typeDef.AsScalar())
	}

	if typeDef.Default != nil {
		return defaultLiteral(typeDef)
	}

	return "Default::default()"
}

// defaultLiteral renders an explicit IR default into a Rust expression, using
// the field's scalar kind to pick the correct literal form.
func defaultLiteral(typeDef ast.Type) string {
	if !typeDef.IsScalar() {
		return formatValue(typeDef.Default)
	}

	kind := typeDef.AsScalar().ScalarKind
	if kind == ast.KindString {
		return fmt.Sprintf("%s.to_string()", formatValue(typeDef.Default))
	}

	return formatScalarValue(typeDef.Default, kind)
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

func formatConstValue(value any, kind ast.ScalarKind) string {
	return formatScalarValue(value, kind)
}

// constScalarLiteral renders a constant scalar field's value for use in a Default
// impl, where String fields need an owned String.
func constScalarLiteral(scalar ast.ScalarType) string {
	if scalar.ScalarKind == ast.KindString {
		return fmt.Sprintf("%s.to_string()", formatValue(scalar.Value))
	}
	return formatScalarValue(scalar.Value, scalar.ScalarKind)
}
