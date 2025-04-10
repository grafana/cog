package typescript

import (
	"fmt"
	"strings"

	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/jennies/common"
	"github.com/grafana/cog/internal/jennies/template"
	"github.com/grafana/cog/internal/languages"
)

type validationMethods struct {
	tmpl            *template.Template
	context         languages.Context
	packageMapper   func(string) string
	apiRefCollector *common.APIReferenceCollector
}

func newValidationMethods(tmpl *template.Template, context languages.Context, packageMapper func(string) string, apiRefCollector *common.APIReferenceCollector) validationMethods {
	return validationMethods{
		tmpl:            tmpl,
		context:         context,
		packageMapper:   packageMapper,
		apiRefCollector: apiRefCollector,
	}
}

func (jenny validationMethods) generateForObject(buffer *strings.Builder, object ast.Object, imports *common.DirectImportMap) error {
	if !object.Type.IsStruct() {
		return nil
	}

	jenny.apiRefCollector.ObjectMethod(object, common.MethodReference{
		Name: "validate",
		Comments: []string{
			fmt.Sprintf("Validate checks all the validation constraints that may be defined on `%s` fields for violations and returns them.", formatObjectName(object.Name)),
		},
		Return: "error",
	})

	return nil
}

func (jenny validationMethods) resolveToConstraints(typeDef ast.Type) bool {
	if typeDef.IsAny() {
		return false
	}

	if typeDef.IsComposableSlot() {
		return true
	}

	if typeDef.IsRef() {
		return jenny.context.ResolveRefs(typeDef).IsStruct()
	}

	if typeDef.IsScalar() {
		return len(typeDef.AsScalar().Constraints) != 0
	}

	if typeDef.IsDisjunction() {
		for _, branch := range typeDef.AsDisjunction().Branches {
			if jenny.resolveToConstraints(branch) {
				return true
			}
		}
	}

	if typeDef.IsIntersection() {
		for _, branch := range typeDef.AsIntersection().Branches {
			if jenny.resolveToConstraints(branch) {
				return true
			}
		}
	}

	if typeDef.IsStruct() {
		for _, field := range typeDef.AsStruct().Fields {
			if jenny.resolveToConstraints(field.Type) {
				return true
			}
		}
	}

	if typeDef.IsMap() {
		return jenny.resolveToConstraints(typeDef.AsMap().ValueType)
	}

	if typeDef.IsArray() {
		return jenny.resolveToConstraints(typeDef.AsArray().ValueType)
	}

	if typeDef.IsConstantRef() {
		obj, _ := jenny.context.LocateObject(typeDef.AsConstantRef().ReferredPkg, typeDef.AsConstantRef().ReferredType)
		return obj.Type.IsEnum()
	}

	return false
}
