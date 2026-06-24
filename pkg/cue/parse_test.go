package cue_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	cuelang "cuelang.org/go/cue"
	"github.com/grafana/cog/pkg/cue"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	schema := `
package sandbox

// Contains things.
Container: {
    str: string
}
`
	req := require.New(t)

	dir := t.TempDir()
	schemaFile := filepath.Join(dir, "schema.cue")

	err := os.WriteFile(schemaFile, []byte(schema), 0600)
	req.NoError(err)

	result, err := cue.Parse(context.TODO(), cue.Input{
		Entrypoint: dir,
	})
	req.NoError(err)

	containerVal := result.Value.LookupPath(cuelang.MakePath(cuelang.Str("Container")))
	req.True(containerVal.Exists())
	req.Equal("sandbox", result.Package)
}
