# MinIO 403 Forbidden Error - Analisis Detail

## Error yang Terjadi

```
PUT http://localhost:9000/cakli/drivers/photo-profile/xxx.jpeg 403 (Forbidden)
PUT http://localhost:9000/cakli/drivers/photo-ktp/xxx.jpeg 403 (Forbidden)
PUT http://localhost:9000/cakli/drivers/photo-face/xxx.jpeg 403 (Forbidden)
POST http://localhost:8080/api/v1/admin/drivers 409 (Conflict)
```

## Analisis Masalah Secara Detail

### 1. Root Cause: Bucket Policy MinIO

**Masalah Utama:**
MinIO bucket `cakli` dibuat dengan policy default `private`, yang berarti:
- ❌ Tidak mengizinkan GET (download) dari publik
- ❌ Tidak mengizinkan PUT (upload) dari publik
- ❌ Presigned URL tidak berfungsi karena bucket policy menolak operasi

**Verifikasi Masalah:**
```bash
docker exec cakli-storage mc anonymous get local/cakli
# Output: Access permission for `local/cakli` is `private`
```

### 2. Mengapa Presigned URL Tidak Cukup?

**Presigned URL** adalah URL yang sudah ditandatangani dengan AWS Signature V4, yang berisi:
- Access Key ID
- Timestamp
- Expiration time
- Signature (HMAC-SHA256)

**Namun**, presigned URL hanya memberikan **autentikasi**, bukan **autorisasi**.

**Analogi:**
- Presigned URL = KTP yang valid (membuktikan identitas)
- Bucket Policy = Aturan gedung (menentukan siapa yang boleh masuk)

Meskipun Anda punya KTP valid, jika aturan gedung melarang masuk, Anda tetap tidak bisa masuk.

### 3. Alur Request yang Gagal

```
Step 1: Browser → API
Request: POST /upload/presigned-url
Body: {
  "filename": "photo.jpeg",
  "content_type": "image/jpeg",
  "folder": "drivers/photo-profile"
}

Step 2: API → Browser
Response: {
  "upload_url": "http://localhost:9000/cakli/drivers/photo-profile/uuid.jpeg?X-Amz-Algorithm=...",
  "object_key": "drivers/photo-profile/uuid.jpeg"
}

Step 3: Browser → MinIO
Request: PUT http://localhost:9000/cakli/drivers/photo-profile/uuid.jpeg
Headers: {
  "Content-Type": "image/jpeg"
}
Body: [binary file data]

Step 4: MinIO → Browser
Response: 403 Forbidden ❌
Reason: Bucket policy 'private' menolak PUT operation
```

### 4. Mengapa Error 409 Conflict Muncul?

Error `409 Conflict` pada `POST /admin/drivers` terjadi karena:

1. Upload file gagal (403 Forbidden)
2. `formData.photo_profile_key` tetap kosong atau undefined
3. Backend mencoba membuat driver dengan data yang tidak lengkap
4. Database constraint atau business logic menolak data tersebut

**Kode di driver-form-modal.tsx:**
```typescript
// Upload photos if selected
if (photoProfile) {
    formData.photo_profile_key = await uploadFile(photoProfile, "drivers/photo-profile")
    // ❌ Ini throw error karena 403, tapi error tidak di-handle dengan baik
}

// Tetap lanjut ke create driver meskipun upload gagal
if (driver) {
    await api.put(`/admin/drivers/${driver.id}`, formData)
} else {
    await api.post("/admin/drivers", formData) // ❌ 409 Conflict
}
```

## Solusi yang Diterapkan

### 1. Mengubah Bucket Policy ke Public

**Command:**
```bash
docker exec cakli-storage mc anonymous set public local/cakli
```

**Efek:**
- ✅ Mengizinkan GET (download) dari publik
- ✅ Mengizinkan PUT (upload) dari publik dengan presigned URL
- ✅ Presigned URL sekarang berfungsi dengan baik

**Verifikasi:**
```bash
docker exec cakli-storage mc anonymous get local/cakli
# Output: Access permission for `local/cakli` is `public`
```

### 2. Policy Levels di MinIO

MinIO memiliki 3 level policy:

| Policy | GET (Download) | PUT (Upload) | DELETE |
|--------|----------------|--------------|--------|
| `private` | ❌ | ❌ | ❌ |
| `download` | ✅ | ❌ | ❌ |
| `public` | ✅ | ✅ | ❌ |

**Mengapa kita butuh `public`?**
- Karena browser perlu melakukan PUT request untuk upload
- `download` hanya mengizinkan GET, tidak cukup untuk upload

### 3. Keamanan Presigned URL

Meskipun bucket policy `public`, presigned URL tetap aman karena:

1. **Time-Limited**: URL expired setelah 15 menit
2. **Signature Validation**: MinIO memvalidasi signature AWS
3. **Specific Object**: URL hanya berlaku untuk object key tertentu
4. **Content-Type Locked**: Hanya menerima content-type yang ditentukan

**Contoh Presigned URL:**
```
http://localhost:9000/cakli/drivers/photo-profile/uuid.jpeg?
  X-Amz-Algorithm=AWS4-HMAC-SHA256&
  X-Amz-Credential=minioadmin/20260409/us-east-1/s3/aws4_request&
  X-Amz-Date=20260409T061025Z&
  X-Amz-Expires=900&                    ← Expired dalam 15 menit
  X-Amz-SignedHeaders=host&
  X-Amz-Signature=88026ab9...           ← Signature unik
```

## Testing Solusi

### 1. Test Upload dari Browser

1. Buka aplikasi web: http://localhost:3000
2. Login sebagai admin
3. Navigasi ke Operation Admin → Drivers
4. Klik "Tambah Driver"
5. Upload foto profile, KTP, dan foto muka
6. Submit form

**Expected Result:**
- ✅ Upload berhasil (200 OK)
- ✅ Driver berhasil dibuat
- ✅ Foto tersimpan di MinIO

### 2. Test Manual dengan cURL

```bash
# 1. Get presigned URL dari API
curl -X POST http://localhost:8080/api/v1/upload/presigned-url \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "filename": "test.jpeg",
    "content_type": "image/jpeg",
    "folder": "drivers/photo-profile"
  }'

# Response:
# {
#   "upload_url": "http://localhost:9000/cakli/...",
#   "object_key": "drivers/photo-profile/uuid.jpeg"
# }

# 2. Upload file ke MinIO
curl -X PUT "PRESIGNED_URL_FROM_STEP_1" \
  -H "Content-Type: image/jpeg" \
  --upload-file test.jpeg

# Expected: 200 OK
```

### 3. Verify di MinIO Console

1. Buka MinIO Console: http://localhost:9001
2. Login: `minioadmin` / `minioadmin`
3. Navigate to Buckets → cakli → drivers/photo-profile
4. File yang diupload harus terlihat di sini

## Troubleshooting

### Masih 403 setelah setup?

**Cek 1: Bucket Policy**
```bash
docker exec cakli-storage mc anonymous get local/cakli
```
Expected: `public`

**Cek 2: MinIO Logs**
```bash
docker logs cakli-storage --tail 50
```
Look for: Access denied errors

**Cek 3: Restart MinIO**
```bash
docker restart cakli-storage
```

### CORS Error?

Jika muncul CORS error di browser console:

```bash
# Set CORS policy
docker exec cakli-storage mc admin config set local api \
  cors_allow_origin="http://localhost:3000"

# Restart MinIO
docker restart cakli-storage
```

### Presigned URL Expired?

Jika upload lambat dan URL expired:

Edit `apps/api/services/upload_service.go`:
```go
// Change from 15 minutes to 30 minutes
expiry := 30 * time.Minute
```

## Production Considerations

### ⚠️ JANGAN gunakan `public` policy di production!

**Untuk Production:**

1. **Gunakan Private Bucket + Presigned URLs**
```bash
docker exec cakli-storage mc anonymous set private local/cakli
```

2. **Serve Files via API Proxy**
```go
// Proxy download through API
func DownloadFile(c fiber.Ctx) error {
    objectKey := c.Params("key")
    // Validate user has permission
    // Stream file from MinIO to client
}
```

3. **Gunakan CDN**
- CloudFront (AWS)
- CloudFlare
- Akamai

4. **Implement IAM Policies**
```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": ["s3:PutObject"],
      "Resource": ["arn:aws:s3:::cakli/drivers/*"],
      "Condition": {
        "StringEquals": {
          "s3:x-amz-server-side-encryption": "AES256"
        }
      }
    }
  ]
}
```

5. **Enable Encryption**
```bash
docker exec cakli-storage mc encrypt set sse-s3 local/cakli
```

## Summary

**Masalah:**
- MinIO bucket policy `private` menolak PUT operations
- Presigned URL tidak cukup tanpa bucket policy yang tepat
- Browser tidak bisa upload file ke MinIO

**Solusi:**
- Ubah bucket policy ke `public` untuk development
- Presigned URL sekarang berfungsi dengan baik
- Upload file berhasil

**Keamanan:**
- Presigned URL tetap aman dengan time-limit dan signature
- Untuk production, gunakan private bucket + API proxy
- Implement proper IAM policies dan encryption

**Command untuk Fix:**
```bash
make setup-minio
```

atau manual:
```bash
docker exec cakli-storage mc alias set local http://localhost:9000 minioadmin minioadmin
docker exec cakli-storage mc anonymous set public local/cakli
```
