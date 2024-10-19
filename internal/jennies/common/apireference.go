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
	"github.com/grafana/cog/internal/tools"
)

type APIReference struct {
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
		files = append(files, *jenny.referenceForBuilder(builder))
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
		buffer.WriteString(fmt.Sprintf(" * %s\n", object.Name))
	})

	buffer.WriteString("## Builders\n\n")

	builders := context.Builders.ByPackage(schema.Package)
	slices.SortFunc(builders, func(builderA, builderB ast.Builder) int {
		return strings.Compare(builderA.Name, builderB.Name)
	})

	for _, builder := range builders {
		// TODO: builder name is language-dependent
		buffer.WriteString(fmt.Sprintf(" * [%[1]s](./%[1]s.md)\n", builder.Name, schema.Package))
	}

	return codejen.NewFile(fmt.Sprintf("docs/%s/index.md", schema.Package), buffer.Bytes(), jenny)
}

func (jenny APIReference) referenceForBuilder(builder ast.Builder) *codejen.File {
	var buffer bytes.Buffer

	buffer.WriteString(fmt.Sprintf("# %s\n\n", builder.Name))

	buffer.WriteString("## Methods\n\n")

	slices.SortFunc(builder.Options, func(optionA, optionB ast.Option) int {
		return strings.Compare(optionA.Name, optionB.Name)
	})

	for _, option := range builder.Options {
		buffer.WriteString(fmt.Sprintf("### %s\n\n", option.Name))

		args := tools.Map(option.Args, func(arg ast.Argument) string {
			return arg.Name
		})

		if len(option.Comments) != 0 {
			buffer.WriteString(strings.Join(option.Comments, "\n\n") + "\n\n")
		}

		// TODO: configurable language tag
		buffer.WriteString("```typescript\n")
		// TODO: language-dependent function definition formatting
		buffer.WriteString(fmt.Sprintf("%[1]s(%[2]s)\n", option.Name, strings.Join(args, ", ")))
		buffer.WriteString("```\n")

		buffer.WriteString("\n")
	}

	return codejen.NewFile(fmt.Sprintf("docs/%s/%s.md", builder.Package, builder.Name), buffer.Bytes(), jenny)
}
