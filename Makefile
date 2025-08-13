# Provider and project information
PROVIDER     ?= pexip
COMPANY      ?= pexip
DOMAIN       ?= pexip.com
NAME         := terraform-provider-$(PROVIDER)
ROOT_DIR     := $(if $(ROOT_DIR),$(ROOT_DIR),$(shell git rev-parse --show-toplevel))
BUILD_DIR    := $(ROOT_DIR)/dist

# Version information
_GIT_VERSION := $(shell git describe --tags --always)
_VERSION_FROM_GIT := $(shell \
	if echo "$(_GIT_VERSION)" | grep -qE '^v[0-9]+\.[0-9]+\.[0-9]+(-[a-zA-Z0-9.-]+)?$$'; then \
		echo "$(_GIT_VERSION)"; \
	else \
		echo "v0.0.1-$(_GIT_VERSION)"; \
	fi)
VERSION      ?= $(_VERSION_FROM_GIT)
VERSION_NO_V := $(patsubst v%,%,$(VERSION))
GITSHA       := $(shell git rev-parse --short HEAD)

# Build information
BUILD_TIME         := $(shell date +'%Y-%m-%d_%T')
GO_OS              := $(if $(GOOS),$(GOOS),$(shell go env GOOS))
GO_ARCH            := $(if $(GOARCH),$(GOARCH),$(shell go env GOARCH))
OS_ARCH            := $(GO_OS)_$(GO_ARCH)
GIT_BRANCH         :=$(shell git rev-parse --abbrev-ref HEAD)
GIT_REVISION       :=$(shell git rev-list -1 HEAD)
GIT_REVISION_DIRTY :=$(shell (git diff-index --quiet HEAD -- . && git diff --staged --quiet -- .) || echo "-dirty")

# Binary and archive names
BINARY_NAME ?= $(NAME)_v$(VERSION_NO_V)
ZIP_ARCHIVE_NAME ?= $(NAME)_$(VERSION_NO_V)_$(OS_ARCH).zip

# Build flags
BUILD_LD_FLAGS := "-X main.commit=$(GIT_BRANCH)@$(GIT_REVISION)$(GIT_REVISION_DIRTY) -X internal/version.appBuildTime=$(BUILD_TIME) -X internal/version.appVersion=$(VERSION_NO_V) -X internal/version.appBuildUser=${USER}"

.PHONY: prepare lint build package install test testacc clean manifest

all: build

prepare:
	@mkdir -p $(BUILD_DIR)

lint:
	$(GO_LINT_HEAD) $(GO_ENV_VARS) golangci-lint run

build: prepare
	@echo "Building $(NAME) version $(VERSION_NO_V) for $(OS_ARCH)..."
	@go build -ldflags $(BUILD_LD_FLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) .

package: build
	@zip -j $(BUILD_DIR)/$(ZIP_ARCHIVE_NAME) $(BUILD_DIR)/$(BINARY_NAME)

install:
	@mkdir -p ~/.terraform.d/plugins/$(DOMAIN)/$(COMPANY)/$(PROVIDER)
	@unzip -o $(BUILD_DIR)/$(ZIP_ARCHIVE_NAME) -d ~/.terraform.d/plugins/$(DOMAIN)/$(COMPANY)/$(PROVIDER)

install2: install
	@mkdir -p ~/.terraform.d/plugins/$(DOMAIN)/$(COMPANY)/$(PROVIDER)/$(VERSION_NO_V)/$(OS_ARCH)
	@cp ~/.terraform.d/plugins/$(DOMAIN)/$(COMPANY)/$(PROVIDER)/$(BINARY_NAME) ~/.terraform.d/plugins/$(DOMAIN)/$(COMPANY)/$(PROVIDER)/$(VERSION_NO_V)/$(OS_ARCH)/$(NAME)

test: prepare
	go test -v -parallel 4 -tags unit -coverprofile=$(BUILD_DIR)/cover.out ./...

testacc: export TF_ACC=true
testacc: prepare
	go test -v -tags integration -coverprofile=$(BUILD_DIR)/cover.out ./...

clean:
	rm -rf $(BUILD_DIR)
	rm -rf ~/.terraform.d/plugins/$(DOMAIN)/$(COMPANY)/$(PROVIDER)

manifest: prepare
	@./generate-manifest.sh $(BUILD_DIR) $(VERSION) $(PROVIDER)

version:
	@echo $(VERSION)