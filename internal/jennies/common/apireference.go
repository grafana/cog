package common

import (
	"bytes"
	"fmt"
	"slices"
	"strings"

	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/ast"
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

type APIReferenceCollector struct {
	objectMethods map[string][]MethodReference
}

func NewAPIReferenceCollector() *APIReferenceCollector {
	return &APIReferenceCollector{
		objectMethods: make(map[string][]MethodReference),
	}
}

func (collector *APIReferenceCollector) RegisterMethodForObject(object ast.Object, methodReference MethodReference) {
	objectRef := object.SelfRef.String()
	collector.objectMethods[objectRef] = append(collector.objectMethods[objectRef], methodReference)
}

func (collector *APIReferenceCollector) MethodsForObject(object ast.Object) []MethodReference {
	return collector.objectMethods[object.SelfRef.String()]
}

type APIReferenceFormatter struct {
	ObjectName       func(object ast.Object) string
	ObjectDefinition func(context languages.Context, object ast.Object) string

	MethodName      func(method MethodReference) string
	MethodSignature func(context languages.Context, method MethodReference) string

	BuilderName          func(builder ast.Builder) string
	ConstructorSignature func(context languages.Context, builder ast.Builder) string
	OptionName           func(option ast.Option) string
	OptionSignature      func(context languages.Context, option ast.Option) string
}

type APIReference struct {
	Collector *APIReferenceCollector
	Language  string
	Formatter APIReferenceFormatter
}

func (jenny APIReference) JennyName() string {
	return "APIReference"
}

func (jenny APIReference) Generate(context languages.Context) (codejen.Files, error) {
	files := make([]codejen.File, 0, len(context.Schemas)+len(context.Builders)+1)

	for _, schema := range context.Schemas {
		files = append(files, jenny.referenceForSchema(context, schema)...)
	}
	for _, builder := range context.Builders {
		files = append(files, jenny.referenceForBuilder(context, builder))
	}

	files = append(files, jenny.index(context))

	return files, nil
}

func (jenny APIReference) index(context languages.Context) codejen.File {
	var buffer bytes.Buffer

	buffer.WriteString("# Packages\n\n")

	slices.SortFunc(context.Schemas, func(schemaA, schemaB *ast.Schema) int {
		return strings.Compare(schemaA.Package, schemaB.Package)
	})

	for _, schema := range context.Schemas {
		buffer.WriteString(fmt.Sprintf(" * [%[1]s](./%[1]s/index.md)\n", schema.Package))
	}

	return *codejen.NewFile("docs/index.md", buffer.Bytes(), jenny)
}

func (jenny APIReference) referenceForSchema(context languages.Context, schema *ast.Schema) codejen.Files {
	files := make([]codejen.File, 0, schema.Objects.Len()+1)

	files = append(files, jenny.schemaIndex(context, schema))

	schema.Objects.Iterate(func(_ string, object ast.Object) {
		files = append(files, jenny.referenceForObject(context, object))
	})

	return files
}

func (jenny APIReference) schemaIndex(context languages.Context, schema *ast.Schema) codejen.File {
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

	return *codejen.NewFile(fmt.Sprintf("docs/%s/index.md", schema.Package), buffer.Bytes(), jenny)
}

func (jenny APIReference) referenceForObject(context languages.Context, object ast.Object) codejen.File {
	var buffer bytes.Buffer

	objectName := jenny.Formatter.ObjectName(object)

	buffer.WriteString(fmt.Sprintf(`---
title: %[2]s %[1]s
---
`, objectName, jenny.kindBadge(object.Type.Kind)))

	buffer.WriteString(fmt.Sprintf("# %[2]s %[1]s\n\n", objectName, jenny.kindBadge(object.Type.Kind)))

	buffer.WriteString("## Definition\n\n")

	buffer.WriteString(fmt.Sprintf("```%s\n", jenny.Language))
	buffer.WriteString(jenny.Formatter.ObjectDefinition(context, object))
	buffer.WriteString("\n```\n")

	if object.Type.IsStruct() {
		jenny.referenceStructMethods(&buffer, context, object)
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

	return *codejen.NewFile(fmt.Sprintf("docs/%s/%s.md", object.SelfRef.ReferredPkg, objectName), buffer.Bytes(), jenny)
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
		buffer.WriteString(jenny.Formatter.MethodSignature(context, method))
		buffer.WriteString("\n```\n")

		buffer.WriteString("\n")
	}

	if len(methods) == 0 {
		buffer.WriteString("No methods.\n")
	}
}

func (jenny APIReference) referenceForType(buffer *bytes.Buffer, context languages.Context, object ast.Object) {
	buffer.WriteString("## Definition\n\n")

	buffer.WriteString(fmt.Sprintf("```%s\n", jenny.Language))
	buffer.WriteString(jenny.Formatter.ObjectDefinition(context, object))
	buffer.WriteString("\n```\n")
}

func (jenny APIReference) referenceForBuilder(context languages.Context, builder ast.Builder) codejen.File {
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
		buffer.WriteString(jenny.Formatter.OptionSignature(context, option))
		buffer.WriteString("\n```\n")

		buffer.WriteString("\n")
	}

	buffer.WriteString("## Examples\n\n")

	// TODO
	buffer.WriteString("TODO\n")

	buffer.WriteString("## See also\n\n")

	if builder.Package == builder.For.SelfRef.ReferredPkg {
		buffer.WriteString(fmt.Sprintf(" * %[2]s [%[1]s](./%[1]s.md)\n", jenny.Formatter.ObjectName(builder.For), jenny.kindBadge(builder.For.Type.Kind)))
	} else {
		buffer.WriteString(fmt.Sprintf(" * %[3]s [%[1]s.%[2]s](../%[1]s/%[2]s.md)\n", builder.For.SelfRef.ReferredPkg, jenny.Formatter.ObjectName(builder.For), jenny.kindBadge(builder.For.Type.Kind)))
	}

	return *codejen.NewFile(fmt.Sprintf("docs/%s/%s.md", builder.Package, builderName), buffer.Bytes(), jenny)
}

func (jenny APIReference) kindBadge(kind ast.Kind) string {
	return fmt.Sprintf("<span class=\"badge object-type-%s\"></span>", kind)
}

func (jenny APIReference) methodBadge() string {
	return "<span class=\"badge object-method\"></span>"
}

func (jenny APIReference) builderBadge() string {
	return "<span class=\"badge builder\"></span>"
}
