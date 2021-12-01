exe = cmd/commander/*
cmd = commander
GIT_RELEASE_TAG ?= "0.0.0"
CWD = $(shell pwd)

.PHONY: init
init: git-hooks

.PHONY: git-hooks
git-hooks:
	$(info INFO: Starting build $@)
	ln -sf ../../.githooks/pre-commit .git/hooks/pre-commit

.PHONY: deps
deps:
	$(info INFO: Starting build $@)
	go mod vendor

.PHONY: build
build:
	$(info INFO: Starting build $@)
	go build $(exe)

.PHONY: lint
lint:
	$(info INFO: Starting build $@)
	golint pkg/ cmd/

.PHONY: test
test:
	$(info INFO: Starting build $@)
	go test ./...

.PHONY: test-coverage
test-coverage:
	$(info INFO: Starting build $@)
	go test -coverprofile c.out ./...

.PHONY: test-coverage-all-dockerized
test-coverage-all-dockerized:
	$(info INFO: Starting build $@)
	./test.sh

.PHONY: test-coverage-all-dockerized-with-codeclimate
test-coverage-all-dockerized-with-codeclimate:
	$(info INFO: Starting build $@)
	CC_TEST_REPORTER_ID=${CC_TEST_REPORTER_ID} ./test.sh

test-coverage-all: export COMMANDER_TEST_ALL = 1
test-coverage-all: export COMMANDER_SSH_TEST = 1
test-coverage-all: export COMMANDER_TEST_SSH_HOST = 172.28.0.2:22
test-coverage-all: export COMMANDER_TEST_SSH_USER = root
test-coverage-all: export COMMANDER_TEST_SSH_IDENTITY_FILE = $(CWD)/integration/containers/ssh/id_rsa
.PHONY: test-coverage-all
test-coverage-all:
	$(info INFO: Starting build $@)
	go test -coverprofile c.out ./...

test-coverage-all: export COMMANDER_TEST_ALL = 1
test-coverage-all: export COMMANDER_SSH_TEST = 1
test-coverage-all: export COMMANDER_TEST_SSH_HOST = 172.28.0.2:22
test-coverage-all: export COMMANDER_TEST_SSH_USER = root
test-coverage-all: export COMMANDER_TEST_SSH_IDENTITY_FILE = $(CWD)/integration/containers/ssh/id_rsa
.PHONY: test-coverage-all-codeclimate
test-coverage-all-codeclimate:
	$(info INFO: Starting build $@)
	./test-reporter before-build
	go test -coverprofile c.out ./...; \
	./test-reporter after-build -t gocov --prefix=github.com/commander-cli/commander/v2 --exit-code $$?

.PHONY: integration-unix
integration-unix: build
	$(info INFO: Starting build $@)
	commander test commander_unix.yaml

.PHONY: integration-linux
integration-linux: build
	$(info INFO: Starting build $@)
	commander test integration/linux/docker.yaml
	DOCKER_HOST=${DOCKER_HOST} DOCKER_CERT_PATH=${DOCKER_CERT_PATH} commander test commander_unix.yaml
	DOCKER_HOST=${DOCKER_HOST} DOCKER_CERT_PATH=${DOCKER_CERT_PATH}  commander test commander_linux.yaml --verbose

.PHONY: integration-linux-dockerized
integration-linux-dockerized:
	$(info INFO: Starting build $@)
	./test.sh integration-linux

.PHONY: integration-windows
integration-windows: build
	$(info INFO: Starting build $@)
	commander test commander_windows.yaml

release-amd64:
	$(info INFO: Starting build $@)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-X main.version=$(GIT_RELEASE_TAG) -s -w" -o release/$(cmd)-linux-amd64 $(exe)

release-arm:
	$(info INFO: Starting build $@)
	CGO_ENABLED=0 GOOS=linux GOARCH=arm go build -ldflags "-X main.version=$(GIT_RELEASE_TAG) -s -w" -o release/$(cmd)-linux-arm $(exe)

release-386:
	$(info INFO: Starting build $@)
	CGO_ENABLED=0 GOOS=linux GOARCH=386 go build -ldflags "-X main.version=$(GIT_RELEASE_TAG) -s -w" -o release/$(cmd)-linux-386 $(exe)

release-darwin-amd64:
	$(info INFO: Starting build $@)
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags "-X main.version=$(GIT_RELEASE_TAG) -s -w" -o release/$(cmd)-darwin-amd64 $(exe)

release-windows-amd64:
	$(info INFO: Starting build $@)
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags "-X main.version=$(GIT_RELEASE_TAG) -s -w" -o release/$(cmd)-windows-amd64.exe $(exe)

release-windows-386:
	$(info INFO: Starting build $@)
	CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -ldflags "-X main.version=$(GIT_RELEASE_TAG) -s -w" -o release/$(cmd)-windows-386.exe $(exe)

release: release-amd64 release-arm release-386 release-darwin-amd64 release-windows-amd64 release-windows-386
