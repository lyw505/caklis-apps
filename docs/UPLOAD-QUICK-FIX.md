# Upload Quick Fix - After Restart

## Problem

Setelah restart device/container, upload file gagal dengan error:
- `500 Internal Server Error` saat request presigned URL
- `403 Forbidden` saat upload ke MinIO

## Root Cause

MinIO bucket policy **tidak persistent** setelah restart. Policy kembali ke default `private`.

## Quick Fix (30 detik)

Jalankan command ini setiap kali setelah restart:

```bash
make fix-upload
```

Atau manual:

```bash
docker exec cakli-storage mc alias set local http://localhost:9000 minioadmin minioadmin
docker exec cakli-storage mc anonymous set public local/cakli
```

## Verify

Check bucket policy:

```bash
docker exec cakli-storage mc anonymous get local/cakli
```

Expected output: `Access permission for 'local/cakli' is 'public'`

## Permanent Solution (Coming Soon)

Saat ini MinIO policy tidak persistent karena:
1. MinIO menyimpan policy di memory, bukan di disk
2. Setiap restart container, policy hilang

**Solusi permanent:**
- Use MinIO with persistent policy configuration
- Or use startup script to auto-configure policy
- Or use Docker entrypoint to set policy on start

## Workflow After Restart

```bash
# 1. Start containers
docker-compose up -d

# 2. Fix upload (set bucket policy)
make fix-upload

# 3. Start development
make dev
```

## Why This Happens

MinIO stores bucket policies in memory by default. When container restarts:
1. MinIO starts fresh
2. Bucket `cakli` exists (data is persistent)
3. But bucket policy is reset to `private`
4. Upload fails because browser can't PUT to private bucket

## Alternative: Run Backend Locally

If you don't want to run `make fix-upload` every time:

```bash
# Stop backend container
docker-compose stop backend-api

# Run backend locally
cd apps/api
go run main.go
```

Backend running on host machine can access MinIO directly and doesn't have the signature mismatch issue.

## Summary

**Quick Fix:** `make fix-upload` after every restart

**Why:** MinIO bucket policy not persistent

**Time:** 30 seconds

**Frequency:** Once per restart
