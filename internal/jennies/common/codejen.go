package common

import (
	"bytes"
	"fmt"
	"path/filepath"
	"text/template"

	"github.com/grafana/codejen"
)

type noopOneToManyJenny[Input any] struct {
}

func (jenny noopOneToManyJenny[Input]) JennyName() string {
	return "noopOneToManyJenny"
}

func (jenny noopOneToManyJenny[Input]) Generate(_ Input) (codejen.Files, error) {
	return nil, nil
}

func If[Input any](condition bool, innerJenny codejen.OneToMany[Input]) codejen.OneToMany[Input] {
	if !condition {
		return noopOneToManyJenny[Input]{}
	}

	return innerJenny
}

// GeneratedCommentHeader produces a FileMapper that injects a comment header onto
// a [codejen.File] indicating  the jenny or jennies that constructed the
// file.
func GeneratedCommentHeader() codejen.FileMapper {
	genHeader := `// Code generated - EDITING IS FUTILE. DO NOT EDIT.
//
// Using jennies:
{{- range .Using }}
//     {{ .JennyName }}
{{- end }}

`

	tmpl, err := template.New("cog").Parse(genHeader)
	if err != nil {
		// not ideal, but also not a big deal: this statement is only reachable when the template is invalid.
		panic(err)
	}

	return func(f codejen.File) (codejen.File, error) {
		// Never inject on certain filetypes, it's never valid
		switch filepath.Ext(f.RelativePath) {
		case ".json", ".yml", ".yaml", ".md":
			return f, nil
		default:
			buf := new(bytes.Buffer)
			if err := tmpl.Execute(buf, map[string]any{
				"Using": f.From,
			}); err != nil {
				return codejen.File{}, fmt.Errorf("failed executing GeneratedCommentHeader() template: %w", err)
			}
			buf.Write(f.Data)

			f.Data = buf.Bytes()
		}
		return f, nil
	}
}
