package simplecue

import (
	"strings"
	"testing"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/load"
	"github.com/grafana/cog/internal/testutils"
	"github.com/grafana/cog/pkg/ir"
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

	objects := []ir.Object{
		ir.NewObject("grafanatest", "Ref", ir.String()),
		ir.NewObject("grafanatest", "Container", ir.NewStruct(
			ir.NewStructField("ref", ir.NewRef("grafanatest", "Ref"), ir.Required()),
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

	objects := []ir.Object{
		ir.NewObject("grafanatest", "Origin", ir.NewStruct(
			ir.NewStructField("creator", ir.String(), ir.Required()),
		)),
		ir.NewObject("grafanatest", "spec", ir.NewStruct(
			ir.NewStructField("title", ir.String(), ir.Required()),
			ir.NewStructField("origin", ir.NewRef("grafanatest", "Origin"), ir.Required()),
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

	objects := []ir.Object{
		ir.NewObject("grafanatest", "schema.#Origin", ir.NewStruct(
			ir.NewStructField("creator", ir.String(), ir.Required()),
		)),
		ir.NewObject("grafanatest", "schema.spec.#Details", ir.NewMap(
			ir.String(),
			ir.Any(),
		)),
		ir.NewObject("grafanatest", "spec", ir.NewStruct(
			ir.NewStructField("title", ir.String(), ir.Required()),
			ir.NewStructField("origin", ir.NewRef("grafanatest", "schema.#Origin"), ir.Required()),
			ir.NewStructField("details", ir.NewRef("grafanatest", "schema.spec.#Details"), ir.Required()),
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

	objects := []ir.Object{
		ir.NewObject("grafanatest", "ValueMap", ir.NewStruct(
			ir.NewStructField("type", ir.NewConstantReferenceType("grafanatest", "MappingType", "value"), ir.Required()),
		)),
		ir.NewObject("grafanatest", "MappingType", ir.NewEnum([]ir.EnumValue{
			{Name: "value", Value: "value", Type: ir.String()},
			{Name: "range", Value: "range", Type: ir.String()},
		})),
		ir.NewObject("grafanatest", "Spec", ir.NewStruct(
			ir.NewStructField("type", ir.NewRef("grafanatest", "ValueMap"), ir.Required()),
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
