.PHONY: build run test fmt

APP=azure-finops

build:
	go build -o bin/$(APP) ./cmd/azure-finops

run:
	go run ./cmd/azure-finops --help

fmt:
	go fmt ./...

test:
	go test ./...
