package typescript

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/jennies/common"
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
			packages[schema.Package] = []string{"types_gen"}
		}
	}

	if jenny.Targets.Builders {
		for _, builder := range context.Builders {
			packages[builder.Package] = append(packages[builder.Package], fmt.Sprintf("%s_builder_gen", strings.ToLower(builder.Name)))
		}
	}

	for pkg, refs := range packages {
		files = append(files, *codejen.NewFile(filepath.Join("src", pkg, "index.ts"), jenny.generateIndex(refs), jenny))
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
