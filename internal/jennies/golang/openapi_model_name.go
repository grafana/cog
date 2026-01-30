package golang

import (
	"fmt"
	"strings"

	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/jennies/common"
	"github.com/grafana/cog/internal/jennies/template"
	"github.com/grafana/cog/internal/languages"
)

type openAPIModelNameMethods struct {
	config          Config
	tmpl            *template.Template
	apiRefCollector *common.APIReferenceCollector
}

func newOpenAPIModelNameMethods(config Config, tmpl *template.Template, apiRefCollector *common.APIReferenceCollector) openAPIModelNameMethods {
	return openAPIModelNameMethods{
		config:          config,
		tmpl:            tmpl,
		apiRefCollector: apiRefCollector,
	}
}

func (jenny openAPIModelNameMethods) generateForObject(buffer *strings.Builder, context languages.Context, object ast.Object, namerFunc func(string) string) error {
	if !object.Type.IsStruct() {
		return nil
	}
	objName := formatObjectName(object.Name)
	modelName := namerFunc(objName)

	jenny.apiRefCollector.ObjectMethod(object, common.MethodReference{
		Name: "OpenAPIModelName",
		Comments: []string{
			fmt.Sprintf("OpenAPIModelName returns the OpenAPI model name for `%s`.", objName),
		},
		Return: "string",
	})

	rendered, err := jenny.tmpl.Render("types/struct_openapi_model_name_method.tmpl", map[string]any{
		"objName": objName,
		"modelName": modelName,
	})
	if err != nil {
		return err
	}

	_, err = buffer.WriteString(rendered)
	return err
}
