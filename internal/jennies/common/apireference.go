package common

import (
	"bytes"
	"fmt"
	"slices"
	"strings"

	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/jennies/template"
	"github.com/grafana/cog/internal/orderedmap"
	"github.com/grafana/cog/pkg/ir"
	"github.com/grafana/cog/pkg/languages"
)

type ArgumentReference struct {
	Name     string
	Type     string
	Comments []string
}

type MethodReference struct {
	ReceiverObject  *ir.Object
	ReceiverBuilder *ir.Builder

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

type VirtualObject struct {
	Object  ir.Object
	Methods []MethodReference
}

type APIReferenceCollector struct {
	virtualObjects   map[string]map[string]VirtualObject
	objectMethods    map[string][]MethodReference
	builderMethods   map[string][]MethodReference
	packageFunctions map[string][]FunctionReference
}

func NewAPIReferenceCollector() *APIReferenceCollector {
	return &APIReferenceCollector{
		virtualObjects:   make(map[string]map[string]VirtualObject),
		objectMethods:    make(map[string][]MethodReference),
		builderMethods:   make(map[string][]MethodReference),
		packageFunctions: make(map[string][]FunctionReference),
	}
}

func (collector *APIReferenceCollector) VirtualObject(object ir.Object) {
	objectRef := object.SelfRef.String()
	pkg := object.SelfRef.ReferredPkg
	if collector.virtualObjects[pkg] == nil {
		collector.virtualObjects[pkg] = make(map[string]VirtualObject)
	}

	if _, ok := collector.virtualObjects[pkg][objectRef]; ok {
		return
	}

	collector.virtualObjects[pkg][objectRef] = VirtualObject{
		Object: object,
	}
}

func (collector *APIReferenceCollector) VirtualObjectMethod(object ir.Object, method MethodReference) {
	pkg := object.SelfRef.ReferredPkg
	objectRef := object.SelfRef.String()

	collector.VirtualObject(object)

	virtualObject := collector.virtualObjects[pkg][objectRef]
	virtualObject.Methods = append(virtualObject.Methods, method)

	collector.virtualObjects[pkg][objectRef] = virtualObject
}

func (collector *APIReferenceCollector) ObjectMethod(object ir.Object, methodReference MethodReference) {
	objectRef := object.SelfRef.String()
	methodReference.ReceiverObject = &object
	collector.objectMethods[objectRef] = append(collector.objectMethods[objectRef], methodReference)
}

func (collector *APIReferenceCollector) methodsForObject(object ir.Object) []MethodReference {
	pkg := object.SelfRef.ReferredPkg
	objectRef := object.SelfRef.String()

	if collector.virtualObjects[pkg] != nil && len(collector.virtualObjects[pkg][objectRef].Methods) != 0 {
		return collector.virtualObjects[pkg][objectRef].Methods
	}

	return collector.objectMethods[objectRef]
}

func (collector *APIReferenceCollector) BuilderMethod(builder ir.Builder, methodReference MethodReference) {
	ref := fmt.Sprintf("%s_%s", builder.Package, builder.Name)
	methodReference.ReceiverBuilder = &builder
	collector.builderMethods[ref] = append(collector.builderMethods[ref], methodReference)
}

func (collector *APIReferenceCollector) methodsForBuilder(builder ir.Builder) []MethodReference {
	ref := fmt.Sprintf("%s_%s", builder.Package, builder.Name)
	return collector.builderMethods[ref]
}

func (collector *APIReferenceCollector) RegisterFunction(pkg string, functionReference FunctionReference) {
	collector.packageFunctions[pkg] = append(collector.packageFunctions[pkg], functionReference)
}

func (collector *APIReferenceCollector) functionsForPackage(pkg string) []FunctionReference {
	return collector.packageFunctions[pkg]
}

type APIReferenceFormatter struct {
	KindName func(kind ir.Kind) string

	FunctionName      func(function FunctionReference) string
	FunctionSignature func(context languages.Context, function FunctionReference) string

	ObjectName       func(object ir.Object) string
	ObjectDefinition func(context languages.Context, object ir.Object) string

	MethodName      func(method MethodReference) string
	MethodSignature func(context languages.Context, method MethodReference) string

	BuilderName          func(builder ir.Builder) string
	ConstructorSignature func(context languages.Context, builder ir.Builder) string
	OptionName           func(option ir.Option) string
	OptionSignature      func(context languages.Context, builder ir.Builder, option ir.Option) string
}

type APIReference struct {
	Collector *APIReferenceCollector
	Language  string
	Formatter APIReferenceFormatter
	Tmpl      *template.Template
}

func (jenny APIReference) JennyName() string {
	return fmt.Sprintf("APIReference[%s]", jenny.Language)
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

	buffer.WriteString("# Packages\n\n")

	slices.SortFunc(context.Schemas, func(schemaA, schemaB *ir.Schema) int {
		return strings.Compare(schemaA.Package, schemaB.Package)
	})

	for _, schema := range context.Schemas {
		badge := jenny.packageBadge(schema)
		if badge != "" {
			badge += " "
		}
		fmt.Fprintf(&buffer, " * %[1]s[%[2]s](./%[2]s/index.md)\n", badge, schema.Package)
	}

	return *codejen.NewFile("docs/Reference/index.md", buffer.Bytes(), jenny), nil
}

func (jenny APIReference) referenceForSchema(context languages.Context, schema *ir.Schema) (codejen.Files, error) {
	files := make([]codejen.File, 0, schema.Objects.Len()+1)

	schemaIndexFile, err := jenny.schemaIndex(context, schema)
	if err != nil {
		return nil, err
	}
	files = append(files, schemaIndexFile)

	var inner error
	schema.Objects.Iterate(func(_ string, object ir.Object) {
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

	virtualObjects := jenny.Collector.virtualObjects[schema.Package]
	for _, virtualObject := range virtualObjects {
		objFile, err := jenny.referenceForObject(context, virtualObject.Object)
		if err != nil {
			inner = err
		}
		files = append(files, objFile)
	}

	return files, nil
}

func (jenny APIReference) schemaIndex(context languages.Context, schema *ir.Schema) (codejen.File, error) {
	var buffer bytes.Buffer

	badge := jenny.packageBadge(schema)
	if badge != "" {
		badge += " "
	}

	fmt.Fprintf(&buffer, "# %s%s\n\n", badge, schema.Package)

	buffer.WriteString("## Objects\n\n")

	objects := orderedmap.New[string, string]()
	schema.Objects.Iterate(func(_ string, object ir.Object) {
		objects.Set(object.Name, fmt.Sprintf(" * %[2]s [%[1]s](./object-%[1]s.md)\n", jenny.Formatter.ObjectName(object), jenny.kindBadge(object.Type.Kind)))
	})
	for _, virtualObject := range jenny.Collector.virtualObjects[schema.Package] {
		objects.Set(virtualObject.Object.Name, fmt.Sprintf(" * %[2]s [%[1]s](./object-%[1]s.md)\n", jenny.Formatter.ObjectName(virtualObject.Object), jenny.kindBadge(virtualObject.Object.Type.Kind)))
	}

	objects.Sort(orderedmap.SortStrings)
	objects.Iterate(func(_ string, value string) {
		buffer.WriteString(value)
	})

	buffer.WriteString("## Builders\n\n")

	builders := context.Builders.GetByPackage(schema.Package)
	slices.SortFunc(builders, func(builderA, builderB ir.Builder) int {
		return strings.Compare(builderA.Name, builderB.Name)
	})

	for _, builder := range builders {
		fmt.Fprintf(&buffer, " * %[2]s [%[1]s](./builder-%[1]s.md)\n", jenny.Formatter.BuilderName(builder), jenny.builderBadge(builder))
	}

	functions := jenny.Collector.functionsForPackage(schema.Package)

	if len(functions) > 0 {
		buffer.WriteString("## Functions\n\n")

		for _, functionReference := range functions {
			fmt.Fprintf(&buffer, "### %[2]s %[1]s\n\n", jenny.Formatter.FunctionName(functionReference), jenny.functionBadge())

			if len(functionReference.Comments) != 0 {
				buffer.WriteString(strings.Join(functionReference.Comments, "\n\n"))
				buffer.WriteString("\n\n")
			}

			fmt.Fprintf(&buffer, "```%s\n", jenny.Language)
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

	return *codejen.NewFile(fmt.Sprintf("docs/Reference/%s/index.md", schema.Package), buffer.Bytes(), jenny), nil
}

func (jenny APIReference) referenceForObject(context languages.Context, object ir.Object) (codejen.File, error) {
	var buffer bytes.Buffer

	objectName := jenny.Formatter.ObjectName(object)

	fmt.Fprintf(&buffer, `---
title: %[2]s %[1]s
---
`, objectName, jenny.objectBadge(object))

	fmt.Fprintf(&buffer, "# %[2]s %[1]s\n\n", objectName, jenny.objectBadge(object))

	if object.DeprecationMessage != "" {
		buffer.WriteString(jenny.deprecatedBadge())
		buffer.WriteString(object.DeprecationMessage)
		buffer.WriteString("\n\n")
	}

	if len(object.Comments) != 0 {
		buffer.WriteString(strings.Join(object.Comments, "\n\n"))
		buffer.WriteString("\n\n")
	}

	buffer.WriteString("## Definition\n\n")

	fmt.Fprintf(&buffer, "```%s\n", jenny.Language)
	buffer.WriteString(jenny.Formatter.ObjectDefinition(context, object))
	buffer.WriteString("\n```\n")

	methods := jenny.Collector.methodsForObject(object)
	if len(methods) != 0 {
		jenny.referenceStructMethods(&buffer, context, methods)
	}

	err := jenny.renderIfExists(&buffer, template.ExtraObjectDocsBlock(object), map[string]any{
		"Object": object,
	})
	if err != nil {
		return codejen.File{}, err
	}

	buildersForObjet := context.Builders.GetAllByObject(object.SelfRef.ReferredPkg, object.SelfRef.ReferredType)
	if len(buildersForObjet) != 0 {
		buffer.WriteString("## See also\n\n")

		slices.SortFunc(buildersForObjet, func(builderA, builderB ir.Builder) int {
			builderAName := fmt.Sprintf("%s.%s", builderA.Package, builderA.Name)
			builderBName := fmt.Sprintf("%s.%s", builderB.Package, builderB.Name)
			return strings.Compare(builderAName, builderBName)
		})
		for _, builder := range buildersForObjet {
			if builder.Package == object.SelfRef.ReferredPkg {
				fmt.Fprintf(&buffer, " * %[2]s [%[1]s](./builder-%[1]s.md)\n", jenny.Formatter.BuilderName(builder), jenny.builderBadge(builder))
			} else {
				fmt.Fprintf(&buffer, " * %[3]s [%[1]s.%[2]s](../%[1]s/builder-%[2]s.md)\n", builder.Package, jenny.Formatter.BuilderName(builder), jenny.builderBadge(builder))
			}
		}
	}

	return *codejen.NewFile(fmt.Sprintf("docs/Reference/%s/object-%s.md", object.SelfRef.ReferredPkg, objectName), buffer.Bytes(), jenny), nil
}

func (jenny APIReference) referenceStructMethods(buffer *bytes.Buffer, context languages.Context, methods []MethodReference) {
	buffer.WriteString("## Methods\n\n")

	for _, method := range methods {
		jenny.formatMethodReference(buffer, context, method)
		buffer.WriteString("\n")
	}

	if len(methods) == 0 {
		buffer.WriteString("No methods.\n")
	}
}

func (jenny APIReference) formatMethodReference(buffer *bytes.Buffer, context languages.Context, method MethodReference) {
	fmt.Fprintf(buffer, "### %[2]s %[1]s\n\n", jenny.Formatter.MethodName(method), jenny.methodBadge())

	if len(method.Comments) != 0 {
		buffer.WriteString(strings.Join(method.Comments, "\n\n"))
		buffer.WriteString("\n\n")
	}

	fmt.Fprintf(buffer, "```%s\n", jenny.Language)
	buffer.WriteString(jenny.Formatter.MethodSignature(context, method))
	buffer.WriteString("\n```\n")
}

func (jenny APIReference) referenceForBuilder(context languages.Context, builder ir.Builder) (codejen.File, error) {
	var buffer bytes.Buffer

	builderName := jenny.Formatter.BuilderName(builder)

	fmt.Fprintf(&buffer, `---
title: %[2]s %[1]s
---
`, builderName, jenny.builderBadge(builder))

	fmt.Fprintf(&buffer, "# %[2]s %[1]s\n\n", builderName, jenny.builderBadge(builder))

	if builder.DeprecationMessage != "" {
		buffer.WriteString(jenny.deprecatedBadge())
		buffer.WriteString(builder.DeprecationMessage)
		buffer.WriteString("\n\n")
	}

	if len(builder.For.Comments) != 0 {
		buffer.WriteString(strings.Join(builder.For.Comments, "\n\n"))
		buffer.WriteString("\n\n")
	}

	buffer.WriteString("## Constructor\n\n")

	fmt.Fprintf(&buffer, "```%s\n", jenny.Language)
	buffer.WriteString(jenny.Formatter.ConstructorSignature(context, builder))
	buffer.WriteString("\n```\n")

	buffer.WriteString("## Methods\n\n")

	builderMethods := jenny.Collector.methodsForBuilder(builder)
	slices.SortFunc(builderMethods, func(methodA, methodB MethodReference) int {
		return strings.Compare(methodA.Name, methodB.Name)
	})

	for _, method := range builderMethods {
		jenny.formatMethodReference(&buffer, context, method)

		buffer.WriteString("\n")
	}

	slices.SortFunc(builder.Options, func(optionA, optionB ir.Option) int {
		return strings.Compare(optionA.Name, optionB.Name)
	})

	for _, option := range builder.Options {
		fmt.Fprintf(&buffer, "### %[2]s %[1]s\n\n", jenny.Formatter.OptionName(option), jenny.methodBadge())

		if len(option.Comments) != 0 {
			buffer.WriteString(strings.Join(option.Comments, "\n\n"))
			buffer.WriteString("\n\n")
		}

		fmt.Fprintf(&buffer, "```%s\n", jenny.Language)
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
		fmt.Fprintf(&buffer, " * %[2]s [%[1]s](./object-%[1]s.md)\n", jenny.Formatter.ObjectName(builder.For), jenny.kindBadge(builder.For.Type.Kind))
	} else {
		fmt.Fprintf(&buffer, " * %[3]s [%[1]s.%[2]s](../%[1]s/object-%[2]s.md)\n", builder.For.SelfRef.ReferredPkg, jenny.Formatter.ObjectName(builder.For), jenny.kindBadge(builder.For.Type.Kind))
	}

	return *codejen.NewFile(fmt.Sprintf("docs/Reference/%s/builder-%s.md", builder.Package, builderName), buffer.Bytes(), jenny), nil
}

func (jenny APIReference) packageBadge(schema *ir.Schema) string {
	if schema.Metadata.Kind == ir.SchemaKindCore {
		return "<span class=\"badge package-core\"></span>"
	}

	if schema.Metadata.Variant == "" {
		return ""
	}

	return fmt.Sprintf("<span class=\"badge package-variant-%s\"></span>", string(schema.Metadata.Variant))
}

func (jenny APIReference) objectBadge(object ir.Object) string {
	badge := jenny.kindBadge(object.Type.Kind)

	if object.DeprecationMessage != "" {
		badge += " " + jenny.deprecatedBadge()
	}

	return badge
}

func (jenny APIReference) kindBadge(kind ir.Kind) string {
	return fmt.Sprintf("<span class=\"badge object-type-%s\"></span>", jenny.Formatter.KindName(kind))
}

func (jenny APIReference) methodBadge() string {
	return "<span class=\"badge object-method\"></span>"
}

func (jenny APIReference) functionBadge() string {
	return "<span class=\"badge function\"></span>"
}

func (jenny APIReference) builderBadge(builder ir.Builder) string {
	badge := "<span class=\"badge builder\"></span>"

	if builder.DeprecationMessage != "" {
		badge += " " + jenny.deprecatedBadge()
	}

	return badge
}

func (jenny APIReference) deprecatedBadge() string {
	return "<span class=\"badge deprecated\"></span>"
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
