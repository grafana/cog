package java

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
	"testing/fstest"

	"github.com/grafana/cog/internal/jennies/common"
	"github.com/grafana/cog/internal/testutils"
	"github.com/grafana/cog/pkg/ir"
	"github.com/grafana/cog/pkg/languages"
	"github.com/grafana/cog/pkg/logs"
	"github.com/stretchr/testify/require"
)

func TestRawTypes_Generate(t *testing.T) {
	test := testutils.GoldenFilesTestSuite[ir.Schema]{
		TestDataRoot: "../../../testdata/jennies/rawtypes",
		Name:         "JavaRawTypes",
	}

	cfg := Config{
		GenerateEqual:          true,
		GenerateJSONMarshaller: true,
	}

	jenny := RawTypes{config: cfg, tmpl: initTemplates(cfg, common.NewAPIReferenceCollector())}
	transforms := New(logs.NoopLogger(), cfg).Transform

	test.Run(t, func(tc *testutils.Test[ir.Schema]) {
		req := require.New(tc)

		// We run the compiler passes defined fo Java since without them, we
		// might not be able to translate some of the IR's semantics into Java.
		// Example: disjunctions.
		schema := tc.UnmarshalJSONInput(testutils.RawTypesIRInputFile)
		processedAsts, err := transforms(ir.Schemas{&schema})
		req.NoError(err)

		req.Len(processedAsts, 1, "we somehow got more ast.Schema than we put in")

		files, err := jenny.Generate(languages.Context{
			Schemas: processedAsts,
		})
		req.NoError(err)

		tc.WriteFiles(files)
	})
}

func TestRawTypes_Generate_CustomObjectMethod(t *testing.T) {
	req := require.New(t)

	widgetMarker := "func-Widget"
	gadgetMarker := "func-Gadget"
	templateContent := `{{ define "object_all_custom_methods" }}
public String customMethod() {
	return "{{ label .Object.Name }}";
}
{{ end }}`

	schema := &ir.Schema{
		Package: "tests",
		Objects: testutils.ObjectsMap(
			ir.NewObject("tests", "Widget", ir.NewStruct(
				ir.NewStructField("name", ir.NewScalar(ir.KindString), ir.Required()),
			)),
			ir.NewObject("tests", "Gadget", ir.NewStruct(
				ir.NewStructField("id", ir.NewScalar(ir.KindString), ir.Required()),
			)),
		),
	}
	runTest := func(config Config) {
		jenny := RawTypes{config: config, tmpl: initTemplates(config, common.NewAPIReferenceCollector())}
		transforms := New(logs.NoopLogger(), config).Transform

		schemas, err := transforms(ir.Schemas{schema})
		req.NoError(err)

		files, err := jenny.Generate(languages.Context{Schemas: schemas})
		req.NoError(err)

		foundWidget := false
		foundGadget := false
		for _, file := range files {
			if bytes.Contains(file.Data, []byte(widgetMarker)) {
				foundWidget = true
			}
			if bytes.Contains(file.Data, []byte(gadgetMarker)) {
				foundGadget = true
			}
		}

		req.True(foundWidget, "expected generated output to include Widget custom method")
		req.True(foundGadget, "expected generated output to include Gadget custom method")
	}

	t.Run("fs", func(t *testing.T) {
		config := Config{
			OverridesTemplatesFS: fstest.MapFS{
				"custom/methods.tmpl": {
					Data: []byte(templateContent),
				},
			},
			OverridesTemplateFuncs: map[string]any{
				"label": func(s string) string {
					return "func-" + s
				},
			},
		}

		runTest(config)
	})

	t.Run("directory", func(t *testing.T) {
		tmpDir := t.TempDir()
		customDir := filepath.Join(tmpDir, "custom")
		err := os.MkdirAll(customDir, 0o755)
		req.NoError(err)

		err = os.WriteFile(filepath.Join(customDir, "methods.tmpl"), []byte(templateContent), 0o600)
		req.NoError(err)

		config := Config{
			OverridesTemplatesDirectories: []string{tmpDir},
			OverridesTemplateFuncs: map[string]any{
				"label": func(s string) string {
					return "func-" + s
				},
			},
		}

		runTest(config)
	})
}
