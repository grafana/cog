package golang

import (
	"testing"

	"github.com/grafana/cog/internal/languages"
	"github.com/stretchr/testify/require"
)

func TestGoMod_Generate(t *testing.T) {
	req := require.New(t)

	jenny := GoMod{
		Config: Config{PackageRoot: "github.com/grafana/heey"},
	}

	files, err := jenny.Generate(languages.Context{})
	req.NoError(err)

	req.Len(files, 1)

	goModFile := files[0]

	req.Equal("go.mod", goModFile.RelativePath)
	req.Equal(`module github.com/grafana/heey

go 1.21

`, string(goModFile.Data))
}
