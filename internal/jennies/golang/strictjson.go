package golang

import (
	"fmt"
	"strings"

	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/jennies/common"
	"github.com/grafana/cog/internal/jennies/template"
	"github.com/grafana/cog/internal/languages"
)

type strictJSONUnmarshal struct {
	tmpl          *template.Template
	imports       *common.DirectImportMap
	packageMapper func(string) string
	typeFormatter *typeFormatter
}

func newStrictJSONUnmarshal(tmpl *template.Template, imports *common.DirectImportMap, packageMapper func(string) string, typeFormatter *typeFormatter) strictJSONUnmarshal {
	return strictJSONUnmarshal{
		tmpl: tmpl.Funcs(template.FuncMap{
			"formatType": typeFormatter.formatType,
			"importStdPkg": func(pkg string) string {
				return imports.Add(pkg, pkg)
			},
			"importPkg": packageMapper,
		}),
		imports:       imports,
		packageMapper: packageMapper,
		typeFormatter: typeFormatter,
	}
}

func (jenny strictJSONUnmarshal) generateForObject(buffer *strings.Builder, context languages.Context, object ast.Object) error {
	if jenny.objectNeedsUnmarshal(object) {
		jsonUnmarshal, err := jenny.renderUnmarshal(context, object)
		if err != nil {
			return err
		}
		buffer.WriteString(jsonUnmarshal)
		buffer.WriteString("\n")
	}

	return nil
}

func (jenny strictJSONUnmarshal) objectNeedsUnmarshal(obj ast.Object) bool {
	return obj.Type.IsStruct()
}

func (jenny strictJSONUnmarshal) renderUnmarshal(context languages.Context, obj ast.Object) (string, error) {
	customUnmarshalTmpl := jenny.customObjectUnmarshalBlock(obj)
	if jenny.tmpl.Exists(customUnmarshalTmpl) {
		return jenny.tmpl.Render(customUnmarshalTmpl, map[string]any{
			"Object": obj,
		})
	}

	return jenny.renderStructUnmarshal(context, obj)
}

func (jenny strictJSONUnmarshal) renderStructUnmarshal(context languages.Context, obj ast.Object) (string, error) {
	return jenny.tmpl.
		Funcs(template.FuncMap{
			"resolvesToScalar": func(typeDef ast.Type) bool {
				return context.ResolveRefs(typeDef).IsScalar()
			},
			"resolvesToStruct": func(typeDef ast.Type) bool {
				return context.ResolveRefs(typeDef).IsStruct()
			},
			"formatRawRef": func(pkg string, ref string) string {
				return jenny.typeFormatter.formatRef(ast.NewRef(pkg, ref), false)
			},
		}).
		Render("types/struct.strict.json_unmarshal.tmpl", map[string]any{
			"def": obj,
		})
}

func (jenny strictJSONUnmarshal) customObjectUnmarshalBlock(obj ast.Object) string {
	return fmt.Sprintf("object_%s_%s_custom_strict_unmarshal", obj.SelfRef.ReferredPkg, obj.SelfRef.ReferredType)
}
