.PHONY: bench test lint

LINE_LENGTH=$(shell yq r .golangci.yml linters-settings.lll.line-length)

test:
	@go test

bench:
	@go test -bench=. -benchmem

 ## Run the linter
lint:
	@golines -m $(LINE_LENGTH) . -w --ignored-dirs=vendor
	@golangci-lint run --fix
