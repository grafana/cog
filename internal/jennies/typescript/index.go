package typescript

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/jennies/context"
)

type Index struct {
}

func (jenny Index) JennyName() string {
	return "TypescriptIndex"
}

func (jenny Index) Generate(context context.Builders) (codejen.Files, error) {
	packages := make(map[string][]string, len(context.Schemas))
	files := codejen.Files{}

	for _, schema := range context.Schemas {
		packages[schema.Package] = []string{"types_gen"}
	}
	for _, builder := range context.Builders {
		packages[builder.Package] = append(packages[builder.Package], fmt.Sprintf("%s_builder_gen", strings.ToLower(builder.For.Name)))
	}

	for pkg, refs := range packages {
		files = append(files, *codejen.NewFile(filepath.Join(pkg, "index.ts"), jenny.generateIndex(refs), jenny))
	}

	return files, nil
}

func (jenny Index) generateIndex(refs []string) []byte {
	output := strings.Builder{}

	for _, ref := range refs {
		output.WriteString(fmt.Sprintf("export * from './%s';\n", ref))
	}

	return []byte(output.String())
}
