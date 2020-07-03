BIN=dots
SRC=$(shell find . -name *.go)

# Use richgo for testing if installed, otherwise fallback to normal go
ifeq (, $(shell which richgo))
	GO_TEST := go
else
	GO_TEST := richgo
endif

.PHONY: install_deps test fmt vet clean

default: all

all: install_deps vet fmt test

vet:
	go vet ./...

test:
	$(GO_TEST) test -v ./...

install_deps:
	go get -v ./...

fmt:
	@test -z $(shell gofmt -l $(SRC)) || (gofmt -d $(SRC); exit 1)

clean:
	rm -rf $(BIN)
