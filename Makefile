# variable definitions
VERSION := $(shell git describe --tags --always --dirty)
GOVERSION := $(shell go version)
BUILDTIME := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
BUILDER := $(shell echo "`git config user.name` <`git config user.email`>")
LDFLAGS := -X 'main.version=$(VERSION)' \
			-X 'main.buildTime=$(BUILDTIME)' \
			-X 'main.builder=$(BUILDER)' \
			-X 'main.goversion=$(GOVERSION)'

.PHONY: test build lint runtime

build:
	go build -o "bin/asml" -v -ldflags "$(LDFLAGS)" ./cmd/asml

test:
	@go test -race $$(go list ./pkg/...)

# https://github.com/golang/lint
# go get github.com/golang/lint/golint
lint:
	@golint ./src/...

runtime:
	go run ./cmd/runtime/make.go -in ./cmd/runtime/runtime.asml -out ./pkg/lexer/runtime.go
