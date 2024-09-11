package common

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/jennies/template"
)

type BuildOptions struct {
	Languages []string
}

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

			tmpl, err := template.New(
				jenny.JennyName(),
				template.Parse(string(templateContent)),
			)
			if err != nil {
				return err
			}

			rendered, err := tmpl.ExecuteAsBytes(jenny.templateData())
			if err != nil {
				return err
			}

			if len(rendered) != 0 {
				files = append(files, *codejen.NewFile(strings.TrimPrefix(path, cleanedRoot), rendered, jenny))
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
