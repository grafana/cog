package python

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"testing/fstest"
	"time"

	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/jennies/common"
	"github.com/grafana/cog/internal/languages"
	"github.com/grafana/cog/internal/testutils"
	"github.com/stretchr/testify/require"
)

func TestRawTypes_Generate(t *testing.T) {
	test := testutils.GoldenFilesTestSuite[ast.Schema]{
		TestDataRoot: "../../../testdata/jennies/rawtypes",
		Name:         "PythonRawTypes",
		Skip: map[string]string{
			"intersections": "Intersections are not implemented",
			"dashboard":     "the dashboard test schema includes a composable slot, which rely on external input to be properly supported",
		},
	}

	config := Config{
		GenerateJSONMarshaller: true,
	}
	jenny := RawTypes{
		config:          config,
		tmpl:            initTemplates(config, common.NewAPIReferenceCollector()),
		apiRefCollector: common.NewAPIReferenceCollector(),
	}
	compilerPasses := New(config).CompilerPasses()

	test.Run(t, func(tc *testutils.Test[ast.Schema]) {
		req := require.New(tc)

		// We run the compiler passes defined fo Python since without them, we
		// might not be able to translate some of the IR's semantics into Python.
		// Example: anonymous objects.
		schema := tc.UnmarshalJSONInput(testutils.RawTypesIRInputFile)
		processedAsts, err := compilerPasses.Process(ast.Schemas{&schema})
		req.NoError(err)

		req.Len(processedAsts, 1, "we somehow got more ast.Schema than we put in")

		files, err := jenny.Generate(languages.Context{Schemas: processedAsts})
		req.NoError(err)

		tc.WriteFiles(files)
	})
}

func TestRawTypes_Generate_CustomObjectMethod(t *testing.T) {
	req := require.New(t)

	customData := fmt.Sprintf("data-%d", time.Now().UnixNano())
	widgetMarker := "func-Widget-" + customData
	gadgetMarker := "func-Gadget-" + customData

	config := Config{
		OverridesTemplatesFS: fstest.MapFS{
			"custom/methods.tmpl": {
				Data: []byte(`{{ define "object_custom_methods_all_python" }}
def custom_method(self) -> str:
    return "{{ label .Object.Name }}-{{ .CustomData }}"
{{ end }}`),
			},
		},
		OverridesTemplatesData: map[string]any{
			"CustomData": customData,
		},
		OverridesTemplateFuncs: map[string]any{
			"label": func(s string) string {
				return "func-" + s
			},
		},
	}

	jenny := RawTypes{
		config:          config,
		tmpl:            initTemplates(config, common.NewAPIReferenceCollector()),
		apiRefCollector: common.NewAPIReferenceCollector(),
	}
	compilerPasses := New(config).CompilerPasses()

	schema := &ast.Schema{
		Package: "tests",
		Objects: testutils.ObjectsMap(
			ast.NewObject("tests", "Widget", ast.NewStruct(
				ast.NewStructField("name", ast.NewScalar(ast.KindString), ast.Required()),
			)),
			ast.NewObject("tests", "Gadget", ast.NewStruct(
				ast.NewStructField("id", ast.NewScalar(ast.KindString), ast.Required()),
			)),
		),
	}
	schemas, err := compilerPasses.Process(ast.Schemas{schema})
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

func TestRawTypes_Generate_CustomObjectMethod_WithDirectory(t *testing.T) {
	req := require.New(t)

	const marker = "directory-based-method"

	tmpDir := t.TempDir()
	customDir := filepath.Join(tmpDir, "custom")
	err := os.MkdirAll(customDir, 0o755)
	req.NoError(err)

	templateContent := []byte(`{{ define "object_custom_methods_all_python" }}
def directory_method(self) -> str:
    return "` + marker + `"
{{ end }}`)
	err = os.WriteFile(filepath.Join(customDir, "methods.tmpl"), templateContent, 0o644)
	req.NoError(err)

	config := Config{
		OverridesTemplatesDirectories: []string{tmpDir},
	}

	jenny := RawTypes{
		config:          config,
		tmpl:            initTemplates(config, common.NewAPIReferenceCollector()),
		apiRefCollector: common.NewAPIReferenceCollector(),
	}
	compilerPasses := New(config).CompilerPasses()

	schema := &ast.Schema{
		Package: "tests",
		Objects: testutils.ObjectsMap(
			ast.NewObject("tests", "Widget", ast.NewStruct(
				ast.NewStructField("name", ast.NewScalar(ast.KindString), ast.Required()),
			)),
		),
	}
	schemas, err := compilerPasses.Process(ast.Schemas{schema})
	req.NoError(err)

	files, err := jenny.Generate(languages.Context{Schemas: schemas})
	req.NoError(err)

	found := false
	for _, file := range files {
		if bytes.Contains(file.Data, []byte(marker)) {
			found = true
			break
		}
	}

	req.True(found, "expected generated output to include method from directory template")
}
