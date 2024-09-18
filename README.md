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
$ devbox run gen-sdk-dev
```

## Development setup

cog relies on [devbox](https://www.jetify.com/devbox/docs/) to manage all
the tools and programming languages it targets.

A shell including all the required tools is accessible via:

```console
$ devbox shell
```

This shell can be exited like any other shell, with `exit` or `CTRL+D`.

One-off commands can be executed within the devbox shell as well:

```console
$ devbox run go version
```

Various cog-specific commands also exist:

```console
$ devbox run
```

Packages can be installed using:

```console
devbox add go@1.23
```

Available packages can be found on the [NixOS package repository](https://search.nixos.org/packages).
