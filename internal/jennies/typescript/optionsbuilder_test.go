package typescript

import (
	"testing"

	"github.com/grafana/cog/internal/jennies/context"
	"github.com/stretchr/testify/require"
)

func TestOptionsBuilder(t *testing.T) {
	req := require.New(t)
	jenny := OptionsBuilder{}

	file, err := jenny.Generate(context.Builders{})
	req.NoError(err)

	req.Equal("options_builder_gen.ts", file.RelativePath)
	req.NotEmpty(file.From)
	req.Equal(`export interface CogOptionsBuilder<T> {
  build: () => T;
}
`, string(file.Data))
}
