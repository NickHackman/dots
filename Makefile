BIN=dots
SRC=$(shell find . -name *.go)

# Use richgo for testing if installed, otherwise fallback to normal go
ifeq (, $(shell which richgo))
	GO_TEST := go
else
	GO_TEST := richgo
endif

.PHONY: install_deps test fmt vet clean run build

default: all

all: install_deps vet fmt test

vet:
	$(info --------------------- vetting ---------------------)
	go vet ./...

test:
	$(info --------------------- testing ---------------------)
	$(GO_TEST) test -v ./...

install_deps:
	$(info --------------------- downloading dependencies ---------------------)
	go get -v ./...

fmt:
	$(info --------------------- checking formatting ---------------------)
	@test -z $(shell gofmt -l $(SRC)) || (gofmt -d $(SRC); exit 1)

run:
	$(info --------------------- running ---------------------)
	go run .

build:
	$(info --------------------- building ---------------------)
	go build ./...

clean:
	rm -rf $(BIN)
