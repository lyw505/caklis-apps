# CAKLI Makefile — Task Automation
# Project: CAKLI Web Admin v1

.PHONY: dev up down build logs seed clean setup-minio

# Default: Show help
help:
	@echo "CAKLI Management Commands:"
	@echo "  make dev         - Start full dev stack (Docker deps + API Air + Frontend local)"
	@echo "  make up          - Start only Docker services (headless)"
	@echo "  make down        - Stop all Docker services"
	@echo "  make build       - Rebuild Docker images"
	@echo "  make logs        - Show follow logs for all Docker services"
	@echo "  make seed        - Run database seeders"
	@echo "  make setup-minio - Configure MinIO bucket policy"
	@echo "  make clean       - Remove build artifacts and temp files"

# ── Development ──

# make dev: Implement Section 4 of PRD
# 1. Start Docker services (postgres, minio, api with air)
# 2. Sync .env to web apps
# 3. Start Next.js frontend locally
dev:
	@echo "🚀 Starting CAKLI development stack..."
	docker compose up -d --build postgres minio backend-api
	@echo "📋 Syncing environment variables to web-admin..."
	powershell -Command "Copy-Item .env ./apps/web/.env.local -Force; (Get-Content ./apps/web/.env.local) -replace 'API_INTERNAL_URL=http://api:8080', 'API_INTERNAL_URL=http://localhost:8080' | Set-Content ./apps/web/.env.local"
	@echo "🌐 Starting Frontend (Next.js)..."
	cd apps/web && npm run dev

# ── Docker Management ──

up:
	docker compose up -d

down:
	docker compose down

build:
	docker compose build

logs:
	docker compose logs -f

# ── Database ──

seed:
	@echo "🌱 Seeding database..."
	docker compose exec backend-api ./api -seed
	@echo "✅ Seeding complete"

# ── MinIO Setup ──

setup-minio:
	@echo "🔧 Configuring MinIO bucket policy..."
	powershell -ExecutionPolicy Bypass -File ./scripts/setup-minio.ps1

# ── Cleanup ──

clean:
	rm -rf ./apps/api/tmp
	rm -rf ./apps/web/.next
	rm -rf ./apps/web/node_modules
	docker system prune -f
