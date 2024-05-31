GOLANGCI_LINT_VERSION?=v1.54.2
DOCKER_TOOLS?=true

GOLANGCI_LINT=docker run --rm -v $(shell pwd):/app -w /app golangci/golangci-lint:${GOLANGCI_LINT_VERSION}-alpine golangci-lint

ifneq ($(DOCKER_TOOLS),true)
  GOLANGCI_LINT=golangci-lint
endif

.PHONY: lint
lint:
	$(GOLANGCI_LINT) run -c .golangci.yaml

.PHONY: tests
tests:
	go test -v ./...

.PHONY: deps
deps:
	go mod vendor

.PHONY: docs
docs:
	go run cmd/compiler-passes-docs/*
	go run cmd/cog-config-schemas/*
