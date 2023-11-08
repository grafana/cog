# cog

`cog` is a CLI tool to generate code from schemas.

## Maturity

> _The code in this repository should be considered experimental. Documentation is only
available alongside the code. It comes with no support, but we are keen to receive
feedback on the product and suggestions on how to improve it, though we cannot commit
to resolution of any particular issue. No SLAs are available. It is not meant to be used
in production environments, and the risks are unknown/high._

Grafana Labs defines experimental features as follows:

> Projects and features in the Experimental stage are supported only by the Engineering
teams; on-call support is not available. Documentation is either limited or not provided
outside of code comments. No SLA is provided.
>
> Experimental projects or features are primarily intended for open source engineers who
want to participate in ensuring systems stability, and to gain consensus and approval
for open source governance projects.
>
> Projects and features in the Experimental phase are not meant to be used in production
environments, and the risks are unknown/high.

## Goals

* Support common schema formats: [JSON Schema](https://json-schema.org/), [OpenAPI](https://www.openapis.org/), ...
* Generate code in a wide range of languages: Golang, Typescript, Java, ...
* Generate *types* described by schemas
* Generate developer-friendly *builder libraries*, allowing the creation of complex Grafana-related objects as-code

## Usage

For specific schemas:

```console
$ go run cmd/cli/main.go generate \
    --output ./generated \
    --kindsys-core ./schemas/kindsys/core/dashboard \
    --cue ./schemas/cue/common \
    --include-cue-import ./schemas/cue/common:github.com/grafana/grafana/packages/grafana-schema/src/common \
    --jsonschema ./schemas/jsonschema/core/playlist/playlist.json \
    --kindsys-composable ./schemas/kindsys/composable/timeseries \
    --kindsys-composable ./schemas/kindsys/composable/logs \
    --kindsys-composable ./schemas/kindsys/composable/prometheus \
    --veneers ./config
```

For the [`kind-registry`](https://github.com/grafana/kind-registry):

```console
$ git clone https://github.com/grafana/kind-registry schemas/kind-registry
$ go run cmd/cli/main.go generate \
    --output ./generated \
    --kind-registry ./schemas/kind-registry \
    --cue ./schemas/cue/common \
    --include-cue-import ./schemas/cue/common:github.com/grafana/grafana/packages/grafana-schema/src/common \
    --veneers ./config
```
