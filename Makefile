# Task runner for the repo. Run `make` on its own to list the targets.
#
# These are thin wrappers over the Go tool — `make build` is exactly
# `go run . -build`. The wrapping earns its keep for the targets Go alone
# cannot express (serve, test, clean) and for not having to remember which
# steps need a GITHUB_TOKEN and which do not.

DIST ?= dist
PORT ?= 8080

# Print the help text when make is run with no target.
.DEFAULT_GOAL := help

.PHONY: help update build check serve test fmt lint clean

help: ## Show this help
	@echo "Usage: make <target>"
	@echo
	@grep -hE '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
		| awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-8s\033[0m %s\n", $$1, $$2}'
	@echo
	@echo "Only 'update' needs network access and a GITHUB_TOKEN."

update: ## Fetch GitHub and refresh README, history and metadata (needs GITHUB_TOKEN)
	@[ -n "$$GITHUB_TOKEN" ] || { \
		echo "GITHUB_TOKEN is not set — see docs/LOCAL_DEV.md for how to get one."; \
		echo "If you only changed tags, notes or the dashboard, use 'make build' instead."; \
		exit 1; \
	}
	go run .

build: ## Render the site into dist/ from committed data (offline)
	go run . -build

check: ## Validate data/agents.yml (offline)
	go run . -check

serve: build ## Build, then preview the dashboard locally (override with PORT=)
	@command -v python3 >/dev/null || { \
		echo "python3 not found — serve $(DIST)/ with any static file server instead."; \
		exit 1; \
	}
	@echo "Serving $(DIST)/ at http://localhost:$(PORT) — Ctrl-C to stop"
	@python3 -m http.server $(PORT) -d $(DIST)

test: ## Run everything CI runs (vet, tests, check, build)
	go vet ./...
	go test ./...
	go run . -check
	go run . -build

fmt: ## Format the Go sources
	gofmt -w .

lint: ## Run golangci-lint if it is installed
	@command -v golangci-lint >/dev/null || { \
		echo "golangci-lint not installed — CI runs it regardless."; \
		echo "Install: https://golangci-lint.run/welcome/install/"; \
		exit 1; \
	}
	golangci-lint run

clean: ## Remove build output
	rm -rf $(DIST)
	rm -f awesome-ai-dev-tools
