 # 📘 API Documentation — CAKLI v1

> **Base URL:** `http://localhost:8080/api/v1`
> **Content-Type:** `application/json`
> **Auth:** JWT Bearer Token (kecuali endpoint publik)
> **Date Format:** ISO 8601 (`2026-04-01T09:00:00Z`)

---

## Daftar Isi

- [Konvensi](#konvensi)
- [1. Authentication](#1-authentication)
  - [1.1 Login](#11-login)
  - [1.2 Refresh Token](#12-refresh-token)
  - [1.3 Logout](#13-logout)
  - [1.4 Get Current Admin](#14-get-current-admin)
- [2. Dashboard](#2-dashboard)
  - [2.1 Get Dashboard Stats](#21-get-dashboard-stats)
- [3. Drivers](#3-drivers)
  - [3.1 List Drivers](#31-list-drivers)
  - [3.2 Get Driver Detail](#32-get-driver-detail)
  - [3.3 Create Driver](#33-create-driver)
  - [3.4 Update Driver](#34-update-driver)
  - [3.5 Delete Driver](#35-delete-driver)
  - [3.6 Update Verification Status](#36-update-verification-status)
- [4. Users](#4-users)
  - [4.1 List Users](#41-list-users)
  - [4.2 Get User Detail](#42-get-user-detail)
  - [4.3 Create User](#43-create-user)
  - [4.4 Update User](#44-update-user)
  - [4.5 Delete User](#45-delete-user)
- [5. Upload (Presigned URL)](#5-upload-presigned-url)
  - [5.1 Generate Upload URL](#51-generate-upload-url)
  - [5.2 Generate View URL](#52-generate-view-url)
- [6. Reference Data](#6-reference-data)
  - [6.1 List Banks](#61-list-banks)

---

## Konvensi

### Response Format

Semua response mengikuti format standar:

**Success Response:**
```json
{
  "success": true,
  "message": "Deskripsi aksi berhasil",
  "data": { ... }
}
```

**Success Response (dengan paginasi):**
```json
{
  "success": true,
  "message": "Data berhasil diambil",
  "data": [ ... ],
  "meta": {
    "page": 1,
    "limit": 10,
    "total": 57,
    "total_pages": 6
  }
}
```

**Error Response:**
```json
{
  "success": false,
  "message": "Deskripsi error yang user-friendly",
  "error": {
    "code": "VALIDATION_ERROR",
    "details": [
      {
        "field": "email",
        "message": "Email sudah terdaftar"
      }
    ]
  }
}
```

### Error Codes

| Code | HTTP Status | Deskripsi |
|------|-------------|-----------|
| `VALIDATION_ERROR` | 400 | Input tidak valid |
| `INVALID_CREDENTIALS` | 401 | Email atau password salah |
| `TOKEN_EXPIRED` | 401 | Access token sudah expired |
| `TOKEN_INVALID` | 401 | Token tidak valid / malformed |
| `REFRESH_TOKEN_EXPIRED` | 401 | Refresh token sudah expired |
| `REFRESH_TOKEN_INVALID` | 401 | Refresh token tidak valid / sudah direvoke |
| `UNAUTHORIZED` | 401 | Tidak memiliki akses |
| `ACCOUNT_INACTIVE` | 403 | Akun admin tidak aktif |
| `FORBIDDEN` | 403 | Tidak memiliki izin untuk aksi ini |
| `NOT_FOUND` | 404 | Resource tidak ditemukan |
| `CONFLICT` | 409 | Data duplikat (email/phone/NIK sudah ada) |
| `RATE_LIMITED` | 429 | Terlalu banyak request |
| `INTERNAL_ERROR` | 500 | Error internal server |

### Auth Header

Semua endpoint yang membutuhkan autentikasi harus menyertakan:
```
Authorization: Bearer <access_token>
```

---

## 1. Authentication

### 1.1 Login

Autentikasi admin menggunakan email dan password. Mengembalikan access token di response body dan menyimpan refresh token di httpOnly cookie.

| | |
|---|---|
| **Endpoint** | `POST /auth/login` |
| **Auth Required** | ❌ Tidak |
| **Rate Limit** | 5 request per menit per IP |

#### Request Body

| Field | Tipe | Wajib | Validasi |
|-------|------|:-----:|----------|
| `email` | string | ✅ | Format email valid |
| `password` | string | ✅ | Tidak boleh kosong |

```json
{
  "email": "admin@cakli.id",
  "password": "admin123"
}
```

#### Success Response — `200 OK`

```json
{
  "success": true,
  "message": "Login berhasil",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG1pbl9pZCI6IjU1MGU4NDAwLWUyOWItNDFkNC1hNzE2LTQ0NjY1NTQ0MDAwMCIsInJvbGUiOiJtYXN0ZXJfYWRtaW4iLCJleHAiOjE3NDM0OTI0MDB9.abc123",
    "expires_in": 900,
    "token_type": "Bearer",
    "admin": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "name": "Super Admin",
      "email": "admin@cakli.id",
      "role": "master_admin"
    }
  }
}
```

**Set-Cookie Header:**
```
Set-Cookie: refresh_token=eyJhbGci...; Path=/api/v1/auth; HttpOnly; Secure; SameSite=Strict; Max-Age=604800
```

#### Error Responses

**`401 Unauthorized` — Kredensial salah:**
```json
{
  "success": false,
  "message": "Email atau password salah",
  "error": {
    "code": "INVALID_CREDENTIALS"
  }
}
```

**`403 Forbidden` — Akun nonaktif:**
```json
{
  "success": false,
  "message": "Akun anda tidak aktif, hubungi administrator",
  "error": {
    "code": "ACCOUNT_INACTIVE"
  }
}
```

**`400 Bad Request` — Validasi gagal:**
```json
{
  "success": false,
  "message": "Input tidak valid",
  "error": {
    "code": "VALIDATION_ERROR",
    "details": [
      {
        "field": "email",
        "message": "Format email tidak valid"
      }
    ]
  }
}
```

**`429 Too Many Requests` — Rate limited:**
```json
{
  "success": false,
  "message": "Terlalu banyak percobaan login, coba lagi dalam 60 detik",
  "error": {
    "code": "RATE_LIMITED"
  }
}
```

---

### 1.2 Refresh Token

Mendapatkan access token baru menggunakan refresh token dari cookie. Refresh token lama direvoke dan diganti baru (rotation).

| | |
|---|---|
| **Endpoint** | `POST /auth/refresh` |
| **Auth Required** | ❌ Tidak (menggunakan cookie) |
| **Cookie Required** | ✅ `refresh_token` (httpOnly) |

#### Request Body

Tidak ada body. Refresh token diambil dari cookie.

#### Success Response — `200 OK`

```json
{
  "success": true,
  "message": "Token berhasil diperbarui",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiJ9.new_token_here.sig",
    "expires_in": 900,
    "token_type": "Bearer"
  }
}
```

**Set-Cookie Header:** Cookie refresh token baru (menggantikan yang lama).

#### Error Responses

**`401 Unauthorized` — Refresh token tidak ada / sudah direvoke:**
```json
{
  "success": false,
  "message": "Sesi tidak valid, silakan login kembali",
  "error": {
    "code": "REFRESH_TOKEN_INVALID"
  }
}
```

**`401 Unauthorized` — Refresh token expired:**
```json
{
  "success": false,
  "message": "Sesi sudah berakhir, silakan login kembali",
  "error": {
    "code": "REFRESH_TOKEN_EXPIRED"
  }
}
```

---

### 1.3 Logout

Merevoke refresh token di database dan menghapus cookie.

| | |
|---|---|
| **Endpoint** | `POST /auth/logout` |
| **Auth Required** | ✅ Bearer Token |
| **Cookie Required** | ✅ `refresh_token` (httpOnly) |

#### Request Body

Tidak ada body.

#### Success Response — `200 OK`

```json
{
  "success": true,
  "message": "Logout berhasil"
}
```

**Set-Cookie Header:** Cookie refresh token dihapus (Max-Age=0).

#### Error Responses

**`401 Unauthorized` — Token tidak valid:**
```json
{
  "success": false,
  "message": "Token tidak valid",
  "error": {
    "code": "TOKEN_INVALID"
  }
}
```

---

### 1.4 Get Current Admin

Mendapatkan profil admin yang sedang login berdasarkan access token.

| | |
|---|---|
| **Endpoint** | `GET /auth/me` |
| **Auth Required** | ✅ Bearer Token |

#### Success Response — `200 OK`

```json
{
  "success": true,
  "message": "Data admin berhasil diambil",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "Super Admin",
    "email": "admin@cakli.id",
    "role": "master_admin",
    "is_active": true,
    "created_at": "2026-04-01T02:00:00Z",
    "updated_at": "2026-04-01T02:00:00Z"
  }
}
```

#### Error Responses

**`401 Unauthorized` — Token expired:**
```json
{
  "success": false,
  "message": "Token sudah expired, silakan refresh",
  "error": {
    "code": "TOKEN_EXPIRED"
  }
}
```

---

## 2. Dashboard

### 2.1 Get Dashboard Stats

Mendapatkan statistik ringkasan untuk dashboard admin.

| | |
|---|---|
| **Endpoint** | `GET /admin/dashboard/stats` |
| **Auth Required** | ✅ Bearer Token |

#### Success Response — `200 OK`

```json
{
  "success": true,
  "message": "Statistik dashboard berhasil diambil",
  "data": {
    "drivers": {
      "total": 150,
      "verified": 120,
      "pending": 25,
      "rejected": 5
    },
    "users": {
      "total": 3500,
      "active": 3200
    }
  }
}
```

#### Error Responses

**`401 Unauthorized`** — Token tidak valid atau expired (lihat error umum di [Konvensi](#error-codes)).

---

## 3. Drivers

### 3.1 List Drivers

Mendapatkan daftar driver dengan paginasi dan pencarian.

| | |
|---|---|
| **Endpoint** | `GET /admin/drivers` |
| **Auth Required** | ✅ Bearer Token |

#### Query Parameters

| Parameter | Tipe | Default | Deskripsi |
|-----------|------|---------|-----------|
| `page` | integer | `1` | Halaman ke-n |
| `limit` | integer | `10` | Jumlah item per halaman (maks 100) |
| `search` | string | `""` | Pencarian di nama, email, NIK, atau no telp |
| `verification_status` | string | `""` | Filter: `pending` \| `accepted` \| `rejected` (kosong = semua) |
| `sort_by` | string | `created_at` | Kolom sort: `name` \| `email` \| `created_at` \| `verification_status` |
| `sort_order` | string | `desc` | Arah sort: `asc` \| `desc` |

#### Success Response — `200 OK`

```json
{
  "success": true,
  "message": "Daftar driver berhasil diambil",
  "data": [
    {
      "id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
      "name": "Sucipto Putra",
      "email": "sucipto@email.com",
      "phone": "081234567890",
      "nik": "3573012345678901",
      "verification_status": "accepted",
      "photo_profile_url": null,
      "is_active": true,
      "created_at": "2026-03-20T10:00:00Z"
    },
    {
      "id": "b2c3d4e5-f6a7-8901-bcde-f12345678901",
      "name": "Ahmad Fauzi",
      "email": "fauzi@email.com",
      "phone": "082345678901",
      "nik": "3573019876543210",
      "verification_status": "pending",
      "photo_profile_url": "https://minio.example.com/cakli/drivers/photo-profile/xxx.jpg?X-Amz-...",
      "is_active": true,
      "created_at": "2026-03-25T14:30:00Z"
    }
  ],
  "meta": {
    "page": 1,
    "limit": 10,
    "total": 57,
    "total_pages": 6
  }
}
```

> **Catatan:** `photo_profile_url` berisi **presigned GET URL** yang di-generate oleh backend saat response. URL ini bersifat sementara (expire 1 jam). Jika driver tidak punya foto profile, nilainya `null`.

---

### 3.2 Get Driver Detail

Mendapatkan detail lengkap seorang driver termasuk semua presigned URL untuk foto.

| | |
|---|---|
| **Endpoint** | `GET /admin/drivers/:id` |
| **Auth Required** | ✅ Bearer Token |

#### Path Parameters

| Parameter | Tipe | Deskripsi |
|-----------|------|-----------|
| `id` | UUID | ID driver |

#### Success Response — `200 OK`

```json
{
  "success": true,
  "message": "Detail driver berhasil diambil",
  "data": {
    "id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
    "name": "Sucipto Putra",
    "email": "sucipto@email.com",
    "phone": "081234567890",
    "nik": "3573012345678901",
    "birth_place": "Malang",
    "birth_date": "1975-05-15",
    "bank": {
      "id": 1,
      "name": "Bank BRI",
      "code": "002"
    },
    "bank_account_number": "1234567890123456",
    "photo_profile_url": "https://minio.example.com/cakli/drivers/photo-profile/a1b2c3d4.jpg?X-Amz-...",
    "photo_ktp_url": "https://minio.example.com/cakli/drivers/photo-ktp/a1b2c3d4.jpg?X-Amz-...",
    "photo_face_url": "https://minio.example.com/cakli/drivers/photo-face/a1b2c3d4.jpg?X-Amz-...",
    "verification_status": "accepted",
    "is_active": true,
    "created_at": "2026-03-20T10:00:00Z",
    "updated_at": "2026-03-28T16:45:00Z"
  }
}
```

#### Error Responses

**`404 Not Found` — Driver tidak ditemukan:**
```json
{
  "success": false,
  "message": "Driver tidak ditemukan",
  "error": {
    "code": "NOT_FOUND"
  }
}
```

---

### 3.3 Create Driver

Membuat data driver baru. Foto-foto sudah di-upload terlebih dahulu via presigned URL, lalu object key dikirim di sini.

| | |
|---|---|
| **Endpoint** | `POST /admin/drivers` |
| **Auth Required** | ✅ Bearer Token |

#### Request Body

| Field | Tipe | Wajib | Validasi |
|-------|------|:-----:|----------|
| `name` | string | ✅ | Maks 255 karakter |
| `email` | string | ✅ | Format email valid, unik |
| `phone` | string | ✅ | Format nomor telepon, unik |
| `password` | string | ✅ | Min 8 karakter |
| `nik` | string | ✅ | Tepat 16 digit angka, unik |
| `birth_place` | string | ✅ | Maks 100 karakter |
| `birth_date` | string | ✅ | Format: `YYYY-MM-DD` |
| `bank_id` | integer | ✅ | ID bank yang valid (dari daftar bank) |
| `bank_account_number` | string | ✅ | Angka saja, maks 30 karakter |
| `photo_profile_key` | string | ❌ | MinIO object key (dari presigned upload) |
| `photo_ktp_key` | string | ❌ | MinIO object key |
| `photo_face_key` | string | ❌ | MinIO object key |
| `verification_status` | string | ❌ | `pending` (default) \| `accepted` \| `rejected` |

```json
{
  "name": "Sucipto Putra",
  "email": "sucipto@email.com",
  "phone": "081234567890",
  "password": "securePass123",
  "nik": "3573012345678901",
  "birth_place": "Malang",
  "birth_date": "1975-05-15",
  "bank_id": 1,
  "bank_account_number": "1234567890123456",
  "photo_profile_key": "drivers/photo-profile/a1b2c3d4-uuid.jpg",
  "photo_ktp_key": "drivers/photo-ktp/a1b2c3d4-uuid.jpg",
  "photo_face_key": "drivers/photo-face/a1b2c3d4-uuid.jpg",
  "verification_status": "pending"
}
```

#### Success Response — `201 Created`

```json
{
  "success": true,
  "message": "Driver berhasil ditambahkan",
  "data": {
    "id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
    "name": "Sucipto Putra",
    "email": "sucipto@email.com",
    "phone": "081234567890",
    "nik": "3573012345678901",
    "birth_place": "Malang",
    "birth_date": "1975-05-15",
    "bank": {
      "id": 1,
      "name": "Bank BRI",
      "code": "002"
    },
    "bank_account_number": "1234567890123456",
    "photo_profile_url": null,
    "photo_ktp_url": null,
    "photo_face_url": null,
    "verification_status": "pending",
    "is_active": true,
    "created_at": "2026-04-01T09:00:00Z",
    "updated_at": "2026-04-01T09:00:00Z"
  }
}
```

#### Error Responses

**`400 Bad Request` — Validasi gagal:**
```json
{
  "success": false,
  "message": "Input tidak valid",
  "error": {
    "code": "VALIDATION_ERROR",
    "details": [
      {
        "field": "nik",
        "message": "NIK harus tepat 16 digit angka"
      },
      {
        "field": "password",
        "message": "Password minimal 8 karakter"
      }
    ]
  }
}
```

**`409 Conflict` — Data duplikat:**
```json
{
  "success": false,
  "message": "Data sudah terdaftar",
  "error": {
    "code": "CONFLICT",
    "details": [
      {
        "field": "email",
        "message": "Email sudah digunakan oleh driver lain"
      }
    ]
  }
}
```

---

### 3.4 Update Driver

Mengupdate data driver yang sudah ada. Hanya kirim field yang ingin diubah.

| | |
|---|---|
| **Endpoint** | `PUT /admin/drivers/:id` |
| **Auth Required** | ✅ Bearer Token |

#### Path Parameters

| Parameter | Tipe | Deskripsi |
|-----------|------|-----------|
| `id` | UUID | ID driver |

#### Request Body

Semua field sama dengan Create Driver, tapi **semua opsional**. Hanya field yang dikirim yang akan diupdate.

| Field | Tipe | Wajib | Catatan |
|-------|------|:-----:|---------|
| `name` | string | ❌ | |
| `email` | string | ❌ | Unik |
| `phone` | string | ❌ | Unik |
| `password` | string | ❌ | Jika kosong / tidak dikirim, password tidak berubah |
| `nik` | string | ❌ | Unik |
| `birth_place` | string | ❌ | |
| `birth_date` | string | ❌ | Format: `YYYY-MM-DD` |
| `bank_id` | integer | ❌ | |
| `bank_account_number` | string | ❌ | |
| `photo_profile_key` | string | ❌ | MinIO object key baru |
| `photo_ktp_key` | string | ❌ | MinIO object key baru |
| `photo_face_key` | string | ❌ | MinIO object key baru |
| `verification_status` | string | ❌ | `pending` \| `accepted` \| `rejected` |

```json
{
  "name": "Sucipto Putra Wijaya",
  "phone": "081234567899",
  "photo_profile_key": "drivers/photo-profile/new-uuid.jpg"
}
```

#### Success Response — `200 OK`

```json
{
  "success": true,
  "message": "Driver berhasil diperbarui",
  "data": {
    "id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
    "name": "Sucipto Putra Wijaya",
    "email": "sucipto@email.com",
    "phone": "081234567899",
    "nik": "3573012345678901",
    "birth_place": "Malang",
    "birth_date": "1975-05-15",
    "bank": {
      "id": 1,
      "name": "Bank BRI",
      "code": "002"
    },
    "bank_account_number": "1234567890123456",
    "photo_profile_url": "https://minio.example.com/cakli/drivers/photo-profile/new-uuid.jpg?X-Amz-...",
    "photo_ktp_url": "https://minio.example.com/cakli/drivers/photo-ktp/a1b2c3d4-uuid.jpg?X-Amz-...",
    "photo_face_url": "https://minio.example.com/cakli/drivers/photo-face/a1b2c3d4-uuid.jpg?X-Amz-...",
    "verification_status": "accepted",
    "is_active": true,
    "created_at": "2026-03-20T10:00:00Z",
    "updated_at": "2026-04-01T09:30:00Z"
  }
}
```

#### Error Responses

**`404 Not Found`** — Driver tidak ditemukan.
**`400 Bad Request`** — Validasi gagal (sama seperti Create).
**`409 Conflict`** — Data duplikat (sama seperti Create).

---

### 3.5 Delete Driver

Menghapus driver (soft delete — set `deleted_at` dan `is_active = false`).

| | |
|---|---|
| **Endpoint** | `DELETE /admin/drivers/:id` |
| **Auth Required** | ✅ Bearer Token |

#### Path Parameters

| Parameter | Tipe | Deskripsi |
|-----------|------|-----------|
| `id` | UUID | ID driver |

#### Success Response — `200 OK`

```json
{
  "success": true,
  "message": "Driver berhasil dihapus"
}
```

#### Error Responses

**`404 Not Found` — Driver tidak ditemukan:**
```json
{
  "success": false,
  "message": "Driver tidak ditemukan",
  "error": {
    "code": "NOT_FOUND"
  }
}
```

---

### 3.6 Update Verification Status

Mengubah status verifikasi driver. Endpoint khusus untuk aksi klik chip status di data table.

| | |
|---|---|
| **Endpoint** | `PATCH /admin/drivers/:id/verification-status` |
| **Auth Required** | ✅ Bearer Token |

#### Path Parameters

| Parameter | Tipe | Deskripsi |
|-----------|------|-----------|
| `id` | UUID | ID driver |

#### Request Body

| Field | Tipe | Wajib | Validasi |
|-------|------|:-----:|----------|
| `verification_status` | string | ✅ | `pending` \| `accepted` \| `rejected` |

```json
{
  "verification_status": "accepted"
}
```

#### Success Response — `200 OK`

```json
{
  "success": true,
  "message": "Status verifikasi driver berhasil diperbarui",
  "data": {
    "id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
    "name": "Sucipto Putra",
    "verification_status": "accepted",
    "updated_at": "2026-04-01T09:35:00Z"
  }
}
```

#### Error Responses

**`400 Bad Request` — Status tidak valid:**
```json
{
  "success": false,
  "message": "Input tidak valid",
  "error": {
    "code": "VALIDATION_ERROR",
    "details": [
      {
        "field": "verification_status",
        "message": "Status harus salah satu dari: pending, accepted, rejected"
      }
    ]
  }
}
```

**`404 Not Found`** — Driver tidak ditemukan.

---

## 4. Users

### 4.1 List Users

Mendapatkan daftar user dengan paginasi dan pencarian.

| | |
|---|---|
| **Endpoint** | `GET /admin/users` |
| **Auth Required** | ✅ Bearer Token |

#### Query Parameters

| Parameter | Tipe | Default | Deskripsi |
|-----------|------|---------|-----------|
| `page` | integer | `1` | Halaman ke-n |
| `limit` | integer | `10` | Jumlah item per halaman (maks 100) |
| `search` | string | `""` | Pencarian di nama, email, atau no telp |
| `sort_by` | string | `created_at` | Kolom sort: `name` \| `email` \| `created_at` |
| `sort_order` | string | `desc` | Arah sort: `asc` \| `desc` |

#### Success Response — `200 OK`

```json
{
  "success": true,
  "message": "Daftar user berhasil diambil",
  "data": [
    {
      "id": "c3d4e5f6-a7b8-9012-cdef-123456789012",
      "name": "Aulia Sukma",
      "email": "aulia@email.com",
      "phone": "085678901234",
      "photo_profile_url": "https://minio.example.com/cakli/users/photo-profile/c3d4e5f6.jpg?X-Amz-...",
      "is_active": true,
      "created_at": "2026-03-15T08:00:00Z"
    },
    {
      "id": "d4e5f6a7-b8c9-0123-def0-234567890123",
      "name": "Budi Santoso",
      "email": "budi@email.com",
      "phone": "087890123456",
      "photo_profile_url": null,
      "is_active": true,
      "created_at": "2026-03-18T12:30:00Z"
    }
  ],
  "meta": {
    "page": 1,
    "limit": 10,
    "total": 35,
    "total_pages": 4
  }
}
```

---

### 4.2 Get User Detail

Mendapatkan detail lengkap seorang user.

| | |
|---|---|
| **Endpoint** | `GET /admin/users/:id` |
| **Auth Required** | ✅ Bearer Token |

#### Path Parameters

| Parameter | Tipe | Deskripsi |
|-----------|------|-----------|
| `id` | UUID | ID user |

#### Success Response — `200 OK`

```json
{
  "success": true,
  "message": "Detail user berhasil diambil",
  "data": {
    "id": "c3d4e5f6-a7b8-9012-cdef-123456789012",
    "name": "Aulia Sukma",
    "email": "aulia@email.com",
    "phone": "085678901234",
    "photo_profile_url": "https://minio.example.com/cakli/users/photo-profile/c3d4e5f6.jpg?X-Amz-...",
    "is_active": true,
    "created_at": "2026-03-15T08:00:00Z",
    "updated_at": "2026-03-30T18:00:00Z"
  }
}
```

#### Error Responses

**`404 Not Found`** — User tidak ditemukan.

---

### 4.3 Create User

Membuat data user baru.

| | |
|---|---|
| **Endpoint** | `POST /admin/users` |
| **Auth Required** | ✅ Bearer Token |

#### Request Body

| Field | Tipe | Wajib | Validasi |
|-------|------|:-----:|----------|
| `name` | string | ✅ | Maks 255 karakter |
| `email` | string | ✅ | Format email valid, unik |
| `phone` | string | ✅ | Format nomor telepon, unik |
| `password` | string | ✅ | Min 8 karakter |
| `photo_profile_key` | string | ❌ | MinIO object key (dari presigned upload) |

```json
{
  "name": "Aulia Sukma",
  "email": "aulia@email.com",
  "phone": "085678901234",
  "password": "userPass123",
  "photo_profile_key": "users/photo-profile/c3d4e5f6-uuid.jpg"
}
```

#### Success Response — `201 Created`

```json
{
  "success": true,
  "message": "User berhasil ditambahkan",
  "data": {
    "id": "c3d4e5f6-a7b8-9012-cdef-123456789012",
    "name": "Aulia Sukma",
    "email": "aulia@email.com",
    "phone": "085678901234",
    "photo_profile_url": null,
    "is_active": true,
    "created_at": "2026-04-01T09:00:00Z",
    "updated_at": "2026-04-01T09:00:00Z"
  }
}
```

#### Error Responses

**`400 Bad Request`** — Validasi gagal.
**`409 Conflict`** — Email atau nomor telepon sudah terdaftar.

---

### 4.4 Update User

Mengupdate data user. Hanya kirim field yang ingin diubah.

| | |
|---|---|
| **Endpoint** | `PUT /admin/users/:id` |
| **Auth Required** | ✅ Bearer Token |

#### Path Parameters

| Parameter | Tipe | Deskripsi |
|-----------|------|-----------|
| `id` | UUID | ID user |

#### Request Body

| Field | Tipe | Wajib | Catatan |
|-------|------|:-----:|---------|
| `name` | string | ❌ | |
| `email` | string | ❌ | Unik |
| `phone` | string | ❌ | Unik |
| `password` | string | ❌ | Jika kosong / tidak dikirim, password tidak berubah |
| `photo_profile_key` | string | ❌ | MinIO object key baru |

```json
{
  "name": "Aulia Sukma Ramadhani",
  "phone": "085678901999"
}
```

#### Success Response — `200 OK`

```json
{
  "success": true,
  "message": "User berhasil diperbarui",
  "data": {
    "id": "c3d4e5f6-a7b8-9012-cdef-123456789012",
    "name": "Aulia Sukma Ramadhani",
    "email": "aulia@email.com",
    "phone": "085678901999",
    "photo_profile_url": "https://minio.example.com/cakli/users/photo-profile/c3d4e5f6.jpg?X-Amz-...",
    "is_active": true,
    "created_at": "2026-03-15T08:00:00Z",
    "updated_at": "2026-04-01T10:00:00Z"
  }
}
```

#### Error Responses

**`404 Not Found`** — User tidak ditemukan.
**`400 Bad Request`** — Validasi gagal.
**`409 Conflict`** — Data duplikat.

---

### 4.5 Delete User

Menghapus user (soft delete).

| | |
|---|---|
| **Endpoint** | `DELETE /admin/users/:id` |
| **Auth Required** | ✅ Bearer Token |

#### Path Parameters

| Parameter | Tipe | Deskripsi |
|-----------|------|-----------|
| `id` | UUID | ID user |

#### Success Response — `200 OK`

```json
{
  "success": true,
  "message": "User berhasil dihapus"
}
```

#### Error Responses

**`404 Not Found`** — User tidak ditemukan.

---

## 5. Upload (Presigned URL)

### 5.1 Generate Upload URL

Menghasilkan presigned PUT URL untuk upload file langsung ke MinIO. Frontend menggunakan URL ini untuk upload file tanpa melalui backend.

| | |
|---|---|
| **Endpoint** | `POST /upload/presigned-url` |
| **Auth Required** | ✅ Bearer Token |

#### Request Body

| Field | Tipe | Wajib | Validasi |
|-------|------|:-----:|----------|
| `filename` | string | ✅ | Nama file asli (untuk extension) |
| `content_type` | string | ✅ | MIME type: `image/jpeg` \| `image/png` |
| `folder` | string | ✅ | Folder tujuan: `drivers/photo-profile` \| `drivers/photo-ktp` \| `drivers/photo-face` \| `users/photo-profile` |

```json
{
  "filename": "ktp-sucipto.jpg",
  "content_type": "image/jpeg",
  "folder": "drivers/photo-ktp"
}
```

#### Success Response — `200 OK`

```json
{
  "success": true,
  "message": "Presigned URL berhasil dibuat",
  "data": {
    "upload_url": "https://minio.example.com/cakli/drivers/photo-ktp/f47ac10b-58cc-4372-a567-0e02b2c3d479.jpg?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=...&X-Amz-Date=...&X-Amz-Expires=900&X-Amz-SignedHeaders=host&X-Amz-Signature=...",
    "object_key": "drivers/photo-ktp/f47ac10b-58cc-4372-a567-0e02b2c3d479.jpg",
    "expires_in": 900
  }
}
```

> **Cara pakai di frontend:**
> ```javascript
> // 1. Dapatkan presigned URL dari API
> const { upload_url, object_key } = response.data;
>
> // 2. Upload langsung ke MinIO
> await fetch(upload_url, {
>   method: 'PUT',
>   headers: { 'Content-Type': file.type },
>   body: file
> });
>
> // 3. Kirim object_key saat create/update driver
> await createDriver({ ...formData, photo_ktp_key: object_key });
> ```

#### Error Responses

**`400 Bad Request` — Content type tidak didukung:**
```json
{
  "success": false,
  "message": "Input tidak valid",
  "error": {
    "code": "VALIDATION_ERROR",
    "details": [
      {
        "field": "content_type",
        "message": "Hanya mendukung image/jpeg dan image/png"
      }
    ]
  }
}
```

**`400 Bad Request` — Folder tidak valid:**
```json
{
  "success": false,
  "message": "Input tidak valid",
  "error": {
    "code": "VALIDATION_ERROR",
    "details": [
      {
        "field": "folder",
        "message": "Folder tidak valid. Pilihan: drivers/photo-profile, drivers/photo-ktp, drivers/photo-face, users/photo-profile"
      }
    ]
  }
}
```

---

### 5.2 Generate View URL

Menghasilkan presigned GET URL untuk melihat/menampilkan file dari MinIO.

| | |
|---|---|
| **Endpoint** | `GET /upload/presigned-view-url` |
| **Auth Required** | ✅ Bearer Token |

#### Query Parameters

| Parameter | Tipe | Wajib | Deskripsi |
|-----------|------|:-----:|-----------|
| `key` | string | ✅ | MinIO object key (contoh: `drivers/photo-ktp/uuid.jpg`) |

#### Success Response — `200 OK`

```json
{
  "success": true,
  "message": "Presigned view URL berhasil dibuat",
  "data": {
    "view_url": "https://minio.example.com/cakli/drivers/photo-ktp/f47ac10b.jpg?X-Amz-Algorithm=AWS4-HMAC-SHA256&...",
    "expires_in": 3600
  }
}
```

#### Error Responses

**`400 Bad Request` — Key tidak valid atau kosong:**
```json
{
  "success": false,
  "message": "Input tidak valid",
  "error": {
    "code": "VALIDATION_ERROR",
    "details": [
      {
        "field": "key",
        "message": "Object key wajib diisi"
      }
    ]
  }
}
```

**`404 Not Found` — File tidak ada di MinIO:**
```json
{
  "success": false,
  "message": "File tidak ditemukan",
  "error": {
    "code": "NOT_FOUND"
  }
}
```

> **Catatan:** Pada endpoint List dan Detail (drivers & users), backend sudah otomatis mengkonversi object key menjadi presigned view URL di response. Endpoint ini hanya digunakan jika frontend perlu me-refresh URL yang sudah expired.

---

## 6. Reference Data

### 6.1 List Banks

Mendapatkan daftar bank yang aktif untuk dropdown pilihan rekening driver.

| | |
|---|---|
| **Endpoint** | `GET /admin/banks` |
| **Auth Required** | ✅ Bearer Token |

#### Success Response — `200 OK`

```json
{
  "success": true,
  "message": "Daftar bank berhasil diambil",
  "data": [
    {
      "id": 1,
      "name": "Bank BRI",
      "code": "002"
    },
    {
      "id": 2,
      "name": "Bank Mandiri",
      "code": "008"
    },
    {
      "id": 3,
      "name": "Bank BNI",
      "code": "009"
    },
    {
      "id": 4,
      "name": "Bank BCA",
      "code": "014"
    },
    {
      "id": 5,
      "name": "Bank BTN",
      "code": "200"
    },
    {
      "id": 6,
      "name": "Bank CIMB Niaga",
      "code": "022"
    },
    {
      "id": 7,
      "name": "Bank Danamon",
      "code": "011"
    },
    {
      "id": 8,
      "name": "Bank Permata",
      "code": "013"
    },
    {
      "id": 9,
      "name": "Bank Panin",
      "code": "019"
    },
    {
      "id": 10,
      "name": "Bank OCBC NISP",
      "code": "028"
    },
    {
      "id": 11,
      "name": "Bank Mega",
      "code": "426"
    },
    {
      "id": 12,
      "name": "Bank Bukopin",
      "code": "441"
    },
    {
      "id": 13,
      "name": "Bank Sinarmas",
      "code": "153"
    },
    {
      "id": 14,
      "name": "Bank Maybank",
      "code": "016"
    },
    {
      "id": 15,
      "name": "Bank BSI",
      "code": "451"
    },
    {
      "id": 16,
      "name": "Bank Jago",
      "code": "542"
    },
    {
      "id": 17,
      "name": "Bank Neo Commerce",
      "code": "490"
    },
    {
      "id": 18,
      "name": "Bank Seabank",
      "code": "535"
    },
    {
      "id": 19,
      "name": "Bank BTPN",
      "code": "213"
    },
    {
      "id": 20,
      "name": "Bank DKI",
      "code": "111"
    },
    {
      "id": 21,
      "name": "Bank Jatim",
      "code": "114"
    },
    {
      "id": 22,
      "name": "Bank Jateng",
      "code": "113"
    }
  ]
}
```

---

## Ringkasan Endpoint

| # | Method | Endpoint | Deskripsi | Auth |
|---|--------|----------|-----------|:----:|
| 1 | `POST` | `/auth/login` | Login admin | ❌ |
| 2 | `POST` | `/auth/refresh` | Refresh access token | ❌* |
| 3 | `POST` | `/auth/logout` | Logout admin | ✅ |
| 4 | `GET` | `/auth/me` | Get current admin profile | ✅ |
| 5 | `GET` | `/admin/dashboard/stats` | Get dashboard statistics | ✅ |
| 6 | `GET` | `/admin/drivers` | List drivers (paginated) | ✅ |
| 7 | `GET` | `/admin/drivers/:id` | Get driver detail | ✅ |
| 8 | `POST` | `/admin/drivers` | Create driver | ✅ |
| 9 | `PUT` | `/admin/drivers/:id` | Update driver | ✅ |
| 10 | `DELETE` | `/admin/drivers/:id` | Delete driver (soft) | ✅ |
| 11 | `PATCH` | `/admin/drivers/:id/verification-status` | Update verification status | ✅ |
| 12 | `GET` | `/admin/users` | List users (paginated) | ✅ |
| 13 | `GET` | `/admin/users/:id` | Get user detail | ✅ |
| 14 | `POST` | `/admin/users` | Create user | ✅ |
| 15 | `PUT` | `/admin/users/:id` | Update user | ✅ |
| 16 | `DELETE` | `/admin/users/:id` | Delete user (soft) | ✅ |
| 17 | `POST` | `/upload/presigned-url` | Generate presigned upload URL | ✅ |
| 18 | `GET` | `/upload/presigned-view-url` | Generate presigned view URL | ✅ |
| 19 | `GET` | `/admin/banks` | List available banks | ✅ |

**Total: 19 endpoints**

> \* Endpoint refresh menggunakan refresh token dari httpOnly cookie, bukan Bearer token.
