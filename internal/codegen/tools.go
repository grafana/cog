package codegen

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/expr-lang/expr"
	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/jennies/common"
	"github.com/grafana/cog/internal/semver"
)

func guessPackageFromFilename(filename string) string {
	pkg := filepath.Base(filepath.Dir(filename))
	if pkg != "." {
		return pkg
	}

	return strings.TrimSuffix(filepath.Base(filename), filepath.Ext(filename))
}

func isFile(path string) (bool, error) {
	stat, err := os.Stat(path)
	//nolint:gocritic
	if err != nil {
		return false, err
	} else if stat.IsDir() {
		return false, nil
	}

	return true, nil
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

	// Necessary for private repositories
	authToken := os.Getenv("AUTH_TOKEN")
	if authToken != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", authToken))
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

func repositoryTemplatesJenny(pipeline *Pipeline) (*codejen.JennyList[common.BuildOptions], error) {
	outputDir, err := pipeline.outputDir(pipeline.currentDirectory)
	if err != nil {
		return nil, err
	}

	repoTemplatesJenny := codejen.JennyListWithNamer[common.BuildOptions](func(_ common.BuildOptions) string {
		return "RepositoryTemplates"
	})
	repoTemplatesJenny.AppendOneToMany(&common.RepositoryTemplate{
		TemplateDir:       pipeline.Output.RepositoryTemplates,
		ExtraData:         pipeline.Output.TemplatesData,
		ReplaceExtensions: pipeline.Output.OutputOptions.ReplaceExtension,
	})
	repoTemplatesJenny.AddPostprocessors(
		common.GeneratedCommentHeader(pipeline.jenniesConfig()),
		common.PathPrefixer(strings.ReplaceAll(outputDir, "%l", ".")),
	)

	return repoTemplatesJenny, nil
}

func runJenny[I any](jenny *codejen.JennyList[I], input I, destinationFS *codejen.FS) error {
	fs, err := jenny.GenerateFS(input)
	if err != nil {
		return err
	}

	return destinationFS.Merge(fs)
}

func evalExpr(in string) (bool, error) {
	if in == "" {
		return true, nil
	}

	env := map[string]any{
		"sprintf": fmt.Sprintf,
		"semver":  semver.ParseTolerant,
	}

	program, err := expr.Compile(in, expr.Env(env))
	if err != nil {
		return false, err
	}

	output, err := expr.Run(program, env)
	if err != nil {
		return false, err
	}

	if _, ok := output.(bool); !ok {
		return false, fmt.Errorf("expected expression to evaluate to a boolean, got %T", output)
	}

	return output.(bool), nil
}
