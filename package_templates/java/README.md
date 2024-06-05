# Grafana Foundation SDK – Java

A set of tools, types and *builder libraries* for building and manipulating Grafana objects in Java.

> [!NOTE]
> This branch contains **types and builders generated for Grafana {{ .Extra.GrafanaVersion }}.**
> Other supported versions of Grafana can be found at [this repository's root](https://github.com/grafana/grafana-foundation-sdk/).

## Installing

```kotlin
implementation("com.grafana:grafana-foundation-sdk:{{ .Extra.GrafanaVersion|registryToSemver }}-{{ .Extra.BuildTimestamp }}")
```

```xml
<dependency>
    <groupId>com.grafana</groupId>
    <artifactId>grafana-foundation-sdk</artifactId>
    <version>{{ .Extra.GrafanaVersion|registryToSemver }}-{{ .Extra.BuildTimestamp }}</version>
</dependency>
```

## Example usage

TODO

## Maturity

> [!WARNING]
> The code in this repository should be considered experimental. Documentation is only
available alongside the code. It comes with no support, but we are keen to receive
feedback and suggestions on how to improve it, though we cannot commit
to resolution of any particular issue.

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

## License

[Apache 2.0 License](./LICENSE)
