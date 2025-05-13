# Makefile – multi-arch Linux & Windows + GitHub release

BINARY_NAME       := trading-bot
BUILD_DIR         := bin
PKG               := trading-bot
CLI               := cmd/cli/main.go

GO                := go
OS_LIST := linux windows darwin freebsd netbsd openbsd dragonfly solaris aix android illumos ios js wasip1 hurd plan9

# allow overriding on the command-line, e.g. make CODEVERSION=1.2.3
CODEVERSION       ?= $(shell ./scripts/next-version.sh)
BUILDDATE         := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
CODEBUILDREVISION := $(shell git rev-parse HEAD)

LDFLAGS := \
  -X $(PKG)/internal/interfaces/cli.CODEVERSION=$(CODEVERSION) \
  -X $(PKG)/internal/interfaces/cli.CODEBUILDDATE=$(BUILDDATE) \
  -X $(PKG)/internal/interfaces/cli.CODEBUILDREVISION=$(CODEBUILDREVISION)

.PHONY: all build-all release clean

all: build-all

build-all:
	@mkdir -p $(BUILD_DIR)
	@echo "Building $(BINARY_NAME) for: $(OS_LIST)"
	@for os in $(OS_LIST); do \
	  arches="$$( $(GO) tool dist list | grep ^$$os/ | cut -d/ -f2 )"; \
	  for arch in $$arches; do \
	    ext=""; if [ "$$os" = "windows" ]; then ext=".exe"; fi; \
	    out="$(BUILD_DIR)/$(BINARY_NAME)-$$os-$$arch$$ext"; \
	    echo " → $$os/$$arch"; \
	    GOOS=$$os GOARCH=$$arch $(GO) build -v -buildvcs=false \
	      -ldflags "$(LDFLAGS)" \
	      -o "$$out" "$(CLI)"; \
	  done; \
	done
	@echo "All builds complete."

release: build-all
	@echo "Creating GitHub release v$(CODEVERSION)…"
	@gh release create v$(CODEVERSION) $(BUILD_DIR)/* \
	  --title "v$(CODEVERSION)" \
	  --notes "Automated release of v$(CODEVERSION)"
	@echo "Release v$(CODEVERSION) created and artifacts uploaded."

clean:
	@rm -rf $(BUILD_DIR)
	@echo "Cleaned all build artifacts."