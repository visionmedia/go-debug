
test:
	@go test

bench:
	@go test -bench=. -benchmem

.PHONY: bench test
