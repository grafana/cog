package common

import (
	"bytes"
	"fmt"
	"slices"
	"strings"

	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/jennies/template"
	"github.com/grafana/cog/internal/languages"
	"github.com/grafana/cog/internal/orderedmap"
)

type ArgumentReference struct {
	Name     string
	Type     string
	Comments []string
}

type MethodReference struct {
	Name      string
	Comments  []string
	Arguments []ArgumentReference
	Return    string
	Static    bool
}

type FunctionReference struct {
	Name      string
	Comments  []string
	Arguments []ArgumentReference
	Return    string
}

type APIReferenceCollector struct {
	objectMethods    map[string][]MethodReference
	packageFunctions map[string][]FunctionReference
}

func NewAPIReferenceCollector() *APIReferenceCollector {
	return &APIReferenceCollector{
		objectMethods:    make(map[string][]MethodReference),
		packageFunctions: make(map[string][]FunctionReference),
	}
}

func (collector *APIReferenceCollector) RegisterMethod(object ast.Object, methodReference MethodReference) {
	objectRef := object.SelfRef.String()
	collector.objectMethods[objectRef] = append(collector.objectMethods[objectRef], methodReference)
}

func (collector *APIReferenceCollector) MethodsForObject(object ast.Object) []MethodReference {
	return collector.objectMethods[object.SelfRef.String()]
}

func (collector *APIReferenceCollector) RegisterFunction(pkg string, functionReference FunctionReference) {
	collector.packageFunctions[pkg] = append(collector.packageFunctions[pkg], functionReference)
}

func (collector *APIReferenceCollector) FunctionsForPackage(pkg string) []FunctionReference {
	return collector.packageFunctions[pkg]
}

type APIReferenceFormatter struct {
	KindName func(kind ast.Kind) string

	FunctionSignature func(context languages.Context, function FunctionReference) string

	ObjectName       func(object ast.Object) string
	ObjectDefinition func(context languages.Context, object ast.Object) string

	MethodName      func(method MethodReference) string
	MethodSignature func(context languages.Context, object ast.Object, method MethodReference) string

	BuilderName          func(builder ast.Builder) string
	ConstructorSignature func(context languages.Context, builder ast.Builder) string
	OptionName           func(option ast.Option) string
	OptionSignature      func(context languages.Context, builder ast.Builder, option ast.Option) string
}

type APIReference struct {
	Collector *APIReferenceCollector
	Language  string
	Formatter APIReferenceFormatter
	Tmpl      *template.Template
}

func (jenny APIReference) JennyName() string {
	return "APIReference"
}

func (jenny APIReference) Generate(context languages.Context) (codejen.Files, error) {
	files := make([]codejen.File, 0, len(context.Schemas)+len(context.Builders)+1)

	for _, schema := range context.Schemas {
		schemaFiles, err := jenny.referenceForSchema(context, schema)
		if err != nil {
			return nil, err
		}

		files = append(files, schemaFiles...)
	}
	for _, builder := range context.Builders {
		builderFile, err := jenny.referenceForBuilder(context, builder)
		if err != nil {
			return nil, err
		}

		files = append(files, builderFile)
	}

	indexFile, err := jenny.index(context)
	if err != nil {
		return nil, err
	}
	files = append(files, indexFile)

	return files, nil
}

func (jenny APIReference) index(context languages.Context) (codejen.File, error) {
	var buffer bytes.Buffer

	err := jenny.renderIfExists(&buffer, template.ExtraIndexDocsBlock(), map[string]any{})
	if err != nil {
		return codejen.File{}, err
	}

	buffer.WriteString("# Packages\n\n")

	slices.SortFunc(context.Schemas, func(schemaA, schemaB *ast.Schema) int {
		return strings.Compare(schemaA.Package, schemaB.Package)
	})

	for _, schema := range context.Schemas {
		buffer.WriteString(fmt.Sprintf(" * [%[1]s](./%[1]s/index.md)\n", schema.Package))
	}

	return *codejen.NewFile("docs/index.md", buffer.Bytes(), jenny), nil
}

func (jenny APIReference) referenceForSchema(context languages.Context, schema *ast.Schema) (codejen.Files, error) {
	files := make([]codejen.File, 0, schema.Objects.Len()+1)

	schemaIndexFile, err := jenny.schemaIndex(context, schema)
	if err != nil {
		return nil, err
	}
	files = append(files, schemaIndexFile)

	var inner error
	schema.Objects.Iterate(func(_ string, object ast.Object) {
		if inner != nil {
			return
		}

		objFile, err := jenny.referenceForObject(context, object)
		if err != nil {
			inner = err
		}
		files = append(files, objFile)
	})
	if inner != nil {
		return nil, inner
	}

	return files, nil
}

func (jenny APIReference) schemaIndex(context languages.Context, schema *ast.Schema) (codejen.File, error) {
	var buffer bytes.Buffer

	buffer.WriteString(fmt.Sprintf("# %s\n\n", schema.Package))

	buffer.WriteString("## Objects\n\n")

	schema.Objects.Sort(orderedmap.SortStrings)
	schema.Objects.Iterate(func(_ string, object ast.Object) {
		buffer.WriteString(fmt.Sprintf(" * %[2]s [%[1]s](./%[1]s.md)\n", jenny.Formatter.ObjectName(object), jenny.kindBadge(object.Type.Kind)))
	})

	buffer.WriteString("## Builders\n\n")

	builders := context.Builders.ByPackage(schema.Package)
	slices.SortFunc(builders, func(builderA, builderB ast.Builder) int {
		return strings.Compare(builderA.Name, builderB.Name)
	})

	for _, builder := range builders {
		buffer.WriteString(fmt.Sprintf(" * %[2]s [%[1]s](./%[1]s.md)\n", jenny.Formatter.BuilderName(builder), jenny.builderBadge()))
	}

	functions := jenny.Collector.FunctionsForPackage(schema.Package)

	if len(functions) > 0 {
		buffer.WriteString("## Functions\n\n")

		for _, functionReference := range functions {
			buffer.WriteString(fmt.Sprintf("### %[2]s %[1]s\n\n", functionReference.Name, jenny.functionBadge()))

			if len(functionReference.Comments) != 0 {
				buffer.WriteString(strings.Join(functionReference.Comments, "\n\n") + "\n\n")
			}

			buffer.WriteString(fmt.Sprintf("```%s\n", jenny.Language))
			buffer.WriteString(jenny.Formatter.FunctionSignature(context, functionReference))
			buffer.WriteString("\n```\n")

			buffer.WriteString("\n")
		}
	}

	err := jenny.renderIfExists(&buffer, template.ExtraPackageDocsBlock(schema), map[string]any{
		"Schema": schema,
	})
	if err != nil {
		return codejen.File{}, err
	}

	return *codejen.NewFile(fmt.Sprintf("docs/%s/index.md", schema.Package), buffer.Bytes(), jenny), nil
}

func (jenny APIReference) referenceForObject(context languages.Context, object ast.Object) (codejen.File, error) {
	var buffer bytes.Buffer

	objectName := jenny.Formatter.ObjectName(object)

	buffer.WriteString(fmt.Sprintf(`---
title: %[2]s %[1]s
---
`, objectName, jenny.kindBadge(object.Type.Kind)))

	buffer.WriteString(fmt.Sprintf("# %[2]s %[1]s\n\n", objectName, jenny.kindBadge(object.Type.Kind)))

	if len(object.Comments) != 0 {
		buffer.WriteString(strings.Join(object.Comments, "\n\n") + "\n\n")
	}

	buffer.WriteString("## Definition\n\n")

	buffer.WriteString(fmt.Sprintf("```%s\n", jenny.Language))
	buffer.WriteString(jenny.Formatter.ObjectDefinition(context, object))
	buffer.WriteString("\n```\n")

	if object.Type.IsStruct() {
		jenny.referenceStructMethods(&buffer, context, object)
	}

	err := jenny.renderIfExists(&buffer, template.ExtraObjectDocsBlock(object), map[string]any{
		"Object": object,
	})
	if err != nil {
		return codejen.File{}, err
	}

	buildersForObjet := context.Builders.LocateAllByObject(object.SelfRef.ReferredPkg, object.SelfRef.ReferredType)
	if len(buildersForObjet) != 0 {
		buffer.WriteString("## See also\n\n")

		slices.SortFunc(buildersForObjet, func(builderA, builderB ast.Builder) int {
			return strings.Compare(builderA.Name, builderB.Name)
		})
		for _, builder := range buildersForObjet {
			if builder.Package == object.SelfRef.ReferredPkg {
				buffer.WriteString(fmt.Sprintf(" * %[2]s [%[1]s](./%[1]s.md)\n", jenny.Formatter.BuilderName(builder), jenny.builderBadge()))
			} else {
				buffer.WriteString(fmt.Sprintf(" * %[3]s [%[1]s.%[2]s](../%[1]s/%[2]s.md)\n", builder.Package, jenny.Formatter.BuilderName(builder), jenny.builderBadge()))
			}
		}
	}

	return *codejen.NewFile(fmt.Sprintf("docs/%s/%s.md", object.SelfRef.ReferredPkg, objectName), buffer.Bytes(), jenny), nil
}

func (jenny APIReference) referenceStructMethods(buffer *bytes.Buffer, context languages.Context, object ast.Object) {
	buffer.WriteString("## Methods\n\n")

	methods := jenny.Collector.MethodsForObject(object)

	for _, method := range methods {
		buffer.WriteString(fmt.Sprintf("### %[2]s %[1]s\n\n", jenny.Formatter.MethodName(method), jenny.methodBadge()))

		if len(method.Comments) != 0 {
			buffer.WriteString(strings.Join(method.Comments, "\n\n") + "\n\n")
		}

		buffer.WriteString(fmt.Sprintf("```%s\n", jenny.Language))
		buffer.WriteString(jenny.Formatter.MethodSignature(context, object, method))
		buffer.WriteString("\n```\n")

		buffer.WriteString("\n")
	}

	if len(methods) == 0 {
		buffer.WriteString("No methods.\n")
	}
}

func (jenny APIReference) referenceForBuilder(context languages.Context, builder ast.Builder) (codejen.File, error) {
	var buffer bytes.Buffer

	builderName := jenny.Formatter.BuilderName(builder)

	buffer.WriteString(fmt.Sprintf(`---
title: %[2]s %[1]s
---
`, builderName, jenny.builderBadge()))

	buffer.WriteString(fmt.Sprintf("# %[2]s %[1]s\n\n", builderName, jenny.builderBadge()))

	buffer.WriteString("## Constructor\n\n")

	buffer.WriteString(fmt.Sprintf("```%s\n", jenny.Language))
	buffer.WriteString(jenny.Formatter.ConstructorSignature(context, builder))
	buffer.WriteString("\n```\n")

	buffer.WriteString("## Methods\n\n")

	slices.SortFunc(builder.Options, func(optionA, optionB ast.Option) int {
		return strings.Compare(optionA.Name, optionB.Name)
	})

	for _, option := range builder.Options {
		buffer.WriteString(fmt.Sprintf("### %[2]s %[1]s\n\n", jenny.Formatter.OptionName(option), jenny.methodBadge()))

		if len(option.Comments) != 0 {
			buffer.WriteString(strings.Join(option.Comments, "\n\n") + "\n\n")
		}

		buffer.WriteString(fmt.Sprintf("```%s\n", jenny.Language))
		buffer.WriteString(jenny.Formatter.OptionSignature(context, builder, option))
		buffer.WriteString("\n```\n")

		buffer.WriteString("\n")
	}

	err := jenny.renderIfExists(&buffer, template.ExtraBuilderDocsBlock(builder), map[string]any{
		"Builder": builder,
	})
	if err != nil {
		return codejen.File{}, err
	}

	buffer.WriteString("## See also\n\n")

	if builder.Package == builder.For.SelfRef.ReferredPkg {
		buffer.WriteString(fmt.Sprintf(" * %[2]s [%[1]s](./%[1]s.md)\n", jenny.Formatter.ObjectName(builder.For), jenny.kindBadge(builder.For.Type.Kind)))
	} else {
		buffer.WriteString(fmt.Sprintf(" * %[3]s [%[1]s.%[2]s](../%[1]s/%[2]s.md)\n", builder.For.SelfRef.ReferredPkg, jenny.Formatter.ObjectName(builder.For), jenny.kindBadge(builder.For.Type.Kind)))
	}

	return *codejen.NewFile(fmt.Sprintf("docs/%s/%s.md", builder.Package, builderName), buffer.Bytes(), jenny), nil
}

func (jenny APIReference) kindBadge(kind ast.Kind) string {
	return fmt.Sprintf("<span class=\"badge object-type-%s\"></span>", jenny.Formatter.KindName(kind))
}

func (jenny APIReference) methodBadge() string {
	return "<span class=\"badge object-method\"></span>"
}

func (jenny APIReference) functionBadge() string {
	return "<span class=\"badge function\"></span>"
}

func (jenny APIReference) builderBadge() string {
	return "<span class=\"badge builder\"></span>"
}

func (jenny APIReference) renderIfExists(buffer *bytes.Buffer, blockName string, data any) error {
	if !jenny.Tmpl.Exists(blockName) {
		return nil
	}

	rendered, err := jenny.Tmpl.Render(blockName, data)
	if err != nil {
		return err
	}

	buffer.WriteString(rendered)

	return nil
}
