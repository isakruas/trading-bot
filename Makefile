# Makefile for trading-bot CLI
# Cross-platform: Linux/macOS & Windows (CMD)

# ——————————————————————————————
# CONFIGURATION
# ——————————————————————————————

BINARY_NAME        := trading-bot
BUILD_DIR          := bin

# Windows needs .exe
ifeq ($(OS),Windows_NT)
  EXT := .exe
else
  EXT :=
endif

BINARY             := $(BUILD_DIR)/$(BINARY_NAME)$(EXT)
PKG                := trading-bot
CLI                := cmd/cli/main.go

GO                 := go
GOOS               := $(shell go env GOOS)
GOARCH             := $(shell go env GOARCH)

# Build date in UTC
ifeq ($(OS),Windows_NT)
  # Powershell call; assuming PowerShell is on the PATH
  BUILDDATE := $(shell powershell -Command "Get-Date -Format yyyy-MM-ddTHH:mm:ssZ")
else
  BUILDDATE := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
endif

CODEVERSION        := 0.0.1

ifeq ($(OS),Windows_NT)
  CODEBUILDREVISION := $(shell powershell -Command "git rev-parse HEAD")
else
  CODEBUILDREVISION := $(shell git rev-parse HEAD)
endif

# Platform‐specific helpers
ifeq ($(OS),Windows_NT)
  MKDIR := if not exist "$(BUILD_DIR)" mkdir "$(BUILD_DIR)"
  RMDIR := if exist "$(BUILD_DIR)" rmdir /S /Q "$(BUILD_DIR)"
  GOENV := set GOOS=$(GOOS)&& set GOARCH=$(GOARCH)&&
else
  MKDIR := mkdir -p $(BUILD_DIR)
  RMDIR := rm -rf $(BUILD_DIR)
  GOENV := GOOS=$(GOOS) GOARCH=$(GOARCH)
endif

LDFLAGS := \
  -X $(PKG)/internal/interfaces/cli.CODEVERSION=$(CODEVERSION) \
  -X $(PKG)/internal/interfaces/cli.CODEBUILDDATE=$(BUILDDATE) \
  -X $(PKG)/internal/interfaces/cli.CODEBUILDREVISION=$(CODEBUILDREVISION)

.PHONY: all build

all: build

build:
	@$(MKDIR)
	@echo Building $(BINARY_NAME) for $(GOOS)/$(GOARCH) version $(CODEVERSION)...
	@$(GOENV) $(GO) build -v -buildvcs=false -ldflags "$(LDFLAGS)" \
	  -o "$(BINARY)" "$(CLI)"
	@echo Built: $(BINARY)
	@echo Build complete.