PROVIDER     ?= pexip
COMPANY      ?= pexip
DOMAIN       ?= pexip.com
NAME         := terraform-provider-$(PROVIDER)
ROOT_DIR     := $(if $(ROOT_DIR),$(ROOT_DIR),$(shell git rev-parse --show-toplevel))
BUILD_DIR    := $(ROOT_DIR)/dist
VERSION      ?= $(shell git describe --tags --always)
VERSION_NO_V := $(patsubst v%,%,$(VERSION))
GITSHA       := $(shell git rev-parse --short HEAD)

BUILD_TIME         := $(shell date +'%Y-%m-%d_%T')
GO_OS              := $(if $(GOOS),$(GOOS),$(shell go env GOOS))
GO_ARCH            := $(if $(GOARCH),$(GOARCH),$(shell go env GOARCH))
OS_ARCH            := $(GO_OS)_$(GO_ARCH)
GIT_BRANCH         :=$(shell git rev-parse --abbrev-ref HEAD)
GIT_REVISION       :=$(shell git rev-list -1 HEAD)
GIT_REVISION_DIRTY :=$(shell (git diff-index --quiet HEAD -- . && git diff --staged --quiet -- .) || echo "-dirty")

.PHONY: prepare lint build package install test testacc clean manifest

all: testacc build

prepare:
	@mkdir -p $(BUILD_DIR)

lint:
	$(GO_LINT_HEAD) $(GO_ENV_VARS) golangci-lint run

build: prepare
	@echo "Building $(NAME) version $(VERSION_NO_V) for $(OS_ARCH)..."
	@go build -ldflags "-X main.commit=$(GIT_BRANCH)@$(GIT_REVISION)$(GIT_REVISION_DIRTY) -X internal/version.appBuildTime=$(BUILD_TIME) -X internal/version.appVersion=$(VERSION_NO_V) -X internal/version.appBuildUser=${USER}" -o $(BUILD_DIR)/$(NAME)_v$(VERSION_NO_V) .

package: build
	@zip -j $(BUILD_DIR)/$(NAME)_$(VERSION_NO_V)_$(OS_ARCH).zip $(BUILD_DIR)/$(NAME)_v$(VERSION_NO_V)

install:
	@mkdir -p ~/.terraform.d/plugins/$(DOMAIN)/$(COMPANY)/$(PROVIDER)
	@unzip -o $(BUILD_DIR)/$(NAME)_$(VERSION)_$(OS_ARCH).zip -d ~/.terraform.d/plugins/$(DOMAIN)/$(COMPANY)/$(PROVIDER)

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