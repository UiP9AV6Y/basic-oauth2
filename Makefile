
GO ?= go
GIT ?= git
DOCKER ?= docker
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
CODE_ORIGIN ?= $(shell $(GIT) config --get remote.origin.url 2>/dev/null || grep module go.mod | cut -d' ' -f2)
COMMIT ?= $(shell $(GIT) rev-parse --short HEAD 2>/dev/null || echo HEAD)
BUILD_DATE ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
SOURCE_DATE_EPOCH ?= $(shell $(GIT) log -1 --format='%ct' 2>/dev/null || echo 0)

GO_MODULE := $(shell $(GO) list -m 2>/dev/null || grep module go.mod | cut -d' ' -f2)
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
DOCKER_REPOSITORY ?= uip9av6y/$(PROJECT_NAME)
DOCKER_TAG ?= latest
DOCKERFILE_PATH ?= Dockerfile
DOCKER_IMAGE ?= $(DOCKER_REGISTRY)/$(DOCKER_REPOSITORY)

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

.PHONY: docker-image
docker-image:
	$(DOCKER) build \
		--build-arg "BUILD_DATE=$(BUILD_DATE)" \
		--build-arg "VERSION=$(VERSION)" \
		--build-arg "VCS_REF=$(COMMIT)" \
		--build-arg "VCS_URL=$(CODE_ORIGIN)" \
		-f $(DOCKERFILE_PATH) \
		-t $(DOCKER_IMAGE):$(DOCKER_TAG) \
		.

.PHONY: docker-deploy
docker-deploy: DOCKER_DEPLOY_TAG := $(DOCKER_TAG)
docker-deploy:
	$(DOCKER) tag \
		$(DOCKER_IMAGE):$(DOCKER_TAG) \
		$(DOCKER_IMAGE):$(DOCKER_DEPLOY_TAG)
	$(DOCKER) push \
		$(DOCKER_IMAGE):$(DOCKER_DEPLOY_TAG)

.PHONY: docker-deploy-version
docker-deploy-version: VERSION_META := $(subst ., ,$(VERSION))
docker-deploy-version:
	$(MAKE) docker-deploy DOCKER_DEPLOY_TAG=$(word 1,$(VERSION_META))
	$(MAKE) docker-deploy DOCKER_DEPLOY_TAG=$(word 1,$(VERSION_META)).$(word 2,$(VERSION_META))

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
