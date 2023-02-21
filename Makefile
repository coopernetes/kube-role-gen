SHELL := /bin/bash

# The name of the executable (default is current directory name)
TARGET := $(shell echo "$${PWD##*/}" )
MAIN := ./cmd/kube-role-gen
.DEFAULT_GOAL: $(TARGET)

# These will be provided to the target
VERSION := 0.0.6
BUILD := `git rev-parse HEAD`

# Use linker flags to provide version/build settings to the target
LDFLAGS=-ldflags "-X=main.Version=$(VERSION) -X=main.Build=$(BUILD)"

# go source files, ignore vendor directory
SRC = $(shell find . -type f -name '*.go' -not -path "./vendor/*")

.PHONY: all build test clean install uninstall fmt simplify check run

all: check install

$(TARGET): $(SRC)
	@go build $(LDFLAGS) -o $(TARGET) $(MAIN)/main.go

build: $(TARGET)
	@true

test:
	go test ./...

e2e:
	tests/e2e_tests.sh

clean:
	@rm -f $(TARGET)

install:
	@go install $(LDFLAGS) $(MAIN)

uninstall: clean
	@rm -f $$(which ${TARGET})

fmt:
	@gofmt -l -w $(SRC)

simplify:
	@gofmt -s -l -w $(SRC)

check:
	@test -z $(shell gofmt -l $(MAIN)/main.go | tee /dev/stderr) || echo "[WARN] Fix formatting issues with 'make fmt'"
	@go vet ./...

run: install
	@$(TARGET)