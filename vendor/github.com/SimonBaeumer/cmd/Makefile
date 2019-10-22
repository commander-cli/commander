exe = cmd/operator/*
cmd = operator
TRAVIS_TAG ?= "0.0.0"

.PHONY: deps lint test integration integration-windows git-hooks init

init: git-hooks

git-hooks:
	$(info INFO: Starting build $@)
	ln -sf ../../.githooks/pre-commit .git/hooks/pre-commit

deps:
	$(info INFO: Starting build $@)
	go mod vendor

build:
	$(info INFO: Starting build $@)
	go build $(exe)

lint:
	$(info INFO: Starting build $@)
	golint pkg/ cmd/

test:
	$(info INFO: Starting build $@)
	go test `go list ./... | grep -v examples`

test-windows:
	$(info INFO: Starting build $@)
	go test .

test-coverage:
	$(info INFO: Starting build $@)
	go test -coverprofile c.out `go list ./... | grep -v examples`
