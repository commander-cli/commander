exe = cmd/commander/*
cmd = commander
TRAVIS_TAG ?= "0.0.0"
CWD = $(shell pwd)

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
	go test ./...

test-coverage:
	$(info INFO: Starting build $@)
	go test -coverprofile c.out ./...


test-coverage-all: export COMMANDER_SSH_TEST = 1
test-coverage-all: export COMMANDER_TEST_SSH_HOST = localhost:2222
#test-coverage-all: export COMMANDER_TEST_SSH_PASS = password
test-coverage-all: export COMMANDER_TEST_SSH_USER = root
test-coverage-all: export COMMANDER_TEST_SSH_IDENTITY_FILE = $(CWD)/integration/containers/ssh/id_rsa
test-coverage-all:
	$(info INFO: Starting build $@)
	./integration/setup_unix.sh
	go test -coverprofile c.out ./...
	./integration/teardown_unix.sh

integration-unix: build
	$(info INFO: Starting build $@)
	./integration/setup_unix.sh
	commander test commander_unix.yaml
	./integration/teardown_unix.sh

integration-linux: build
	$(info INFO: Starting build $@)
	./integration/setup_unix.sh
	commander test commander_unix.yaml
	commander test commander_linux.yaml
	./integration/teardown_unix.sh

integration-windows: build
	$(info INFO: Starting build $@)
	commander test commander_windows.yaml

release-amd64:
	$(info INFO: Starting build $@)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-X main.version=$(TRAVIS_TAG) -s -w" -o release/$(cmd)-linux-amd64 $(exe)

release-arm:
	$(info INFO: Starting build $@)
	CGO_ENABLED=0 GOOS=linux GOARCH=arm go build -ldflags "-X main.version=$(TRAVIS_TAG) -s -w" -o release/$(cmd)-linux-arm $(exe)

release-386:
	$(info INFO: Starting build $@)
	CGO_ENABLED=0 GOOS=linux GOARCH=386 go build -ldflags "-X main.version=$(TRAVIS_TAG) -s -w" -o release/$(cmd)-linux-386 $(exe)

release-darwin-amd64:
	$(info INFO: Starting build $@)
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags "-X main.version=$(TRAVIS_TAG) -s -w" -o release/$(cmd)-darwin-amd64 $(exe)

release-darwin-386:
	$(info INFO: Starting build $@)
	CGO_ENABLED=0 GOOS=darwin GOARCH=386 go build -ldflags "-X main.version=$(TRAVIS_TAG) -s -w" -o release/$(cmd)-darwin-386 $(exe)

release-windows-amd64:
	$(info INFO: Starting build $@)
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags "-X main.version=$(TRAVIS_TAG) -s -w" -o release/$(cmd)-windows-amd64.exe $(exe)

release-windows-386:
	$(info INFO: Starting build $@)
	CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -ldflags "-X main.version=$(TRAVIS_TAG) -s -w" -o release/$(cmd)-windows-386.exe $(exe)


release: release-amd64 release-arm release-386 release-darwin-amd64 release-darwin-386 release-windows-amd64 release-windows-386