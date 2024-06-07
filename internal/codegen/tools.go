package codegen

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/jennies/common"
	"github.com/grafana/cog/internal/languages"
	"github.com/grafana/cog/internal/tools"
)

func guessPackageFromFilename(filename string) string {
	pkg := filepath.Base(filepath.Dir(filename))
	if pkg != "." {
		return pkg
	}

	return strings.TrimSuffix(filepath.Base(filename), filepath.Ext(filename))
}

func dirExists(dir string) (bool, error) {
	stat, err := os.Stat(dir)
	//nolint:gocritic
	if os.IsNotExist(err) {
		return false, nil
	} else if err != nil {
		return false, err
	} else if !stat.IsDir() {
		return false, fmt.Errorf("'%s' is not a directory", dir)
	}

	return true, nil
}

func loadURL(ctx context.Context, url string) (io.ReadCloser, error) {
	client := http.DefaultClient

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("expecting 200 when loading '%s', got %d", url, resp.StatusCode)
	}

	return resp.Body, nil
}

func repositoryTemplatesJenny(pipeline *Pipeline) (*codejen.JennyList[languages.BuildOptions], error) {
	outputDir, err := pipeline.outputDir(pipeline.currentDirectory)
	if err != nil {
		return nil, err
	}

	repoTemplatesJenny := codejen.JennyListWithNamer[languages.BuildOptions](func(_ languages.BuildOptions) string {
		return "RepositoryTemplates"
	})
	repoTemplatesJenny.AppendOneToMany(&common.RepositoryTemplate{
		TemplateDir: pipeline.Output.RepositoryTemplates,
		ExtraData:   pipeline.Output.TemplatesData,
	})
	repoTemplatesJenny.AddPostprocessors(
		common.GeneratedCommentHeader(pipeline.jenniesConfig()),
		common.PathPrefixer(strings.ReplaceAll(outputDir, "%l", ".")),
	)

	return repoTemplatesJenny, nil
}

func packageTemplatesJenny(pipeline *Pipeline, language string) (*codejen.JennyList[languages.Context], error) {
	outputDir, err := pipeline.languageOutputDir(pipeline.currentDirectory, language)
	if err != nil {
		return nil, err
	}

	pkgTemplatesJenny := codejen.JennyListWithNamer[languages.Context](func(_ languages.Context) string {
		return "PackageTemplates" + tools.UpperCamelCase(language)
	})
	pkgTemplatesJenny.AppendOneToMany(&common.PackageTemplate{
		Language:    language,
		TemplateDir: pipeline.Output.PackageTemplates,
		ExtraData:   pipeline.Output.TemplatesData,
	})
	pkgTemplatesJenny.AddPostprocessors(common.PathPrefixer(outputDir))

	return pkgTemplatesJenny, nil
}

func runJenny[I any](jenny *codejen.JennyList[I], input I, destinationFS *codejen.FS) error {
	fs, err := jenny.GenerateFS(input)
	if err != nil {
		return err
	}

	return destinationFS.Merge(fs)
}
