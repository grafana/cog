package codegen

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/jennies/common"
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

func repositoryTemplatesJenny(pipeline *Pipeline) (*codejen.JennyList[common.BuildOptions], error) {
	outputDir, err := pipeline.outputDir(pipeline.currentDirectory)
	if err != nil {
		return nil, err
	}

	repoTemplatesJenny := codejen.JennyListWithNamer(func(_ common.BuildOptions) string {
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

func createInterpolator(parameters map[string]string) ParametersInterpolator {
	interpolateFun := func(in string) string {
		interpolated := in

		for key, value := range parameters {
			interpolated = strings.ReplaceAll(interpolated, "%"+key+"%", value)
		}

		return interpolated
	}

	return func(input string) string {
		finalOutput := input
		out := interpolateFun(finalOutput)
		for out != finalOutput {
			finalOutput = out
			out = interpolateFun(finalOutput)
		}

		return finalOutput
	}
}
