SHELL := /bin/bash 
GO111MODULE=on

# Build ldflags
VERSION ?= "v0.0.0"
GITCOMMIT=$(shell git rev-parse HEAD)
BUILDDATE=$(shell date -u +'%Y-%m-%dT%H:%M:%SZ')
PKG_PATH=github.com/aws-controllers-k8s/dev-tools/pkg
GO_LDFLAGS=-ldflags "-X $(PKG_PATH)/version.GitVersion=$(VERSION) \
			-X $(PKG_PATH)/version.GitCommit=$(GITCOMMIT) \
			-X $(PKG_PATH)/version.BuildDate=$(BUILDDATE)"

build:
	go build ${GO_LDFLAGS} -o ./bin/ackdev ./cmd/ackdev/main.go

install: build
	cp ./bin/ackdev $(shell go env GOPATH)/bin/ackdev

test:
	go test -v ./...

.PHONY: test install mocks

mocks:
	@echo -n "building mocks for pkg/git ... "
	@mockery --quiet --name OpenCloner --tags=codegen --case=underscore --output=mocks --dir=pkg/git
	@echo "ok."
	@echo -n "building mocks for pkg/github ... "
	@mockery --quiet --all --tags=codegen --case=underscore --output=mocks --dir=pkg/github
	@echo "ok."
