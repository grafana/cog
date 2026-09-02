package apiref

import (
	"github.com/grafana/cog/internal/tools"
	"github.com/grafana/cog/pkg/ir"
	"github.com/grafana/cog/pkg/template"
)

func TemplateHelpers(apiRefCollector *APIReferenceCollector) template.FuncMap {
	return template.FuncMap{
		"apiDeclareFunction": func(data map[string]any) string {
			apiRefCollector.RegisterFunction(maybeGet[string](data, "pkg"), FunctionReference{
				Name:      maybeGet[string](data, "name"),
				Comments:  maybeGet[[]string](data, "comments"),
				Arguments: tools.Map(maybeGet[[]map[string]any](data, "arguments"), dataToArgumentRef),
				Return:    maybeGet[string](data, "return"),
			})

			return ""
		},
		"apiDeclareMethod": func(data map[string]any) string {
			apiRefCollector.ObjectMethod(maybeGet[ir.Object](data, "object"), MethodReference{
				Name:      maybeGet[string](data, "name"),
				Comments:  maybeGet[[]string](data, "comments"),
				Arguments: tools.Map(maybeGet[[]map[string]any](data, "arguments"), dataToArgumentRef),
				Return:    maybeGet[string](data, "return"),
				Static:    false,
			})

			return ""
		},
	}
}

func dataToArgumentRef(data map[string]any) ArgumentReference {
	return ArgumentReference{
		Name:     maybeGet[string](data, "name"),
		Type:     maybeGet[string](data, "type"),
		Comments: maybeGet[[]string](data, "comments"),
	}
}

func maybeGet[T any](data map[string]any, key string) T {
	var result T
	if data[key] == nil {
		return result
	}

	return data[key].(T)
}
