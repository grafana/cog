package common

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/jennies/template"
	"github.com/grafana/cog/internal/languages"
)

type PackageTemplate struct {
	Language    string
	TemplateDir string
	ExtraData   map[string]string
}

func (jenny PackageTemplate) JennyName() string {
	return fmt.Sprintf("PackageTemplate[%s]", jenny.Language)
}

func (jenny PackageTemplate) Generate(context languages.Context) (codejen.Files, error) {
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

		tmpl, err := template.New(
			jenny.JennyName(),
			template.Funcs(template.FuncMap{
				"registryToSemver": jenny.registryToSemver,
			}),
			template.Parse(string(templateContent)),
		)
		if err != nil {
			return err
		}

		rendered, err := tmpl.ExecuteAsBytes(jenny.templateData(context))
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

	return files, nil
}

func (jenny PackageTemplate) templateData(context languages.Context) map[string]any {
	packages := make([]string, 0, len(context.Schemas))
	for _, schema := range context.Schemas {
		packages = append(packages, schema.Package)
	}

	sort.Strings(packages)

	extra := map[string]string{}
	if jenny.ExtraData != nil {
		extra = jenny.ExtraData
	}

	return map[string]any{
		"Packages": packages,
		"Extra":    extra,
	}
}

// registryToSemver turns a "v10.2.x" input (version string coming from
// the kind-registry) into a semver-compatible version "10.2.0"
func (jenny PackageTemplate) registryToSemver(registryVersion string) string {
	semver := strings.TrimPrefix(registryVersion, "v")

	if strings.HasSuffix(semver, "x") {
		semver = semver[:len(semver)-1] + "0"
	}

	return semver
}
