.DEFAULT_GOAL = test
.PHONY: FORCE

build: fmt vet chip8
.PHONY: build

fmt:
	go fmt ./...
.PHONY: fmt

vet:
	go vet ./...
.PHONY: vet

test: build
	go test ./...
.PHONY: test

lint:
	golangci-lint run
.PHONY: lint

chip8: FORCE
	go build ./...
.PHONY: chip8

go.mod: FORCE
	go mod tidy
	go mod verify
go.sum: go.mod
