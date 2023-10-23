package typescript

import (
	"embed"
	"encoding/json"
	"text/template"

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
	base.Funcs(map[string]any{
		"jsonEncode":     mustJSONEncode,
		"upperCamelCase": tools.UpperCamelCase,
		"setOperator":    setOperator,
		"setCloser":      setCloser,
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

func setOperator(valueType Type) string {
	switch valueType {
	case TypeInterface, TypeEnum:
		return "{\n"
	case TypeConst, TypeType:
		return "="
	}
	return ""
}

func setCloser(valueType Type) string {
	switch valueType {
	case TypeInterface, TypeEnum:
		return "\n}"
	case TypeConst, TypeType:
		return ";"
	}
	return ""
}
