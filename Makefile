# ---- Config ----
APP_NAME    ?= url-shortener 
MAIN_PKG    ?= ./cmd/api      # path to your main package
PKG         ?= ./...
BUILD_DIR   ?= bin
GO          ?= go

# Versioning (from git)
GIT_TAG     := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
GIT_COMMIT  := $(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
BUILD_DATE  := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

# ldflags (change import path to match your version vars if you have them)
LDFLAGS     := -X main.version=$(GIT_TAG) -X main.commit=$(GIT_COMMIT) -X main.date=$(BUILD_DATE)
GOFLAGS     ?=

# Platform-aware binary extension
OS := $(shell uname -s)
EXT :=
ifeq ($(OS),Windows_NT)
  EXT := .exe
endif

BIN := $(BUILD_DIR)/$(APP_NAME)$(EXT)

# Environment for run
PORT ?= 8080
ENV  ?= development

# Tools (optional: install separately as needed)
GOLANGCI_LINT ?= golangci-lint
AIR           ?= air

# ---- Helpers ----
.PHONY: help
help: ## Show this help
	@grep -E '^[a-zA-Z0-9_\-\/]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-18s\033[0m %s\n", $$1, $$2}'

# ---- Dev tasks ----
.PHONY: run
run: ## Run the server (go run) with ENV/PORT
	ENV=$(ENV) PORT=$(PORT) $(GO) run $(GOFLAGS) $(MAIN_PKG)

.PHONY: watch
watch: ## Live-reload with air (if installed)
	@[ -x "$$(command -v $(AIR))" ] || (echo "air not found. Install: go install github.com/air-verse/air@latest" && exit 1)
	ENV=$(ENV) PORT=$(PORT) $(AIR)

.PHONY: build
build: $(BIN) ## Build the binary

$(BIN):
	@mkdir -p $(BUILD_DIR)
	$(GO) build $(GOFLAGS) -ldflags '$(LDFLAGS)' -o $(BIN) $(MAIN_PKG)

.PHONY: clean
clean: ## Remove build artifacts
	@rm -rf $(BUILD_DIR) coverage.out

# ---- Code quality ----
.PHONY: fmt
fmt: ## go fmt + goimports (if installed)
	$(GO) fmt $(PKG)
	@if command -v goimports >/dev/null 2>&1; then \
		goimports -w . ; \
	else \
		echo "goimports not found (optional): go install golang.org/x/tools/cmd/goimports@latest"; \
	fi

.PHONY: vet
vet: ## go vet
	$(GO) vet $(PKG)

.PHONY: lint
lint: ## Run golangci-lint (if installed)
	@[ -x "$$(command -v $(GOLANGCI_LINT))" ] || (echo "golangci-lint not found. Install: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest" && exit 1)
	$(GOLANGCI_LINT) run

# ---- Modules / deps ----
.PHONY: tidy
tidy: ## go mod tidy
	$(GO) mod tidy

.PHONY: vendor
vendor: ## vendor deps
	$(GO) mod vendor

# ---- Tests ----
.PHONY: test
test: ## Run tests (short)
	$(GO) test -race -vet=off -timeout=3m $(PKG)

.PHONY: test-cover
test-cover: ## Run tests with coverage
	$(GO) test -race -covermode=atomic -coverprofile=coverage.out $(PKG)
	@$(GO) tool cover -func=coverage.out | tail -n 1

.PHONY: bench
bench: ## Run benchmarks
	$(GO) test -bench=. -benchmem $(PKG)

# ---- Docker ----
IMAGE ?= $(APP_NAME):$(GIT_TAG)

.PHONY: docker-build
docker-build: ## Build Docker image
	docker build --build-arg GIT_TAG=$(GIT_TAG) --build-arg GIT_COMMIT=$(GIT_COMMIT) --build-arg BUILD_DATE=$(BUILD_DATE) -t $(IMAGE) .

.PHONY: docker-run
docker-run: ## Run Docker image mapping $(PORT)
	docker run --rm -e ENV=$(ENV) -e PORT=$(PORT) -p $(PORT):$(PORT) $(IMAGE)

# ---- Misc ----
.PHONY: info
info: ## Print build info
	@echo "App:        $(APP_NAME)"
	@echo "Main pkg:   $(MAIN_PKG)"
	@echo "Binary:     $(BIN)"
	@echo "Git tag:    $(GIT_TAG)"
	@echo "Commit:     $(GIT_COMMIT)"
	@echo "Build date: $(BUILD_DATE)"

