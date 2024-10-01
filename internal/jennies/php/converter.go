package php

import (
	"fmt"
	"path/filepath"

	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/jennies/template"
	"github.com/grafana/cog/internal/languages"
)

type Converter struct {
	config         Config
	nullableConfig languages.NullableConfig
	tmpl           *template.Template
}

func (jenny *Converter) JennyName() string {
	return "PHPConverter"
}

func (jenny *Converter) Generate(context languages.Context) (codejen.Files, error) {
	files := codejen.Files{}

	for _, builder := range context.Builders {
		output, err := jenny.generateConverter(context, builder)
		if err != nil {
			return nil, err
		}

		filename := filepath.Join(
			"src",
			formatPackageName(builder.Package),
			fmt.Sprintf("%sConverter.php", formatObjectName(builder.Name)),
		)

		files = append(files, *codejen.NewFile(filename, output, jenny))
	}

	return files, nil
}

func (jenny *Converter) generateConverter(context languages.Context, builder ast.Builder) ([]byte, error) {
	converter := languages.NewConverterGenerator(jenny.nullableConfig).FromBuilder(context, builder)
	formatter := builderTypeFormatter(jenny.config, context)

	return jenny.tmpl.
		Funcs(map[string]any{
			"resolvesToEnum": func(typeDef ast.Type) bool {
				return context.ResolveRefs(typeDef).IsEnum()
			},
			"formatType": builderTypeFormatter(jenny.config, context).formatType,
			"formatPath": formatFieldPath,
			"formatRawRef": func(pkg string, ref string) string {
				return formatter.formatRef(ast.NewRef(pkg, ref), false)
			},
			"disjunctionCaseForType": func(input string, typeDef ast.Type) string {
				// TODO: shaky at best
				if typeDef.IsAnyOf(ast.KindArray, ast.KindMap) {
					return fmt.Sprintf("is_array(%s)", input)
				}

				if typeDef.IsScalar() {
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

					testFunc := testMap[typeDef.Scalar.ScalarKind]
					if testFunc == "" {
						return "/* unhandled scalar type */"
					}

					return fmt.Sprintf("%s(%s)", testFunc, input)
				}

				if typeDef.IsRef() {
					return fmt.Sprintf("%s instanceof %s", input, formatter.formatRef(typeDef.Ref.AsType(), false))
				}

				return "/* unhandled type */"
			},
		}).
		RenderAsBytes("converters/converter.tmpl", map[string]any{
			"NamespaceRoot": jenny.config.NamespaceRoot,
			"Converter":     converter,
		})
}
