package main

import (
	"context"
	"fmt"

	"cuelang.org/go/cue/cuecontext"
	"github.com/grafana/cog"
)

const val = `
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

func main() {
	v := cuecontext.New().CompileString(val)
	if v.Err() != nil {
		panic(v.Err())
	}

	files, err := cog.TypesFromSchema().
		CUEValue("sandbox", v).
		Golang(cog.GoConfig{}).
		//Typescript().
		SchemaTransformations(
			cog.AppendCommentToObjects("+k8s:openapi-gen=true"),
			cog.PrefixObjectsNames("Foo"),
		).
		Run(context.Background())
	if err != nil {
		panic(err)
	}

	fmt.Println(string(files[0].Data))
}
