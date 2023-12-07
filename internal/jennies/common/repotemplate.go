package common

import (
	"bytes"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/grafana/codejen"
)

type RepositoryTemplate struct {
	TemplateDir string
	ExtraData   map[string]string
}

func (jenny RepositoryTemplate) JennyName() string {
	return "RepositoryTemplate"
}

func (jenny RepositoryTemplate) Generate(buildOpts BuildOptions) (codejen.Files, error) {
	var files codejen.Files

	for _, language := range buildOpts.Languages {
		templateRoot := filepath.Join(jenny.TemplateDir, language)
		cleanedRoot := filepath.Clean(templateRoot) + string(filepath.Separator)

		err := filepath.WalkDir(templateRoot, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if d.IsDir() {
				return nil
			}

			templateContent, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			tmpl, err := template.New(jenny.JennyName()).
				Option("missingkey=error").
				Funcs(sprig.FuncMap()).
				Parse(string(templateContent))
			if err != nil {
				return err
			}

			buf := new(bytes.Buffer)
			if err := tmpl.Execute(buf, jenny.templateData()); err != nil {
				return err
			}

			if buf.Len() != 0 {
				files = append(files, *codejen.NewFile(strings.TrimPrefix(path, cleanedRoot), buf.Bytes(), jenny))
			}

			return nil
		})
		if err != nil && !errors.Is(err, os.ErrNotExist) {
			return nil, err
		}
	}

	return files, nil
}

func (jenny RepositoryTemplate) templateData() map[string]any {
	extra := map[string]string{}
	if jenny.ExtraData != nil {
		extra = jenny.ExtraData
	}

	return map[string]any{
		"Extra": extra,
	}
}
