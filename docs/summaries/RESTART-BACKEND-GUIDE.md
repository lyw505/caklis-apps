# Backend Restart Guide

## Why Restart?

The MinIO upload fix has been applied to `apps/api/services/upload_service.go`. The backend needs to be restarted to apply these changes.

## What Was Fixed?

The presigned URLs now use `localhost:9000` instead of `minio:9000`, which allows the browser to access MinIO correctly.

## How to Restart Backend

### Option 1: Using Air (Hot Reload)

If you're using Air for development:

```bash
cd apps/api
air
```

Air will automatically detect changes and restart the server.

### Option 2: Manual Restart

1. Stop the current backend process (Ctrl+C)
2. Rebuild and run:

```bash
cd apps/api
go run main.go
```

### Option 3: Using Docker (if applicable)

```bash
docker-compose restart api
```

## Verify Backend is Running

Check that the backend is running on port 8080:

```bash
curl http://localhost:8080/api/health
```

Expected response:
```json
{"status": "ok"}
```

## Test Upload Functionality

1. Login to the application
2. Navigate to Driver Management or User Management
3. Click "Tambah Driver" or "Tambah User"
4. Fill in the form and upload photos
5. Click "Simpan"

The upload should now work without the `ERR_NAME_NOT_RESOLVED` error.

## Troubleshooting

### Still Getting ERR_NAME_NOT_RESOLVED?

1. Verify backend was restarted
2. Check browser console for the upload URL - it should contain `localhost:9000`, not `minio:9000`
3. Verify MinIO is running:
   ```bash
   docker ps | grep minio
   ```
4. Test MinIO access directly:
   ```bash
   curl http://localhost:9000/minio/health/live
   ```

### Upload Fails with 403 Forbidden?

Check MinIO credentials in `apps/api/.env`:
```
MINIO_ENDPOINT=localhost:9000
MINIO_ACCESS_KEY=minioadmin
MINIO_SECRET_KEY=minioadmin
MINIO_USE_SSL=false
```

### Upload Fails with Network Error?

Ensure MinIO container is running:
```bash
docker-compose up -d minio
```

## Next Steps

After restarting the backend:

1. Test driver creation with photo uploads
2. Test user creation with photo upload
3. Verify photos are stored in MinIO
4. Check that photos can be viewed in the table

## Summary

The fix is already in place - you just need to restart the backend to apply it. Once restarted, the upload functionality should work correctly.
