package java

import (
	gotemplate "text/template"

	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/jennies/template"
)

type JSONMarshaller struct {
	config        Config
	tmpl          *gotemplate.Template
	typeFormatter *typeFormatter
}

func (j JSONMarshaller) genToJSONFunction(t ast.Type) string {
	if !j.config.generateBuilders || j.config.SkipRuntime {
		return ""
	}

	j.typeFormatter.packageMapper(fasterXMLPackageName, "core.JsonProcessingException")
	j.typeFormatter.packageMapper(fasterXMLPackageName, "databind.ObjectMapper")
	j.typeFormatter.packageMapper(fasterXMLPackageName, "databind.ObjectWriter")
	if t.IsStructGeneratedFromDisjunction() {
		if t.IsStruct() && (t.HasHint(ast.HintDiscriminatedDisjunctionOfRefs) || t.HasHint(ast.HintDisjunctionOfScalars)) {
			rendered, _ := template.Render(j.tmpl, "marshalling/disjunctions.json_marshall.tmpl", map[string]any{
				"Fields": t.AsStruct().Fields,
			})
			return rendered
		}
	}

	rendered, _ := template.Render(j.tmpl, "marshalling/marshalling.tmpl", map[string]any{})
	return rendered
}

func (j JSONMarshaller) annotation(t ast.Type) string {
	if !j.config.generateBuilders || j.config.SkipRuntime {
		return ""
	}

	if t.IsStructGeneratedFromDisjunction() && t.IsStruct() {
		j.typeFormatter.packageMapper(fasterXMLPackageName, "annotation.JsonUnwrapped")
		return "@JsonUnwrapped"
	}

	j.typeFormatter.packageMapper(fasterXMLPackageName, "annotation.JsonProperty")
	return "@JsonProperty(%#v)"
}
