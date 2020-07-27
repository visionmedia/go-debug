.PHONY: help test report bench lint

LINE_LENGTH=$(shell yq r .golangci.yml linters-settings.lll.line-length)

help: ## Generates this help message
	@grep -E '^[0-9a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

test: ## unit tests
	@go test -failfast -covermode=atomic -coverprofile=coverage.txt

report: ## coverage report
	@go tool cover -html=coverage.txt

bench: ## benchmark tests
	@go test -bench=. -benchmem

lint: ## Run the linter
	@golines -m $(LINE_LENGTH) . -w --ignored-dirs=vendor
	@golangci-lint run --fix
