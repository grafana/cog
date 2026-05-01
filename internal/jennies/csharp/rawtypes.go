package csharp

import (
	"path/filepath"
	"sort"

	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/jennies/template"
	"github.com/grafana/cog/internal/languages"
)

// RawTypes is the C# equivalent of [java.RawTypes]: it walks every
// schema in the input context and emits one `.cs` file per object
// (struct, enum, ref, intersection) plus an aggregated `Constants.cs`
// per package collecting concrete scalar values.
//
// Phase 2 scope: structs, enums, intersections-as-classes,
// concrete-scalar constants. Disjunctions, JSON converters and
// equality overrides are deliberately deferred.
type RawTypes struct {
	config Config
	tmpl   *template.Template

	// Reset per-file in genFilesForObject.
	imports       *importMap
	typeFormatter *typeFormatter
}

func (jenny RawTypes) JennyName() string {
	return "CSharpRawTypes"
}

func (jenny RawTypes) Generate(context languages.Context) (codejen.Files, error) {
	files := make(codejen.Files, 0)
	jenny.typeFormatter = newTypeFormatter(context, jenny.config, nil)

	for _, schema := range context.Schemas {
		out, err := jenny.genFilesForSchema(schema)
		if err != nil {
			return nil, err
		}
		files = append(files, out...)
	}

	return files, nil
}

func (jenny RawTypes) genFilesForSchema(schema *ast.Schema) (codejen.Files, error) {
	files := make(codejen.Files, 0)
	scalars := make(map[string]ast.ScalarType)
	pkgFolder := formatPackageName(schema.Package)
	namespace := jenny.config.formatNamespace(schema.Package)

	var iterErr error
	schema.Objects.Iterate(func(_ string, object ast.Object) {
		if iterErr != nil {
			return
		}

		// Concrete scalar definitions are aggregated into a single
		// Constants.cs per package.
		if object.Type.IsScalar() {
			if object.Type.AsScalar().IsConcrete() {
				scalars[object.Name] = object.Type.AsScalar()
			}
			return
		}

		output, err := jenny.generateObject(namespace, object)
		if err != nil {
			iterErr = err
			return
		}
		if output == nil {
			return
		}

		filename := filepath.Join(jenny.config.ProjectPath, pkgFolder, formatObjectName(object.Name)+".cs")
		files = append(files, *codejen.NewFile(filename, output, jenny))
	})
	if iterErr != nil {
		return nil, iterErr
	}

	if len(scalars) > 0 {
		output, err := jenny.formatConstants(namespace, scalars)
		if err != nil {
			return nil, err
		}
		filename := filepath.Join(jenny.config.ProjectPath, pkgFolder, "Constants.cs")
		files = append(files, *codejen.NewFile(filename, output, jenny))
	}

	return files, nil
}

func (jenny RawTypes) generateObject(namespace string, object ast.Object) ([]byte, error) {
	switch object.Type.Kind {
	case ast.KindStruct:
		return jenny.formatStruct(namespace, object)
	case ast.KindEnum:
		return jenny.formatEnum(namespace, object)
	case ast.KindRef:
		return jenny.formatRefAlias(namespace, object)
	case ast.KindIntersection:
		return jenny.formatIntersection(namespace, object)
	}
	return nil, nil
}

// renderClass renders the class template against the given model.
// All field/arg type expressions in `data` must already be formatted
// strings — the template no longer calls formatType so that imports
// are guaranteed to be populated before the {{ .Imports }} block runs.
func (jenny RawTypes) renderClass(data classTemplate) ([]byte, error) {
	return jenny.tmpl.RenderAsBytes("types/class.tmpl", data)
}

func (jenny RawTypes) formatStruct(namespace string, object ast.Object) ([]byte, error) {
	jenny.imports = newImportMap(jenny.config.namespaceRoot(), namespace)
	jenny.typeFormatter = jenny.typeFormatter.withImports(jenny.imports)

	fields := object.Type.AsStruct().Fields
	cls := jenny.buildClassTemplate(namespace, formatObjectName(object.Name), object.Comments, nil, fields)

	return jenny.renderClass(cls)
}

// formatRefAlias emits a thin C# subclass for a top-level type that is
// just a ref to another type. Mirrors the Java jenny's behaviour.
func (jenny RawTypes) formatRefAlias(namespace string, object ast.Object) ([]byte, error) {
	jenny.imports = newImportMap(jenny.config.namespaceRoot(), namespace)
	jenny.typeFormatter = jenny.typeFormatter.withImports(jenny.imports)

	ref := object.Type.AsRef()
	parent := jenny.typeFormatter.formatReference(ref)

	cls := classTemplate{
		Namespace: namespace,
		Name:      formatObjectName(object.Name),
		Imports:   jenny.imports,
		Comments:  object.Comments,
		Extends:   []string{parent},
	}
	return jenny.renderClass(cls)
}

func (jenny RawTypes) formatIntersection(namespace string, object ast.Object) ([]byte, error) {
	jenny.imports = newImportMap(jenny.config.namespaceRoot(), namespace)
	jenny.typeFormatter = jenny.typeFormatter.withImports(jenny.imports)

	intersection := object.Type.AsIntersection()
	extends := make([]string, 0)
	fields := make([]ast.StructField, 0)
	inheritedFieldNames := make(map[string]bool)

	for _, branch := range intersection.Branches {
		if branch.IsRef() {
			if obj, found := jenny.typeFormatter.context.LocateObjectByRef(branch.AsRef()); found && obj.Type.IsStruct() {
				for _, field := range obj.Type.AsStruct().Fields {
					inheritedFieldNames[field.Name] = true
				}
			}
		}
	}

	for _, branch := range intersection.Branches {
		switch branch.Kind {
		case ast.KindRef:
			extends = append(extends, jenny.typeFormatter.formatReference(branch.AsRef()))
		case ast.KindStruct:
			for _, field := range branch.AsStruct().Fields {
				if inheritedFieldNames[field.Name] {
					continue
				}
				fields = append(fields, field)
			}
		}
	}

	cls := jenny.buildClassTemplate(namespace, formatObjectName(object.Name), object.Comments, extends, fields)
	return jenny.renderClass(cls)
}

// buildClassTemplate populates the template-facing classTemplate struct
// from a struct's fields, generating the field declarations and the
// two constructors' assignment lists.
func (jenny RawTypes) buildClassTemplate(namespace string, name string, comments []string, extends []string, fields []ast.StructField) classTemplate {
	clsFields := make([]classField, 0, len(fields))
	defaults := make([]assignment, 0, len(fields))
	args := make([]classArg, 0, len(fields))
	argAssignments := make([]assignment, 0, len(fields))

	for _, field := range fields {
		fieldName := formatFieldName(field.Name)
		argName := formatArgName(field.Name)
		typeExpr := jenny.typeFormatter.formatFieldType(field.Type)

		clsFields = append(clsFields, classField{
			Name:     fieldName,
			Type:     typeExpr,
			Comments: field.Comments,
		})

		// Default-constructor assignment: explicit default, or empty value.
		var defaultExpr string
		if field.Type.Default != nil {
			defaultExpr = jenny.formatDefault(field.Type, field.Type.Default)
		} else if field.Required && !field.Type.Nullable {
			defaultExpr = jenny.typeFormatter.emptyValueForType(field.Type)
		}
		if defaultExpr != "" {
			defaults = append(defaults, assignment{Name: fieldName, Value: defaultExpr})
		}

		args = append(args, classArg{Name: argName, Type: typeExpr})
		argAssignments = append(argAssignments, assignment{
			Name:         fieldName,
			ValueFromArg: argName,
		})
	}

	return classTemplate{
		Namespace:          namespace,
		Name:               name,
		Imports:            jenny.imports,
		Comments:           comments,
		Extends:            extends,
		Fields:             clsFields,
		DefaultAssignments: defaults,
		HasArgsConstructor: len(args) > 0,
		Args:               args,
		ArgAssignments:     argAssignments,
	}
}

// formatDefault renders a literal expression for a field's default value.
// Only scalar defaults are supported in Phase 2 — refs and enums fall
// back to the type's "empty" value.
func (jenny RawTypes) formatDefault(t ast.Type, value any) string {
	if t.IsScalar() {
		return formatScalarValue(t.AsScalar().ScalarKind, value)
	}
	return jenny.typeFormatter.emptyValueForType(t)
}

func (jenny RawTypes) formatEnum(namespace string, object ast.Object) ([]byte, error) {
	enumDef := object.Type.AsEnum()
	if len(enumDef.Values) == 0 {
		return nil, nil
	}

	isString := enumDef.Values[0].Type.IsScalar() && enumDef.Values[0].Type.AsScalar().ScalarKind == ast.KindString

	values := make([]EnumValue, 0, len(enumDef.Values))
	for _, v := range enumDef.Values {
		var raw any
		if isString {
			s, _ := v.Value.(string)
			raw = s
		} else {
			raw = formatScalarValue(v.Type.AsScalar().ScalarKind, v.Value)
		}
		values = append(values, EnumValue{
			Name:     formatEnumMemberName(v.Name),
			RawValue: raw,
		})
	}

	data := enumTemplate{
		Namespace:       namespace,
		Name:            formatObjectName(object.Name),
		Comments:        object.Comments,
		IsString:        isString,
		NeedsEnumMember: isString,
		Values:          values,
	}

	return jenny.tmpl.RenderAsBytes("types/enum.tmpl", data)
}

func (jenny RawTypes) formatConstants(namespace string, scalars map[string]ast.ScalarType) ([]byte, error) {
	consts := make([]constant, 0, len(scalars))
	for name, scalar := range scalars {
		consts = append(consts, constant{
			Name:  name,
			Type:  formatScalarType(scalar),
			Value: formatScalarValue(scalar.ScalarKind, scalar.Value),
		})
	}
	sort.SliceStable(consts, func(i, j int) bool { return consts[i].Name < consts[j].Name })

	return jenny.tmpl.RenderAsBytes("types/constants.tmpl", constantsTemplate{
		Namespace: namespace,
		Name:      "Constants",
		Constants: consts,
	})
}

// formatEnumMemberName converts a raw enum-value name into a valid C#
// PascalCased identifier.
func formatEnumMemberName(name string) string {
	if name == "" {
		return "Default"
	}
	formatted := pascalCase(name)
	// Identifier cannot start with a digit.
	if formatted == "" || (formatted[0] >= '0' && formatted[0] <= '9') {
		formatted = "V" + formatted
	}
	return escapeVarName(formatted)
}
