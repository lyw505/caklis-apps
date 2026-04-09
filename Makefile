# CAKLI Makefile — Task Automation
# Project: CAKLI Web Admin v1

.PHONY: dev up down build logs seed clean setup-minio db-init

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
	$(MAKE) db-init
	$(MAKE) seed
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
	docker cp ./docs/agile-development/v1/seed-multi-role-admins.sql cakli-db:/seed.sql
	docker exec cakli-db psql -U cakli -d cakli_db -f /seed.sql
	docker exec cakli-db rm /seed.sql
	@echo "✅ Seeding complete"

# ── MinIO Setup ──

setup-minio:
	@echo "🔧 Configuring MinIO bucket policy..."
	powershell -ExecutionPolicy Bypass -File ./scripts/setup-minio.ps1

db-init:
	@echo "⏳ Waiting for database to be ready..."
	@powershell -Command "$$max=30; for($$i=0;$$i -lt $$max;$$i++){ docker exec cakli-db pg_isready -U cakli -d cakli_db -q; if($$LASTEXITCODE -eq 0){ break }; Start-Sleep -Seconds 1 }"
	@echo "🗄️ Initializing database schema (if needed)..."
	docker cp ./docs/agile-development/v1/db.sql cakli-db:/db.sql
	docker exec cakli-db sh -lc "psql -U cakli -d cakli_db -tAc \"SELECT to_regclass('public.admins') IS NOT NULL\" | grep -qi t || psql -U cakli -d cakli_db -v ON_ERROR_STOP=0 -f /db.sql"
	docker exec cakli-db rm /db.sql

# ── Cleanup ──

clean:
	rm -rf ./apps/api/tmp
	rm -rf ./apps/web/.next
	rm -rf ./apps/web/node_modules
	docker system prune -f
