# MinIO Upload 403 Forbidden - Fix Guide

## Problem Explanation

### What Happened?
You encountered a `403 Forbidden` error when trying to upload files (driver photos) to MinIO from the browser:

```
PUT http://localhost:9000/cakli/drivers/photo-profile/xxx.jfif 403 (Forbidden)
PUT http://localhost:9000/cakli/drivers/photo-ktp/xxx.jpeg 403 (Forbidden)
PUT http://localhost:9000/cakli/drivers/photo-face/xxx.png 403 (Forbidden)
```

### Root Cause
The MinIO bucket `cakli` was created without proper access policies. By default, MinIO buckets are private and don't allow:
- Direct PUT requests from browsers (even with presigned URLs)
- Public access to objects
- CORS requests from web applications

### Why Presigned URLs Weren't Enough
Even though the backend generates presigned URLs with proper AWS signatures, MinIO still needs:
1. Bucket policy that allows the operations
2. CORS configuration to accept browser requests
3. Proper permissions on the bucket itself

## Solution

### Quick Fix (Run This Now)

1. Make sure your Docker containers are running:
```bash
docker compose ps
```

2. Run the MinIO setup script:
```bash
make setup-minio
```

Or manually:
```powershell
powershell -ExecutionPolicy Bypass -File ./scripts/setup-minio.ps1
```

### What the Script Does

The setup script performs these actions:

1. **Configures MinIO Client (mc)**
   - Sets up alias to connect to local MinIO instance
   - Uses admin credentials (minioadmin/minioadmin)

2. **Creates Bucket**
   - Ensures `cakli` bucket exists
   - Skips if already created

3. **Sets Bucket Policy**
   - Enables public download access
   - Allows presigned URLs to work properly
   - Permits browser-based uploads

### Verify the Fix

1. Open MinIO Console: http://localhost:9001
   - Username: `minioadmin`
   - Password: `minioadmin`

2. Navigate to Buckets → cakli → Access Policy
   - Should show "Public" or "Download" access

3. Try uploading a driver photo again from the web interface

## Technical Details

### Upload Flow
```
Browser → API (POST /upload/presigned-url)
       ← Presigned PUT URL + Object Key
       
Browser → MinIO (PUT with presigned URL)
       ← 200 OK (if policy is correct)
       ← 403 Forbidden (if policy is missing)
```

### MinIO Bucket Policy
The script sets this policy:
```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {"AWS": ["*"]},
      "Action": ["s3:GetObject"],
      "Resource": ["arn:aws:s3:::cakli/*"]
    }
  ]
}
```

### Why This is Safe for Development
- Only allows GET (download) operations publicly
- PUT operations still require presigned URLs with valid signatures
- Presigned URLs are time-limited (15 minutes for uploads)
- In production, you should use more restrictive policies

## Production Considerations

For production deployment, consider:

1. **Use IAM Policies**
   - Create specific MinIO users for the API
   - Don't use root credentials (minioadmin)

2. **Restrict Bucket Access**
   - Use private buckets
   - Serve files through CDN or API proxy
   - Implement proper authentication

3. **Enable HTTPS**
   - Set `MINIO_USE_SSL=true`
   - Configure SSL certificates
   - Update `MINIO_ENDPOINT` to use HTTPS

4. **CORS Configuration**
   - Restrict allowed origins to your domain
   - Don't use wildcard (*) in production

## Troubleshooting

### Still Getting 403?

1. Check MinIO logs:
```bash
docker logs cakli-storage
```

2. Verify bucket policy:
```bash
docker exec cakli-storage mc anonymous get local/cakli
```

3. Test presigned URL manually:
```bash
curl -X PUT -H "Content-Type: image/jpeg" --upload-file test.jpg "PRESIGNED_URL_HERE"
```

### CORS Errors?

If you see CORS errors in browser console:
```bash
docker exec cakli-storage mc admin config set local api cors_allow_origin="http://localhost:3000"
docker restart cakli-storage
```

## References

- [MinIO Bucket Policy Guide](https://min.io/docs/minio/linux/administration/identity-access-management/policy-based-access-control.html)
- [Presigned URLs Documentation](https://min.io/docs/minio/linux/developers/go/API.html#presigned-put-object)
- [MinIO Client (mc) Commands](https://min.io/docs/minio/linux/reference/minio-mc.html)
