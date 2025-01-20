# Codegen pipelines

Code generation is configured via a *codegen pipeline* describing:

* the schemas used as input
* possible transformations applied to the schemas
* possible transformations applied to the builders
* the desired outputs

The diagram below describes – from a high-level perspective – how `cog` runs such a pipeline:

```mermaid
flowchart LR
    schemas@{ shape: docs, label: "Schemas\n(CUE, OpenAPI, Jsonschema)"}
    types_ir@{shape: diamond, label: "Types IR"}
    builders_ir@{shape: diamond, label: "Builders IR"}

    parsers@{label: "Parsers"}
    jennies@{label: "Jennies"}

    compiler_passes@{shape: lean-l, label: "Schema transformations"}
    veneers@{shape: lean-l, label: "Builder transformations"}

    output_go@{shape: docs, label: "Go types & builders"}
    output_ts@{shape: docs, label: "Typescript types & builders"}
    output_etc@{shape: docs, label: "…"}

    schemas --> parsers
    compiler_passes -.-> |Modifies|types_ir
    veneers -.-> |Modifies|builders_ir
    types_ir -.-> builders_ir

    subgraph cog [Cog]
        parsers --> types_ir

        subgraph jennies_ctx [Intermediate representations]
            types_ir
            builders_ir
        end

        jennies_ctx --> jennies
    end

    jennies --> output_go & output_ts & output_etc
```