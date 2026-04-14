# Makefile — shortcuts for common development tasks
#
# Usage: make <target>
# Run `make help` to see all available commands

# Variables
APP_NAME = urlshortener
BINARY   = server
CMD_PATH = ./cmd/main.go

# ── Help ───────────────────────────────────────────────────────────────────
.PHONY: help
help: ## Show this help message
	@echo ""
	@echo "URL Shortener — Available Commands"
	@echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'
	@echo ""

# ── Local Development ─────────────────────────────────────────────────────
.PHONY: run
run: ## Run the server locally (requires local PostgreSQL)
	go run $(CMD_PATH)

.PHONY: build
build: ## Compile the binary to ./server
	CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o $(BINARY) $(CMD_PATH)
	@echo "✅ Binary built: ./$(BINARY)"

.PHONY: deps
deps: ## Download and tidy Go dependencies
	go mod tidy
	@echo "✅ Dependencies ready"

# ── Docker ────────────────────────────────────────────────────────────────
.PHONY: docker-up
docker-up: ## Build and start the full Docker stack (app + postgres + pgadmin)
	docker compose up --build

.PHONY: docker-up-d
docker-up-d: ## Start Docker stack in background (detached)
	docker compose up --build -d
	@echo "✅ Stack running in background"
	@echo "   API:     http://localhost:3000"
	@echo "   pgAdmin: http://localhost:5050"

.PHONY: docker-down
docker-down: ## Stop and remove Docker containers
	docker compose down

.PHONY: docker-down-v
docker-down-v: ## Stop containers AND delete volumes (wipes the database!)
	docker compose down -v
	@echo "⚠️  All volumes deleted — database is wiped"

.PHONY: docker-logs
docker-logs: ## Tail logs from the app container
	docker compose logs -f app

.PHONY: docker-ps
docker-ps: ## Show status of all containers
	docker compose ps

.PHONY: docker-restart
docker-restart: ## Restart only the app container (not postgres)
	docker compose restart app

# ── Database ──────────────────────────────────────────────────────────────
.PHONY: psql
psql: ## Open a psql shell in the running postgres container
	docker compose exec postgres psql -U postgres -d urlshortener

# ── Testing the API ───────────────────────────────────────────────────────
.PHONY: test-create
test-create: ## Test creating a short URL
	curl -s -X POST http://localhost:3000/api/urls/ \
		-H "Content-Type: application/json" \
		-d '{"original_url": "https://google.com", "custom_code": "goog"}' | jq .

.PHONY: test-list
test-list: ## Test listing all URLs
	curl -s http://localhost:3000/api/urls/ | jq .

.PHONY: test-stats
test-stats: ## Test getting stats for "goog"
	curl -s http://localhost:3000/api/urls/goog | jq .

.PHONY: test-redirect
test-redirect: ## Test the redirect for "goog" (follows redirect)
	curl -Ls -o /dev/null -w "Redirected to: %{url_effective}\n" http://localhost:3000/goog

.PHONY: test-delete
test-delete: ## Test deleting "goog"
	curl -s -X DELETE http://localhost:3000/api/urls/goog | jq .

.PHONY: test-health
test-health: ## Test the health check endpoint
	curl -s http://localhost:3000/health | jq .

.PHONY: test-all
test-all: ## Run all API tests in sequence
	@echo "\n── Health Check ──────────────────"
	@make test-health
	@echo "\n── Create Short URL ──────────────"
	@make test-create
	@echo "\n── List All URLs ─────────────────"
	@make test-list
	@echo "\n── Get Stats ─────────────────────"
	@make test-stats
	@echo "\n── Test Redirect ─────────────────"
	@make test-redirect
	@echo "\n── Delete URL ────────────────────"
	@make test-delete

# ── Cleanup ───────────────────────────────────────────────────────────────
.PHONY: clean
clean: ## Remove compiled binary
	rm -f $(BINARY)
	@echo "✅ Cleaned"