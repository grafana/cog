package typescript

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/jennies/common"
	"github.com/grafana/cog/internal/tools"
)

type Index struct {
	Targets common.Config
}

func (jenny Index) JennyName() string {
	return "TypescriptIndex"
}

func (jenny Index) Generate(context common.Context) (codejen.Files, error) {
	packages := make(map[string][]string, len(context.Schemas))
	files := codejen.Files{}

	if jenny.Targets.Types {
		for _, schema := range context.Schemas {
			packages[schema.Package] = []string{"types.gen"}
		}
	}

	if jenny.Targets.Builders {
		for _, builder := range context.Builders {
			packages[builder.Package] = append(packages[builder.Package], fmt.Sprintf("%sBuilder.gen", tools.LowerCamelCase(builder.Name)))
		}
	}

	for pkg, refs := range packages {
		files = append(files, *codejen.NewFile(filepath.Join("src", formatPackageName(pkg), "index.ts"), jenny.generateIndex(refs), jenny))
	}

	return files, nil
}

func (jenny Index) generateIndex(refs []string) []byte {
	output := strings.Builder{}

	for _, ref := range refs {
		output.WriteString(fmt.Sprintf("export * from './%s';\n", ref))
		output.WriteString(fmt.Sprintf("export type * from './%s';\n", ref))
	}

	return []byte(output.String())
}
