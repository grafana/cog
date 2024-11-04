package cog

import (
	"context"
	"fmt"

	"cuelang.org/go/cue/cuecontext"
)

// ExampleCueValue demonstrates how to generate types from a CUE value.
func Example_cueValue() {
	schema := `
// Contains things.
Container: {
    str: string
    trueOrFalse: bool
    anything: {...}
    data: bytes

    num_unit8: uint8
    num_int8: int8
    num_uint16: uint16
    num_int16: int16
    num_uint32: uint32
    num_int32: int32
    num_uint64: uint64
    num_int64: int64

	disjunction: Bar | Baz
}

// This is a bar.
Bar: {
       type: "bar"
       foo: string
}

// This is a baz.
Baz: {
       type: "baz"
       boo: string
}
`

	cueValue := cuecontext.New().CompileString(schema)
	if cueValue.Err() != nil {
		panic(cueValue.Err())
	}

	types, err := TypesFromSchema().
		CUEValue("sandbox", cueValue).
		Typescript().
		Run(context.Background())
	if err != nil {
		panic(err)
	}

	fmt.Println(string(types))
}
