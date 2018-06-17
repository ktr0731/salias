VERSION := 0.1.0
REVISION := $(shell git rev-parse --short HEAD)

# Show version
version:
	@echo "Version: $(VERSION)($(REVISION))"

# Install dep
dep: 
ifeq ($(shell which dep 2>/dev/null),)
	curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
endif

# Install Go dependencies
deps: dep
	dep ensure

# Test 
test: deps
	go test -v
	
# Build app
build: deps test
	go build
