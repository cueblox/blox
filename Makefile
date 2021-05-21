SOURCE_FILES?=./...
TEST_PATTERN?=.
TEST_OPTIONS?=
TEST_TIMEOUT?=15m
TEST_PARALLEL?=2
DOCKER_BUILDKIT?=1
export DOCKER_BUILDKIT

export PATH := ./bin:$(PATH)
export GO111MODULE := on

# Install all the build and lint dependencies
setup:
	go mod tidy
	git config core.hooksPath .githooks
.PHONY: setup

test:
	go test $(TEST_OPTIONS) -p $(TEST_PARALLEL) -v -failfast -race -coverpkg=./... -covermode=atomic -coverprofile=coverage.txt $(SOURCE_FILES) -run $(TEST_PATTERN) -timeout=$(TEST_TIMEOUT)
.PHONY: test

cover: test
	go tool cover -html=coverage.txt
.PHONY: cover

fmt:
	gofumpt -w .
.PHONY: fmt

lint: check
	golangci-lint run
.PHONY: check

ci: lint test
.PHONY: ci

build:
	go build -o blox ./cmd/blox/main.go
	./scripts/cmd_docs.sh
.PHONY: build

deps:
	go get -u
	go mod tidy
	go mod verify
.PHONY: deps

serve:
	@docker run --rm -it -p 8000:8000 -v ${PWD}/dogfood/sites/docs:/docs squidfunk/mkdocs-material
.PHONY: serve

todo:
	@grep \
		--exclude-dir=vendor \
		--exclude-dir=node_modules \
		--exclude-dir=bin \
		--exclude=Makefile \
		--text \
		--color \
		-nRo -E ' TODO:.*|SkipNow' .
.PHONY: todo

.DEFAULT_GOAL := build
