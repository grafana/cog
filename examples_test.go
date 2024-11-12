package cog

import (
	"context"
	"fmt"

	"cuelang.org/go/cue/cuecontext"
)

// ExampleSchemaTransformations demonstrates how apply transformations to input schemas.
func Example_schemaTransformations() {
	schema := `
// Contains things.
Container: {
    str: string
}

// This is a bar.
Bar: {
       type: "bar"
       foo: string
}
`

	cueValue := cuecontext.New().CompileString(schema)
	if cueValue.Err() != nil {
		panic(cueValue.Err())
	}

	files, err := TypesFromSchema().
		CUEValue("sandbox", cueValue).
		Golang(GoConfig{}).
		SchemaTransformations(
			AppendCommentToObjects("Transformed by cog."),
			PrefixObjectsNames("Example"),
		).
		Run(context.Background())
	if err != nil {
		panic(err)
	}

	if len(files) != 1 {
		panic("expected a single file :(")
	}

	fmt.Println(string(files[0].Data))
}

// ExampleCueValue demonstrates how to generate types from a CUE value.
func Example_cueValue() {
	schema := `
// Contains things.
Container: {
    str: string
}

// This is a bar.
Bar: {
       type: "bar"
       foo: string
}
`

	cueValue := cuecontext.New().CompileString(schema)
	if cueValue.Err() != nil {
		panic(cueValue.Err())
	}

	files, err := TypesFromSchema().
		CUEValue("sandbox", cueValue).
		Golang(GoConfig{}).
		Run(context.Background())
	if err != nil {
		panic(err)
	}

	if len(files) != 1 {
		panic("expected a single file :(")
	}

	fmt.Println(string(files[0].Data))
}

// ExampleCueModule demonstrates how to generate types from a CUE module living on the filesystem.
func Example_cueModule() {
	files, err := TypesFromSchema().
		CUEModule("/path/to/cue/module").
		Golang(GoConfig{}).
		Run(context.Background())
	if err != nil {
		panic(err)
	}

	if len(files) != 1 {
		panic("expected a single file :(")
	}

	fmt.Println(string(files[0].Data))
}

// ExampleGoOutput demonstrates how to generate Go types from a CUE value.
func Example_goOutput() {
	schema := `
// Contains things.
Container: {
    str: string
}

// This is a bar.
Bar: {
       type: "bar"
       foo: string
}
`

	cueValue := cuecontext.New().CompileString(schema)
	if cueValue.Err() != nil {
		panic(cueValue.Err())
	}

	files, err := TypesFromSchema().
		CUEValue("sandbox", cueValue).
		Golang(GoConfig{}).
		Run(context.Background())
	if err != nil {
		panic(err)
	}

	if len(files) != 1 {
		panic("expected a single file :(")
	}

	fmt.Println(string(files[0].Data))
}

// ExampleTypescriptOutput demonstrates how to generate Typescript types from a CUE value.
func Example_typescriptOutput() {
	schema := `
// Contains things.
Container: {
    str: string
}

// This is a bar.
Bar: {
       type: "bar"
       foo: string
}
`

	cueValue := cuecontext.New().CompileString(schema)
	if cueValue.Err() != nil {
		panic(cueValue.Err())
	}

	files, err := TypesFromSchema().
		CUEValue("sandbox", cueValue).
		Typescript(TypescriptConfig{}).
		Run(context.Background())
	if err != nil {
		panic(err)
	}

	if len(files) != 1 {
		panic("expected a single file :(")
	}

	fmt.Println(string(files[0].Data))
}
