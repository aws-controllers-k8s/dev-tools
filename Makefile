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
	go build ${GO_LDFLAGS} -o ackdev ./cmd/ackdev/main.go

install: build
	cp ./ackdev $(shell go env GOPATH)/bin/ackdev

cpconfig:
	cp config.yaml ~/.ackdev.yaml

all: cpconfig build list