package typescript

import (
	"testing"

	"github.com/grafana/cog/internal/ast"
	"github.com/stretchr/testify/require"
)

func TestOptionsBuilder(t *testing.T) {
	req := require.New(t)
	jenny := OptionsBuilder{}

	file, err := jenny.Generate([]*ast.Schema{})
	req.NoError(err)

	req.Equal("options_builder_gen.ts", file.RelativePath)
	req.NotEmpty(file.From)
	req.Equal(`export interface CogOptionsBuilder<T> {
  build: () => T;
}
`, string(file.Data))
}
