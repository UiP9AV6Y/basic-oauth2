
GO ?= go
GOLINT ?= golangci-lint
GOFMT ?= goimports
INSTALL ?= install

ifeq ($(VERSION),)
VCS_META := $(subst -, ,$(shell $(GIT) describe --abbrev --tags 2>/dev/null || echo 0.0.0-1-HEAD))
VERSION := $(word 1,$(VCS_META))
ifneq ($(word 2,$(VCS_META)),0)
VERSION := $(VERSION)-dev
endif
endif
COMMIT ?= $(shell $(GIT) rev-parse --short HEAD 2>/dev/null || echo HEAD)
BUILD_DATE ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
SOURCE_DATE_EPOCH ?= $(shell $(GIT) log -1 --format='%ct' 2>/dev/null || echo 0)

GO_MODULE := $(shell $(GO) list -m)
GO_SOURCES := $(shell find . -path './.*' -prune -o -name '*.go' -print)
GO_CMDS := $(notdir $(wildcard ./cmd/*))

GO_LDFLAGS := '-extldflags "-static"
GO_LDFLAGS += -X $(GO_MODULE)/pkg/version.version=$(VERSION)
GO_LDFLAGS += -X $(GO_MODULE)/pkg/version.commit=$(COMMIT)
GO_LDFLAGS += -X $(GO_MODULE)/pkg/version.date=$(BUILD_DATE)
GO_LDFLAGS += -w -s # Drop debugging symbols.
GO_LDFLAGS += -buildid= # Reproducable builds
GO_LDFLAGS += '

PROJECT_NAME ?= $(notdir $(GO_MODULE))
BUILD_DIR ?= out

ifeq ($(GOOS),windows)
PROGRAMS := $(addsuffix .exe,$(GO_CMDS))
else
PROGRAMS := $(GO_CMDS)
endif

PREFIX ?= /usr
DOCKER_REGISTRY ?= docker.io
DOCKER_REPOSITORY ?= UiP9AV6Y/$(PROJECT_NAME)

.PHONY: default
default: all

.PHONY: all
all: lint build

.PHONY: clean
clean:
	$(RM) -f $(BUILD_DIR)

.PHONY: lint
lint: $(GO_SOURCES)
	$(GOLINT) run -v ./...

.PHONY: format
format: $(GO_SOURCES)
	$(GOFMT) -w $(GO_SOURCES)

.PHONY: test
test: $(GO_SOURCES)
	$(GO) test ./...

.PHONY: build
build: $(addprefix $(BUILD_DIR)/,$(GO_CMDS))

.PHONY: build-deps
build-deps:
	# go install -v -tags tools ./...
	$(GO) list -tags tools -f '{{range .Imports}}{{ . }} {{end}}' ./tools \
	| xargs -n1 $(GO) install -v

$(BUILD_DIR)/%: $(GO_SOURCES)
	GO111MODULE=on CGO_ENABLED=0 $(GO) build \
		-ldflags=$(GO_LDFLAGS) \
		-o $@ ./cmd/$*

.PHONY: install
install: build
	$(INSTALL) -d -m 755 $(DESTDIR)$(PREFIX)/bin
	$(INSTALL) -m 755 -t $(DESTDIR)$(PREFIX)/bin/ $(BUILD_DIR)/*
