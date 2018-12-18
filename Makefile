VERSION := 0.1.0
REVISION := $(shell git rev-parse --short HEAD)

version:
	@echo "Version: $(VERSION)($(REVISION))"

test: 
	go test -v -race ./...

build: 
	go build ./...
