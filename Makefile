build:
	$(info INFO: Starting build $@)
	go build cmd/commander/commander.go

test:
	$(info INFO: Starting build $@)
	go test ./...

test-integration: build
	$(info INFO: Starting build $@)
	./commander test