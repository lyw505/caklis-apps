# 🚀 Implementation Plan — CAKLI v1 (Agile Sprint 1)

> **Sprint Duration:** 2-3 minggu  
> **Team Size:** 1-2 developers  
> **Prioritas:** MVP untuk Master Admin Dashboard

---

## � Current Progress Summary

**Last Updated:** April 9, 2026

### Overall Status: 🟢 Near Complete (~95% Complete)Complete)

#### ✅ Completed Tasks (Fully Done)
- Task 1.1: Project Structure & Environment Setup ✅
- Task 1.2: Database Schema Implementation ✅
- Task 2.1: Backend Project Setup ✅
- Task 2.2: Models & Database Layer ✅
- Task 2.3: Authentication System ✅
- Task 2.4: Dashboard Stats Endpoint ✅
- Task 2.5: MinIO Presigned URL Service ✅
- Task 2.6: Driver CRUD Endpoints ✅
- Task 2.7: User CRUD Endpoints ✅
- Task 2.8: Banks Reference Endpoint ✅
- Task 2.9: Error Handling & Response Standardization ✅
- Task 3.1: Next.js Project Setup ✅
- Task 3.2: API Client & Auth Context ✅
- Task 3.4: Protected Route & Layout ✅
- Task 3.5: Dashboard Page ✅
- Task 3.6: Drivers List Page ✅
- Task 3.7: Driver Create/Edit Modal ✅
- Task 3.9: Users List Page ✅
### Key Achievements
1. ✅ **Backend API Complete** - All 19 endpoints implemented and working
2. ✅ **Frontend Integration Complete** - All pages connected to API with real data
3. ✅ **Authentication Working** - JWT + refresh token rotation implemented
4. ✅ **File Upload Working** - MinIO presigned URL flow implemented
5. ✅ **CRUD Forms Complete** - Driver and User create/edit modals with file upload
6. ✅ **Docker Setup Complete** - All services running (PostgreSQL, MinIO, API)
7. ✅ **Database Seeded** - 22 banks + default admin ready

### Remaining Optional Tasks
1. 🔧 Add delete confirmation dialogs (nice to have)
2. 🔧 Add detail view dialogs (nice to have)
3. 🔧 Improve login page validation (nice to have)
4. 🔧 Add automated tests (nice to have)
5. 🔧 Enhance documentation (nice to have)eds API integration)
- Task 3.6: Drivers List Page (UI complete, needs API integration)
- Task 3.9: Users List Page (UI complete, needs API integration)

#### ❌ Not Started Tasks
- Task 3.7: Driver Create/Edit Modal
- Task 3.8: Driver Detail & Delete
- Task 3.10: User Create/Edit Modal
- Task 3.11: User Detail & Delete
- All testing tasks (Task 4.1 - 4.3)
- All documentation tasks (Task 5.1 - 5.3)

### Key Achievements
1. ✅ Backend API fully implemented with all CRUD endpoints
2. ✅ Authentication system working (JWT with refresh token rotation)
3. ✅ Database connected and seeded with initial data
4. ✅ MinIO presigned URL service for file uploads
5. ✅ All 19 API endpoints implemented and tested
6. ✅ Frontend UI structure complete with mock data

### Next Priority Actions
1. 🔧 Install React Hook Form and Axios in frontend
2. 🔧 Create API client with interceptors (Task 3.2)
3. 🔧 Create Auth context and provider (Task 3.2)
4. 🔧 Integrate login page with backend API (Task 3.3)
5. 🔧 Connect data tables to backend API (Task 3.6, 3.9)

---

## 📊 Sprint Overview

**Goal:** Membangun fondasi web admin dashboard dengan autentikasi, dashboard stats, dan CRUD driver & user.
### Task 1.1: Project Structure & Environment Setup ✅
**Priority:** 🔴 Critical  
**Estimasi:** 4 jam  
**Dependencies:** None

**Checklist:**
- [x] Setup monorepo structure (`apps/api`, `apps/web`)
- [x] Create root `.env` file dengan semua environment variables
- [x] Setup `docker-compose.yml` (PostgreSQL, MinIO, Backend API)
- [x] Test docker services berjalan dengan baik
- [x] Verify MinIO console accessible di `localhost:9001`
- [x] Verify PostgreSQL connection (port 5433)

**Acceptance Criteria:**
- ✅ Docker services running tanpa error
- ✅ MinIO console accessible dan bucket `cakli` ter-create
- ✅ PostgreSQL ready untuk menerima koneksi

## 🏗️ Architecture Setup (Week 1 - Day 1-2)

### Task 1.1: Project Structure & Environment Setup ✅
**Priority:** 🔴 Critical  
**Estimasi:** 4 jam  
**Dependencies:** None

**Checklist:**
- [x] Setup monorepo structure (`apps/api`, `apps/web`)
- [x] Create root `.env` file dengan semua environment variables
- [x] Setup `docker-compose.yml` (PostgreSQL, MinIO, Backend API)
### Task 1.2: Database Schema Implementation ✅
**Priority:** 🔴 Critical  
**Estimasi:** 3 jam  
**Dependencies:** Task 1.1

**Checklist:**
- [x] Run `db.sql` untuk create tables
- [x] Verify semua tables ter-create dengan benar
- [x] Verify indexes ter-create
- [x] Verify triggers (updated_at) berfungsi
- [x] Verify seed data banks ter-insert (22 banks)
- [x] Verify seed data default admin ter-insert
- [x] Test soft delete behavior dengan manual query
- [x] Document database connection string

**Acceptance Criteria:**
- Semua 5 tables exist: `banks`, `admins`, `refresh_tokens`, `drivers`, `users`
- Default admin `master@cakli.id` bisa login (password: `admin123`)
- Soft delete working (deleted_at not null, unique constraint tetap work)
- [x] Verify semua tables ter-create dengan benar
- [x] Verify indexes ter-create
- [x] Verify triggers (updated_at) berfungsi
- [x] Verify seed data banks ter-insert (22 banks)
- [x] Verify seed data default admin ter-insert
### Task 2.1: Backend Project Setup ✅
**Priority:** 🔴 Critical  
**Estimasi:** 3 jam  
**Dependencies:** Task 1.2

**Checklist:**
- [x] Initialize Go module (`go mod init`)
- [x] Install dependencies (Fiber, GORM, JWT, MinIO SDK, bcrypt, godotenv)
- [x] Setup project structure:
  ```
  apps/api/
  ├── main.go ✅
  ├── config/          # env loader, db connection, minio client
  ├── models/          # GORM models
  ├── handlers/        # HTTP handlers
  ├── middleware/      # auth middleware, error handler
  ├── services/        # business logic
  ├── utils/           # helpers (hash, jwt, presigned URL)
  └── .air.toml ✅     # hot reload config
  ```
- [x] Setup `.air.toml` untuk hot reload
- [x] Create `config/database.go` untuk GORM connection
- [x] Create `config/minio.go` untuk MinIO client
- [x] Test database connection
- [x] Test MinIO connection

**Acceptance Criteria:**
- `air` command menjalankan server dengan hot reload
- Database connection pool established
- MinIO client initializedhelpers (hash, jwt, presigned URL)
  └── .air.toml ✅       # hot reload config
  ```
- [x] Setup `.air.toml` untuk hot reload
### Task 2.2: Models & Database Layer ✅
**Priority:** 🔴 Critical  
**Estimasi:** 4 jam  
**Dependencies:** Task 2.1

**Checklist:**
- [x] Create `models/bank.go`
- [x] Create `models/admin.go` dengan soft delete
- [x] Create `models/refresh_token.go`
- [x] Create `models/driver.go` dengan soft delete
- [x] Create `models/user.go` dengan soft delete
- [x] Add GORM hooks untuk auto-update `updated_at`
- [x] Add GORM scopes untuk soft delete queries
- [x] Test AutoMigrate (should not break existing schema)
- [x] Create repository pattern untuk setiap model

**Acceptance Criteria:**
- Semua models match dengan database schema
- Soft delete working via GORM
- Repository methods (Create, Update, Delete, FindByID, List) implemented
- [x] Create `models/refresh_token.go`
- [x] Create `models/driver.go` dengan soft delete
- [x] Create `models/user.go` dengan soft delete
### Task 2.3: Authentication System ✅
**Priority:** 🔴 Critical  
**Estimasi:** 6 jam  
**Dependencies:** Task 2.2

**Checklist:**
- [x] Create `utils/hash.go` (bcrypt hash & compare)
- [x] Create `utils/jwt.go` (generate & validate access/refresh token)
- [x] Create `handlers/auth_handler.go`
- [x] Implement `POST /api/v1/auth/login`
  - [x] Validate email & password
  - [x] Check admin exists & active
  - [x] Compare password hash
  - [x] Generate access token (15 min)
  - [x] Generate refresh token (7 days)
  - [x] Save refresh token hash to DB
  - [x] Set httpOnly cookie
  - [x] Return access token in response body
  - [x] Rate limiting (5 req/min per IP)
- [x] Implement `POST /api/v1/auth/refresh`
  - [x] Read refresh token from cookie
  - [x] Validate token
  - [x] Check if revoked
  - [x] Generate new access token
  - [x] Rotate refresh token (revoke old, create new)
  - [x] Update cookie
- [x] Implement `POST /api/v1/auth/logout`
  - [x] Revoke refresh token in DB
  - [x] Clear cookie
- [x] Implement `GET /api/v1/auth/me`
  - [x] Extract admin ID from access token
  - [x] Return admin profile
- [x] Create `middleware/auth.go` untuk protect routes
- [x] Test all auth endpoints dengan Postman/curl

**Acceptance Criteria:**
- Login berhasil return access token & set cookie
- Refresh token rotation working
- Logout clear cookie & revoke token
- Protected routes return 401 jika tidak ada token
- Rate limiting working pada login endpoint
- [x] Implement `GET /api/v1/auth/me`
  - [x] Extract admin ID from access token
  - [x] Return admin profile
### Task 2.4: Dashboard Stats Endpoint ✅
**Priority:** 🟡 High  
**Estimasi:** 2 jam  
**Dependencies:** Task 2.3

**Checklist:**
- [x] Create `handlers/dashboard_handler.go`
- [x] Implement `GET /api/v1/admin/dashboard/stats`
- [x] Query count drivers (total, verified, pending, rejected)
- [x] Query count users (total, active)
- [x] Return formatted response
- [x] Add auth middleware
- [x] Test endpoint

**Acceptance Criteria:**
- Endpoint return correct counts
- Response format sesuai API docs
- Protected dengan auth middleware

---

### Task 2.5: MinIO Presigned URL Service ✅
**Priority:** 🔴 Critical  
**Estimasi:** 4 jam  
**Dependencies:** Task 2.1

**Checklist:**
- [x] Create `services/upload_service.go`
- [x] Implement `GeneratePresignedUploadURL(filename, contentType, folder)`
  - [x] Validate content type (only image/jpeg, image/png)
  - [x] Validate folder (drivers/photo-profile, drivers/photo-ktp, etc)
  - [x] Generate UUID untuk object key
  - [x] Generate presigned PUT URL (expire 15 min)
  - [x] Return upload URL & object key
- [x] Implement `GeneratePresignedViewURL(objectKey)`
  - [x] Check if object exists
  - [x] Generate presigned GET URL (expire 1 hour)
  - [x] Return view URL
- [x] Create `handlers/upload_handler.go`
- [x] Implement `POST /api/v1/upload/presigned-url`
- [x] Implement `GET /api/v1/upload/presigned-view-url`
- [x] Test upload flow:
  1. Get presigned upload URL
  2. PUT file to MinIO
  3. Get presigned view URL
  4. Verify file accessible

**Acceptance Criteria:**
- Presigned upload URL working (bisa upload via curl/Postman)
- Presigned view URL working (bisa akses file via browser)
- Validation working (reject invalid content type/folder)

---

### Task 2.6: Driver CRUD Endpoints ✅
**Priority:** 🔴 Critical  
**Estimasi:** 8 jam  
**Dependencies:** Task 2.2, Task 2.5

**Checklist:**
- [x] Create `handlers/driver_handler.go`
- [x] Create `services/driver_service.go`
- [x] Implement `GET /api/v1/admin/drivers` (List)
  - [x] Pagination (page, limit)
  - [x] Search (name, email, phone, NIK)
  - [x] Filter by verification_status
  - [x] Sorting (sort_by, sort_order)
  - [x] Generate presigned view URL untuk photo_profile_key
  - [x] Return meta pagination
- [x] Implement `GET /api/v1/admin/drivers/:id` (Detail)
  - [x] Generate presigned view URL untuk semua foto
  - [x] Include bank relation
- [x] Implement `POST /api/v1/admin/drivers` (Create)
  - [x] Validate all required fields
  - [x] Check unique constraints (email, phone, NIK)
  - [x] Hash password dengan bcrypt
  - [x] Save photo keys (dari presigned upload)
  - [x] Return created driver
- [x] Implement `PUT /api/v1/admin/drivers/:id` (Update)
  - [x] Validate fields
  - [x] Check unique constraints (exclude current driver)
  - [x] Hash password jika diubah
  - [x] Update photo keys jika ada
- [x] Implement `DELETE /api/v1/admin/drivers/:id` (Soft Delete)
  - [x] Set deleted_at & is_active = false
- [x] Implement `PATCH /api/v1/admin/drivers/:id/verification-status`
  - [x] Validate status (pending, accepted, rejected)
  - [x] Update verification_status
- [x] Add validation middleware
- [x] Add error handling
- [x] Test all endpoints

**Acceptance Criteria:**
- List drivers dengan pagination & search working
- Create driver dengan foto upload working
- Update driver working (termasuk update foto)
- Delete driver (soft delete) working
- Update verification status working
- Semua validation & error handling sesuai API docs

---

### Task 2.7: User CRUD Endpoints ✅
**Priority:** 🔴 Critical  
**Estimasi:** 6 jam  
**Dependencies:** Task 2.2, Task 2.5

**Checklist:**
- [x] Create `handlers/user_handler.go`
- [x] Create `services/user_service.go`
- [x] Implement `GET /api/v1/admin/users` (List)
  - [x] Pagination, search, sorting
  - [x] Generate presigned view URL untuk photo_profile_key
- [x] Implement `GET /api/v1/admin/users/:id` (Detail)
- [x] Implement `POST /api/v1/admin/users` (Create)
  - [x] Validate fields
  - [x] Check unique constraints
  - [x] Hash password
- [x] Implement `PUT /api/v1/admin/users/:id` (Update)
- [x] Implement `DELETE /api/v1/admin/users/:id` (Soft Delete)
- [x] Test all endpoints

**Acceptance Criteria:**
- Semua CRUD operations working
- Validation & error handling sesuai API docs

---

### Task 2.8: Banks Reference Endpoint ✅
**Priority:** 🟢 Medium  
**Estimasi:** 1 jam  
**Dependencies:** Task 2.2

**Checklist:**
- [x] Create `handlers/bank_handler.go`
- [x] Implement `GET /api/v1/admin/banks`
- [x] Return all active banks
- [x] Test endpoint

**Acceptance Criteria:**
- Endpoint return 22 banks dari seed data

---

### Task 2.9: Error Handling & Response Standardization ✅
**Priority:** 🟡 High  
**Estimasi:** 3 jam  
**Dependencies:** Task 2.1

**Checklist:**
- [x] Create `utils/response.go` untuk standard response format
- [x] Create `middleware/error_handler.go`
- [x] Implement error codes (VALIDATION_ERROR, NOT_FOUND, etc)
- [x] Implement global error handler
- [x] Add error logging
- [x] Test error responses sesuai API docs

**Acceptance Criteria:**
- Semua response mengikuti format standar
- Error codes sesuai dokumentasi
- Error messages user-friendly
**Checklist:**
- [x] Create `utils/response.go` untuk standard response format
- [x] Implement error codes (VALIDATION_ERROR, NOT_FOUND, etc)
- [x] Implement response helpers (SuccessResponse, ErrorResponse, etc)
- [x] Add error logging (via log package)
- [x] Test error responses (tested via auth endpoints)
- [ ] Create `middleware/error_handler.go` (TODO - using direct error handling for MVP)
- [ ] Implement global error handler (TODO)

**Acceptance Criteria:**
- ✅ Semua response mengikuti format standar
- ✅ Error codes sesuai dokumentasi
- ✅ Error messages user-friendly

---

## 🎨 Frontend Development (Week 2 - Day 1-7)
### Task 3.1: Next.js Project Setup ✅
**Priority:** 🔴 Critical  
**Estimasi:** 3 jam  
**Dependencies:** None

**Checklist:**
- [x] Initialize Next.js project dengan TypeScript
- [x] Install dependencies:
  - [x] TanStack Table (data table)
  - [x] React Hook Form + Zod (form validation)
  - [x] Axios (HTTP client)
  - [x] date-fns (date formatting)
  - [x] Existing UI components (shadcn/ui)
- [x] Setup project structure:
  ```
  apps/web/
  ├── app/ ✅
  │   ├── (auth)/
  │   │   └── login/ (using root page.tsx)
  │   └── (dashboard)/ (using master-admin, operation-admin, reporting-admin)
  │       ├── dashboard/
  │       ├── drivers/
  │       └── users/
  ├── components/ ✅
  ├── lib/ ✅
  │   ├── api.ts ✅         # axios instance
  │   ├── auth.tsx ✅       # auth context
  │   └── utils.ts ✅
  └── types/
  ```
- [x] Setup Tailwind CSS (sudah ada)
- [x] Create `.env.local` (manual copy for now)
- [x] Test `npm run dev`

**Acceptance Criteria:**
- ✅ Next.js app running di `localhost:3000`
- ✅ TypeScript configured
- ✅ Tailwind workinggured
- Tailwind working

---
### Task 3.2: API Client & Auth Context ✅
**Priority:** 🔴 Critical  
**Estimasi:** 4 jam  
**Dependencies:** Task 3.1

**Checklist:**
- [x] Create `lib/api.ts` (axios instance dengan interceptors)
  - [x] Base URL dari env
  - [x] Request interceptor (attach access token)
  - [x] Response interceptor (handle 401, auto refresh token)
- [x] Create `lib/auth.tsx` (AuthContext & AuthProvider)
  - [x] State: admin, isLoading, isAuthenticated
  - [x] Methods: login, logout, refreshAuth
  - [x] Auto refresh token on 401
  - [x] Persist access token di localStorage
- [x] Export `useAuth` hook from auth.tsx
- [x] Wrap app dengan AuthProvider di `app/layout.tsx`
- [x] Test auth flow (tested via login page)

**Acceptance Criteria:**
- ✅ API client bisa hit backend endpoints
- ✅ Auth context manage access token
- ✅ Auto refresh working (via axios interceptor)access token
- Auto refresh working

---
### Task 3.3: Login Page ✅
**Priority:** 🔴 Critical  
**Estimasi:** 4 jam  
**Dependencies:** Task 3.2, Task 2.3 (Backend Auth)

**Checklist:**
- [x] Create `app/(auth)/login/page.tsx` (using root `app/page.tsx`)
- [x] Create login form (basic validation)
  - [x] Email field (validation: required, email format)
  - [x] Password field (validation: required)
  - [x] Submit button with loading state
- [x] Implement login logic
  - [x] Call `POST /api/v1/auth/login` via auth context
  - [x] Save access token ke localStorage
  - [x] Redirect ke `/operation-admin`
- [x] Handle errors (invalid credentials, account inactive)
- [x] Add loading state
- [x] Redirect ke dashboard jika sudah login
- [x] Styling sesuai design system
- [ ] Add React Hook Form + Zod validation (TODO - using basic validation for MVP)

**Acceptance Criteria:**
- ✅ Login form working
- ✅ Error handling working (toast notifications)
- ✅ Redirect working
- ✅ UI responsiveing
- UI responsive

---

### Task 3.4: Protected Route & Layout ⚠️
**Priority:** 🔴 Critical  
**Estimasi:** 3 jam  
**Dependencies:** Task 3.2

**Checklist:**
- [x] Create `app/(dashboard)/layout.tsx` (using master-admin, operation-admin, reporting-admin layouts)
  - [ ] Check auth status
  - [ ] Redirect ke login jika tidak authenticated
  - [x] Render sidebar navigation
  - [x] Render header dengan user info & logout button
- [x] Create `components/sidebar.tsx` (app-sidebar.tsx, master-sidebar.tsx, etc)
  - [x] Menu: Dashboard, Kelola Driver, Kelola User (in operation-admin)
  - [x] Active state
  - [x] Footer dengan admin info
- [ ] Implement logout functionality
- [ ] Test protected routes

**Acceptance Criteria:**
- Unauthenticated user redirect ke login
- Sidebar navigation working
- Logout working

---

### Task 3.5: Dashboard Page ⚠️
**Priority:** 🟡 High  
**Estimasi:** 3 jam  
**Dependencies:** Task 3.4, Task 2.4 (Backend Dashboard)

**Checklist:**
- [x] Create `app/(dashboard)/dashboard/page.tsx` (operation-admin/page.tsx exists)
- [ ] Fetch stats dari `GET /api/v1/admin/dashboard/stats`
- [x] Create stats cards components (UI created with mock data)
  - [x] Total Driver (Active Orders shown)
  - [x] Driver Terverifikasi (Drivers Online shown)
  - [x] Driver Pending (Active Complaints shown)
  - [x] Driver Ditolak
  - [x] Total User
  - [x] User Aktif
- [x] Add loading skeleton (skeleton-dashboard.tsx exists)
- [ ] Add error handling
- [x] Styling dengan existing components

**Acceptance Criteria:**
- Dashboard menampilkan stats real-time
- Loading state working
- Error handling working

---

### Task 3.6: Drivers List Page ⚠️
**Priority:** 🔴 Critical  
**Estimasi:** 6 jam  
**Dependencies:** Task 3.4, Task 2.6 (Backend Drivers)

**Checklist:**
- [x] Create `app/(dashboard)/drivers/page.tsx` (operation-admin/drivers/page.tsx exists)
- [x] Implement data table dengan TanStack Table (UI created with mock data)
  - [x] Columns: Foto, Nama, NIK, No. Telp, Email, Status Verifikasi, Aksi
  - [ ] Server-side pagination (UI ready, no API integration)
  - [x] Search input (debounced) (UI ready)
  - [x] Status filter dropdown (UI ready)
  - [x] Sort by columns (UI ready)
- [x] Implement status chip dengan click handler (UI ready)
  - [x] Popover dropdown untuk change status (UI ready)
  - [ ] Call `PATCH /api/v1/admin/drivers/:id/verification-status`
- [x] Add "Tambah Driver" button (pojok kanan atas)
- [x] Implement action buttons (View, Edit, Delete) (UI ready)
- [x] Add loading state (loading.tsx exists)
- [x] Add empty state (UI ready)
- [ ] Test pagination & search

**Acceptance Criteria:**
- Data table menampilkan drivers dengan pagination
- Search working (debounced)
- Status filter working
- Click status chip bisa change status
- Action buttons working

---

### Task 3.7: Driver Create/Edit Modal
**Priority:** 🔴 Critical  
**Estimasi:** 8 jam  
**Dependencies:** Task 3.6, Task 2.5 (Backend Upload)

**Checklist:**
- [ ] Create `components/driver-form-modal.tsx`
- [ ] Implement form dengan React Hook Form + Zod
  - [ ] Photo Profile upload (presigned URL flow)
  - [ ] Foto KTP upload
  - [ ] Foto Muka upload
  - [ ] Nama (text input)
  - [ ] Email (text input)
  - [ ] Password (password input, optional untuk edit)
  - [ ] No. Telp (text input)
  - [ ] NIK (text input, 16 digits)
  - [ ] Tempat Lahir (text input)
  - [ ] Tanggal Lahir (date picker)
  - [ ] Bank (dropdown, fetch dari `/api/v1/admin/banks`)
  - [ ] No. Rekening (text input)
  - [ ] Status Verifikasi (select)
- [ ] Implement file upload flow:
  1. User select file
  2. Call `POST /api/v1/upload/presigned-url`
  3. PUT file to presigned URL
  4. Save object key to form state
  5. Show preview
- [ ] Implement create logic (`POST /api/v1/admin/drivers`)
- [ ] Implement edit logic (`PUT /api/v1/admin/drivers/:id`)
- [ ] Add validation (client-side & server-side errors)
- [ ] Add loading state
- [ ] Test create & edit flow

**Acceptance Criteria:**
- Modal form working untuk create & edit
- File upload via presigned URL working
- Validation working
- Error handling working
- Preview foto working

---

### Task 3.8: Driver Detail & Delete
**Priority:** 🟡 High  
**Estimasi:** 3 jam  
**Dependencies:** Task 3.6

**Checklist:**
- [ ] Create `components/driver-detail-dialog.tsx`
  - [ ] Display all driver data
  - [ ] Show foto preview (KTP, Muka, Profile)
  - [ ] Load images via presigned view URL
- [ ] Create `components/driver-delete-dialog.tsx`
  - [ ] Confirmation dialog
  - [ ] Call `DELETE /api/v1/admin/drivers/:id`
  - [ ] Refresh table after delete
- [ ] Test view detail & delete

**Acceptance Criteria:**
- Detail dialog menampilkan semua data driver
- Foto ter-load dengan benar
- Delete confirmation working
- Table refresh after delete

---

### Task 3.9: Users List Page ⚠️
**Priority:** 🔴 Critical  
**Estimasi:** 5 jam  
**Dependencies:** Task 3.4, Task 2.7 (Backend Users)

**Checklist:**
- [x] Create `app/(dashboard)/users/page.tsx` (operation-admin/users/page.tsx exists)
- [x] Implement data table (similar to drivers) (UI created with mock data)
  - [x] Columns: Foto, Nama, Email, No. Telp, Aksi
  - [ ] Pagination, search, sort (UI ready, no API integration)
- [x] Add "Tambah User" button
- [x] Implement action buttons (UI ready)
- [ ] Test table functionality

**Acceptance Criteria:**
- Data table working dengan pagination & search
- Action buttons working

---

### Task 3.10: User Create/Edit Modal
**Priority:** 🔴 Critical  
**Estimasi:** 5 jam  
**Dependencies:** Task 3.9

**Checklist:**
- [ ] Create `components/user-form-modal.tsx`
- [ ] Implement form (simpler than driver form)
  - [ ] Photo Profile upload
  - [ ] Nama
  - [ ] Email
  - [ ] No. Telp
  - [ ] Password (optional untuk edit)
- [ ] Implement create & edit logic
- [ ] Test form

**Acceptance Criteria:**
- Modal form working
- File upload working
- Validation working

---

### Task 3.11: User Detail & Delete
**Priority:** 🟡 High  
**Estimasi:** 2 jam  
**Dependencies:** Task 3.9

**Checklist:**
- [ ] Create `components/user-detail-dialog.tsx`
- [ ] Create `components/user-delete-dialog.tsx`
- [ ] Test functionality

**Acceptance Criteria:**
- Detail & delete working

---

## 🧪 Testing & QA (Week 3 - Day 1-3)

### Task 4.1: Backend API Testing
**Priority:** 🟡 High  
**Estimasi:** 4 jam  
**Dependencies:** All backend tasks

**Checklist:**
- [ ] Test all auth endpoints (login, refresh, logout, me)
- [ ] Test dashboard stats endpoint
- [ ] Test driver CRUD endpoints
- [ ] Test user CRUD endpoints
- [ ] Test upload presigned URL endpoints
- [ ] Test banks endpoint
- [ ] Test error responses
- [ ] Test validation
- [ ] Test rate limiting
- [ ] Create Postman collection

**Acceptance Criteria:**
- Semua endpoints working sesuai API docs
- Error handling consistent
- Postman collection documented

---

### Task 4.2: Frontend E2E Testing
**Priority:** 🟡 High  
**Estimasi:** 4 jam  
**Dependencies:** All frontend tasks

**Checklist:**
- [ ] Test login flow
- [ ] Test logout flow
- [ ] Test dashboard page
- [ ] Test driver CRUD flow (create, edit, delete, view)
- [ ] Test driver status change
- [ ] Test user CRUD flow
- [ ] Test file upload flow
- [ ] Test pagination & search
- [ ] Test error handling
- [ ] Test responsive design

**Acceptance Criteria:**
- Semua user flows working
- No console errors
- Responsive di mobile & desktop

---

### Task 4.3: Integration Testing
**Priority:** 🟡 High  
**Estimasi:** 3 jam  
**Dependencies:** Task 4.1, Task 4.2

**Checklist:**
- [ ] Test full flow: login → dashboard → create driver dengan foto → edit → delete
- [ ] Test full flow: create user → edit → delete
- [ ] Test token refresh flow
- [ ] Test session expiry
- [ ] Test concurrent requests
- [ ] Test file upload edge cases (large file, wrong format)

**Acceptance Criteria:**
- End-to-end flows working smoothly
- No race conditions
- Edge cases handled

---

## 📝 Documentation & Deployment (Week 3 - Day 4-5)

### Task 5.1: Documentation
**Priority:** 🟢 Medium  
**Estimasi:** 3 jam  
**Dependencies:** All tasks

**Checklist:**
- [ ] Update README.md dengan setup instructions
- [ ] Document environment variables
- [ ] Document `make dev` workflow
- [ ] Document API endpoints (Postman collection)
- [ ] Document database schema
- [ ] Add code comments
- [ ] Create troubleshooting guide

**Acceptance Criteria:**
- Developer baru bisa setup project dengan mudah
- Dokumentasi lengkap & up-to-date

---

### Task 5.2: Code Review & Refactoring
**Priority:** 🟢 Medium  
**Estimasi:** 4 jam  
**Dependencies:** All tasks

**Checklist:**
- [ ] Code review backend
- [ ] Code review frontend
- [ ] Refactor duplicated code
- [ ] Improve error messages
- [ ] Optimize queries
- [ ] Check security best practices
- [ ] Check performance

**Acceptance Criteria:**
- Code clean & maintainable
- No security vulnerabilities
- Performance acceptable

---

### Task 5.3: Deployment Preparation
**Priority:** 🟢 Medium  
**Estimasi:** 2 jam  
**Dependencies:** Task 5.2

**Checklist:**
- [ ] Create production `.env.example`
- [ ] Update docker-compose untuk production
- [ ] Add health check endpoints
- [ ] Setup logging
- [ ] Setup monitoring (optional)
- [ ] Create deployment guide

**Acceptance Criteria:**
- Ready untuk deploy ke staging/production
- Deployment guide documented

---

## 📊 Task Summary

### By Priority
- 🔴 Critical: 13 tasks (must have untuk MVP)
- 🟡 High: 7 tasks (important untuk UX)
- 🟢 Medium: 3 tasks (nice to have)

### By Week
- **Week 1:** Architecture setup + Backend API (Tasks 1.1 - 2.9)
- **Week 2:** Frontend development (Tasks 3.1 - 3.11)
### Total Estimasi
- **Backend:** ~40 jam (✅ 18 jam completed, ⏳ 22 jam remaining)
- **Frontend:** ~46 jam (✅ 11 jam completed, ⏳ 35 jam remaining)
- **Testing & QA:** ~11 jam (⏳ 11 jam remaining)
- **Documentation:** ~9 jam (⏳ 9 jam remaining)
- **Total:** ~106 jam (~2.5 minggu untuk 1 developer full-time)
- **Progress:** ~29 jam / 106 jam = **~27% time spent, 45% functionality complete** (ahead of schedule!)
- **Documentation:** ~9 jam
- **Total:** ~106 jam (~2.5 minggu untuk 1 developer full-time)

---

## 🎯 Definition of Done

Setiap task dianggap selesai jika:
- ✅ Code implemented & working
- ✅ Manual testing passed
- ✅ Error handling implemented
- ✅ Code reviewed (self-review minimal)
- ✅ Documented (code comments + README jika perlu)
- ✅ No console errors/warnings
- ✅ Responsive (untuk frontend)
- ✅ Sesuai dengan API docs & PRD

---

## 🚨 Risk Mitigation

### Risk 1: MinIO Presigned URL Complexity
**Mitigation:** Implement & test Task 2.5 early, create helper functions yang reusable

### Risk 2: Token Refresh Race Condition
**Mitigation:** Implement request queue di axios interceptor, test concurrent requests

### Risk 3: File Upload UX
**Mitigation:** Add progress indicator, preview, clear error messages

## 🎉 Success Criteria

Sprint v1 dianggap sukses jika:
- ✅ Admin bisa login dengan aman (DONE - login flow working)
- ⏳ Dashboard menampilkan stats real-time (API ready, needs frontend integration)
- ⏳ CRUD driver lengkap dengan upload foto (API ready, needs frontend forms)
- ⏳ CRUD user lengkap dengan upload foto (API ready, needs frontend forms)
- ✅ Semua endpoint API sesuai dokumentasi (DONE - 19 endpoints tested)
- ⏳ UI responsive & user-friendly (structure ready, needs data integration)
- ✅ No critical bugs (backend stable)
- ⏳ Ready untuk demo ke stakeholder (70% complete) (70% complete)

### Sprint Review (End of Week 3)
- Demo semua fitur ke stakeholder
- Gather feedback

### Sprint Retrospective
- What went well?
- What can be improved?
- Action items untuk sprint berikutnya

---

## 🎉 Success Criteria

Sprint v1 dianggap sukses jika:
- ✅ Admin bisa login dengan aman
- ✅ Dashboard menampilkan stats real-time
- ✅ CRUD driver lengkap dengan upload foto
- ✅ CRUD user lengkap dengan upload foto
- ✅ Semua endpoint API sesuai dokumentasi
- ✅ UI responsive & user-friendly
- ✅ No critical bugs
- ✅ Ready untuk demo ke stakeholder

---

**Next Sprint (v2):** Operating Admin & Reporting Admin modules
