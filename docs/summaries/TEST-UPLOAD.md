# Test Upload Functionality

## Quick Test

Setelah menjalankan `make setup-minio`, test upload dengan langkah berikut:

### 1. Test dari Web Interface

1. Buka browser: http://localhost:3000
2. Login dengan credentials admin
3. Navigate: Operation Admin → Drivers
4. Klik "Tambah Driver"
5. Fill form dengan data:
   - Nama: Test Driver
   - Email: test@driver.com
   - Password: password123
   - Phone: 081234567890
   - NIK: 1234567890123456
6. Upload 3 foto (profile, KTP, muka)
7. Klik "Simpan"

**Expected Result:**
- ✅ Upload progress indicator muncul
- ✅ "Driver berhasil ditambahkan" toast notification
- ✅ Driver muncul di tabel
- ✅ Foto bisa dilihat di MinIO Console

### 2. Verify di MinIO Console

1. Buka: http://localhost:9001
2. Login: `minioadmin` / `minioadmin`
3. Navigate: Buckets → cakli
4. Check folders:
   - drivers/photo-profile/
   - drivers/photo-ktp/
   - drivers/photo-face/

**Expected Result:**
- ✅ File dengan UUID filename terlihat
- ✅ File size sesuai dengan yang diupload
- ✅ Content-Type: image/jpeg atau image/png

### 3. Test Manual dengan cURL

```bash
# Step 1: Login dan dapatkan token
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@cakli.com",
    "password": "admin123"
  }'

# Save the access_token from response

# Step 2: Request presigned URL
curl -X POST http://localhost:8080/api/v1/upload/presigned-url \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -d '{
    "filename": "test.jpeg",
    "content_type": "image/jpeg",
    "folder": "drivers/photo-profile"
  }'

# Save the upload_url from response

# Step 3: Upload file
curl -X PUT "UPLOAD_URL_FROM_STEP_2" \
  -H "Content-Type: image/jpeg" \
  --upload-file /path/to/test.jpeg

# Expected: Empty response with 200 OK status
```

## Troubleshooting

### Still Getting 403?

Run these commands:

```bash
# Check bucket policy
docker exec cakli-storage mc anonymous get local/cakli

# Should output: Access permission for `local/cakli` is `public`

# If not, run setup again
make setup-minio
```

### Upload Timeout?

Check if MinIO is running:

```bash
docker ps | grep cakli-storage

# Should show: cakli-storage container running
```

### Network Error?

Check if ports are accessible:

```bash
# Test MinIO API
curl http://localhost:9000/minio/health/live

# Test MinIO Console
curl http://localhost:9001

# Test Backend API
curl http://localhost:8080/api/v1/health
```

## Success Indicators

✅ Bucket policy is `public`
✅ MinIO container is running
✅ Backend API is running
✅ Web frontend is running
✅ Upload returns 200 OK
✅ Files visible in MinIO Console
✅ Driver created successfully in database

## Common Issues

### Issue 1: 403 Forbidden
**Cause:** Bucket policy not set
**Fix:** Run `make setup-minio`

### Issue 2: 409 Conflict
**Cause:** Upload failed but form still submitted
**Fix:** Check upload error handling in driver-form-modal.tsx

### Issue 3: CORS Error
**Cause:** Browser blocking cross-origin request
**Fix:** MinIO CORS not needed with presigned URLs (already handled)

### Issue 4: Presigned URL Expired
**Cause:** Upload took longer than 15 minutes
**Fix:** Increase expiry in upload_service.go

### Issue 5: Network Error
**Cause:** MinIO not accessible from browser
**Fix:** Check docker-compose.yml ports mapping
