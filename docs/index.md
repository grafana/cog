# Cog

`cog` is a **Co**de **g**enerator created with the following objectives in mind:

* Support multiple schema formats: [CUE](https://cuelang.org/), [JSON Schema](https://json-schema.org/), [OpenAPI](https://www.openapis.org/), ...
* Generate code in a wide range of languages: Golang, Java, PHP, Python, Typescript, …
* Generate *types* described by schemas
* Generate developer-friendly *builder libraries*, allowing the creation of complex objects as-code

## Usage


=== "As a CLI"

    Download the `cog` binary from our [releases](https://github.com/grafana/cog/releases),
    and run the codegen pipeline:

    ```console
    cog generate --config ./cog-pipeline.yaml
    ```

=== "As a Go Library"

    See the [Go documentation](https://pkg.go.dev/github.com/grafana/cog) for more example and a complete API reference.

    ```go
    package main

    import (
        "context"
        "fmt"

	    "github.com/grafana/cog"
    )

    func main() {
        files, err := cog.TypesFromSchema().
            CUEModule("/path/to/cue/module").
            SchemaTransformations(
                cog.AppendCommentToObjects("Transformed by cog."),
                cog.PrefixObjectsNames("Example"),
            ).
            Golang(cog.GoConfig{}).
            Run(context.Background())
        if err != nil {
            panic(err)
        }

        if len(files) != 1 {
            panic("expected a single file :(")
        }

        fmt.Println(string(files[0].Data))
    }
    ```

## Maturity

Cog should be considered as "public preview". While it is used by Grafana Labs in production, it still is under active development.

Additional information can be found in [Release life cycle for Grafana Labs](https://grafana.com/docs/release-life-cycle/).

!!! note

    Bugs and issues are handled solely by Engineering teams. On-call support or SLAs are not available.