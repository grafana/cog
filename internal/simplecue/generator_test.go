package simplecue

import (
	"strings"
	"testing"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/load"
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/testutils"
	"github.com/stretchr/testify/require"
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

func TestGenerateAST_withPackageOverride(t *testing.T) {
	req := require.New(t)
	schema := `
package foo

#Ref: string
Container: {
  ref: #Ref
}
`

	cueVal := cuecontext.New().CompileString(schema)

	schemaAst, err := GenerateAST(cueVal, Config{Package: "grafanatest"})
	req.NoError(err)
	require.NotNil(t, schemaAst)

	objects := []ast.Object{
		ast.NewObject("grafanatest", "Ref", ast.String()),
		ast.NewObject("grafanatest", "Container", ast.NewStruct(
			ast.NewStructField("ref", ast.NewRef("grafanatest", "Ref"), ast.Required()),
		)),
	}

	req.Equal(testutils.ObjectsMap(objects...), schemaAst.Objects)
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

func TestGenerateAST_withCustomNameFunc(t *testing.T) {
	req := require.New(t)
	schema := `
schema: {
  #Origin: { creator: string }
  spec: {
    title: string
    origin: #Origin
    details: #Details
    #Details: {
      [string]: _
    }
  }
}
`

	nameFunc := func(_ cue.Value, path cue.Path) string {
		return strings.Trim(path.String(), "?#")
	}

	cueVal := cuecontext.New().CompileString(schema)
	specCueVal := cueVal.LookupPath(cue.ParsePath("schema.spec"))

	schemaAst, err := GenerateAST(specCueVal, Config{Package: "grafanatest", ForceNamedEnvelope: "spec", NameFunc: nameFunc})
	req.NoError(err)
	require.NotNil(t, schemaAst)

	objects := []ast.Object{
		ast.NewObject("grafanatest", "schema.#Origin", ast.NewStruct(
			ast.NewStructField("creator", ast.String(), ast.Required()),
		)),
		ast.NewObject("grafanatest", "schema.spec.#Details", ast.NewMap(
			ast.String(),
			ast.Any(),
		)),
		ast.NewObject("grafanatest", "spec", ast.NewStruct(
			ast.NewStructField("title", ast.String(), ast.Required()),
			ast.NewStructField("origin", ast.NewRef("grafanatest", "schema.#Origin"), ast.Required()),
			ast.NewStructField("details", ast.NewRef("grafanatest", "schema.spec.#Details"), ast.Required()),
		)),
	}

	req.Equal(testutils.ObjectsMap(objects...), schemaAst.Objects)
}

func TestGenerateAST_withEnvelopeAndConstantRef(t *testing.T) {
	req := require.New(t)
	schema := `
Spec: {
	type: ValueMap
}

MappingType: "value" | "range"

ValueMap: {
	type: MappingType & "value"
}
`

	cueVal := cuecontext.New().CompileString(schema)
	specCueVal := cueVal.LookupPath(cue.ParsePath("Spec"))

	schemaAst, err := GenerateAST(specCueVal, Config{Package: "grafanatest", ForceNamedEnvelope: "Spec"})
	req.NoError(err)
	require.NotNil(t, schemaAst)

	objects := []ast.Object{
		ast.NewObject("grafanatest", "ValueMap", ast.NewStruct(
			ast.NewStructField("type", ast.NewConstantReferenceType("grafanatest", "MappingType", "value"), ast.Required()),
		)),
		ast.NewObject("grafanatest", "MappingType", ast.NewEnum([]ast.EnumValue{
			{Name: "value", Value: "value", Type: ast.String()},
			{Name: "range", Value: "range", Type: ast.String()},
		})),
		ast.NewObject("grafanatest", "Spec", ast.NewStruct(
			ast.NewStructField("type", ast.NewRef("grafanatest", "ValueMap"), ast.Required()),
		)),
	}

	req.Equal(testutils.ObjectsMap(objects...), schemaAst.Objects)
}

func txtarTestToCueInstance(tc *testutils.Test[string]) cue.Value {
	tc.Helper()

	return bytesToCueValue(tc.T, tc.ReadInput("schema.cue"))
}

func bytesToCueValue(t *testing.T, input []byte) cue.Value {
	t.Helper()

	overlay := map[string]load.Source{
		"/schema.cue": load.FromBytes(input),
	}

	bis := load.Instances([]string{"/schema.cue"}, &load.Config{
		Overlay:    overlay,
		ModuleRoot: "/",
	})
	values, err := cuecontext.New().BuildInstances(bis)
	require.NoError(t, err)

	return values[0]
}
