GOPATH ?= $(shell go env GOPATH)
GOBIN ?= $(or $(shell go env GOBIN),$(GOPATH)/bin)
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)


BUILD_PACKAGE = $(REPOPATH)/cmd/skaffold

GO_FILES = $(shell find . -type f -name "*.go" -not -path "./pkg/diag/*")

VERSION_PACKAGE = $(REPOPATH)/pkg/version
COMMIT = $(shell git rev-parse HEAD)

ifeq "$(strip $(VERSION))" ""
	override VERSION = $(shell git describe --always --tags --dirty)
endif

DATE_FMT = +%Y-%m-%dT%H:%M:%SZ
ifdef SOURCE_DATE_EPOCH
    BUILD_DATE ?= $(shell date -u -d "@$(SOURCE_DATE_EPOCH)" "$(DATE_FMT)" 2>/dev/null || date -u -r "$(SOURCE_DATE_EPOCH)" "$(DATE_FMT)" 2>/dev/null || date -u "$(DATE_FMT)")
else
    BUILD_DATE ?= $(shell date "$(DATE_FMT)")
endif

.PHONY: $(BUILD_DIR)/VERSION
$(BUILD_DIR)/VERSION: $(BUILD_DIR)
		@ echo $(VERSION) > $@

$(BUILD_DIR):
	mkdir -p $(BUILD_DIR)

.PHONY: fmt
fmt: $(BUILD_DIR)
		@ ./hack/linters.sh

