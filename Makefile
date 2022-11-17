GOFILES := $(shell find . -name "*.go")
GO ?= go
GOFMT ?= gofmt "-s"
GO_VERSION=$(shell $(GO) version | cut -c 14- | cut -d' ' -f1 | cut -d'.' -f2)
PACKAGES ?= $(shell $(GO) list ./...)

.PHONY: fmt
fmt:
	$(GOFMT) -w $(GOFILES)


.PHONY: run
run:
	make install; ~/go/bin/gtdbot


.PHONY: install
install:
	$(GO) install $(GOFILES)



