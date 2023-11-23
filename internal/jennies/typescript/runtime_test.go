package typescript

import (
	"testing"

	"github.com/grafana/cog/internal/jennies/common"
	"github.com/stretchr/testify/require"
)

func TestRuntime(t *testing.T) {
	req := require.New(t)
	jenny := Runtime{}

	files, err := jenny.Generate(common.Context{})
	req.NoError(err)

	req.Len(files, 3)
}
