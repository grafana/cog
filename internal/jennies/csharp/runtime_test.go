package csharp

import (
	"strings"
	"testing"

	"github.com/grafana/cog/internal/jennies/common"
	"github.com/grafana/cog/internal/languages"
	"github.com/stretchr/testify/require"
)

// TestRuntime_Generate sanity-checks that the runtime jenny emits the
// IBuilder<T> interface at the expected location with the expected
// namespace. Goldens aren't used because the runtime is a single,
// trivial file — a smoke test is enough.
func TestRuntime_Generate(t *testing.T) {
	cfg := Config{GenerateBuilders: true}
	cfg.InterpolateParameters(func(s string) string { return s })

	jenny := Runtime{
		config: cfg,
		tmpl:   initTemplates(cfg, common.NewAPIReferenceCollector()),
	}
	files, err := jenny.Generate(languages.Context{})
	require.NoError(t, err)
	require.Len(t, files, 1)

	require.Equal(t, "src/Grafana/Foundation/Cog/IBuilder.cs", files[0].RelativePath)
	body := string(files[0].Data)
	require.True(t, strings.Contains(body, "namespace Grafana.Foundation.Cog;"), body)
	require.True(t, strings.Contains(body, "public interface IBuilder<out T>"), body)
	require.True(t, strings.Contains(body, "T Build();"), body)
}
