COVERAGE_FILE ?= coverage.out

.PHONY: build
build: build_bot build_scrapper

.PHONY: build_bot
build_bot:
	@echo "Выполняется go build для таргета bot"
	@mkdir -p bin
	@go build -o ./bin/bot ./cmd/bot

.PHONY: build_scrapper
build_scrapper:
	@echo "Выполняется go build для таргета scrapper"
	@mkdir -p bin
	@go build -o ./bin/scrapper ./cmd/scrapper

## test: run all tests
.PHONY: test
test:
	@go test -coverpkg='github.com/es-debug/backend-academy-2024-go-template/...' --race -count=1 -coverprofile='$(COVERAGE_FILE)' ./...
	@go tool cover -func='$(COVERAGE_FILE)' | grep ^total | tr -s '\t'
	@go tool cover -html='$(COVERAGE_FILE)' -o coverage.html && xdg-open coverage.html

.PHONY: lint
lint: lint-golang lint-proto

.PHONY: lint-golang
lint-golang:
	@if ! command -v 'golangci-lint' &> /dev/null; then \
  		echo "Please install golangci-lint!"; exit 1; \
  	fi;
	@golangci-lint -v run --fix ./...

.PHONY: lint-proto
lint-proto:
	@if ! command -v 'easyp' &> /dev/null; then \
  		echo "Please install easyp!"; exit 1; \
	fi;
	@easyp lint

.PHONY: generate
generate: generate_proto generate_openapi

.PHONY: generate_openapi
generate_openapi: generate_bot generate_scrapper

.PHONY: generate_proto
generate_proto:
	@if ! command -v 'easyp' &> /dev/null; then \
		echo "Please install easyp!"; exit 1; \
	fi;
	@easyp generate

.PHONY: generate_bot
generate_bot:
	@if ! command -v 'oapi-codegen' &> /dev/null; then \
		echo "Please install oapi-codegen!"; exit 1; \
	fi;
	@mkdir -p internal/api/openapi/v1/clients/bot
	@mkdir -p internal/api/openapi/v1/servers/bot
	@oapi-codegen --config oapi.bot-client.yaml api/openapi/v1/bot-api.yaml
	@oapi-codegen --config oapi.bot-server.yaml api/openapi/v1/bot-api.yaml

.PHONY: generate_scrapper
generate_scrapper:
	@if ! command -v 'oapi-codegen' &> /dev/null; then \
		echo "Please install oapi-codegen!"; exit 1; \
	fi;
	@mkdir -p internal/api/openapi/v1/clients/scrapper
	@mkdir -p internal/api/openapi/v1/servers/scrapper
	@oapi-codegen --config oapi.scrapper-client.yaml api/openapi/v1/scrapper-api.yaml
	@oapi-codegen --config oapi.scrapper-server.yaml api/openapi/v1/scrapper-api.yaml

.PHONY: clean
clean:
	@rm -rf bin
