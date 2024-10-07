# Releasing `grafana-foundation-sdk`

Releases to [`grafana-foundation-sdk`](https://github.com/grafana/grafana-foundation-sdk/)
are handled by the `./scripts/release-all.sh` and `./scripts/release-version.sh`
scripts.

## Requirements

* [`go`](https://go.dev/doc/install)
* [`gh`](https://cli.github.com/) (GitHub CLI), with authentication configured using `gh auth login`

Before releasing, ensure that cog's `main` branch is checked out and up-to-date:

```console
git checkout main
git pull --ff-only origin main
```

## Releasing all supported Grafana versions

The `release-all.sh` script is a thin wrapper that calls `release-version.sh`
in a loop, for a pre-configured list of supported versions.

As such, it's options and behavior are identical to what was described previously.

```console
DRY_RUN=no ./scripts/release-all.sh
```

## Releasing a specific Grafana version

Start a dry-run release for the desired version of the schemas:

```console
./scripts/release-version.sh "next;v11.2.x;v11.1.x;v11.0.x;v10.4.x;v10.3.x;v10.2.x;v10.1.x" v10.2.x
```

This will perform the release process without pushing any change to give a safe opportunity to review the release.
Details on where to find the generated code and inspect it will be written to the standard output.

> [!NOTE]
> To generate accurate READMEs, a list of all supported Grafana versions needs to be given to the script.
> An accurate list can be found in the `./scripts/release-all.sh` file.

If everything looks good, proceed for real:

```console
DRY_RUN=no ./scripts/release-version.sh "next;v11.2.x;v11.1.x;v11.0.x;v10.4.x;v10.3.x;v10.2.x;v10.1.x" v10.2.x
```

Disabling dry-run will allow the script to push changes to a *release preview branch* and open a pull-request onto the
actual *release branch*. To finalise the release, review this PR and merge it.

## Options

The following environment variables can be used to alter the behavior of the release script:

* `DRY_RUN`: flag indicating whether the release should actually be published or not – **dry-run is ON by default**.
* `LOG_LEVEL`: number indicating how verbose the script should be. 7 = debug -> 0 = emergency, defaults to 6.
* `COG_CMD`: command used to run `cog`.
* `GH_CLI_CMD`: command used to run `gh` (GitHub CLI).
* `KIND_REGISTRY_PATH`: path to the `kind-registry` repository. If it doesn't exist, the [`kind-registry`](https://github.com/grafana/kind-registry/) will be cloned at that path.
* `FOUNDATION_SDK_PATH`: path to the `grafana-foundation-sdk` repository. If it doesn't exist, the [`grafana-foundation-sdk`](https://github.com/grafana/grafana-foundation-sdk/) will be cloned at that path.
