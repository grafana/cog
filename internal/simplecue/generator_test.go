package simplecue

import (
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/load"
	"fmt"
	"github.com/stretchr/testify/require"
	"io"
	"io/fs"
	"path/filepath"
	"testing"
)

// ToOverlay converts a fs.FS into a CUE loader overlay.
func toCueOverlay(prefix string, vfs fs.FS, overlay map[string]load.Source) error {
	// TODO why not just stick the prefix on automatically...?
	if !filepath.IsAbs(prefix) {
		return fmt.Errorf("must provide absolute path prefix when generating cue overlay, got %q", prefix)
	}
	err := fs.WalkDir(vfs, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		f, err := vfs.Open(path)
		if err != nil {
			return err
		}
		defer func() { _ = f.Close() }()

		b, err := io.ReadAll(f)
		if err != nil {
			return err
		}

		overlay[filepath.Join(prefix, path)] = load.FromBytes(b)

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func TestTypes_TestCueGenerator(t *testing.T) {
	//libFS := &merged_fs.EmptyFS{}
	//overlay := make(map[string]load.Source)
	//err := toCueOverlay("/", libFS, overlay)
	//require.NoError(t, err)

	bis := load.Instances([]string{"testdata/dashboard.cue"}, &load.Config{
		//Overlay:    overlay,
		ModuleRoot: "/",
	})
	values, err := cuecontext.New().BuildInstances(bis)
	require.NoError(t, err)

	schemaAst, err := GenerateAST(values[0], Config{Package: "grafanatest"})
	require.NoError(t, err)
	require.NotNil(t, schemaAst)
	//
	//testCases := []struct {
	//	description string
	//	expected    bool
	//}{
	//	{
	//		description: "foo",
	//		expected:    false,
	//	},
	//}
	//
	//for _, testCase := range testCases {
	//	tc := testCase
	//
	//	t.Run(tc.description, func(t *testing.T) {
	//		req := require.New(t)
	//
	//		req.Equal(tc.expected, false)
	//	})
	//}
}
