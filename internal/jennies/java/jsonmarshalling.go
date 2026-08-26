package java

import (
	"github.com/grafana/cog/internal/ir"
	"github.com/grafana/cog/internal/jennies/template"
)

type JSONMarshaller struct {
	config        Config
	tmpl          *template.Template
	typeFormatter *typeFormatter
}

func (j JSONMarshaller) genToJSONFunction() string {
	if !j.config.GenerateJSONMarshaller || !j.config.GenerateBuilders || j.config.SkipRuntime {
		return ""
	}

	j.typeFormatter.packageMapper(fasterXMLPackageName, "core.JsonProcessingException")
	j.typeFormatter.packageMapper(fasterXMLPackageName, "databind.ObjectMapper")
	j.typeFormatter.packageMapper(fasterXMLPackageName, "databind.ObjectWriter")
	rendered, _ := j.tmpl.Render("marshalling/marshalling.tmpl", map[string]any{})
	return rendered
}

func (j JSONMarshaller) annotation(t ir.Type) string {
	if !j.config.GenerateJSONMarshaller || !j.config.GenerateBuilders || j.config.SkipRuntime {
		return ""
	}

	if t.IsStructGeneratedFromDisjunction() && t.IsStruct() {
		j.typeFormatter.packageMapper(fasterXMLPackageName, "annotation.JsonUnwrapped")
		return "@JsonUnwrapped"
	}

	j.typeFormatter.packageMapper(fasterXMLPackageName, "annotation.JsonProperty")
	return "@JsonProperty(%#v)"
}
