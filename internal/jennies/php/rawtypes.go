package php

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/jennies/template"
	"github.com/grafana/cog/internal/languages"
	"github.com/grafana/cog/internal/orderedmap"
	"github.com/grafana/cog/internal/tools"
)

type RawTypes struct {
	config Config
	tmpl   *template.Template

	typeFormatter *typeFormatter
	shaper        *shape
}

func (jenny RawTypes) JennyName() string {
	return "PHPRawTypes"
}

func (jenny RawTypes) Generate(context languages.Context) (codejen.Files, error) {
	var err error
	files := make(codejen.Files, 0, len(context.Schemas))

	jenny.shaper = &shape{context: context}

	// generate typehints with a compiler pass
	context.Schemas, err = (&AddTypehintsComments{config: jenny.config}).Process(context.Schemas)
	if err != nil {
		return nil, err
	}

	for _, schema := range context.Schemas {
		schemaFiles, err := jenny.generateSchema(context, schema)
		if err != nil {
			return nil, err
		}

		files = append(files, schemaFiles...)
	}

	return files, nil
}

func (jenny RawTypes) generateSchema(context languages.Context, schema *ast.Schema) (codejen.Files, error) {
	var err error

	files := make(codejen.Files, 0, schema.Objects.Len())
	schema.Objects.Iterate(func(_ string, object ast.Object) {
		// Constants are handled separately
		if object.Type.IsConcreteScalar() {
			return
		}

		file, innerErr := jenny.formatObject(context, schema, object)
		if innerErr != nil {
			err = innerErr
			return
		}

		files = append(files, file)
	})
	if err != nil {
		return nil, err
	}

	constants := schema.Objects.Filter(func(_ string, object ast.Object) bool {
		return object.Type.IsConcreteScalar()
	})
	if constants.Len() != 0 {
		files = append(files, jenny.generateConstants(schema, constants))
	}

	if schema.Metadata.Kind == ast.SchemaKindComposable && schema.Metadata.Variant == ast.SchemaVariantPanel {
		files = append(files, jenny.generatePanelCfgVariantConfigFunc(schema))
	}

	if file := jenny.generateDataqueryVariantConfig(context, schema); file != nil {
		files = append(files, *file)
	}

	return files, nil
}

func (jenny RawTypes) generateConstants(schema *ast.Schema, objects *orderedmap.Map[string, ast.Object]) codejen.File {
	constants := make([]string, 0, objects.Len())

	objects.Iterate(func(_ string, object ast.Object) {
		name := formatConstantName(object.Name)
		value := formatValue(object.Type.Scalar.Value)

		constant := fmt.Sprintf("const %s = %s;", name, value)
		if len(object.Comments) != 0 {
			constant = formatCommentsBlock(object.Comments) + constant
		}

		constants = append(constants, tools.Indent(constant, 4))
	})

	content := fmt.Sprintf(`<?php

namespace %[1]s;

final class Constants
{
%[2]s
}`, jenny.config.fullNamespace(formatPackageName(schema.Package)), strings.Join(constants, "\n"))

	filename := filepath.Join(
		"src",
		formatPackageName(schema.Package),
		"Constants.php",
	)

	return *codejen.NewFile(filename, []byte(content), jenny)
}

func (jenny RawTypes) generatePanelCfgVariantConfigFunc(schema *ast.Schema) codejen.File {
	identifier := schema.Metadata.Identifier

	options := "null"
	if _, hasOptions := schema.LocateObject("Options"); hasOptions {
		options = `[` + jenny.config.fullNamespaceRef(fmt.Sprintf("%s\\Options", formatPackageName(schema.Package))) + `::class, 'fromArray']`
	}

	fieldConfig := "null"
	if _, hasFieldConfig := schema.LocateObject("FieldConfig"); hasFieldConfig {
		fieldConfig = `[` + jenny.config.fullNamespaceRef(fmt.Sprintf("%s\\FieldConfig", formatPackageName(schema.Package))) + `::class, 'fromArray']`
	}

	panelcfgConfigRef := jenny.config.fullNamespaceRef("Cog\\PanelcfgConfig")
	convert := ""
	if jenny.config.converters {
		convert = fmt.Sprintf("\n            convert: [\\%[1]s\\PanelConverter::class, 'convert'],\n", jenny.config.fullNamespace(formatPackageName(schema.Package)))
	}

	content := fmt.Sprintf(`<?php

namespace %[5]s;

final class VariantConfig
{
    public static function get(): %[1]s
    {
        return new %[1]s(
            identifier: '%[2]s',
            optionsFromArray: %[3]s,
            fieldConfigFromArray: %[4]s,%[6]s
        );
    }
}`, panelcfgConfigRef, identifier, options, fieldConfig, jenny.config.fullNamespace(formatPackageName(schema.Package)), convert)

	filename := filepath.Join(
		"src",
		formatPackageName(schema.Package),
		"VariantConfig.php",
	)

	return *codejen.NewFile(filename, []byte(content), jenny)
}

func (jenny RawTypes) formatObject(context languages.Context, schema *ast.Schema, def ast.Object) (codejen.File, error) {
	var buffer strings.Builder

	jenny.typeFormatter = defaultTypeFormatter(jenny.config, context)

	defName := formatObjectName(def.Name)

	comments := def.Comments
	if jenny.config.debug {
		passesTrail := tools.Map(def.PassesTrail, func(trail string) string {
			return fmt.Sprintf("Modified by compiler pass '%s'", trail)
		})
		comments = append(comments, passesTrail...)
	}

	buffer.WriteString(formatCommentsBlock(comments))

	switch def.Type.Kind {
	case ast.KindEnum:
		enum, err := jenny.formatEnumDef(def)
		if err != nil {
			return codejen.File{}, err
		}

		buffer.WriteString(enum)
	case ast.KindRef:
		buffer.WriteString(fmt.Sprintf("class %s extends %s {}", defName, jenny.typeFormatter.formatType(def.Type)))
	case ast.KindStruct:
		buffer.WriteString(jenny.formatStructDef(context, schema, def))
	default:
		return codejen.File{}, fmt.Errorf("unhandled type def kind: %s", def.Type.Kind)
	}

	buffer.WriteString("\n")

	filename := filepath.Join(
		"src",
		formatPackageName(schema.Package),
		fmt.Sprintf("%s.php", defName),
	)

	output := fmt.Sprintf("<?php\n\nnamespace %s;\n\n", jenny.config.fullNamespace(formatPackageName(schema.Package)))
	output += buffer.String()

	return *codejen.NewFile(filename, []byte(output), jenny), nil
}

func (jenny RawTypes) formatStructDef(context languages.Context, schema *ast.Schema, def ast.Object) string {
	var buffer strings.Builder

	variant := ""
	if def.Type.ImplementsVariant() {
		variant = ", " + jenny.config.fullNamespaceRef("Cog\\"+formatObjectName(def.Type.ImplementedVariant()))
	}

	buffer.WriteString(fmt.Sprintf("class %s implements \\JsonSerializable%s\n{\n", formatObjectName(def.Name), variant))

	for _, fieldDef := range def.Type.Struct.Fields {
		buffer.WriteString(tools.Indent(jenny.typeFormatter.formatField(fieldDef), 4))
		buffer.WriteString("\n\n")
	}

	buffer.WriteString(tools.Indent(jenny.generateConstructor(context, def), 4))
	buffer.WriteString("\n\n")

	buffer.WriteString(tools.Indent(jenny.generateFromJSON(context, def), 4))
	buffer.WriteString("\n\n")

	buffer.WriteString(tools.Indent(jenny.generateJSONSerialize(def), 4))

	if def.Type.IsDataqueryVariant() {
		buffer.WriteString("\n\n")
		buffer.WriteString(tools.Indent(jenny.generateDataqueryType(schema), 4))
	}

	buffer.WriteString("\n}")

	return buffer.String()
}

func (jenny RawTypes) generateDataqueryVariantConfig(context languages.Context, schema *ast.Schema) *codejen.File {
	if schema.Metadata.Variant != ast.SchemaVariantDataQuery || schema.EntryPoint == "" {
		return nil
	}

	dataqueryConfigRef := jenny.config.fullNamespaceRef("Cog\\DataqueryConfig")
	var fromArrayCallable string

	_, entryPointFound := schema.LocateObject(schema.EntryPoint)

	switch {
	case !entryPointFound && schema.EntryPointType.Kind == "": // no entrypoint at all
		return nil
	case !entryPointFound && schema.EntryPointType.IsDisjunction(): // the entrypoint is a disjunction that was inlined (its object no longer exists)
		fromArrayCallable = jenny.unmarshalDisjunctionFunc(context, schema.EntryPointType.AsDisjunction())
	case entryPointFound: // the entrypoint is a valid reference to an object
		fromArrayCallable = `[` + jenny.config.fullNamespaceRef(fmt.Sprintf("%s\\%s", formatPackageName(schema.Package), formatObjectName(schema.EntryPoint))) + `::class, 'fromArray']`
	default: // No valid entrypoint found
		return nil
	}

	converterCallable := ""
	if jenny.config.converters {
		converterCallable = `[` + jenny.config.fullNamespaceRef(fmt.Sprintf("%s\\%s", formatPackageName(schema.Package), formatObjectName(schema.EntryPoint))) + `Converter::class, 'convert']`
		if !entryPointFound && schema.EntryPointType.IsDisjunction() {
			converterCallable = jenny.convertDisjunctionFunc(schema.EntryPointType.AsDisjunction())
		}

		converterCallable = fmt.Sprintf("\n            convert: %s,", converterCallable)
	}

	content := fmt.Sprintf(`<?php

namespace %[4]s;

final class VariantConfig
{
    public static function get(): %[1]s
    {
        return new %[1]s(
            identifier: "%[2]s",
            fromArray: %[3]s,%[5]s
        );
    }
}`, dataqueryConfigRef, schema.Metadata.Identifier, fromArrayCallable, jenny.config.fullNamespace(formatPackageName(schema.Package)), converterCallable)

	filename := filepath.Join(
		"src",
		formatPackageName(schema.Package),
		"VariantConfig.php",
	)

	return codejen.NewFile(filename, []byte(content), jenny)
}

func (jenny RawTypes) convertDisjunctionFunc(disjunction ast.DisjunctionType) string {
	decodingSwitch := "switch (true) {\n"
	discriminators := tools.Keys(disjunction.DiscriminatorMapping)
	sort.Strings(discriminators) // to ensure a deterministic output
	for _, discriminator := range discriminators {
		if discriminator == ast.DiscriminatorCatchAll {
			continue
		}

		objectRef := disjunction.DiscriminatorMapping[discriminator]
		decodingSwitch += fmt.Sprintf(`    case $input instanceof %[1]s:
        return %[1]sConverter::convert($input);
`, objectRef)
	}

	if defaultBranchType, ok := disjunction.DiscriminatorMapping[ast.DiscriminatorCatchAll]; ok {
		decodingSwitch += fmt.Sprintf(`    default:
        return %[1]sConverter::convert($input);
`, defaultBranchType)
	} else {
		decodingSwitch += `    default:
        throw new \ValueError('can not convert unknown disjunction branch');
`
	}

	decodingSwitch += "}"

	dataqueryRef := jenny.config.fullNamespaceRef("Cog\\Dataquery")

	return fmt.Sprintf(`(function(%s $input) {

    %s
})`, dataqueryRef, decodingSwitch)
}

func (jenny RawTypes) generateConstructor(context languages.Context, def ast.Object) string {
	var buffer strings.Builder
	hinter := typehints{config: jenny.config, context: context}

	var typeAnnotations []string
	var args []string
	var assignments []string

	for _, field := range def.Type.AsStruct().Fields {
		fieldName := formatFieldName(field.Name)
		defaultValue := (any)(nil)

		// set for default values for fields that need one or have one
		if !field.Type.Nullable || field.Type.Default != nil {
			var defaultsOverrides map[string]any
			if overrides, ok := field.Type.Default.(map[string]interface{}); ok {
				defaultsOverrides = overrides
			}

			defaultValue = defaultValueForType(jenny.config, context.Schemas, field.Type, orderedmap.FromMap(defaultsOverrides))
		}

		// initialize constant fields
		if field.Type.IsConcreteScalar() {
			assignments = append(assignments, fmt.Sprintf("    $this->%s = %s;\n", fieldName, formatValue(field.Type.AsScalar().Value)))
			continue
		}

		argType := field.Type.DeepCopy()
		argType.Nullable = true

		args = append(args, fmt.Sprintf("%s $%s = null", jenny.typeFormatter.formatType(argType), fieldName))
		typeAnnotation := hinter.paramAnnotationForType(fieldName, argType)
		if typeAnnotation != "" {
			typeAnnotations = append(typeAnnotations, typeAnnotation)
		}

		if field.Type.Nullable {
			assignments = append(assignments, fmt.Sprintf("    $this->%[1]s = $%[1]s;", fieldName))
		} else {
			assignments = append(assignments, fmt.Sprintf("    $this->%[1]s = $%[1]s ?: %[2]s;", fieldName, formatValue(defaultValue)))
		}
	}

	if len(typeAnnotations) != 0 {
		buffer.WriteString(formatCommentsBlock(typeAnnotations))
	}

	buffer.WriteString(fmt.Sprintf("public function __construct(%s)\n", strings.Join(args, ", ")))
	buffer.WriteString("{\n")

	buffer.WriteString(strings.Join(assignments, "\n"))

	buffer.WriteString("\n}")

	return buffer.String()
}

func (jenny RawTypes) generateFromJSON(context languages.Context, def ast.Object) string {
	var buffer strings.Builder

	var constructorArgs []string

	for _, field := range def.Type.AsStruct().Fields {
		// No need to unmarshal constant scalar fields since they're set in
		// the object's constructor
		if field.Type.IsConcreteScalar() {
			continue
		}

		var value string

		fieldName := formatFieldName(field.Name)
		inputVar := fmt.Sprintf(`$data["%[1]s"]`, field.Name)

		switch {
		// Special case to properly parse dashboard.Panel options
		case def.SelfRef.ReferredPkg == "dashboard" && strings.EqualFold(def.Name, "panel") && field.Name == "options":
			decodingFunc := jenny.unmarshalDashboardOptionsFunc()

			value = fmt.Sprintf(`isset(%[1]s) ? %[2]s($data) : null`, inputVar, decodingFunc)

			// Special case to properly parse dashboard.Panel fieldConfig
		case def.SelfRef.ReferredPkg == "dashboard" && strings.EqualFold(def.Name, "panel") && field.Name == "fieldConfig":
			decodingFunc := jenny.unmarshalDashboardFieldConfigFunc(context, field)

			value = fmt.Sprintf(`isset(%[1]s) ? %[2]s($data) : null`, inputVar, decodingFunc)
		default:
			value = jenny.unmarshalForType(context, def, field.Type, inputVar)
		}

		constructorArgs = append(constructorArgs, fmt.Sprintf("        %s: %s,\n", fieldName, value))
	}

	buffer.WriteString("/**\n")
	buffer.WriteString(" * @param array<string, mixed> $inputData\n")
	buffer.WriteString(" */\n")
	buffer.WriteString("public static function fromArray(array $inputData): self\n")
	buffer.WriteString("{\n")
	if len(constructorArgs) != 0 {
		buffer.WriteString(fmt.Sprintf("    /** @var %s $inputData */\n", jenny.shaper.typeShape(def.Type)))
		buffer.WriteString("    $data = $inputData;\n")
	}
	buffer.WriteString("    return new self(\n")
	buffer.WriteString(strings.Join(constructorArgs, ""))
	buffer.WriteString("    );\n")
	buffer.WriteString("}")

	return buffer.String()
}

func (jenny RawTypes) unmarshalForType(context languages.Context, object ast.Object, def ast.Type, inputVar string) string {
	if _, ok := context.ResolveToComposableSlot(def); ok {
		return jenny.unmarshalComposableSlot(context, object, def, inputVar)
	}

	switch {
	case def.IsRef():
		return fmt.Sprintf(`isset(%[2]s) ? %[1]s(%[2]s) : null`, jenny.unmarshalRefFunc(context, def), inputVar)
	case def.IsArray() && def.Array.ValueType.IsRef():
		return fmt.Sprintf(`array_filter(array_map(%s, %s ?? []))`, jenny.unmarshalRefFunc(context, def.Array.ValueType), inputVar)
	case def.IsArray() && def.Array.ValueType.IsDisjunction():
		disjunctionType := def.Array.ValueType.AsDisjunction()
		decodingFunc := jenny.unmarshalDisjunctionFunc(context, disjunctionType)

		return fmt.Sprintf(`!empty(%[1]s) ? array_map(%[2]s, %[1]s) : null`, inputVar, decodingFunc)
	case def.IsDisjunction():
		decodingFunc := jenny.unmarshalDisjunctionFunc(context, def.AsDisjunction())

		return fmt.Sprintf(`isset(%[1]s) ? %[2]s(%[1]s) : null`, inputVar, decodingFunc)
	case def.IsMap():
		return jenny.unmarshalMap(context, def.AsMap(), inputVar)
	default:
		return fmt.Sprintf(`%[1]s ?? null`, inputVar)
	}
}

func (jenny RawTypes) unmarshalDashboardOptionsFunc() string {
	cog := jenny.config.fullNamespaceRef("Cog")

	return fmt.Sprintf(`(function($panel) {
    /** @var array<string, mixed> $options */
    $options = $panel["options"];

    if (!%[1]s\Runtime::get()->panelcfgVariantExists($panel["type"] ?? "")) {
        return $options;
    }

    $config = %[1]s\Runtime::get()->panelcfgVariantConfig($panel["type"] ?? "");
    if ($config->optionsFromArray === null) {
        return $options;
    }

	return ($config->optionsFromArray)($options);
})`, cog)
}

func (jenny RawTypes) unmarshalDashboardFieldConfigFunc(context languages.Context, field ast.StructField) string {
	cog := jenny.config.fullNamespaceRef("Cog")

	fieldConfigShape := jenny.shaper.typeShape(context.ResolveRefs(field.Type))

	return fmt.Sprintf(`(function($panel) {
    /** @var %[2]s */
    $fieldConfigData = $panel["fieldConfig"];
    $fieldConfig = FieldConfigSource::fromArray($fieldConfigData);

    if (!%[1]s\Runtime::get()->panelcfgVariantExists($panel["type"] ?? "")) {
        return $fieldConfig;
    }

    $config = %[1]s\Runtime::get()->panelcfgVariantConfig($panel["type"] ?? "");
    if ($config->fieldConfigFromArray === null) {
        return $fieldConfig;
    }

    if (!isset($fieldConfigData["defaults"])) {
		return $fieldConfig;
    }
    /** @var array{custom?: array<string, mixed>}*/
    $defaults = $fieldConfigData["defaults"];
    if (!isset($defaults["custom"])) {
		return $fieldConfig;
    }

	$fieldConfig->defaults->custom = ($config->fieldConfigFromArray)($defaults["custom"]);

    return $fieldConfig;
})`, cog, fieldConfigShape)
}

func (jenny RawTypes) unmarshalMap(context languages.Context, mapDef ast.MapType, inputVar string) string {
	if !mapDef.ValueType.IsRef() {
		return fmt.Sprintf("%s ?? null", inputVar)
	}

	decodeRef := jenny.unmarshalRefFunc(context, mapDef.ValueType)

	return fmt.Sprintf(`isset(%[1]s) ? array_map(%[2]s, %[1]s) : null`, inputVar, decodeRef)
}

func (jenny RawTypes) unmarshalRefFunc(context languages.Context, refDef ast.Type) string {
	referredObject, found := context.LocateObjectByRef(refDef.AsRef())
	formattedRef := jenny.typeFormatter.formatRef(refDef, false)

	if found && referredObject.Type.IsStruct() {
		assignment := fmt.Sprintf("/** @var %s */\n", jenny.shaper.typeShape(referredObject.Type))
		assignment += "$val = $input;"

		return fmt.Sprintf(`(function($input) {
	%[2]s
	return %[1]s::fromArray($val);
})`, formattedRef, assignment)
	} else if found && referredObject.Type.IsEnum() {
		return fmt.Sprintf(`(function($input) { return %[1]s::fromValue($input); })`, formattedRef)
	}

	// TODO: should not happen?
	return `/* ref to a non-struct, non-enum, this should have been inlined */ (function(array $input) { return $input; })`
}

func (jenny RawTypes) unmarshalComposableSlot(context languages.Context, parentObject ast.Object, def ast.Type, inputVar string) string {
	slotType, _ := context.ResolveToComposableSlot(def)

	if slotType.ComposableSlot.Variant == ast.SchemaVariantDataQuery {
		return jenny.renderUnmarshalDataqueryField(parentObject, def, inputVar)
	}

	// TODO
	return inputVar
}

func (jenny RawTypes) renderUnmarshalDataqueryField(parentObject ast.Object, def ast.Type, inputVar string) string {
	dataqueryHint := `""`

	// First: try to locate a field that would contain the type of datasource being used.
	// We're looking for a field defined as a reference to the `DataSourceRef` type.
	for _, candidate := range parentObject.Type.AsStruct().Fields {
		if !candidate.Type.IsRef() {
			continue
		}
		if candidate.Type.AsRef().ReferredType != "DataSourceRef" {
			continue
		}

		dataqueryHint = fmt.Sprintf(`(isset($in["%[1]s"], $in["%[1]s"]["type"]) && is_string($in["%[1]s"]["type"])) ? $in["%[1]s"]["type"] : ""`, candidate.Name)
		break
	}

	runtimeRef := jenny.config.fullNamespaceRef("Cog\\Runtime")

	if def.IsArray() {
		return fmt.Sprintf(`isset(%[1]s) ? (function ($in) {
	/** @var array{datasource?: array{type?: mixed}} $in */
    $hint = %[3]s;
    /** @var array<array<string, mixed>> $in */
    return %[2]s::get()->dataqueriesFromArray($in, $hint);
})(%[1]s): null`, inputVar, runtimeRef, dataqueryHint)
	}

	return fmt.Sprintf(`isset(%[1]s) ? (function($in) {
	/** @var array{datasource?: array{type?: mixed}} $in */
    $hint = %[3]s;
    /** @var array<string, mixed> $in */
    return %[2]s::get()->dataqueryFromArray($in, $hint);
})(%[1]s): null`, inputVar, runtimeRef, dataqueryHint)
}

func (jenny RawTypes) unmarshalDisjunctionFunc(context languages.Context, disjunction ast.DisjunctionType) string {
	// this potentially generates incorrect code, but there isn't much we can do without more information.
	if disjunction.Discriminator == "" || disjunction.DiscriminatorMapping == nil {
		decodingSwitch := "switch (true) {\n"

		var ignoredBranches []ast.Type
		for _, branch := range disjunction.Branches {
			if branch.IsScalar() {
				testMap := map[ast.ScalarKind]string{
					ast.KindBytes:   "is_string",
					ast.KindString:  "is_string",
					ast.KindFloat32: "is_float",
					ast.KindFloat64: "is_float",
					ast.KindUint8:   "is_int",
					ast.KindUint16:  "is_int",
					ast.KindUint32:  "is_int",
					ast.KindUint64:  "is_int",
					ast.KindInt8:    "is_int",
					ast.KindInt16:   "is_int",
					ast.KindInt32:   "is_int",
					ast.KindInt64:   "is_int",
					ast.KindBool:    "is_bool",
				}

				testFunc := testMap[branch.Scalar.ScalarKind]
				if testFunc == "" {
					ignoredBranches = append(ignoredBranches, branch)
					continue
				}

				decodingSwitch += fmt.Sprintf(`    case %[1]s($input):
        return $input;
`, testFunc)
				continue
			}

			ignoredBranches = append(ignoredBranches, branch)
		}

		//nolint:gocritic
		if len(ignoredBranches) == 1 && ignoredBranches[0].IsRef() {
			ref := ignoredBranches[0].AsRef()
			referredObject, found := context.LocateObjectByRef(ref)
			formattedRef := jenny.typeFormatter.formatRef(ignoredBranches[0], false)

			value := "$input"
			if found && referredObject.Type.IsStruct() {
				value = fmt.Sprintf(`%[1]s::fromArray($input)`, formattedRef)
			} else if found && referredObject.Type.IsEnum() {
				value = fmt.Sprintf(`%[1]s::fromValue($input)`, formattedRef)
			}

			decodingSwitch += fmt.Sprintf(`    default:
        /** @var %[2]s $input */
        return %[1]s;
`, value, jenny.shaper.typeShape(referredObject.Type))
		} else if len(ignoredBranches) >= 1 {
			decodingSwitch += `    default:
        return $input;
`
		} else if len(ignoredBranches) == 0 {
			decodingSwitch += `    default:
        throw new \ValueError('incorrect value for disjunction');
`
		}

		decodingSwitch += "}"

		return fmt.Sprintf(`(function($input) {
    %s
})`, decodingSwitch)
	}

	decodingSwitch := fmt.Sprintf("switch ($input[\"%s\"]) {\n", disjunction.Discriminator)
	discriminators := tools.Keys(disjunction.DiscriminatorMapping)
	sort.Strings(discriminators) // to ensure a deterministic output
	for _, discriminator := range discriminators {
		if discriminator == ast.DiscriminatorCatchAll {
			continue
		}

		objectRef := disjunction.DiscriminatorMapping[discriminator]
		decodingSwitch += fmt.Sprintf(`    case "%[1]s":
        return %[2]s::fromArray($input);
`, discriminator, objectRef)
	}

	if defaultBranchType, ok := disjunction.DiscriminatorMapping[ast.DiscriminatorCatchAll]; ok {
		decodingSwitch += fmt.Sprintf(`    default:
        return %[1]s::fromArray($input);
`, defaultBranchType)
	} else {
		decodingSwitch += `    default:
        throw new \ValueError('can not parse disjunction from array');
`
	}

	decodingSwitch += "}"

	return fmt.Sprintf(`(function($input) {
    \assert(is_array($input), 'expected disjunction value to be an array');

    %s
})`, decodingSwitch)
}

func (jenny RawTypes) generateJSONSerialize(def ast.Object) string {
	var buffer strings.Builder

	buffer.WriteString("/**\n")
	buffer.WriteString(" * @return array<string, mixed>\n")
	buffer.WriteString(" */\n")
	buffer.WriteString("public function jsonSerialize(): array\n")
	buffer.WriteString("{\n")

	buffer.WriteString("    $data = [\n")

	for _, field := range def.Type.AsStruct().Fields {
		if field.Type.Nullable {
			continue
		}

		buffer.WriteString(fmt.Sprintf(`        "%s" => $this->%s,`+"\n", field.Name, formatFieldName(field.Name)))
	}

	buffer.WriteString("    ];\n")

	for _, field := range def.Type.AsStruct().Fields {
		if !field.Type.Nullable {
			continue
		}

		fieldName := formatFieldName(field.Name)

		buffer.WriteString(fmt.Sprintf("    if (isset($this->%s)) {\n", fieldName))
		buffer.WriteString(fmt.Sprintf(`        $data["%s"] = $this->%s;`+"\n", field.Name, fieldName))
		buffer.WriteString("    }\n")
	}

	buffer.WriteString("    return $data;\n")

	buffer.WriteString("}")

	return buffer.String()
}

func (jenny RawTypes) generateDataqueryType(schema *ast.Schema) string {
	var buffer strings.Builder

	buffer.WriteString("public function dataqueryType(): string\n")
	buffer.WriteString("{\n")
	buffer.WriteString(fmt.Sprintf("    return \"%s\";\n", strings.ToLower(schema.Metadata.Identifier)))
	buffer.WriteString("}")

	return buffer.String()
}

func (jenny RawTypes) formatEnumDef(def ast.Object) (string, error) {
	return jenny.tmpl.
		Funcs(map[string]any{
			"formatType": jenny.typeFormatter.formatType,
		}).
		Render("types/enum.tmpl", map[string]any{
			"Object":   def,
			"EnumType": def.Type.Enum.Values[0].Type,
		})
}
