GOFILES := $(shell find . -name "*.go")
GO ?= go
GOFMT ?= gofmt "-s"
GO_VERSION=$(shell $(GO) version | cut -c 14- | cut -d' ' -f1 | cut -d'.' -f2)
PACKAGES ?= $(shell $(GO) list ./...)
COMPILED_EXEC_FILE = "./gtdbot"

.PHONY: fmt
fmt:
	$(GOFMT) -w $(GOFILES)


.PHONY: run
run:
	make build; $(COMPILED_EXEC_FILE)


.PHONY: install
install:
	$(GO) install $(GOFILES)


.PHONY: build
build:
	$(GO) build
