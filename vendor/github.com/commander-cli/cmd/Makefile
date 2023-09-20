
.PHONY: deps lint test test-coverage

init: git-hooks

git-hooks:
	$(info INFO: Starting build $@)
	ln -sf ../../.githooks/pre-commit .git/hooks/pre-commit

deps:
	$(info INFO: Starting build $@)
	go mod vendor

test:
	$(info INFO: Starting build $@)
	go test `go list ./... | grep -v examples`

test-coverage:
	$(info INFO: Starting build $@)
	go test -coverprofile c.out `go list ./... | grep -v examples`
