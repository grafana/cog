package simplecue

import (
	"fmt"
	"io"
	"io/fs"
	"path/filepath"
	"testing"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/load"
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/testutils"
	"github.com/stretchr/testify/require"
	"github.com/yalue/merged_fs"
)

func TestGenerateAST(t *testing.T) {
	test := testutils.GoldenFilesTestSuite[string]{
		TestDataRoot: "../../testdata/simplecue",
		Name:         "GenerateAST",
	}

	test.Run(t, func(tc *testutils.Test[string]) {
		req := require.New(tc)

		schemaAst, err := GenerateAST(txtarTestToCueInstance(tc), Config{Package: "grafanatest"})
		req.NoError(err)
		require.NotNil(t, schemaAst)

		tc.WriteJSON(testutils.GeneratorOutputFile, schemaAst)
	})
}

func TestGenerateAST_withOutOfRootReference(t *testing.T) {
	req := require.New(t)
	schema := `
schema: {
  #Origin: { creator: string }
  spec: {
    title: string
    origin: #Origin
  }
}
`

	cueVal := cuecontext.New().CompileString(schema)
	specCueVal := cueVal.LookupPath(cue.ParsePath("schema.spec"))

	schemaAst, err := GenerateAST(specCueVal, Config{Package: "grafanatest", ForceNamedEnvelope: "spec"})
	req.NoError(err)
	require.NotNil(t, schemaAst)

	objects := []ast.Object{
		ast.NewObject("grafanatest", "Origin", ast.NewStruct(
			ast.NewStructField("creator", ast.String(), ast.Required()),
		)),
		ast.NewObject("grafanatest", "spec", ast.NewStruct(
			ast.NewStructField("title", ast.String(), ast.Required()),
			ast.NewStructField("origin", ast.NewRef("grafanatest", "Origin"), ast.Required()),
		)),
	}

	req.Equal(testutils.ObjectsMap(objects...), schemaAst.Objects)
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

func txtarTestToCueInstance(tc *testutils.Test[string]) cue.Value {
	tc.Helper()

	return bytesToCueValue(tc.T, tc.ReadInput("schema.cue"))
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
