package typescript

import (
	"testing"

	"github.com/grafana/cog/internal/jennies/context"
	"github.com/stretchr/testify/require"
)

func TestRuntime(t *testing.T) {
	req := require.New(t)
	jenny := Runtime{}

	files, err := jenny.Generate(context.Builders{})
	req.NoError(err)

	req.Len(files, 3)
}
