package golang

import (
	"strings"

	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/jennies/common"
	"github.com/grafana/cog/internal/jennies/template"
	"github.com/grafana/cog/internal/languages"
)

type validationMethods struct {
	tmpl          *template.Template
	packageMapper func(string) string
}

func newValidationMethods(tmpl *template.Template, packageMapper func(string) string) validationMethods {
	return validationMethods{
		tmpl:          tmpl,
		packageMapper: packageMapper,
	}
}

func (jenny validationMethods) generateForObject(buffer *strings.Builder, context languages.Context, schema *ast.Schema, object ast.Object, imports *common.DirectImportMap) error {
	if !object.Type.IsStruct() {
		return nil
	}

	var resolvesToConstraints func(typeDef ast.Type) bool
	resolvesToConstraints = func(typeDef ast.Type) bool {
		if typeDef.IsAny() {
			return false
		}

		if typeDef.IsComposableSlot() {
			return true
		}

		if typeDef.IsRef() {
			return context.ResolveRefs(typeDef).IsStruct()
		}

		if typeDef.IsScalar() {
			return len(typeDef.AsScalar().Constraints) != 0
		}

		if typeDef.IsDisjunction() {
			for _, branch := range typeDef.AsDisjunction().Branches {
				if resolvesToConstraints(branch) {
					return true
				}
			}
		}

		if typeDef.IsIntersection() {
			for _, branch := range typeDef.AsIntersection().Branches {
				if resolvesToConstraints(branch) {
					return true
				}
			}
		}

		if typeDef.IsStruct() {
			for _, field := range typeDef.AsStruct().Fields {
				if resolvesToConstraints(field.Type) {
					return true
				}
			}
		}

		if typeDef.IsMap() {
			return resolvesToConstraints(typeDef.AsMap().ValueType)
		}

		if typeDef.IsArray() {
			return resolvesToConstraints(typeDef.AsArray().ValueType)
		}

		return false
	}

	tmpl := jenny.tmpl.
		Funcs(common.TypeResolvingTemplateHelpers(context)).
		Funcs(template.FuncMap{
			"resolvesToConstraints": resolvesToConstraints,
			"importPkg":             jenny.packageMapper,
			"importStdPkg": func(pkg string) string {
				return imports.Add(pkg, pkg)
			},
		})

	rendered, err := tmpl.Render("types/struct_validation_method.tmpl", map[string]any{
		"def": object,
	})
	if err != nil {
		return err
	}
	buffer.WriteString(rendered)

	return nil
}
