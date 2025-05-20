package cog

import (
	"context"
	"testing"

	"cuelang.org/go/cue/cuecontext"
	"github.com/grafana/codejen"
	"github.com/stretchr/testify/require"

	"github.com/grafana/cog/internal/testutils"
)

func TestTypesFromSchema_OutputGo(t *testing.T) {
	schema := `
// Contains things.
Container: {
    str: string
}
`

	runCodegen := func(t *testing.T, config GoConfig) codejen.Files {
		t.Helper()

		req := require.New(t)

		cueValue := cuecontext.New().CompileString(schema)
		req.NoError(cueValue.Err())

		files, err := TypesFromSchema().
			CUEValue("sandbox", cueValue).
			Golang(config).
			Run(context.Background())
		req.NoError(err)

		return files
	}

	t.Run("generating go from cue - default options", func(t *testing.T) {
		req := require.New(t)

		expectedCode := `// Code generated - EDITING IS FUTILE. DO NOT EDIT.

package sandbox

// Contains things.
type Container struct {
	Str string ` + "`json:\"str\"`" + `
}

// NewContainer creates a new Container object.
func NewContainer() *Container {
	return &Container{}
}`

		cueValue := cuecontext.New().CompileString(schema)
		req.NoError(cueValue.Err())

		files := runCodegen(t, GoConfig{})

		req.Len(files, 1, "expected a single file")
		testutils.TrimSpacesDiffComparator(t, []byte(expectedCode), files[0].Data, "test.go")
	})

	t.Run("generating go from cue - with Equal() functions", func(t *testing.T) {
		req := require.New(t)

		expectedCode := `// Code generated - EDITING IS FUTILE. DO NOT EDIT.

package sandbox

// Contains things.
type Container struct {
	Str string ` + "`json:\"str\"`" + `
}

// NewContainer creates a new Container object.
func NewContainer() *Container {
	return &Container{}
}

// Equals tests the equality of two ` + "`Container`" + ` objects.
func (resource Container) Equals(other Container) bool {
	if resource.Str != other.Str {
		return false
	}

	return true
}`

		files := runCodegen(t, GoConfig{
			GenerateEqual: true,
		})

		req.Len(files, 1, "expected a single file")
		testutils.TrimSpacesDiffComparator(t, []byte(expectedCode), files[0].Data, "test.go")
	})
}

func TestTypesFromSchema_OutputTypescript(t *testing.T) {
	schema := `
// Contains things.
Container: {
    str: string
}
`

	expectedCode := `// Code generated - EDITING IS FUTILE. DO NOT EDIT.

// Contains things.
export interface Container {
	str: string;
}

export const defaultContainer = (): Container => ({
	str: "",
});`

	t.Run("generating typescript from cue", func(t *testing.T) {
		req := require.New(t)

		cueValue := cuecontext.New().CompileString(schema)
		req.NoError(cueValue.Err())

		files, err := TypesFromSchema().
			CUEValue("sandbox", cueValue).
			Typescript(TypescriptConfig{}).
			Run(context.Background())
		req.NoError(err)

		req.Len(files, 1, "expected a single file")
		testutils.TrimSpacesDiffComparator(t, []byte(expectedCode), files[0].Data, "test.go")
	})
}

func TestTypesFromSchema_OutputOpenAPI(t *testing.T) {
	schema := `
// Contains things.
Container: {
    str: string
}
`

	expectedCode := `{
  "openapi": "3.0.0",
  "info": {
    "title": "sandbox",
    "version": "0.0.0",
    "x-schema-identifier": "",
    "x-schema-kind": ""
  },
  "paths": {},
  "components": {
    "schemas": {
      "Container": {
        "type": "object",
        "additionalProperties": false,
        "required": [
          "str"
        ],
        "properties": {
          "str": {
            "type": "string"
          }
        },
        "description": "Contains things."
      }
    }
  }
}`

	t.Run("generating openAPI from cue", func(t *testing.T) {
		req := require.New(t)

		cueValue := cuecontext.New().CompileString(schema)
		req.NoError(cueValue.Err())

		files, err := TypesFromSchema().
			CUEValue("sandbox", cueValue).
			GenerateOpenAPI(OpenAPIGenerationConfig{}).
			Run(context.Background())
		req.NoError(err)

		req.Len(files, 1, "expected a single file")
		testutils.TrimSpacesDiffComparator(t, []byte(expectedCode), files[0].Data, "test.json")
	})
}

func TestTypesFromSchema_SchemaTransformations(t *testing.T) {
	schema := `
// Contains things.
Container: {
    str: string
}
`

	expectedCode := `// Code generated - EDITING IS FUTILE. DO NOT EDIT.

package sandbox

// Contains things.
// Transformed by cog.
type ExampleContainer struct {
	Str string ` + "`json:\"str\"`" + `
}

// NewExampleContainer creates a new ExampleContainer object.
func NewExampleContainer() *ExampleContainer {
	return &ExampleContainer{}
}`

	t.Run("transformations can be applied on the input schema", func(t *testing.T) {
		req := require.New(t)

		cueValue := cuecontext.New().CompileString(schema)
		req.NoError(cueValue.Err())

		files, err := TypesFromSchema().
			CUEValue("sandbox", cueValue).
			Golang(GoConfig{}).
			SchemaTransformations(
				AppendCommentToObjects("Transformed by cog."),
				PrefixObjectsNames("Example"),
			).
			Run(context.Background())
		req.NoError(err)

		req.Len(files, 1, "expected a single file")
		testutils.TrimSpacesDiffComparator(t, []byte(expectedCode), files[0].Data, "test.go")
	})
}
