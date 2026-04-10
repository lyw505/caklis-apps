# Run Backend Locally - Upload Fix

## Problem

Upload gagal dengan 403 Forbidden karena signature mismatch antara Docker container dan browser.

## Solution: Run Backend Locally

Jalankan backend di host machine (bukan di Docker) untuk menghindari masalah networking.

## Steps

### 1. Stop Backend Container

```bash
docker-compose stop backend-api
```

### 2. Verify Backend Stopped

```bash
docker ps | grep cakli-api
# Should show nothing
```

### 3. Start Backend Locally

Open terminal dan jalankan:

```bash
cd apps/api
go run main.go
```

Atau dengan Air (hot reload):

```bash
cd apps/api
air
```

### 4. Verify Backend Running

```bash
curl http://localhost:8080/api/v1/health
```

Expected: `{"success":true,"message":"CAKLI v1 API is running",...}`

### 5. Test Upload

1. Open http://localhost:3000
2. Login: `master@cakli.id` / `admin123`
3. Navigate to Operation Admin → Users
4. Click "Tambah User"
5. Upload profile photo
6. Submit form

**Expected:** Upload berhasil, foto masuk ke MinIO, user tersimpan dengan photo key.

## Why This Works

### Docker Container (Problem)

```
Backend Container:
  - Generate presigned URL with: minio:9000
  - Signature for: minio:9000
  - Replace hostname: minio:9000 → localhost:9000
  - Browser access: localhost:9000
  - Signature mismatch! ❌
```

### Local Backend (Solution)

```
Backend on Host:
  - Generate presigned URL with: localhost:9000
  - Signature for: localhost:9000
  - No hostname replacement needed
  - Browser access: localhost:9000
  - Signature match! ✅
```

## Configuration

Backend `.env` already configured for local run:

```env
# apps/api/.env
DB_HOST=localhost
DB_PORT=5433              # Docker PostgreSQL port
MINIO_ENDPOINT=localhost:9000
MINIO_ACCESS_KEY=minioadmin
MINIO_SECRET_KEY=minioadmin
```

## Workflow

### Daily Development

```bash
# Terminal 1: Start Docker services
docker-compose up -d postgres minio

# Terminal 2: Run backend locally
cd apps/api
air

# Terminal 3: Run frontend
cd apps/web
npm run dev
```

### Or Use Make (Modified)

Update your workflow:

```bash
# 1. Start only Docker services (no backend)
docker-compose up -d postgres minio
make fix-upload

# 2. Run backend locally (Terminal 1)
cd apps/api
air

# 3. Run frontend (Terminal 2)
cd apps/web
npm run dev
```

## Benefits

✅ **Upload works** - No signature mismatch
✅ **Hot reload** - Air auto-reloads on code changes
✅ **Easy debugging** - See logs directly in terminal
✅ **Faster iteration** - No Docker rebuild needed
✅ **Simple networking** - Everything on localhost

## Troubleshooting

### Backend Can't Connect to PostgreSQL

**Error:** `connection refused`

**Solution:** Check PostgreSQL port in `.env`:

```env
DB_PORT=5433  # Not 5432!
```

### Backend Can't Connect to MinIO

**Error:** `connection refused`

**Solution:** Make sure MinIO container is running:

```bash
docker ps | grep cakli-storage
```

If not running:

```bash
docker-compose up -d minio
```

### Port 8080 Already in Use

**Error:** `bind: address already in use`

**Solution:** Stop backend container:

```bash
docker-compose stop backend-api
```

Or kill process using port 8080:

```bash
# Windows
netstat -ano | findstr :8080
taskkill /PID <PID> /F
```

## Alternative: Fix Docker Backend

If you really want to use Docker backend, you need to:

1. Update code to use proper hostname replacement
2. Rebuild container: `docker-compose up -d --build backend-api`
3. Restart: `make down && make dev`

But running locally is simpler and more reliable for development.

## Summary

**Problem:** Upload 403 Forbidden (signature mismatch)

**Solution:** Run backend locally

**Steps:**
1. `docker-compose stop backend-api`
2. `cd apps/api && air`
3. Test upload

**Result:** Upload works! ✅
