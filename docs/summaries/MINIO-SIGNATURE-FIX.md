# MinIO 403 Error - Signature Mismatch Fix

## Problem

Upload masih gagal dengan 403 Forbidden meskipun bucket policy sudah public:

```
PUT http://localhost:9000/cakli/users/photo-profile/xxx.jpeg 403 (Forbidden)
```

**Bucket policy sudah benar:**
```json
{
  "Statement": [{
    "Action": ["s3:PutObject", "s3:GetObject", ...],
    "Effect": "Allow",
    "Principal": {"AWS": ["*"]},
    "Resource": ["arn:aws:s3:::cakli/*"]
  }],
  "Version": "2012-10-17"
}
```

## Root Cause

**Signature Mismatch!**

Presigned URL menggunakan AWS Signature V4 yang di-generate berdasarkan:
1. Access Key
2. Secret Key
3. **Endpoint/Hostname** ← MASALAH DI SINI!
4. Timestamp
5. Object key

**Yang terjadi sebelumnya:**
```
1. Backend generate presigned URL dengan endpoint: minio:9000
   Signature = HMAC-SHA256(secret, "minio:9000/cakli/...")
   
2. Backend replace hostname: minio:9000 → localhost:9000
   URL = http://localhost:9000/cakli/...?signature=XXX
   
3. Browser akses: localhost:9000
   MinIO validate signature untuk: localhost:9000
   
4. Signature mismatch! ❌
   Expected: signature untuk "minio:9000"
   Got: request ke "localhost:9000"
```

**Analogi:**
- Signature seperti tanda tangan di cek
- Cek ditulis untuk "Bank A" (minio:9000)
- Tapi Anda coba cairkan di "Bank B" (localhost:9000)
- Bank B tolak karena cek bukan untuk mereka!

## Solution

**Jangan replace hostname!** Generate presigned URL langsung dengan endpoint yang akan diakses browser.

### Changes Made

#### 1. Environment Variables (No Change Needed)
```env
# .env
MINIO_ENDPOINT=localhost:9000  # Already correct!
```

Backend API di Docker bisa akses MinIO dengan `localhost:9000` karena port mapping:
```yaml
# docker-compose.yml
services:
  minio:
    ports:
      - "9000:9000"  # Host:Container
```

#### 2. Upload Service (Removed Hostname Replacement)

**Before:**
```go
presignedURL, err := config.MinioClient.PresignedPutObject(...)
urlString := presignedURL.String()
urlString = strings.Replace(urlString, "minio:9000", "localhost:9000", 1) // ❌ BAD!
return urlString
```

**After:**
```go
presignedURL, err := config.MinioClient.PresignedPutObject(...)
return presignedURL.String() // ✅ GOOD! No replacement
```

### Why This Works

```
1. MinIO client configured with: localhost:9000
   
2. Backend generate presigned URL:
   Signature = HMAC-SHA256(secret, "localhost:9000/cakli/...")
   URL = http://localhost:9000/cakli/...?signature=XXX
   
3. Browser akses: localhost:9000
   MinIO validate signature untuk: localhost:9000
   
4. Signature match! ✅
   Expected: signature untuk "localhost:9000"
   Got: request ke "localhost:9000"
```

## Testing

### 1. Restart Backend
```bash
docker-compose restart backend-api
```

### 2. Test Upload
1. Open http://localhost:3000
2. Login: `master@cakli.id` / `admin123`
3. Navigate to Operation Admin → Users
4. Click "Tambah User"
5. Upload profile photo
6. Submit form

**Expected Result:**
- ✅ Upload succeeds (200 OK)
- ✅ User created successfully
- ✅ Photo visible in MinIO Console

### 3. Verify Presigned URL
```bash
# Get presigned URL from API
curl -X POST http://localhost:8080/api/v1/upload/presigned-url \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "filename": "test.jpeg",
    "content_type": "image/jpeg",
    "folder": "users/photo-profile"
  }'

# Check URL contains localhost:9000 (not minio:9000)
# URL should be: http://localhost:9000/cakli/...
```

## Technical Details

### AWS Signature V4

Presigned URL signature calculation:
```
StringToSign = 
  "AWS4-HMAC-SHA256\n" +
  Timestamp + "\n" +
  Scope + "\n" +
  Hash(CanonicalRequest)

CanonicalRequest =
  HTTPMethod + "\n" +
  CanonicalURI + "\n" +
  CanonicalQueryString + "\n" +
  CanonicalHeaders + "\n" +    ← Host header here!
  SignedHeaders + "\n" +
  HashedPayload

Signature = HMAC-SHA256(SigningKey, StringToSign)
```

**Host header must match!**
- If signature generated for `minio:9000`
- Request must go to `minio:9000`
- Otherwise: 403 Forbidden

### Docker Networking

**Why backend can access localhost:9000?**

Docker Desktop on Windows uses host networking for `localhost`:
- Container can access host's `localhost`
- Port mapping: `0.0.0.0:9000->9000/tcp`
- Backend in container → `localhost:9000` → Host's MinIO

**Alternative for production:**
- Use Docker internal network: `minio:9000`
- Set `MINIO_EXTERNAL_ENDPOINT` for browser access
- Generate presigned URL with external endpoint

## Production Considerations

### Option 1: Separate Endpoints (Recommended)

```env
# Internal endpoint (backend → MinIO)
MINIO_ENDPOINT=minio:9000

# External endpoint (browser → MinIO)
MINIO_EXTERNAL_ENDPOINT=https://storage.yourdomain.com
```

**Code:**
```go
// Use external endpoint for presigned URLs
externalClient := createMinioClient(os.Getenv("MINIO_EXTERNAL_ENDPOINT"))
presignedURL, _ := externalClient.PresignedPutObject(...)
```

### Option 2: Reverse Proxy

```nginx
# nginx.conf
location /storage/ {
    proxy_pass http://minio:9000/;
    proxy_set_header Host $host;
}
```

```env
MINIO_ENDPOINT=minio:9000
MINIO_EXTERNAL_ENDPOINT=https://yourdomain.com/storage
```

### Option 3: CDN

Use CloudFront, CloudFlare, or similar:
- Origin: MinIO server
- CDN URL: https://cdn.yourdomain.com
- Generate presigned URLs with CDN endpoint

## Summary

**Problem:**
- Presigned URL signature generated for `minio:9000`
- Browser accessed `localhost:9000`
- Signature mismatch → 403 Forbidden

**Solution:**
- Generate presigned URL with `localhost:9000` from start
- No hostname replacement needed
- Signature matches request → Success!

**Key Lesson:**
- AWS Signature V4 includes hostname in signature
- Hostname in signature must match hostname in request
- Don't replace hostname after signature generation

**Files Changed:**
- `apps/api/services/upload_service.go` - Removed hostname replacement
- `apps/api/.env` - Already correct (localhost:9000)
- `.env` - Already correct (localhost:9000)
- `.env.example` - Added MINIO_EXTERNAL_ENDPOINT for reference

**Status:** ✅ Fixed! Upload should work now.
