.PHONY: default setup build test lint vet format tidy clean

GO ?= go

default: help

## Install development and CI tooling
setup:
	$(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	$(GO) install mvdan.cc/gofumpt@latest
	$(GO) install golang.org/x/tools/cmd/goimports@latest

## Build the CLI binary
build:
	$(GO) build ./...

run:
	@$(GO) run ./cmd/dcalc/main.go

## Run unit tests
test:
	$(GO) test ./...

## Run golangci-lint and go vet
lint: vet
	golangci-lint run ./...

## Run go vet static analysis
vet:
	$(GO) vet ./...

## Run gofumpt + goimports on all source files
format:
	gofumpt -w $$(find . -name '*.go' -not -path './vendor/*')
	goimports -w $$(find . -name '*.go' -not -path './vendor/*')

## Update go.mod/go.sum
tidy:
	$(GO) mod tidy

## Clean build artifacts
clean:
	$(GO) clean

# Utils
# ============================================================================

.PHONY = help vars _print-var

## This help screen
help:
	@printf "Available targets:\n\n"
	@awk '/^[a-zA-Z\-\_0-9%:\\]+/ { \
         helpMessage = match(lastLine, /^## (.*)/); \
         if (helpMessage) { \
             helpCommand = $$1; \
             helpMessage = substr(lastLine, RSTART + 3, RLENGTH); \
             gsub("\\\\", "", helpCommand); \
             gsub(":+$$", "", helpCommand); \
             printf "  \x1b[32;01m%-35s\x1b[0m %s\n", helpCommand, helpMessage; \
         } \
    } \
    { lastLine = $$0 }' $(MAKEFILE_LIST) | sort -u
	@printf "\n"

## Show the variables used in the Makefile and their values
vars:
	@printf "Variable values:\n\n"
	@awk 'BEGIN { FS = "[:?]?="; } /^[A-Za-z0-9_]+[[:space:]]*[:?]?=/ { \
		if ($$0 ~ /\?=/) operator = "?="; \
		else if ($$0 ~ /:=/) operator = ":="; \
		else operator = "="; \
		print $$1, operator; \
		} \
		{ lastLine = $$0 }' $(MAKEFILE_LIST) | \
		while read var op; do \
			value=$$(make --no-print-directory -f $(MAKEFILE_LIST) _print-var VAR=$$var); \
			printf "  \x1b[32;01m%-35s\x1b[0m%2s \x1b[34;01m%s\x1b[0m\n" "$$var" "$$op" "$$value"; \
			done
		@printf "\n"

_print-var:
	@echo $($(VAR))
