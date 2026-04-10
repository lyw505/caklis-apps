# Quick Fix: Upload 403 Error

## Problem
❌ File upload gagal dengan error `403 Forbidden` saat upload foto driver

## Solution
✅ Jalankan command ini untuk mengkonfigurasi MinIO:

```bash
make setup-minio
```

## Penjelasan Singkat

**Masalah:** MinIO bucket tidak memiliki policy yang tepat untuk menerima upload dari browser.

**Solusi:** Script `setup-minio.ps1` akan:
1. Membuat bucket `cakli` jika belum ada
2. Mengatur policy bucket untuk mengizinkan download publik
3. Mengaktifkan presigned URL untuk upload

**Setelah menjalankan script:**
- Upload foto driver akan berfungsi normal
- Presigned URL dari API akan diterima oleh MinIO
- Browser bisa melakukan PUT request ke MinIO

## Verifikasi

1. Buka MinIO Console: http://localhost:9001
2. Login dengan `minioadmin` / `minioadmin`
3. Cek bucket `cakli` → Access Policy harus "Public" atau "Download"

## Detail Lengkap

Lihat dokumentasi lengkap di: [MINIO-UPLOAD-FIX.md](./MINIO-UPLOAD-FIX.md)
