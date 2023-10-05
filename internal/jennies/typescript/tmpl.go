package typescript

import (
	"embed"
	"strings"
	"text/template"
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
		"simpleQuotes": simpleQuotes,
	})
	templates = template.Must(base.ParseFS(templatesFS, "templates/*.tmpl"))
}

func simpleQuotes(s string) string {
	splitString := strings.Split(s, " ")
	if strings.HasSuffix(splitString[0], "\"") && strings.HasPrefix(splitString[0], "\"") {
		newString := s[1 : len(s)-1]
		return "'" + newString + "'"

	}

	return s
}
