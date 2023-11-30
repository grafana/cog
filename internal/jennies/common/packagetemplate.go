package common

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/grafana/codejen"
)

type PackageTemplate struct {
	Language    string
	TemplateDir string
}

func (jenny PackageTemplate) JennyName() string {
	return fmt.Sprintf("PackageTemplate[%s]", jenny.Language)
}

func (jenny PackageTemplate) Generate(context Context) (codejen.Files, error) {
	var files codejen.Files

	templateRoot := filepath.Join(jenny.TemplateDir, jenny.Language)
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
			Funcs(sprig.FuncMap()).
			Parse(string(templateContent))
		if err != nil {
			return err
		}

		buf := new(bytes.Buffer)
		if err := tmpl.Execute(buf, jenny.templateData(context)); err != nil {
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

	return files, nil
}

func (jenny PackageTemplate) templateData(context Context) map[string]any {
	packages := make([]string, 0, len(context.Schemas))
	for _, schema := range context.Schemas {
		packages = append(packages, schema.Package)
	}

	return map[string]any{
		"Packages": packages,
	}
}
