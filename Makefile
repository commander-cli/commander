build:
	$(info INFO: Starting build $@)
	go build cmd/commander/commander.go

lint:
	$(info INFO: Starting build $@)
	golint pkg/ cmd/

test:
	$(info INFO: Starting build $@)
	go test ./...

test-coverage:
	$(info INFO: Starting build $@)
	go test -coverprofile coverage.out ./...

test-integration: build
	$(info INFO: Starting build $@)
	./commander test
