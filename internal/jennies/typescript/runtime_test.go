package typescript

import (
	"testing"

	"github.com/grafana/cog/internal/languages"
	"github.com/stretchr/testify/require"
)

func TestRuntime(t *testing.T) {
	req := require.New(t)
	jenny := Runtime{
		tmpl:    initTemplates(nil),
		targets: languages.Config{Converters: true},
	}

	files, err := jenny.Generate(languages.Context{})
	req.NoError(err)

	req.Len(files, 4)
}
