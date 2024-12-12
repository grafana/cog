package common

import (
	"errors"
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
	TmplFuncs           template.FuncMap
	TemplateDirectories []string
	Data                map[string]any
	ExtraData           map[string]string
}

func (jenny PackageTemplate) JennyName() string {
	return "PackageTemplate"
}

func (jenny PackageTemplate) Generate(context languages.Context) (codejen.Files, error) {
	var allFiles codejen.Files

	for _, templatesDir := range jenny.TemplateDirectories {
		files, err := jenny.generateForTemplatesDirectory(context, templatesDir)
		if err != nil {
			return nil, err
		}

		allFiles = append(allFiles, files...)
	}

	return allFiles, nil
}

func (jenny PackageTemplate) generateForTemplatesDirectory(context languages.Context, directory string) (codejen.Files, error) {
	var files codejen.Files

	templateRoot := directory
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
			path,
			template.Funcs(TypeResolvingTemplateHelpers(context)),
			template.Funcs(template.FuncMap{
				"registryToSemver": jenny.registryToSemver,
			}),
			template.Funcs(jenny.TmplFuncs),
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

	return map[string]any{
		"Context":  context,
		"Packages": packages,
		"Data":     jenny.Data,
		"Extra":    jenny.ExtraData,
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
