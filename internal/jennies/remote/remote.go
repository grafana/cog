package remote

import (
	"fmt"

	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/tools"
	"github.com/grafana/cog/pkg/languages"
	"github.com/grafana/cog/pkg/plugins"
)

type remote struct {
	globalConfig languages.Config
	config       map[string]any
	plugin       plugins.Language
}

func (jenny remote) JennyName() string {
	return "Remote"
}

func (jenny remote) Generate(context languages.Context) (codejen.Files, error) {
	files, err := jenny.plugin.Generate(jenny.globalConfig, jenny.config, context)
	if err != nil {
		return nil, fmt.Errorf("error while generating context: %w", err)
	}

	return tools.Map(files, func(file codejen.File) codejen.File {
		file.From = append(file.From, jenny)
		return file
	}), nil
}
