package typescript

import (
	"embed"
	"encoding/json"
	"text/template"

	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/tools"
)

//nolint:gochecknoglobals
var templates *template.Template

//go:embed templates/*.tmpl
//nolint:gochecknoglobals
var templatesFS embed.FS

//nolint:gochecknoinits
func init() {
	base := template.New("ts")
	base.
		Option("missingkey=error").
		Funcs(map[string]any{
			// placeholder functions, will be overridden by jennies
			"typeHasBuilder":           func(_ ast.Type) bool { return false },
			"resolvesToComposableSlot": func(_ ast.Type) bool { return false },

			"jsonEncode":     mustJSONEncode,
			"upperCamelCase": tools.UpperCamelCase,
			"lowerCamelCase": tools.LowerCamelCase,
			"formatScalar":   formatScalar,
		})
	templates = template.Must(base.ParseFS(templatesFS, "templates/*.tmpl"))
}

func mustJSONEncode(val any) string {
	encoded, err := json.MarshalIndent(val, "", "  ")
	if err != nil {
		panic(err)
	}

	return string(encoded)
}
