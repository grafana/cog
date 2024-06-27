package java

import (
	"strings"

	"github.com/grafana/cog/internal/ast"
)

func genToJSONFunction(t ast.Type) string {
	var buffer strings.Builder
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
