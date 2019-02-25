build:
	$(info INFO: Starting build $@)
	go build cmd/commander/commander.go

lint:
	$(info INFO: Starting build $@)
	golint pkg/ cmd/

test:
	$(info INFO: Starting build $@)
	go test ./...

test-integration: build
	$(info INFO: Starting build $@)
	./commander test
