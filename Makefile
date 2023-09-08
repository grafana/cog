GOLANGCI_LINT_VERSION='v1.54.2'

.PHONY: lint
lint:
	docker run --rm -v $(shell pwd):/app -w /app golangci/golangci-lint:${GOLANGCI_LINT_VERSION}-alpine golangci-lint run -c .golangci.yaml

.PHONY: tests
tests:
	go test -v ./...

.PHONY: deps
deps:
	go mod vendor
