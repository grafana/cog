package java

import (
	"strings"

	"github.com/grafana/cog/internal/ast"
)

type JSONMarshaller struct {
	config        Config
	typeFormatter *typeFormatter
}

func (j JSONMarshaller) genToJSONFunction(t ast.Type) string {
	if !j.config.generateBuilders || j.config.SkipRuntime {
		return ""
	}

	var buffer strings.Builder
	j.typeFormatter.packageMapper(fasterXmlPackageName, "core.JsonProcessingException")
	j.typeFormatter.packageMapper(fasterXmlPackageName, "databind.ObjectMapper")
	j.typeFormatter.packageMapper(fasterXmlPackageName, "databind.ObjectWriter")
	if t.IsStructGeneratedFromDisjunction() {
		if t.IsStruct() && (t.HasHint(ast.HintDiscriminatedDisjunctionOfRefs) || t.HasHint(ast.HintDisjunctionOfScalars)) {
			_ = templates.ExecuteTemplate(&buffer, "marshalling/disjunctions.json_marshall.tmpl", map[string]any{
				"Fields": t.AsStruct().Fields,
			})
			return buffer.String()
		}
	}

	_ = templates.ExecuteTemplate(&buffer, "marshalling/marshalling.tmpl", map[string]any{})
	return buffer.String()
}

func (j JSONMarshaller) annotation(t ast.Type) string {
	if !j.config.generateBuilders || j.config.SkipRuntime {
		return ""
	}

	if t.IsStructGeneratedFromDisjunction() && t.IsStruct() {
		j.typeFormatter.packageMapper(fasterXmlPackageName, "annotation.JsonUnwrapped")
		return "@JsonUnwrapped"
	}

	j.typeFormatter.packageMapper(fasterXmlPackageName, "annotation.JsonProperty")
	return "@JsonProperty(%#v)"
}
