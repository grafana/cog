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
func GeneratedCommentHeader(config Config) codejen.FileMapper {
	genHeader := `{{ .Leader }} Code generated - EDITING IS FUTILE. DO NOT EDIT.
{{- with .Using }}
{{ $.Leader }}
{{ $.Leader }} Using jennies:
{{- range . }}
{{ $.Leader }}     {{ .JennyName }}
{{- end }}
{{- end }}

`

	tmpl, err := template.New("cog").Parse(genHeader)
	if err != nil {
		// not ideal, but also not a big deal: this statement is only reachable when the template is invalid.
		panic(err)
	}

	return func(f codejen.File) (codejen.File, error) {
		var leader string
		switch filepath.Ext(f.RelativePath) {
		case ".ts", ".go":
			leader = "//"
		case ".yml", ".yaml", ".py":
			leader = "#"
		default:
			leader = ""
		}

		if leader == "" {
			return f, nil
		}

		var from []codejen.NamedJenny
		if config.Debug {
			from = f.From
		}

		buf := new(bytes.Buffer)
		if err := tmpl.Execute(buf, map[string]any{
			"Using":  from,
			"Leader": leader,
		}); err != nil {
			return codejen.File{}, fmt.Errorf("failed executing GeneratedCommentHeader() template: %w", err)
		}
		buf.Write(f.Data)

		f.Data = buf.Bytes()

		return f, nil
	}
}

// PathPrefixer returns a FileMapper that injects the provided path prefix to files
// passed through it.
func PathPrefixer(prefix string) codejen.FileMapper {
	return func(f codejen.File) (codejen.File, error) {
		f.RelativePath = filepath.Join(prefix, f.RelativePath)
		return f, nil
	}
}
