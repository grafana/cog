package php

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/languages"
	"github.com/grafana/cog/internal/orderedmap"
	"github.com/grafana/cog/internal/tools"
)

type RawTypes struct {
	config Config

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

	if schema.Metadata.Kind == ast.SchemaKindComposable && schema.Metadata.Variant == ast.SchemaVariantPanel {
		files = append(files, jenny.generatePanelCfgVariantConfigFunc(schema))
	}

	return files, nil
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

	content := fmt.Sprintf(`<?php

namespace %[5]s;

final class VariantConfig
{
    public static function get(): %[1]s
    {
        return new %[1]s(
            identifier: '%[2]s',
            optionsFromArray: %[3]s,
            fieldConfigFromArray: %[4]s
        );
    }
}`, panelcfgConfigRef, identifier, options, fieldConfig, jenny.config.fullNamespace(formatPackageName(schema.Package)))

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
	case ast.KindScalar:
		scalarType := def.Type.AsScalar()

		//nolint: gocritic
		if scalarType.IsConcrete() {
			buffer.WriteString(fmt.Sprintf("const %s = %s;", formatConstantName(def.Name), formatValue(scalarType.Value)))
		} else {
			return codejen.File{}, fmt.Errorf("type aliases on scalar types is not supported")
		}
	case ast.KindRef:
		buffer.WriteString(fmt.Sprintf("class %s extends %s {}", defName, jenny.typeFormatter.formatType(def.Type)))
	case ast.KindStruct:
		buffer.WriteString(jenny.formatStructDef(context, def))
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

func (jenny RawTypes) formatStructDef(context languages.Context, def ast.Object) string {
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

	buffer.WriteString("\n}")

	return buffer.String()
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

		// Special cases to properly parse dashboard.Panel options
		if def.SelfRef.ReferredPkg == "dashboard" && strings.EqualFold(def.Name, "panel") && field.Name == "options" {
			decodingFunc := jenny.unmarshalDashboardOptionsFunc()

			value = fmt.Sprintf(`isset(%[1]s) ? %[2]s($data) : null`, inputVar, decodingFunc)

			// Special cases to properly parse dashboard.Panel fieldConfig
		} else if def.SelfRef.ReferredPkg == "dashboard" && strings.EqualFold(def.Name, "panel") && field.Name == "fieldConfig" {
			decodingFunc := jenny.unmarshalDashboardFieldConfigFunc(context, field)

			value = fmt.Sprintf(`isset(%[1]s) ? %[2]s($data) : null`, inputVar, decodingFunc)
		} else {
			value = jenny.unmarshalForType(context, def, field.Type, inputVar)
		}

		constructorArgs = append(constructorArgs, fmt.Sprintf("        %s: %s,\n", fieldName, value))
	}

	buffer.WriteString("/**\n")
	buffer.WriteString(" * @param array<string, mixed> $inputData\n")
	buffer.WriteString(" */\n")
	buffer.WriteString("public static function fromArray(array $inputData): self\n")
	buffer.WriteString("{\n")
	buffer.WriteString(fmt.Sprintf("    /** @var %s $inputData */\n", jenny.shaper.typeShape(def.Type)))
	buffer.WriteString("    $data = $inputData;\n")
	buffer.WriteString("    return new self(\n")
	buffer.WriteString(strings.Join(constructorArgs, ""))
	buffer.WriteString("    );\n")
	buffer.WriteString("}")

	return buffer.String()
}

func (jenny RawTypes) unmarshalForType(context languages.Context, object ast.Object, def ast.Type, inputVar string) string {
	if _, ok := context.ResolveToComposableSlot(def); ok {
		return jenny.unmarshalComposableSlot(context, object, def, inputVar)
	} else if def.IsRef() {
		return fmt.Sprintf(`isset(%[2]s) ? %[1]s(%[2]s) : null`, jenny.unmarshalRefFunc(context, def), inputVar)
	} else if def.IsArray() && def.Array.ValueType.IsRef() {
		return fmt.Sprintf(`array_filter(array_map(%s, %s ?? []))`, jenny.unmarshalRefFunc(context, def.Array.ValueType), inputVar)
	} else if def.IsArray() && def.Array.ValueType.IsDisjunction() {
		disjunctionType := def.Array.ValueType.AsDisjunction()
		decodingFunc := jenny.unmarshalDisjunctionFunc(context, disjunctionType)

		return fmt.Sprintf(`!empty(%[1]s) ? array_map(%[2]s, %[1]s) : null`, inputVar, decodingFunc)
	} else if def.IsDisjunction() {
		decodingFunc := jenny.unmarshalDisjunctionFunc(context, def.AsDisjunction())

		return fmt.Sprintf(`isset(%[1]s) ? %[2]s(%[1]s) : null`, inputVar, decodingFunc)
	} else if def.IsMap() {
		return jenny.unmarshalMap(context, def.AsMap(), inputVar)
	}

	return fmt.Sprintf(`%[1]s ?? null`, inputVar)
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

		if len(ignoredBranches) == 1 && ignoredBranches[0].IsRef() {
			ref := ignoredBranches[0].AsRef()
			referredObject, found := context.LocateObjectByRef(ref)
			formattedRef := jenny.typeFormatter.formatRef(ignoredBranches[0], false)

			var value string
			if found && referredObject.Type.IsStruct() {
				value = fmt.Sprintf(`%[1]s::fromArray($input)`, formattedRef)
			} else if found && referredObject.Type.IsEnum() {
				value = fmt.Sprintf(`%[1]s::fromValue($input)`, formattedRef)
			} else {
				// TODO: should not happen?
				// ref to a non-struct, non-enum, this should have been inlined
				value = "$input"
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
	for discriminator, objectRef := range disjunction.DiscriminatorMapping {
		if discriminator == ast.DiscriminatorCatchAll {
			continue
		}

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
    if (!is_array($input)) {
        throw new \ValueError('expected disjunction value to be an array');
    }

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

func (jenny RawTypes) formatEnumDef(def ast.Object) (string, error) {
	enumType := def.Type.Enum.Values[0].Type

	buf := bytes.Buffer{}
	if err := templates.
		Funcs(map[string]any{
			"formatType": jenny.typeFormatter.formatType,
		}).
		ExecuteTemplate(&buf, "types/enum.tmpl", map[string]any{
			"Object":   def,
			"EnumType": enumType,
		}); err != nil {
		return "", fmt.Errorf("failed executing template: %w", err)
	}

	return buf.String(), nil
}
