VERSION := 0.1.0
REVISION := $(shell git rev-parse --short HEAD)

# Show version
version:
	@echo "Version: $(VERSION)($(REVISION))"

# Install Glide
glide: 
ifeq ($(shell which glide 2>/dev/null),)
	mkdir -p $(GOPATH)/bin
	curl -s https://glide.sh/get | sh
endif

# Install Go dependencies
deps: glide
	glide install

# Test 
test: deps
	go test -v
	
# Build app
build: deps test
	go build
