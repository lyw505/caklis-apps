# 🔧 Troubleshooting Guide - CAKLI v1

## Common Issues and Solutions

### 1. Port Already Allocated Error

**Error Message:**
```
Error response from daemon: driver failed programming external connectivity on endpoint cakli-storage: 
Bind for 0.0.0.0:9000 failed: port is already allocated
```

**Cause:**
Port 9000 (MinIO) atau port lainnya sudah digunakan oleh container lama atau aplikasi lain.

**Solution:**

#### Step 1: Check Running Containers
```bash
docker ps -a
```

#### Step 2: Stop and Remove Old Containers
```bash
# Stop specific container
docker stop cakli-minio
docker rm cakli-minio

# Or stop all CAKLI containers
docker stop cakli-postgres cakli-minio cakli-api
docker rm cakli-postgres cakli-minio cakli-api
```

#### Step 3: Restart Services
```bash
docker-compose down
docker-compose up -d --build
```

#### Step 4: Verify
```bash
docker ps
curl http://localhost:8080/api/v1/health
```

---

### 2. Database Connection Failed

**Error Message:**
```
Failed to connect to database: failed to connect to `user=cakli database=cakli_db`
```

**Cause:**
- PostgreSQL container tidak running
- Port salah (5432 vs 5433)
- Credentials salah

**Solution:**

#### Check Container Status
```bash
docker ps | grep postgres
```

#### Check Environment Variables
Pastikan `.env` dan `apps/api/.env` memiliki konfigurasi yang benar:
```env
DB_HOST=localhost  # atau 'postgres' jika di Docker
DB_PORT=5432
DB_USER=cakli
DB_PASSWORD=cakli_secret
DB_NAME=cakli_db
```

#### Restart Database
```bash
docker-compose restart postgres
```

---

### 3. MinIO Connection Failed

**Error Message:**
```
Failed to check bucket: dial tcp [::1]:9000: connect: connection refused
```

**Cause:**
- MinIO container tidak running
- Environment variable salah
- Network issue antar container

**Solution:**

#### Update docker-compose.yml
Pastikan environment variables untuk MinIO sudah di-set:
```yaml
backend-api:
  environment:
    - MINIO_ENDPOINT=minio:9000  # Gunakan nama service, bukan localhost
    - MINIO_ACCESS_KEY=minioadmin
    - MINIO_SECRET_KEY=minioadmin
    - MINIO_BUCKET=cakli
    - MINIO_USE_SSL=false
```

#### Rebuild Containers
```bash
docker-compose down
docker-compose up -d --build
```

---

### 4. Frontend Build Error

**Error Message:**
```
Parsing ecmascript source code failed
Unterminated regexp literal
```

**Cause:**
- Syntax error di file TypeScript/React
- Karakter tersembunyi
- Import statement salah

**Solution:**

#### Check Import Statements
Pastikan semua import menggunakan format yang benar:
```typescript
// ✅ Correct
import { api } from "@/lib/api"

// ❌ Wrong (if api is default export)
import api from "@/lib/api"
```

#### Run Dev Mode First
```bash
cd apps/web
npm run dev
```

Dev mode lebih toleran dan akan menunjukkan error yang lebih jelas.

---

### 5. API Export Not Found

**Error Message:**
```
Export api doesn't exist in target module
Did you mean to import default?
```

**Cause:**
File `lib/api.ts` hanya export default, tapi di-import sebagai named export.

**Solution:**

Update `apps/web/lib/api.ts`:
```typescript
export default api;
export { api };  // Add this line
```

---

### 6. Authentication Failed

**Error Message:**
```
401 Unauthorized
Token tidak valid
```

**Cause:**
- Access token expired
- Refresh token expired
- Cookie tidak ter-set

**Solution:**

#### Clear Browser Storage
```javascript
// Di browser console
localStorage.clear()
// Reload page
```

#### Check Cookie Settings
Pastikan `withCredentials: true` di axios config:
```typescript
const api = axios.create({
  baseURL: process.env.NEXT_PUBLIC_API_URL,
  withCredentials: true,  // Important!
});
```

---

## Quick Commands

### Reset Everything
```bash
# Stop all containers
docker-compose down

# Remove volumes (WARNING: This deletes all data!)
docker-compose down -v

# Rebuild and start
docker-compose up -d --build

# Check status
docker ps
docker logs cakli-api
docker logs cakli-db
docker logs cakli-storage
```

### Check Logs
```bash
# API logs
docker logs cakli-api --tail 50 -f

# Database logs
docker logs cakli-db --tail 50 -f

# MinIO logs
docker logs cakli-storage --tail 50 -f
```

### Database Access
```bash
# Connect to PostgreSQL
docker exec -it cakli-db psql -U cakli -d cakli_db

# List tables
\dt

# Check data
SELECT * FROM admins;
SELECT * FROM drivers;
SELECT * FROM users;
```

### MinIO Access
- Console: http://localhost:9001
- Username: minioadmin
- Password: minioadmin

---

## Environment Variables Checklist

### Root `.env`
```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=cakli
DB_PASSWORD=cakli_secret
DB_NAME=cakli_db
DB_SSLMODE=disable

MINIO_ENDPOINT=localhost:9000
MINIO_ACCESS_KEY=minioadmin
MINIO_SECRET_KEY=minioadmin
MINIO_BUCKET=cakli
MINIO_USE_SSL=false

JWT_ACCESS_SECRET=your-access-secret-key-change-in-production
JWT_REFRESH_SECRET=your-refresh-secret-key-change-in-production
JWT_ACCESS_EXPIRY=15m
JWT_REFRESH_EXPIRY=168h

API_PORT=8080
API_ENV=development

NEXT_PUBLIC_API_URL=http://localhost:8080/api/v1
NEXT_PUBLIC_APP_NAME=CAKLI Admin
```

### `apps/api/.env`
Same as root `.env`

### `apps/web/.env.local`
```env
NEXT_PUBLIC_API_URL=http://localhost:8080/api/v1
NEXT_PUBLIC_APP_NAME=CAKLI Admin
```

---

## Health Check Commands

```bash
# Check all services
curl http://localhost:8080/api/v1/health
curl http://localhost:9000/minio/health/live
psql -h localhost -U cakli -d cakli_db -c "SELECT 1"

# Check ports
netstat -an | findstr "8080"
netstat -an | findstr "5432"
netstat -an | findstr "9000"
```

---

## Getting Help

If you encounter issues not covered here:

1. Check container logs: `docker logs <container-name>`
2. Check Docker status: `docker ps -a`
3. Verify environment variables: `docker exec cakli-api env`
4. Check network: `docker network ls` and `docker network inspect caklis-apps_default`

---

**Last Updated:** April 9, 2026


---

## 🆕 Recent Issues

### Issue: MinIO Upload Error - ERR_NAME_NOT_RESOLVED

**Symptom:**
```
PUT http://minio:9000/cakli/drivers/photo-profile/xxx.jpg net::ERR_NAME_NOT_RESOLVED
Error: "gagal menyimpan driver"
```

**Cause:** 
Frontend (browser) mencoba mengakses MinIO menggunakan hostname internal Docker (`minio:9000`) yang tidak bisa diakses dari browser. Browser hanya bisa akses melalui `localhost:9000`.

**Solution:**

Backend sudah diperbaiki untuk mengganti hostname di presigned URL. Restart backend untuk apply changes:

```bash
# Stop backend (Ctrl+C)

# Start again
cd apps/api
go run main.go
```

**Verification:**

1. Test presigned URL generation:
```bash
curl -X POST http://localhost:8080/api/v1/upload/presigned-url \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "filename": "test.jpg",
    "content_type": "image/jpeg",
    "folder": "drivers/photo-profile"
  }'
```

2. Check response - URL should contain `localhost:9000` NOT `minio:9000`:
```json
{
  "success": true,
  "data": {
    "upload_url": "http://localhost:9000/cakli/...",  // ✅ localhost
    "object_key": "drivers/photo-profile/xxx.jpg"
  }
}
```

3. Try upload driver/user with photo again - should work now!

**Technical Details:**

The fix replaces internal Docker hostname with localhost in presigned URLs:

```go
// apps/api/services/upload_service.go
urlString := presignedURL.String()
urlString = strings.Replace(urlString, "minio:9000", "localhost:9000", 1)
```

This ensures browser can access MinIO from outside Docker network.

---

**Last Updated:** April 9, 2026
