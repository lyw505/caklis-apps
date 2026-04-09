# CAKLI Setup Summary

## ✅ Files Created/Updated

### New Files
1. `.gitignore` - Git ignore rules untuk project
2. `.env.example` - Template environment variables
3. `scripts/setup-minio.ps1` - Script setup MinIO untuk Windows
4. `scripts/setup-minio.sh` - Script setup MinIO untuk Linux/Mac
5. `docs/MINIO-403-ERROR-DETAILED.md` - Analisis detail error upload
6. `docs/QUICK-FIX-UPLOAD.md` - Quick reference upload fix
7. `docs/TEST-UPLOAD.md` - Panduan testing upload
8. `docs/SETUP-SUMMARY.md` - File ini

### Updated Files
1. `README.md` - Added MinIO setup instructions & troubleshooting
2. `Makefile` - Added `setup-minio` command
3. `docker-compose.yml` - Added healthcheck for MinIO

## 🚀 Quick Start (For New Developers)

```bash
# 1. Clone repository
git clone <repo-url>
cd caklis-apps

# 2. Copy environment file
cp .env.example .env

# 3. Start all services
make dev

# 4. Configure MinIO (first time only)
make setup-minio

# 5. Open browser
# http://localhost:3000
```

## 📋 Available Make Commands

```bash
make help          # Show all available commands
make dev           # Start full development stack
make up            # Start Docker services only
make down          # Stop all services
make build         # Rebuild Docker images
make logs          # Show logs
make seed          # Seed database
make setup-minio   # Configure MinIO bucket
make clean         # Clean build artifacts
```

## 🔧 MinIO Setup (Important!)

After first `make dev`, you MUST run:

```bash
make setup-minio
```

This configures MinIO bucket to allow file uploads. Without this:
- ❌ Driver photo upload will fail (403 Forbidden)
- ❌ User profile picture upload will fail (403 Forbidden)

## 🎯 What's in .gitignore

The `.gitignore` file excludes:
- Environment files (`.env`, `.env.local`)
- Node modules (`node_modules/`)
- Build artifacts (`.next/`, `tmp/`, `bin/`)
- IDE files (`.vscode/`, `.idea/`)
- OS files (`.DS_Store`, `Thumbs.db`)
- Docker volumes data
- Logs and temp files
- Sensitive files (`.pem`, `.key`)

## 📦 Project Structure

```
caklis-apps/
├── .env                    # Environment variables (gitignored)
├── .env.example            # Template for .env
├── .gitignore              # Git ignore rules
├── Makefile                # Development commands
├── docker-compose.yml      # Docker services
├── README.md               # Main documentation
├── apps/
│   ├── api/               # Go backend
│   └── web/               # Next.js frontend
├── docs/                  # Documentation
├── scripts/               # Setup scripts
└── infra/                 # Infrastructure configs
```

## 🔐 Default Credentials

### Application Login
- Master Admin: `master@cakli.id` / `admin123`
- Operation Admin: `operation@cakli.id` / `admin123`
- Reporting Admin: `reporting@cakli.id` / `admin123`

### MinIO Console
- URL: http://localhost:9001
- Username: `minioadmin`
- Password: `minioadmin`

### PostgreSQL
- Host: `localhost:5432`
- Database: `cakli_db`
- Username: `cakli`
- Password: `cakli_secret`

## 🐛 Common Issues

### 1. Upload fails with 403
```bash
make setup-minio
```

### 2. Port already in use
```bash
# Windows
netstat -ano | findstr :8080
taskkill /PID <PID> /F

# Linux/Mac
lsof -ti:8080 | xargs kill -9
```

### 3. Docker won't start
```bash
docker-compose down -v
docker-compose up -d
```

## 📚 Documentation Links

- [README.md](../README.md) - Main documentation
- [TROUBLESHOOTING.md](./TROUBLESHOOTING.md) - Troubleshooting guide
- [TESTING-GUIDE.md](./TESTING-GUIDE.md) - Testing instructions
- [MINIO-403-ERROR-DETAILED.md](./MINIO-403-ERROR-DETAILED.md) - Upload error details

## ✨ Next Steps

1. ✅ Run `make dev`
2. ✅ Run `make setup-minio`
3. ✅ Open http://localhost:3000
4. ✅ Login and test features
5. ✅ Read documentation for advanced features

## 🎉 You're Ready!

Everything is set up and ready to go. Happy coding! 🚀
