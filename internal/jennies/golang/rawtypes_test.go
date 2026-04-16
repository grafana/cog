package golang

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
	"testing/fstest"

	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/jennies/common"
	"github.com/grafana/cog/internal/languages"
	"github.com/grafana/cog/internal/testutils"
	"github.com/stretchr/testify/require"
)

func TestRawTypes_Generate(t *testing.T) {
	test := testutils.GoldenFilesTestSuite[ast.Schema]{
		TestDataRoot: "../../../testdata/jennies/rawtypes",
		Name:         "GoRawTypes",
		Skip: map[string]string{
			"dashboard": "the dashboard test schema includes a composable slot, which rely on external input to be properly supported",
		},
	}

	config := Config{
		PackageRoot:                "github.com/grafana/cog/generated",
		GenerateEqual:              true,
		GenerateJSONMarshaller:     true,
		GenerateStrictUnmarshaller: true,
		GenerateValidate:           true,
	}
	jenny := RawTypes{
		config:          config,
		tmpl:            initTemplates(config, common.NewAPIReferenceCollector()),
		apiRefCollector: common.NewAPIReferenceCollector(),
	}
	compilerPasses := New(config).CompilerPasses()

	test.Run(t, func(tc *testutils.Test[ast.Schema]) {
		req := require.New(tc)

		// We run the compiler passes defined fo Go since without them, we
		// might not be able to translate some of the IR's semantics into Go.
		// Example: disjunctions.
		schema := tc.UnmarshalJSONInput(testutils.RawTypesIRInputFile)
		processedAsts, err := compilerPasses.Process(ast.Schemas{&schema})
		req.NoError(err)

		req.Len(processedAsts, 1, "we somehow got more ast.Schema than we put in")

		files, err := jenny.Generate(languages.Context{
			Schemas: processedAsts,
		})
		req.NoError(err)

		tc.WriteFiles(files)
	})
}

func TestRawTypes_Generate_WithUndiscriminatedDisjunctions(t *testing.T) {
	test := testutils.GoldenFilesTestSuite[ast.Schema]{
		TestDataRoot: "../../../testdata/jennies/rawtypes",
		Name:         "GoRawTypesUndiscriminated",
		Skip: map[string]string{
			"arrays":                           "not relevant for undiscriminated disjunctions",
			"constant_reference_as_default":    "not relevant for undiscriminated disjunctions",
			"constant_reference_discriminator": "not relevant for undiscriminated disjunctions",
			"constant_references":              "not relevant for undiscriminated disjunctions",
			"constraints":                      "not relevant for undiscriminated disjunctions",
			"dashboard":                        "the dashboard test schema includes a composable slot, which rely on external input to be properly supported",
			"disjunction_of_numbers":           "not relevant for undiscriminated disjunctions",
			"disjunctions":                     "not relevant for undiscriminated disjunctions",
			"disjunctions_of_scalars_and_refs": "not relevant for undiscriminated disjunctions",
			"enums":                            "not relevant for undiscriminated disjunctions",
			"field_with_struct_with_defaults":  "not relevant for undiscriminated disjunctions",
			"intersections":                    "not relevant for undiscriminated disjunctions",
			"maps":                             "not relevant for undiscriminated disjunctions",
			"nullable_fields":                  "not relevant for undiscriminated disjunctions",
			"package-with-dashes":              "not relevant for undiscriminated disjunctions",
			"refs":                             "not relevant for undiscriminated disjunctions",
			"scalars":                          "not relevant for undiscriminated disjunctions",
			"struct_with_complex_fields":       "not relevant for undiscriminated disjunctions",
			"struct_with_defaults":             "not relevant for undiscriminated disjunctions",
			"struct_with_optional_fields":      "not relevant for undiscriminated disjunctions",
			"struct_with_scalar_fields":        "not relevant for undiscriminated disjunctions",
			"time_hint":                        "not relevant for undiscriminated disjunctions",
			"variant_dataquery":                "not relevant for undiscriminated disjunctions",
			"variant_panelcfg_full":            "not relevant for undiscriminated disjunctions",
			"variant_panelcfg_only_options":    "not relevant for undiscriminated disjunctions",
		},
	}

	config := Config{
		PackageRoot:                         "github.com/grafana/cog/generated",
		GenerateJSONMarshaller:              true,
		GenerateStrictUnmarshaller:          true,
		GenerateValidate:                    true,
		GenerateUndiscriminatedDisjunctions: true,
	}
	jenny := RawTypes{
		config:          config,
		tmpl:            initTemplates(config, common.NewAPIReferenceCollector()),
		apiRefCollector: common.NewAPIReferenceCollector(),
	}
	compilerPasses := New(config).CompilerPasses()

	test.Run(t, func(tc *testutils.Test[ast.Schema]) {
		req := require.New(tc)

		schema := tc.UnmarshalJSONInput(testutils.RawTypesIRInputFile)
		processedAsts, err := compilerPasses.Process(ast.Schemas{&schema})
		req.NoError(err)

		req.Len(processedAsts, 1, "we somehow got more ast.Schema than we put in")

		files, err := jenny.Generate(languages.Context{
			Schemas: processedAsts,
		})
		req.NoError(err)

		tc.WriteFiles(files)
	})
}

func TestRawTypes_Generate_AllowMarshalEmptyDisjunctions(t *testing.T) {
	req := require.New(t)

	config := Config{
		PackageRoot:            "github.com/grafana/cog/generated",
		GenerateJSONMarshaller: true,
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
			ast.NewObject("tests", "DisjunctionOfScalars", ast.NewDisjunction(ast.Types{ast.String(), ast.Bool()})),
		),
	}
	schemas, err := compilerPasses.Process([]*ast.Schema{schema})
	req.NoError(err)

	context := languages.Context{
		Schemas: schemas,
	}

	result, err := jenny.generateSchema(context, schemas[0])
	req.NoError(err)

	req.NotContains(string(result), "no value for disjunction of scalars")
}

func TestRawTypes_Generate_CustomObjectMethod(t *testing.T) {
	req := require.New(t)

	widgetMarker := "func-Widget"
	gadgetMarker := "func-Gadget"
	templateContent := `{{ define "object_all_custom_methods" }}
func (resource {{ .Object.Name }}) CustomMethod() string {
	return "{{ label .Object.Name }}"
}
{{ end }}`

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
	runTest := func(config Config) {
		jenny := RawTypes{
			config:          config,
			tmpl:            initTemplates(config, common.NewAPIReferenceCollector()),
			apiRefCollector: common.NewAPIReferenceCollector(),
		}
		compilerPasses := New(config).CompilerPasses()

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

	t.Run("fs", func(t *testing.T) {
		config := Config{
			PackageRoot: "github.com/grafana/cog/generated",
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
			PackageRoot:                   "github.com/grafana/cog/generated",
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
