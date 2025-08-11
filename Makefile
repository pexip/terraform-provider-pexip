NAME       := terraform-provider-pexip
ROOT_DIR   := $(if $(ROOT_DIR),$(ROOT_DIR),$(shell git rev-parse --show-toplevel))
BUILD_DIR  := $(ROOT_DIR)/dist
VERSION    ?= 0.0.1
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
	go build -ldflags "-X main.commit=$(GIT_BRANCH)@$(GIT_REVISION)$(GIT_REVISION_DIRTY) -X internal/version.appBuildTime=$(BUILD_TIME) -X internal/version.appVersion=$(VERSION) -X internal/version.appBuildUser=${USER}" -o $(BUILD_DIR)/$(NAME)_$(VERSION) .
	chmod +x $(BUILD_DIR)/$(NAME)_$(VERSION)
	zip -j $(BUILD_DIR)/$(NAME)_$(VERSION)_$(OS_ARCH).zip $(BUILD_DIR)/$(NAME)_$(VERSION)

install:
	mkdir -p ~/.terraform.d/plugins/pexip.com/pexip/pexip
	unzip -o $(BUILD_DIR)/$(NAME)_$(VERSION)_$(OS_ARCH).zip -d ~/.terraform.d/plugins/pexip.com/pexip/pexip

test: prepare
	go test -v -parallel 4 -tags unit -coverprofile=$(BUILD_DIR)/cover.out ./...

testacc: export TF_ACC=true
testacc: prepare
	go test -v -tags integration -coverprofile=$(BUILD_DIR)/cover.out ./...

fmt:
	go fmt ./...

clean:
	rm -rf $(BUILD_DIR)
	rm -rf ~/.terraform.d/plugins/pexip.com/pexip/pexip

.PHONY: generate-manifest
generate-manifest: prepare
	cd $(BUILD_DIR) && \
	VERSION_VALUE="$(VERSION)"; \
	VERSION_VALUE="$${VERSION_VALUE#v}"; \
	PROVIDER_VALUE="$(PROVIDER_NAME)"; \
	MANIFEST="terraform-provider-$${PROVIDER_VALUE}_$${VERSION_VALUE}_manifest.json"; \
	PLATFORMS_JSON="["; \
	first=1; \
	for f in terraform-provider-$${PROVIDER_VALUE}_$${VERSION_VALUE}_*.zip; do \
	  sha=$$(shasum -a 256 "$$f" | awk '{print $$1}'); \
	  base=$${f%.zip}; \
	  os_arch=$${base#terraform-provider-$${PROVIDER_VALUE}_$${VERSION_VALUE}_}; \
	  os=$${os_arch%_*}; \
	  arch=$${os_arch#$$os_}; \
	  entry=$$(jq -n --arg os "$$os" --arg arch "$$arch" --arg filename "$$f" --arg shasum "$$sha" '{os:$os, arch:$arch, filename:$filename, shasum:$shasum}'); \
	  if [ $$first -eq 1 ]; then PLATFORMS_JSON="[$$entry"; first=0; else PLATFORMS_JSON="$$PLATFORMS_JSON,$$entry"; fi; \
	done; \
	PLATFORMS_JSON="$$PLATFORMS_JSON]"; \
	jq -n --arg version "$$VERSION_VALUE" --argjson platforms "$$PLATFORMS_JSON" '{version:$version, protocols:["5.0"], platforms:$platforms}' > "$$MANIFEST"
