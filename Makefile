NAME       := terraform-provider-pexip
ROOT_DIR   := $(if $(ROOT_DIR),$(ROOT_DIR),$(shell git rev-parse --show-toplevel))
BUILD_DIR  := $(ROOT_DIR)/dist
VERSION    := $(shell cat VERSION)
GITSHA     := $(shell git rev-parse --short HEAD)


BUILD_TIME         := $(shell date +'%Y-%m-%d_%T')
GO_OS              := $(if $(GOHOSTOS),$(GOHOSTOS),$(shell go env GOHOSTOS))
GO_ARCH            := $(if $(GOHOSTARCH),$(GOHOSTARCH),$(shell go env GOHOSTARCH))
OS_ARCH            := $(GO_OS)_$(GO_ARCH)
GIT_BRANCH         :=$(shell git rev-parse --abbrev-ref HEAD)
GIT_REVISION       :=$(shell git rev-list -1 HEAD)
GIT_REVISION_DIRTY :=$(shell (git diff-index --quiet HEAD -- . && git diff --staged --quiet -- .) || echo "-dirty")

.PHONY: prepare lint check build-dev build install test testacc fmt clean

all: testacc build

prepare:
	mkdir -p $(BUILD_DIR)

lint:
	$(GO_LINT_HEAD) $(GO_ENV_VARS) golangci-lint run

check: lint test

build-dev:
	go build -ldflags "-X main.commit=$(GIT_BRANCH)@$(GIT_REVISION)$(GIT_REVISION_DIRTY) -X internal/version.appBuildTime=$(BUILD_TIME) -X internal/version.appVersion=$(VERSION) -X internal/version.appBuildUser=${USER}" -o ~/.terraform.d/plugins/$(NAME)_$(VERSION) .

build: prepare
	go build -ldflags "-X main.commit=$(GIT_BRANCH)@$(GIT_REVISION)$(GIT_REVISION_DIRTY) -X internal/version.appBuildTime=$(BUILD_TIME) -X internal/version.appVersionn=$(VERSION) -X internal/version.appBuildUser=${USER}" -o $(BUILD_DIR)/$(NAME)_$(VERSION) .

install:
	mkdir -p ~/.terraform.d/plugins/pexip.com/pexip/pexip/$(VERSION)/$(OS_ARCH)
	mv $(BUILD_DIR)/$(NAME)_$(VERSION) ~/.terraform.d/plugins/pexip.com/pexip/pexip/$(VERSION)/$(OS_ARCH)/$(NAME)
	chmod +x ~/.terraform.d/plugins/pexip.com/pexip/pexip/$(VERSION)/$(OS_ARCH)/$(NAME)
	cp ~/.terraform.d/plugins/pexip.com/pexip/pexip/$(VERSION)/$(OS_ARCH)/$(NAME) ~/.terraform.d/plugins/pexip.com/pexip/pexip

test: prepare
	go test -v -tags unit -coverprofile=$(BUILD_DIR)/cover.out ./...

testacc: export TF_ACC=true
testacc: prepare
	go test -v -tags integration -coverprofile=$(BUILD_DIR)/cover.out ./...

fmt:
	go fmt ./...

clean:
	rm -rf $(BUILD_DIR)
	rm -rf ~/.terraform.d/plugins/pexip.com/pexip/pexip
