# CAKLI Makefile — Task Automation
# Project: CAKLI Web Admin v1

.PHONY: dev dev-local up down build logs seed db-seed clean setup-minio fix-upload db-init

# Default: Show help
help:
	@echo "CAKLI Management Commands:"
	@echo "  make dev         - Start full dev stack (Docker + MinIO setup + Frontend)"
	@echo "  make dev-local   - Start Docker services only (run backend manually)"
	@echo "  make up          - Start only Docker services (headless)"
	@echo "  make down        - Stop all Docker services"
	@echo "  make build       - Rebuild Docker images"
	@echo "  make logs        - Show follow logs for all Docker services"
	@echo "  make db-init     - Initialize database schema (run once)"
	@echo "  make db-seed     - Seed database with default data (admins, banks)"
	@echo "  make setup-minio - Configure MinIO bucket policy (first time)"
	@echo "  make fix-upload  - Fix upload after restart (run after every restart)"
	@echo "  make clean       - Remove build artifacts and temp files"

# ── Development ──

# make dev: Start development environment
# 1. Start Docker services (postgres, minio, backend-api)
# 2. Wait for services to be ready
# 3. Initialize database schema (if needed)
# 4. Configure MinIO bucket policy (auto-fix upload)
# 5. Sync .env to web apps
# 6. Start Next.js frontend locally
dev:
	@echo "🚀 Starting CAKLI development stack..."
	docker compose up -d --build postgres minio backend-api
	@echo "⏳ Waiting for services to be ready..."
	@powershell -Command "& { $$count = 0; while ($$count -lt 60) { $$status = (docker inspect -f '{{.State.Running}}' cakli-storage); if ($$status -eq 'true') { break }; Start-Sleep -Seconds 1; $$count++ }; if ($$status -ne 'true') { Write-Error 'cakli-storage container did not start in time.' } }"
	@powershell -Command "& { $$count = 0; while ($$count -lt 30) { docker exec cakli-db pg_isready -U cakli -d cakli_db -q; if ($$LASTEXITCODE -eq 0){ break }; Start-Sleep -Seconds 1 }; $$count++ }"
	$(MAKE) db-init
	@echo "🔧 Configuring MinIO bucket policy..."
	@powershell -Command "docker exec cakli-storage mc alias set local http://localhost:9000 minioadmin minioadmin | Out-Null"
	@powershell -Command "docker exec cakli-storage mc mb local/cakli --ignore-existing | Out-Null"
	@powershell -Command "docker exec cakli-storage mc anonymous set public local/cakli | Out-Null"
	@echo "✅ MinIO configured!"
	@echo "📋 Syncing environment variables to web-admin..."
	@powershell -Command "Copy-Item .env ./apps/web/.env.local -Force; (Get-Content ./apps/web/.env.local) -replace 'API_INTERNAL_URL=http://api:8080', 'API_INTERNAL_URL=http://localhost:8080' | Set-Content ./apps/web/.env.local"
	@echo ""
	@echo "✅ Development stack ready!"
	@echo "💡 Tip: Run 'make db-seed' to seed database with default data"
	@echo ""
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

db-init:
	@echo "⏳ Waiting for database container to be ready..."
	@powershell -Command "$$count = 0; while ($$count -lt 30) { $$status = (docker inspect -f '{{.State.Running}}' cakli-db); if ($$status -eq 'true') { break }; Start-Sleep -Seconds 1; $$count++ }; if ($$status -ne 'true') { Write-Error 'cakli-db container did not start in time.' }"
	@echo "⏳ Waiting for database to be ready..."
	@powershell -Command "$$count = 0; while ($$count -lt 30) { docker exec cakli-db pg_isready -U cakli -d cakli_db -q; if ($$LASTEXITCODE -eq 0) { break }; Start-Sleep -Seconds 1; $$count++ }"
	@echo "🗄️ Initializing database schema (if needed)..."
	docker cp ./docs/agile-development/v1/db.sql cakli-db:/db.sql
	docker exec cakli-db sh -lc "psql -U cakli -d cakli_db -tAc \"SELECT to_regclass('public.admins') IS NOT NULL\" | grep -qi t || psql -U cakli -d cakli_db -v ON_ERROR_STOP=0 -f /db.sql"
	docker exec cakli-db rm /db.sql
	@echo "✅ Database schema initialized"
	@echo "✅ Database schema initialized"

db-seed:
	@echo "🌱 Seeding database with default data..."
	@echo "   - 3 admin accounts (master, operation, reporting)"
	@echo "   - 22 Indonesian banks"
	docker cp ./docs/agile-development/v1/seed-multi-role-admins.sql cakli-db:/seed.sql
	docker exec cakli-db psql -U cakli -d cakli_db -f /seed.sql
	docker exec cakli-db rm /seed.sql
	@echo "✅ Database seeded successfully!"
	@echo ""
	@echo "📋 Default Admin Accounts:"
	@echo "   Master Admin:    master@cakli.id / admin123"
	@echo "   Operation Admin: operation@cakli.id / admin123"
	@echo "   Reporting Admin: reporting@cakli.id / admin123"

# Legacy alias for backward compatibility
seed: db-seed

# ── MinIO Setup ──

setup-minio:
	@echo "🔧 Configuring MinIO bucket policy..."
	powershell -ExecutionPolicy Bypass -File ./scripts/setup-minio.ps1

fix-upload:
	@echo "🔧 Fixing upload issue (setting bucket policy)..."
	@docker exec cakli-storage mc alias set local http://localhost:9000 minioadmin minioadmin >nul 2>&1
	@docker exec cakli-storage mc mb local/cakli --ignore-existing >nul 2>&1
	@docker exec cakli-storage mc anonymous set public local/cakli >nul 2>&1
	@echo "✅ Upload should work now!"



# ── Cleanup ──

clean:
	rm -rf ./apps/api/tmp
	rm -rf ./apps/web/.next
	rm -rf ./apps/web/node_modules
	docker system prune -f
