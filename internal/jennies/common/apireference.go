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

type APIReferenceFormatter struct {
	ObjectName      func(object ast.Object) string
	BuilderName     func(builder ast.Builder) string
	OptionName      func(option ast.Option) string
	OptionSignature func(context languages.Context, option ast.Option) string
}

type APIReference struct {
	Language  string
	Formatter APIReferenceFormatter
}

func (jenny APIReference) JennyName() string {
	return "APIReference"
}

func (jenny APIReference) Generate(context languages.Context) (codejen.Files, error) {
	// one file per package + one file per builder + an index
	files := make([]codejen.File, 0, len(context.Schemas)+len(context.Builders)+1)

	for _, schema := range context.Schemas {
		files = append(files, *jenny.referenceForSchema(context, schema))
	}
	for _, builder := range context.Builders {
		files = append(files, *jenny.referenceForBuilder(context, builder))
	}

	files = append(files, *jenny.index(context))

	return files, nil
}

func (jenny APIReference) index(context languages.Context) *codejen.File {
	var buffer bytes.Buffer

	buffer.WriteString("# Packages\n\n")

	slices.SortFunc(context.Schemas, func(schemaA, schemaB *ast.Schema) int {
		return strings.Compare(schemaA.Package, schemaB.Package)
	})

	for _, schema := range context.Schemas {
		buffer.WriteString(fmt.Sprintf(" * [%[1]s](./%[1]s/index.md)\n", schema.Package))
	}

	return codejen.NewFile("docs/index.md", buffer.Bytes(), jenny)
}

func (jenny APIReference) referenceForSchema(context languages.Context, schema *ast.Schema) *codejen.File {
	var buffer bytes.Buffer

	buffer.WriteString(fmt.Sprintf("# %s\n\n", schema.Package))

	buffer.WriteString("## Objects\n\n")

	schema.Objects.Sort(orderedmap.SortStrings)
	schema.Objects.Iterate(func(_ string, object ast.Object) {
		buffer.WriteString(fmt.Sprintf(" * %s\n", jenny.Formatter.ObjectName(object)))
	})

	buffer.WriteString("## Builders\n\n")

	builders := context.Builders.ByPackage(schema.Package)
	slices.SortFunc(builders, func(builderA, builderB ast.Builder) int {
		return strings.Compare(builderA.Name, builderB.Name)
	})

	for _, builder := range builders {
		buffer.WriteString(fmt.Sprintf(" * [%[1]s](./%[1]s.md)\n", jenny.Formatter.BuilderName(builder)))
	}

	return codejen.NewFile(fmt.Sprintf("docs/%s/index.md", schema.Package), buffer.Bytes(), jenny)
}

func (jenny APIReference) referenceForBuilder(context languages.Context, builder ast.Builder) *codejen.File {
	var buffer bytes.Buffer

	builderName := jenny.Formatter.BuilderName(builder)

	buffer.WriteString(fmt.Sprintf("# %s\n\n", builderName))

	buffer.WriteString("## Methods\n\n")

	slices.SortFunc(builder.Options, func(optionA, optionB ast.Option) int {
		return strings.Compare(optionA.Name, optionB.Name)
	})

	for _, option := range builder.Options {
		buffer.WriteString(fmt.Sprintf("### %s\n\n", jenny.Formatter.OptionName(option)))

		if len(option.Comments) != 0 {
			buffer.WriteString(strings.Join(option.Comments, "\n\n") + "\n\n")
		}

		buffer.WriteString(fmt.Sprintf("```%s\n", jenny.Language))
		buffer.WriteString(jenny.Formatter.OptionSignature(context, option))
		buffer.WriteString("\n```\n")

		buffer.WriteString("\n")
	}

	return codejen.NewFile(fmt.Sprintf("docs/%s/%s.md", builder.Package, builderName), buffer.Bytes(), jenny)
}
