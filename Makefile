# Within devbox
ifneq "$(MISE_SHELL)" ""
    MISE_EXEC:=
else # Normal shell
    MISE_EXEC:=mise exec --
endif

##@ General

# The help target prints out all targets with their descriptions organized
# beneath their categories. The categories are represented by '##@' and the
# target descriptions by '##'. The awk commands is responsible for reading the
# entire set of makefiles included in this invocation, looking for lines of the
# file as xyz: ## something, and then pretty-format the target and help. Then,
# if there's a line with ##@ something, that gets pretty-printed as a category.
# More info on the usage of ANSI control characters for terminal formatting:
# https://en.wikipedia.org/wiki/ANSI_escape_code#SGR_parameters
# More info on the awk command:
# http://linuxcommand.org/lc3_adv_awk.php

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development

.PHONY: lint
lint: dev-env-check-binaries ## Lints the code base.
	$(MISE_EXEC) golangci-lint run -c .golangci.yaml

.PHONY: tests
tests: dev-env-check-binaries gen-tests ## Runs the tests.
	$(MISE_EXEC) go test ./...

.PHONY: deps
deps: dev-env-check-binaries ## Installs the dependencies.
	$(MISE_EXEC) go mod vendor
	$(MISE_EXEC) composer install --no-interaction --quiet
	$(MISE_EXEC) pip install -qq -r requirements.txt

.PHONY: docs
docs: dev-env-check-binaries ## Generates the documentation.
	@$(MISE_EXEC) go run cmd/compiler-passes-docs/*
	@$(MISE_EXEC) go run cmd/veneers-docs/*
	@$(MISE_EXEC) go run cmd/cog-config-schemas/*
	@$(MISE_EXEC) zensical build -f ./mkdocs-github.yml

.PHONY: serve-docs
serve-docs: dev-env-check-binaries ## Builds and serves the documentation.
	$(MISE_EXEC) zensical serve -f ./mkdocs-github.yml

.PHONY: fetch-foundation-sdk
fetch-foundation-sdk: ## Fetches the Foundation SDK.
	@./scripts/fetch-repo.sh ./build/foundation-sdk https://github.com/grafana/grafana-foundation-sdk.git

.PHONY: fetch-kind-registry
fetch-kind-registry: ## Fetches the kind-registry.
	@./scripts/fetch-repo.sh ./build/kind-registry https://github.com/grafana/kind-registry.git

.PHONY: gen-sdk-dev
gen-sdk-dev: dev-env-check-binaries fetch-foundation-sdk fetch-kind-registry ## Generates a dev version of the Foundation SDK.
	@rm -rf generated
	$(MISE_EXEC) go run cmd/cli/main.go generate \
		--config ./build/foundation-sdk/.cog/config.yaml \
		--debug \
		--parameters "output_dir=./generated/%l,kind_registry_path=./build/kind-registry,go_package_root=github.com/grafana/cog/generated/go"

.PHONY: gen-tests
gen-tests: dev-env-check-binaries ## Generates the code described by tests schemas.
	$(MISE_EXEC) go run ./cmd/cli/ generate \
		--config ./config/foundation_sdk.tests.yaml

.PHONY: dev-env-check-binaires
dev-env-check-binaries: ## Check that the required binary are present.
	@mise version >/dev/null 2>&1 || (echo "ERROR: mise is required. See https://mise.jdx.dev/getting-started.html"; exit 1)

.PHONY: test-inspect
test-inspect:
	$(MISE_EXEC) go run cmd/cli/main.go inspect \
		--config config/inspect.test.yaml \
		--parameters "cue_path=$(CUE_PATH)"
