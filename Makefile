.PHONY: help build test fmt tidy clean

GO ?= go

help:
	@echo "Common targets:"
	@echo "  make build   # Compile the CLI binary"
	@echo "  make test    # Run unit tests"
	@echo "  make fmt     # gofmt project source files"
	@echo "  make tidy    # Update go.mod/go.sum"
	@echo "  make clean   # Remove build artifacts"

build:
	$(GO) build ./...

test:
	$(GO) test ./...

fmt:
	gofmt -w $$(find . -name '*.go' -not -path './vendor/*')

tidy:
	$(GO) mod tidy

clean:
	$(GO) clean
