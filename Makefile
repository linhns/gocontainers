.DEFAULT_GOAL := help

# ====================== #
# HELPERS
# ====================== #

.PHONY: help confirm no-dirty
## help: display help message
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

confirm:
	@echo "Are you sure? [y/N] " && read ans && [ $${ans:-N} = y ]

no-dirty:
	@test -z "$(shell git status --porcelain)"

# ====================== #
# DEVELOPMENT
# ====================== #

.PHONY: build tidy clean
## tidy: tidy go modules
tidy:
	go mod tidy -v
	go fmt ./...

## build: build the binary
build:
	go build ./...

## clean: remove build artifacts
clean:
	go clean

# ====================== #
# QC
# ====================== #

.PHONY: test audit cover vet lint check-vulnerabilities
## test: run tests
test:
	go test -v -race -buildvcs=true ./...

verify-deps:
	go mod tidy -diff
	go mod verify

check-format:
	@test -z "$(shell gofmt -l .)"

## vet: check common mistakes
vet:
	go vet ./...

lint:
	go run honnef.co/go/tools/cmd/staticcheck@latest ./...
	go run github.com/mgechev/revive@latest ./...

check-vulnerabilities:
	go run golang.org/x/vuln/cmd/govulncheck@latest ./...

## audit: run quality control checks
audit: test verify-deps check-format vet

## cover: run code coverage
cover:
	go test -v -race -buildvcs -coverprofile=/tmp/coverage.out ./...
	go tool cover -html=/tmp/coverage.out

# ====================== #
# DEPLOY
# ====================== #

.PHONY: push
## push: push changes to remote
push: confirm audit no-dirty
	git push
