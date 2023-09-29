package simplecue

import (
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"path/filepath"
	"testing"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/load"
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/txtartest"
	"github.com/stretchr/testify/require"
	"github.com/yalue/merged_fs"
)

func TestGenerateAST(t *testing.T) {
	test := txtartest.TxTarTest{
		Root: "../../testdata/simplecue",
		Name: "simplecue/GenerateAST",
	}

	test.Run(t, func(tc *txtartest.Test) {
		req := require.New(tc)

		schemaAst, err := GenerateAST(txtarTestToCueInstance(tc), Config{Package: "grafanatest"})
		req.NoError(err)
		require.NotNil(t, schemaAst)

		writeIR(schemaAst, tc)
	})
}

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

func writeIR(irFile *ast.File, tc *txtartest.Test) {
	tc.Helper()

	// TODO: use FileWriter instead in separate CL.
	_, err := fmt.Fprintln(tc, "==", "ir.json")
	require.NoError(tc, err)

	marshaledIR, err := json.MarshalIndent(irFile, "", "  ")
	require.NoError(tc, err)

	_, err = tc.Write(marshaledIR)
	require.NoError(tc, err)
}

func txtarTestToCueInstance(tc *txtartest.Test) cue.Value {
	tc.Helper()

	for _, f := range tc.Archive.Files {
		if f.Name != "schema.cue" {
			continue
		}

		return bytesToCueValue(tc.T, f.Data)
	}

	// the ir.json file wasn't found, let's fail hard.
	tc.Fatal("could not load types IR: file 'ir.json' not found in test archive")

	return cue.Value{}
}

func bytesToCueValue(t *testing.T, input []byte) cue.Value {
	t.Helper()

	req := require.New(t)

	libFS := &merged_fs.EmptyFS{}
	overlay := make(map[string]load.Source)
	err := toCueOverlay("/", libFS, overlay)
	req.NoError(err)

	someSource := load.FromBytes(input)

	overlay["/schema.cue"] = someSource

	bis := load.Instances([]string{"/schema.cue"}, &load.Config{
		Overlay:    overlay,
		ModuleRoot: "/",
	})
	values, err := cuecontext.New().BuildInstances(bis)
	req.NoError(err)

	return values[0]
}
