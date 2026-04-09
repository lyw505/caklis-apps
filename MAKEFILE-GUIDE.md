# 📘 Makefile Guide - CAKLI Admin System

Panduan lengkap menggunakan Makefile untuk development CAKLI Admin System.

---

## 🚀 Quick Start

```bash
# Install dependencies
make install

# Start everything (Docker + Backend + Frontend)
make dev
```

Setelah `make dev`:
- Docker services (PostgreSQL + MinIO) akan start
- Backend API akan start di PowerShell window terpisah (minimized)
- Frontend akan start di terminal saat ini
- Browser akan otomatis membuka http://localhost:3000

---

## 📋 Semua Perintah

### Development

| Perintah | Deskripsi |
|----------|-----------|
| `make dev` | Start full stack (Docker + API + Frontend) |
| `make dev-api` | Start hanya backend API |
| `make dev-web` | Start hanya frontend |
| `make dev-stop` | Stop semua processes (Go, Node, Docker) |

### Docker Management

| Perintah | Deskripsi |
|----------|-----------|
| `make up` | Start Docker services |
| `make down` | Stop Docker services |
| `make restart` | Restart Docker services |
| `make build` | Rebuild Docker images |

### Logs

| Perintah | Deskripsi |
|----------|-----------|
| `make logs` | Show all Docker logs |
| `make logs-db` | Show PostgreSQL logs |
| `make logs-minio` | Show MinIO logs |

### Database

| Perintah | Deskripsi |
|----------|-----------|
| `make db-init` | Initialize database schema |
| `make db-seed` | Seed database with test data |
| `make db-reset` | Reset database (drop + init + seed) |
| `make db-connect` | Connect to PostgreSQL CLI |

### Health & Testing

| Perintah | Deskripsi |
|----------|-----------|
| `make health` | Check all services health |
| `make test` | Run tests |

### Cleanup

| Perintah | Deskripsi |
|----------|-----------|
| `make clean` | Remove build artifacts |
| `make install` | Install dependencies |

---

## 🎯 Workflow Umum

### Workflow 1: Start Development (Otomatis)

```bash
# Start everything dengan satu perintah
make dev
```

Ini akan:
1. ✅ Start Docker (PostgreSQL + MinIO)
2. ✅ Initialize database
3. ✅ Sync environment variables
4. ✅ Start backend di PowerShell window terpisah
5. ✅ Start frontend di terminal saat ini

**Backend logs:** Lihat di PowerShell window yang minimized  
**Frontend logs:** Lihat di terminal saat ini

**Stop:** Tekan `Ctrl+C` untuk stop frontend, lalu jalankan `make dev-stop` untuk stop semua

### Workflow 2: Start Development (Manual Control)

```bash
# Terminal 1: Start Docker
make up

# Terminal 2: Start Backend
make dev-api

# Terminal 3: Start Frontend
make dev-web
```

Dengan cara ini, Anda punya kontrol penuh dan bisa melihat logs di setiap terminal.

**Stop:** Tekan `Ctrl+C` di setiap terminal

### Workflow 3: Database Reset

```bash
# Reset database lengkap
make db-reset

# Atau step by step:
make db-init    # Initialize schema
make db-seed    # Seed data
```

---

## 🔍 Troubleshooting

### Backend tidak start?

**Cek PowerShell window:**
- `make dev` membuka PowerShell window minimized
- Klik taskbar untuk melihat window
- Lihat error di window tersebut

**Atau start manual:**
```bash
make dev-api
```

### Frontend tidak bisa connect ke backend?

**Check backend health:**
```bash
make health
```

**Atau manual:**
```bash
curl http://localhost:8080/api/v1/health
```

### Database error "already exists"?

Ini normal jika database sudah di-initialize sebelumnya. Makefile akan skip error ini.

**Jika ingin reset:**
```bash
make db-reset
```

### Port sudah digunakan?

**Stop semua processes:**
```bash
make dev-stop
```

**Atau manual:**
```bash
# Stop Docker
make down

# Kill Go processes
taskkill /F /IM go.exe

# Kill Node processes
taskkill /F /IM node.exe
```

### Docker tidak running?

**Start Docker Desktop terlebih dahulu**, lalu:
```bash
make up
```

---

## 💡 Tips & Tricks

### 1. Lihat Semua Perintah
```bash
make help
```

### 2. Check Health Semua Services
```bash
make health
```

Output:
```
1. PostgreSQL: ✅ Ready
2. MinIO: ✅ Healthy
3. Backend API: ✅ Healthy
4. Frontend: ✅ Running
```

### 3. Connect ke Database
```bash
make db-connect
```

Kemudian di PostgreSQL CLI:
```sql
-- List tables
\dt

-- Check admins
SELECT * FROM admins;

-- Check drivers
SELECT * FROM drivers;

-- Exit
\q
```

### 4. Clean Build Artifacts
```bash
make clean
```

Ini akan hapus:
- `apps/api/tmp/`
- `apps/api/api.exe`
- `apps/web/.next/`
- `apps/web/node_modules/`

### 5. Reinstall Dependencies
```bash
make clean
make install
```

---

## 🎨 Customization

### Mengubah Port

Edit `.env`:
```env
API_PORT=8080        # Backend port
```

Edit `apps/web/package.json`:
```json
{
  "scripts": {
    "dev": "next dev -p 3001"  // Frontend port
  }
}
```

### Mengubah Database Credentials

Edit `.env`:
```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=cakli
DB_PASSWORD=cakli_secret
DB_NAME=cakli_db
```

Kemudian:
```bash
make down
make up
make db-init
```

---

## 📊 Service URLs

| Service | URL | Credentials |
|---------|-----|-------------|
| Frontend | http://localhost:3000 | See login credentials below |
| Backend API | http://localhost:8080 | N/A |
| MinIO Console | http://localhost:9001 | minioadmin/minioadmin |
| PostgreSQL | localhost:5432 | cakli/cakli_secret |

### Login Credentials

| Role | Email | Password |
|------|-------|----------|
| Master Admin | master@cakli.id | admin123 |
| Operation Admin | operation@cakli.id | admin123 |
| Reporting Admin | reporting@cakli.id | admin123 |

---

## 🔄 Development Cycle

```bash
# 1. Start development
make dev

# 2. Make code changes
# ... edit files ...

# 3. Backend auto-reload (if using Air)
# Or restart manually:
# Ctrl+C in backend terminal, then: make dev-api

# 4. Frontend auto-reload (Next.js hot reload)
# Just save files, browser will auto-refresh

# 5. Test changes
make health

# 6. Stop when done
make dev-stop
```

---

## 🆘 Common Issues

### Issue 1: "make: command not found"

**Solution:** Install Make for Windows
- Download from: http://gnuwin32.sourceforge.net/packages/make.htm
- Or use Chocolatey: `choco install make`

### Issue 2: "docker-compose: command not found"

**Solution:** Install Docker Desktop
- Download from: https://www.docker.com/products/docker-desktop

### Issue 3: "go: command not found"

**Solution:** Install Go
- Download from: https://go.dev/dl/

### Issue 4: "npm: command not found"

**Solution:** Install Node.js
- Download from: https://nodejs.org/

### Issue 5: Backend tidak start di `make dev`

**Solution:** Start manual di terminal terpisah
```bash
make dev-api
```

---

## 📚 Related Documentation

- [START-HERE.md](START-HERE.md) - Getting started guide
- [QUICK-START.md](QUICK-START.md) - Quick reference
- [TROUBLESHOOTING.md](docs/TROUBLESHOOTING.md) - Detailed troubleshooting
- [README.md](README.md) - Project documentation

---

**Last Updated:** April 9, 2026
