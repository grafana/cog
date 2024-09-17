.PHONY: lint
lint:
	golangci-lint run -c .golangci.yaml

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

.PHONY: gen-sdk-dev
gen-sdk-dev:
	rm -rf generated
	go run cmd/cli/main.go generate \
		--config ./config/foundation_sdk.dev.yaml \
		--parameters kind_registry_version=next,grafana_version=main

.PHONY: run-go-example
run-go-example:
	go run ./examples/_go/*

.PHONY: run-php-example
run-php-example:
	cd ./examples/php && \
	composer install && \
	php index.php

.PHONY: run-ts-example
run-ts-example:
	ts-node examples/typescript
