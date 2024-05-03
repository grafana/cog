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

For the [`kind-registry`](https://github.com/grafana/kind-registry):

```console
$ git clone https://github.com/grafana/kind-registry ../kind-registry
$ go run cmd/cli/main.go generate \
    --config ./config/foundation_sdk.dev.yaml \
    --parameters kind_registry_version=v10.4.x
```
