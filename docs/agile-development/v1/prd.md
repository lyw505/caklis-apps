# 📋 PRD — CAKLI Web Admin v1

> **Sprint:** v1 (MVP Agile)  
> **Tanggal:** 1 April 2026  
> **Fokus:** Web Admin — Master Admin (Login, Dashboard, Kelola Driver, Kelola User)

---

## 1. Ringkasan

CAKLI v1 adalah **web admin dashboard** untuk mengelola data driver dan user platform transportasi becak listrik CAKLI. Sprint pertama ini berfokus pada fondasi sistem: autentikasi admin, dashboard statistik sederhana, serta CRUD data driver dan user.

---

## 2. Tech Stack

| Layer | Teknologi |
|-------|-----------|
| **Frontend** | Next.js (TypeScript) |
| **Backend API** | Go (Fiber + GORM) |
| **Database** | PostgreSQL |
| **Object Storage** | MinIO (presigned URL strategy) |
| **Auth** | JWT dual token (access token + refresh token) |
| **Container** | Docker + Docker Compose |
| **Dev Tools** | Makefile, Air (Go hot reload) |
| **Env Strategy** | Single `.env` di root, diinject via docker-compose untuk services, di-copy ke direktori web admin saat `make dev` |

---

## 3. Arsitektur Monorepo

```
caklis/
├── .env                          # Single env file (root)
├── Makefile                      # Perintah development
├── docker-compose.yml            # PostgreSQL, MinIO, Backend API
├── backend/                      # Go Fiber API
│   ├── main.go
│   ├── .air.toml
│   └── ...
├── web-admin/                    # Next.js Frontend
│   └── ...
└── docs/
    └── agile-development/v1/
        ├── prd.md
        ├── db.sql
        └── api.md
```

---

## 4. Development Workflow

### `make dev`

Menjalankan **seluruh stack** untuk development:

1. **Build & start Docker services** (PostgreSQL, MinIO, Backend API dengan Air hot reload)
2. **Copy `.env` dari root** ke direktori `web-admin/` sebagai `.env.local`
3. **Jalankan `npm run dev`** di direktori `web-admin/`

Urutan eksekusi:
```
make dev
  ├─ docker compose up -d --build   (postgres, minio, backend-api)
  ├─ cp .env ./web-admin/.env.local
  └─ cd web-admin && npm run dev
```

---

## 5. Strategi Object Storage (MinIO Presigned URL)

Semua upload file (foto KTP, foto muka, photo profile) menggunakan **presigned URL** agar tidak ada beban berlebih di backend.

### Flow Upload

```
1. Frontend → POST /api/v1/upload/presigned-url { filename, contentType, folder }
2. Backend  → Generate presigned PUT URL dari MinIO (expire 15 menit)
3. Backend  → Response { uploadUrl, objectKey }
4. Frontend → PUT file langsung ke MinIO via uploadUrl
5. Frontend → Kirim objectKey ke backend saat create/update driver/user
6. Saat menampilkan gambar → GET /api/v1/upload/presigned-view-url?key=xxx
7. Backend  → Generate presigned GET URL dari MinIO (expire 1 jam)
8. Frontend → Load image dari URL tersebut
```

**Keuntungan:**
- Backend tidak handle binary file → ringan
- Upload langsung ke MinIO → cepat
- URL sementara → aman

---

## 6. Autentikasi

### 6.1 Admin Login

- Admin login menggunakan **email** dan **password**
- Hanya role **Master Admin** untuk v1
- Password di-hash menggunakan **bcrypt**

### 6.2 JWT Dual Token

| Token | Lifetime | Deskripsi |
|-------|----------|-----------|
| **Access Token** | 15 menit | Dikirim di header `Authorization: Bearer <token>`, digunakan untuk setiap request API yang membutuhkan autentikasi |
| **Refresh Token** | 7 hari | Disimpan di **httpOnly secure cookie**, digunakan untuk mendapatkan access token baru |

### 6.3 Security Best Practices

- Access token **tidak disimpan di cookie**, hanya di memory frontend (state/context)
- Refresh token di **httpOnly cookie** (tidak bisa diakses JavaScript)
- Refresh token di-hash sebelum disimpan di database
- Saat refresh, token lama direvoke dan diganti baru (**rotation**)
- Logout menghapus refresh token dari database dan clear cookie
- Rate limiting pada endpoint login (max 5 percobaan per menit per IP)

---

## 7. Fitur v1

### 7.1 Halaman Login Admin

**Route:** `/login`

- Input email dan password
- Tombol "Masuk"
- Validasi frontend: email format, password required
- Error handling: kredensial salah, akun nonaktif
- Redirect ke `/dashboard` setelah login berhasil
- Jika sudah login (token valid), redirect otomatis ke `/dashboard`

---

### 7.2 Dashboard

**Route:** `/dashboard`

**Deskripsi:** Halaman utama setelah login. Menampilkan statistik ringkasan data driver dan user.

#### Stats Cards
| Stat | Deskripsi |
|------|-----------|
| **Total Driver** | Jumlah seluruh driver terdaftar |
| **Driver Terverifikasi** | Jumlah driver dengan status verifikasi "diterima" |
| **Driver Pending** | Jumlah driver dengan status verifikasi "pending" |
| **Driver Ditolak** | Jumlah driver dengan status verifikasi "ditolak" |
| **Total User** | Jumlah seluruh user terdaftar |
| **User Aktif** | Jumlah user dengan status aktif |

---

### 7.3 Kelola Driver

**Route:** `/drivers`

**Deskripsi:** Halaman manajemen data driver dengan data table dan fitur CRUD.

#### Data Table

**Kolom yang ditampilkan (ringkas, tidak semua field):**

| Kolom | Deskripsi |
|-------|-----------|
| **Foto** | Thumbnail photo profile (circle avatar) |
| **Nama** | Nama lengkap driver |
| **NIK** | Nomor Induk Kependudukan (16 digit) |
| **No. Telp** | Nomor telepon |
| **Email** | Alamat email |
| **Status Verifikasi** | Chip/badge berwarna: 🟡 Pending / 🟢 Diterima / 🔴 Ditolak |
| **Aksi** | Tombol: View Detail, Edit, Delete |

**Fitur Data Table:**
- Pencarian (search by nama, NIK, email, no telp)
- Paginasi (server-side, 10 item per halaman)
- **Chip Status Verifikasi:** Klik chip → popover dropdown untuk mengganti status verifikasi (Pending / Diterima / Ditolak). Perubahan langsung ter-save ke backend.

**Tombol Tambah Driver:** Pojok kanan atas halaman → buka modal.

#### Modal Tambah Driver

| Field | Tipe | Validasi |
|-------|------|----------|
| **Photo Profile** | File upload (presigned URL) | JPG/PNG, maks 2MB |
| **Foto KTP** | File upload (presigned URL) | JPG/PNG, maks 5MB |
| **Foto Muka Driver** | File upload (presigned URL) | JPG/PNG, maks 5MB |
| **Nama** | Text input | Wajib, maks 255 karakter |
| **Email** | Text input | Wajib, format email valid, unik |
| **Password** | Password input | Wajib, min 8 karakter |
| **No. Telp** | Text input | Wajib, format nomor telepon, unik |
| **NIK** | Text input | Wajib, tepat 16 digit angka, unik |
| **Tempat Lahir** | Text input | Wajib |
| **Tanggal Lahir** | Date picker | Wajib |
| **Bank** | Dropdown select | Wajib, pilih dari daftar bank |
| **No. Rekening** | Text input | Wajib, angka saja |
| **Status Verifikasi** | Select | Pilih: Pending (default) / Diterima / Ditolak |

#### Modal Edit Driver
Sama dengan modal Tambah tapi pre-filled dengan data driver yang dipilih. Password opsional (kosong = tidak berubah).

#### Dialog View Detail Driver
Menampilkan semua data driver secara lengkap termasuk preview foto KTP, foto muka, dan photo profile (loaded via presigned GET URL).

#### Dialog Konfirmasi Delete
Konfirmasi sebelum menghapus driver (soft delete → set `is_active = false` dan `deleted_at = NOW()`).

---

### 7.4 Kelola User

**Route:** `/users`

**Deskripsi:** Halaman manajemen data user dengan data table dan fitur CRUD.

#### Data Table

**Kolom yang ditampilkan:**

| Kolom | Deskripsi |
|-------|-----------|
| **Foto** | Thumbnail photo profile (circle avatar) |
| **Nama** | Nama lengkap user |
| **Email** | Alamat email |
| **No. Telp** | Nomor telepon |
| **Aksi** | Tombol: View Detail, Edit, Delete |

**Fitur Data Table:**
- Pencarian (search by nama, email, no telp)
- Paginasi (server-side, 10 item per halaman)

**Tombol Tambah User:** Pojok kanan atas halaman → buka modal.

#### Modal Tambah User

| Field | Tipe | Validasi |
|-------|------|----------|
| **Photo Profile** | File upload (presigned URL) | JPG/PNG, maks 2MB |
| **Nama** | Text input | Wajib, maks 255 karakter |
| **Email** | Text input | Wajib, format email valid, unik |
| **No. Telp** | Text input | Wajib, format nomor telepon, unik |
| **Password** | Password input | Wajib, min 8 karakter |

#### Modal Edit User
Sama dengan modal Tambah tapi pre-filled. Password opsional.

#### Dialog View Detail User
Menampilkan semua data user lengkap termasuk preview photo profile.

#### Dialog Konfirmasi Delete
Konfirmasi sebelum menghapus user (soft delete).

---

## 8. Sidebar Navigasi

| Menu | Route | Ikon |
|------|-------|------|
| Dashboard | `/dashboard` | LayoutDashboard |
| Kelola Driver | `/drivers` | Car |
| Kelola User | `/users` | Users |

**Footer sidebar:** Info admin yang sedang login (nama + email) + tombol Logout.

---

## 9. Environment Variables

```env
# ===== Database =====
DB_HOST=localhost
DB_PORT=5432
DB_USER=cakli
DB_PASSWORD=cakli_secret
DB_NAME=cakli_db
DB_SSLMODE=disable

# ===== MinIO =====
MINIO_ENDPOINT=localhost:9000
MINIO_ACCESS_KEY=minioadmin
MINIO_SECRET_KEY=minioadmin
MINIO_BUCKET=cakli
MINIO_USE_SSL=false

# ===== JWT =====
JWT_ACCESS_SECRET=your-access-secret-key-change-in-production
JWT_REFRESH_SECRET=your-refresh-secret-key-change-in-production
JWT_ACCESS_EXPIRY=15m
JWT_REFRESH_EXPIRY=168h

# ===== Backend API =====
API_PORT=8080
API_ENV=development

# ===== Frontend =====
NEXT_PUBLIC_API_URL=http://localhost:8080/api/v1
NEXT_PUBLIC_APP_NAME=CAKLI Admin
```

---

## 10. Scope Exclusion (Tidak Termasuk v1)

Yang **TIDAK** termasuk dalam v1:
- ❌ Operating Admin & Reporting Admin
- ❌ Manajemen zona & tarif
- ❌ Manajemen pesanan/order
- ❌ Keluhan & sengketa
- ❌ Voucher management
- ❌ Peta operasional
- ❌ Log audit
- ❌ Kebijakan mitra
- ❌ Analitik & laporan detail
- ❌ Notification system
- ❌ Mobile app integration
- ❌ MFA / 2FA
- ❌ Real-time features (WebSocket)

Semua fitur di atas akan diimplementasikan di sprint berikutnya.
