package java

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/jennies/common"
	"github.com/grafana/cog/internal/jennies/template"
	"github.com/grafana/cog/internal/languages"
	"github.com/grafana/cog/internal/tools"
)

type RawTypes struct {
	config  Config
	tmpl    *template.Template
	imports *common.DirectImportMap

	typeFormatter  *typeFormatter
	jsonMarshaller JSONMarshaller
}

func (jenny RawTypes) JennyName() string {
	return "JavaRawTypes"
}

func (jenny RawTypes) Generate(context languages.Context) (codejen.Files, error) {
	files := make(codejen.Files, 0)
	jenny.imports = NewImportMap(jenny.config.PackagePath)
	jenny.typeFormatter = createFormatter(context, jenny.config)
	jenny.jsonMarshaller = JSONMarshaller{
		config:        jenny.config,
		tmpl:          jenny.tmpl,
		typeFormatter: jenny.typeFormatter,
	}

	for _, schema := range context.Schemas {
		output, err := jenny.genFilesForSchema(schema)
		if err != nil {
			return nil, err
		}

		files = append(files, output...)
	}

	return files, nil
}

func (jenny RawTypes) getTemplate() *template.Template {
	return jenny.tmpl.Funcs(map[string]any{
		"formatBuilderFieldType":        jenny.typeFormatter.formatBuilderFieldType,
		"formatType":                    jenny.typeFormatter.formatFieldType,
		"typeHasBuilder":                jenny.typeFormatter.typeHasBuilder,
		"fillNullableAnnotationPattern": jenny.typeFormatter.fillNullableAnnotationPattern,
	})
}

func (jenny RawTypes) genFilesForSchema(schema *ast.Schema) (codejen.Files, error) {
	var err error
	files := make(codejen.Files, 0)
	scalars := make(map[string]ast.ScalarType)

	packageMapper := func(pkg string, class string) string {
		if jenny.imports.IsIdentical(pkg, schema.Package) {
			return ""
		}

		return jenny.imports.Add(class, pkg)
	}

	jenny.typeFormatter = jenny.typeFormatter.withPackageMapper(packageMapper)

	schema.Objects.Iterate(func(_ string, object ast.Object) {
		jenny.imports = NewImportMap(jenny.config.PackagePath)
		if object.Type.IsScalar() {
			if object.Type.AsScalar().IsConcrete() {
				scalars[object.Name] = object.Type.AsScalar()
			}
			return
		}

		pkg := formatPackageName(schema.Package)
		output, innerErr := jenny.generateSchema(pkg, schema.Metadata.Identifier, object)
		if innerErr != nil {
			err = innerErr
			return
		}

		filename := filepath.Join(jenny.config.ProjectPath, pkg, fmt.Sprintf("%s.java", tools.UpperCamelCase(object.Name)))

		files = append(files, *codejen.NewFile(filename, output, jenny))
	})

	if err != nil {
		return nil, err
	}

	if len(scalars) > 0 {
		output, err := jenny.formatScalars(schema.Package, scalars)
		if err != nil {
			return nil, err
		}

		filename := filepath.Join(jenny.config.ProjectPath, strings.ToLower(schema.Package), "Constants.java")
		files = append(files, *codejen.NewFile(filename, output, jenny))
	}

	return files, nil
}

func (jenny RawTypes) generateSchema(pkg string, identifier string, object ast.Object) ([]byte, error) {
	switch object.Type.Kind {
	case ast.KindStruct:
		return jenny.formatStruct(pkg, identifier, object)
	case ast.KindEnum:
		return formatEnum(jenny.config.formatPackage(pkg), object, jenny.getTemplate())
	case ast.KindRef:
		return jenny.formatReference(pkg, identifier, object)
	case ast.KindIntersection:
		return jenny.formatIntersection(pkg, identifier, object)
	}

	return nil, nil
}

func (jenny RawTypes) formatStruct(pkg string, identifier string, object ast.Object) ([]byte, error) {
	return jenny.getTemplate().RenderAsBytes("types/class.tmpl", ClassTemplate{
		Package:                 jenny.config.formatPackage(pkg),
		RawPackage:              pkg,
		Imports:                 jenny.imports,
		Name:                    tools.UpperCamelCase(object.Name),
		Fields:                  object.Type.AsStruct().Fields,
		Comments:                object.Comments,
		Variant:                 jenny.getVariant(object.Type),
		Identifier:              identifier,
		Annotation:              jenny.jsonMarshaller.annotation(object.Type),
		ToJSONFunction:          jenny.jsonMarshaller.genToJSONFunction(object.Type),
		ShouldAddSerializer:     jenny.typeFormatter.objectNeedsCustomSerializer(object),
		ShouldAddDeserializer:   jenny.typeFormatter.objectNeedsCustomDeserializer(object),
		ShouldAddFactoryMethods: object.Type.IsStructGeneratedFromDisjunction(),
		Constructors:            jenny.constructors(object),
	})
}

func (jenny RawTypes) formatScalars(pkg string, scalars map[string]ast.ScalarType) ([]byte, error) {
	constants := make([]Constant, 0)
	for name, scalar := range scalars {
		constants = append(constants, Constant{
			Name:  name,
			Type:  formatScalarType(scalar),
			Value: scalar.Value,
		})
	}

	// To ensure deterministic output
	sort.SliceStable(constants, func(i, j int) bool {
		return constants[i].Name < constants[j].Name
	})

	return jenny.getTemplate().RenderAsBytes("types/constants.tmpl", ConstantTemplate{
		Package:   jenny.config.formatPackage(pkg),
		Name:      "Constants",
		Constants: constants,
	})
}

func (jenny RawTypes) formatReference(pkg string, identifier string, object ast.Object) ([]byte, error) {
	ref := object.Type.AsRef()
	reference := fmt.Sprintf("%s.%s", jenny.config.formatPackage(formatPackageName(ref.ReferredPkg)), formatObjectName(ref.ReferredType))

	return jenny.getTemplate().RenderAsBytes("types/class.tmpl", ClassTemplate{
		Package:    jenny.config.formatPackage(pkg),
		Imports:    jenny.imports,
		Name:       formatObjectName(object.Name),
		Extends:    []string{reference},
		Comments:   object.Comments,
		Variant:    jenny.getVariant(object.Type),
		Identifier: identifier,
	})
}

func (jenny RawTypes) formatIntersection(pkg string, identifier string, object ast.Object) ([]byte, error) {
	intersection := object.Type.AsIntersection()
	extensions := make([]string, 0)
	fields := make([]ast.StructField, 0)

	for _, branch := range intersection.Branches {
		switch branch.Kind {
		case ast.KindRef:
			extensions = append(extensions, jenny.typeFormatter.formatReference(branch.AsRef()))
		case ast.KindStruct:
			fields = append(fields, branch.AsStruct().Fields...)
		}
	}

	return jenny.getTemplate().RenderAsBytes("types/class.tmpl", ClassTemplate{
		Package:    jenny.config.formatPackage(pkg),
		Imports:    jenny.imports,
		Name:       object.Name,
		Extends:    extensions,
		Comments:   object.Comments,
		Fields:     fields,
		Variant:    jenny.getVariant(object.Type),
		Identifier: identifier,
	})
}

func (jenny RawTypes) getVariant(t ast.Type) string {
	variant := ""
	if t.ImplementsVariant() {
		variant = fmt.Sprintf("cog.variants.%s", tools.UpperCamelCase(t.ImplementedVariant()))
		variant = jenny.config.formatPackage(variant)
	}
	return variant
}

func (jenny RawTypes) constructors(object ast.Object) []ConstructorTemplate {
	if object.Type.IsStructGeneratedFromDisjunction() {
		return nil
	}

	fields := object.Type.AsStruct().Fields

	args := make([]ast.Argument, 0)
	assignments := make([]ConstructorAssignmentTemplate, 0)
	defaultConstructorAssignments := make([]ConstructorAssignmentTemplate, 0)
	for _, field := range fields {
		name := tools.LowerCamelCase(escapeVarName(field.Name))
		if field.Type.IsConstantRef() {
			assign := ConstructorAssignmentTemplate{
				Name:  name,
				Type:  field.Type,
				Value: jenny.typeFormatter.constantRefValue(field.Type.AsConstantRef()),
			}
			assignments = append(assignments, assign)
			defaultConstructorAssignments = append(defaultConstructorAssignments, assign)
			continue
		}

		args = append(args, ast.Argument{
			Name: name,
			Type: field.Type,
		})

		assignments = append(assignments, ConstructorAssignmentTemplate{
			Name:         name,
			Type:         field.Type,
			ValueFromArg: name,
		})

		if field.Type.Default != nil {
			defaultConstructorAssignments = append(defaultConstructorAssignments, ConstructorAssignmentTemplate{
				Name:  name,
				Type:  field.Type,
				Value: jenny.genDefaultForType(field.Type, field.Type.Default),
			})
		} else if field.Required && !field.Type.Nullable { // Fields without an explicit default, but that aren't allowed to be null
			defaultConstructorAssignments = append(defaultConstructorAssignments, ConstructorAssignmentTemplate{
				Name:  name,
				Type:  field.Type,
				Value: jenny.typeFormatter.emptyValueForType(field.Type, false),
			})
		}
	}

	constructors := []ConstructorTemplate{
		// Default constructor
		{
			Args:        []ast.Argument{},
			Assignments: defaultConstructorAssignments,
		},
	}

	if len(args) > 0 {
		constructors = append(constructors, ConstructorTemplate{
			Args:        args,
			Assignments: assignments,
		})
	}

	return constructors
}

func (jenny RawTypes) genDefaultForType(t ast.Type, value any) string {
	switch t.Kind {
	case ast.KindScalar:
		return formatType(t.AsScalar().ScalarKind, value)
	case ast.KindRef:
		return jenny.formatReferenceDefaults(t, value)
	case ast.KindArray:
		if value == nil {
			return "List.of()"
		}
		return fmt.Sprintf("List.of(%s)", jenny.genDefaultForType(t.AsArray().ValueType, value))
	}

	return ""
}

func (jenny RawTypes) formatReferenceDefaults(ref ast.Type, value any) string {
	// Enums
	if _, ok := value.(map[string]interface{}); !ok {
		jenny.typeFormatter.packageMapper(ref.AsRef().ReferredPkg, ref.AsRef().ReferredType)
		return jenny.typeFormatter.formatRefType(ref, value)
	}

	obj, ok := jenny.typeFormatter.context.LocateObjectByRef(ref.AsRef())
	if !ok {
		return ""
	}

	defaultValues := value.(map[string]interface{})
	objectFields := obj.Type.AsStruct().Fields

	args := make([]string, len(objectFields))
	for i, f := range objectFields {
		if v, ok := defaultValues[f.Name]; ok {
			args[i] = jenny.genDefaultForType(f.Type, v)
		} else {
			args[i] = jenny.typeFormatter.emptyValueForType(f.Type, false)
		}
	}

	class := fmt.Sprintf("%s.%s", ref.AsRef().ReferredPkg, ref.AsRef().ReferredType)
	if ref.AsRef().ReferredPkg == obj.SelfRef.ReferredPkg {
		class = ref.AsRef().ReferredType
	}

	jenny.typeFormatter.packageMapper(ref.AsRef().ReferredPkg, ref.AsRef().ReferredType)
	return fmt.Sprintf("new %s(%s)", class, strings.Join(args, ", "))
}
