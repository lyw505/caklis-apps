# CAKLI Admin Dashboard v1 🚀

Platform admin dashboard untuk mengelola sistem transportasi becak listrik CAKLI.

**Status:** 🟢 MVP Complete (100%)

**📖 New to this project? Start here:** [START-HERE.md](START-HERE.md)

---

## 📋 Table of Contents

- [Features](#-features)
- [Tech Stack](#-tech-stack)
- [Quick Start](#-quick-start)
- [Multi-Role Admin](#-multi-role-admin)
- [Project Structure](#-project-structure)
- [API Endpoints](#-api-endpoints)
- [Development](#-development)
- [Testing](#-testing)
- [Troubleshooting](#-troubleshooting)
- [Documentation](#-documentation)

---

## ✨ Features

### Completed ✅

- **Multi-Role Authentication System**
  - Login with email/password
  - JWT dual token (access + refresh)
  - Automatic token refresh
  - Secure logout with token revocation
  - Protected routes
  - 3 admin roles: Master, Operation, Reporting
  - Role-based access control
  - Auto-redirect based on role

- **Dashboard**
  - Real-time statistics (drivers, users)
  - Responsive design
  - Loading states

- **Driver Management**
  - List drivers with pagination
  - Search by name, email, phone, NIK
  - Filter by verification status
  - Create/Edit drivers with form modal
  - Upload photos (profile, KTP, face) via MinIO
  - Soft delete

- **User Management**
  - List users with pagination
  - Search by name, email, phone
  - Create/Edit users with form modal
  - Upload profile photo
  - Soft delete

- **File Upload**
  - MinIO presigned URL flow
  - Support JPG/PNG images
  - Preview before upload
  - Secure upload with expiring URLs

- **Infrastructure**
  - Docker setup (PostgreSQL, MinIO)
  - Database seeded with 22 banks + default admin
  - Hot reload for development
  - CORS configured

---

## 🛠️ Tech Stack

### Backend
- **Go 1.25+** - Programming language
- **Fiber v3** - Web framework
- **GORM** - ORM for PostgreSQL
- **JWT** - Authentication
- **bcrypt** - Password hashing
- **MinIO SDK** - Object storage

### Frontend
- **Next.js 16** - React framework
- **React 19** - UI library
- **TypeScript** - Type safety
- **Tailwind CSS** - Styling
- **shadcn/ui** - UI components
- **TanStack Table** - Data tables
- **Axios** - HTTP client
- **Sonner** - Toast notifications

### Infrastructure
- **PostgreSQL 15** - Database
- **MinIO** - Object storage
- **Docker** - Containerization

---

## 🚀 Quick Start

### Prerequisites

- Go 1.25+
- Node.js 18+
- Docker & Docker Compose
- Make (optional, for convenience)

### 1. Clone Repository

```bash
git clone <repo-url>
cd caklis-apps
```

### 2. Start All Services (Recommended)

```bash
# Start everything with one command
make dev
```

This will:
1. Start PostgreSQL and MinIO containers
2. Start backend API with hot reload
3. Start frontend development server

### 3. Configure MinIO (First Time Only)

After starting services for the first time, configure MinIO bucket:

```bash
make setup-minio
```

This will:
- Create `cakli` bucket if not exists
- Set bucket policy to allow uploads
- Enable presigned URL functionality

**Note:** This step is required for file upload features (driver photos, user profile pictures) to work properly.

### 3. Manual Setup (Alternative)

If you prefer manual setup:

```bash
# Start Docker services
docker-compose up -d

# Wait for services to be ready
sleep 10

# Configure MinIO bucket (first time only)
make setup-minio

# Start backend
cd apps/api
cp ../../.env .env
go run main.go

# In another terminal, start frontend
cd apps/web
npm install
cp ../../.env .env.local
npm run dev
```

### 4. Access the Application

- **Frontend:** http://localhost:3000
- **Backend API:** http://localhost:8080
- **MinIO Console:** http://localhost:9001 (minioadmin / minioadmin)
- **PostgreSQL:** localhost:5432 (cakli / cakli_secret)

### 5. Login

**Multi-Role Admin Credentials:**

| Role | Email | Password | Dashboard |
|------|-------|----------|-----------|
| Master Admin | `master@cakli.id` | `admin123` | /master-admin |
| Operation Admin | `operation@cakli.id` | `admin123` | /operation-admin |
| Reporting Admin | `reporting@cakli.id` | `admin123` | /reporting-admin |

**Note:** Sistem akan otomatis redirect ke dashboard sesuai role admin setelah login.

---

## 👥 Multi-Role Admin

CAKLI menggunakan sistem multi-role admin dengan 3 tingkat akses:

### 1. Master Admin
**Dashboard:** `/master-admin`
**Access:** Full system access
- ✅ Kelola Driver & User
- ✅ Tariff Management
- ✅ Area & Zone Management
- ✅ Admin Roles Management
- ✅ Audit Logs
- ✅ Partner Policies

### 2. Operation Admin
**Dashboard:** `/operation-admin`
**Access:** Operational management
- ✅ Kelola Driver & User
- ✅ Live Map Monitoring
- ✅ Order Management
- ✅ Activity Monitoring
- ✅ Complaints & Disputes

### 3. Reporting Admin
**Dashboard:** `/reporting-admin`
**Access:** Read-only analytics
- ✅ Analytics Dashboard
- ✅ Revenue Reports
- ✅ Driver Performance Reports
- ✅ Order History
- ✅ Cancellation Reports
- ❌ Cannot manage drivers/users

### Role-Based Access Control

Backend menggunakan middleware untuk mengontrol akses:
- `RequireMasterOrOperation()` - Driver & User management
- `RequireMasterOnly()` - Master-only features
- `RequireOperationOnly()` - Operation-only features
- `RequireReportingOnly()` - Reporting-only features

**See:** [Multi-Role Admin Guide](docs/MULTI-ROLE-ADMIN-GUIDE.md) for detailed information.

---

## 📁 Project Structure

```
caklis-apps/
├── apps/
│   ├── api/                    # Go Fiber Backend
│   │   ├── config/            # Database & MinIO config
│   │   │   ├── database.go    # PostgreSQL connection
│   │   │   └── minio.go       # MinIO client
│   │   ├── models/            # GORM models
│   │   │   ├── admin.go       # Admin model
│   │   │   ├── driver.go      # Driver model
│   │   │   ├── user.go        # User model
│   │   │   ├── bank.go        # Bank model
│   │   │   └── refresh_token.go
│   │   ├── handlers/          # HTTP handlers
│   │   │   ├── auth_handler.go      # Authentication
│   │   │   ├── driver_handler.go    # Driver CRUD
│   │   │   ├── user_handler.go      # User CRUD
│   │   │   ├── upload_handler.go    # File upload
│   │   │   ├── bank_handler.go      # Banks reference
│   │   │   └── dashboard_handler.go # Dashboard stats
│   │   ├── middleware/        # Middleware
│   │   │   ├── auth.go        # JWT authentication
│   │   │   └── cors.go        # CORS configuration
│   │   ├── services/          # Business logic
│   │   │   └── upload_service.go    # Presigned URL service
│   │   ├── utils/             # Utilities
│   │   │   ├── jwt.go         # JWT generation/validation
│   │   │   ├── hash.go        # Password hashing
│   │   │   └── response.go    # Standard responses
│   │   └── main.go            # Entry point
│   └── web/                   # Next.js Frontend
│       ├── app/               # App router pages
│       │   ├── page.tsx       # Login page
│       │   └── operation-admin/
│       │       ├── layout.tsx # Protected layout
│       │       ├── page.tsx   # Dashboard
│       │       ├── drivers/   # Driver management
│       │       └── users/     # User management
│       ├── components/        # UI components
│       │   ├── ui/            # shadcn/ui components
│       │   ├── driver-form-modal.tsx
│       │   ├── user-form-modal.tsx
│       │   └── app-sidebar.tsx
│       └── lib/               # Libraries
│           ├── api.ts         # Axios client
│           └── auth.tsx       # Auth context
├── docs/
│   ├── agile-development/v1/
│   │   ├── api.md             # API documentation
│   │   ├── prd.md             # Product requirements
│   │   ├── db.sql             # Database schema
│   │   └── implementation-plan.md
│   ├── TROUBLESHOOTING.md     # Common issues & solutions
│   └── TESTING-GUIDE.md       # Testing instructions
├── .env                       # Environment variables
├── docker-compose.yml         # Docker services
├── Makefile                   # Development commands
└── README.md                  # This file
```

---

## 📡 API Endpoints

### Authentication (Public)

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/auth/login` | Login with email/password |
| POST | `/api/v1/auth/refresh` | Refresh access token |
| POST | `/api/v1/auth/logout` | Logout and revoke token |

### Authentication (Protected)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/auth/me` | Get current admin profile |

### Dashboard (Protected)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/admin/dashboard/stats` | Get dashboard statistics |

### Drivers (Protected)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/admin/drivers` | List drivers (pagination, search, filter) |
| GET | `/api/v1/admin/drivers/:id` | Get driver detail |
| POST | `/api/v1/admin/drivers` | Create new driver |
| PUT | `/api/v1/admin/drivers/:id` | Update driver |
| DELETE | `/api/v1/admin/drivers/:id` | Soft delete driver |
| PATCH | `/api/v1/admin/drivers/:id/verification-status` | Update verification status |

### Users (Protected)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/admin/users` | List users (pagination, search) |
| GET | `/api/v1/admin/users/:id` | Get user detail |
| POST | `/api/v1/admin/users` | Create new user |
| PUT | `/api/v1/admin/users/:id` | Update user |
| DELETE | `/api/v1/admin/users/:id` | Soft delete user |

### Upload (Protected)

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/upload/presigned-url` | Generate presigned upload URL |
| GET | `/api/v1/upload/presigned-view-url` | Generate presigned view URL |

### Banks (Protected)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/admin/banks` | List all active banks |

### Health Check (Public)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/health` | API health check |

**Total:** 19 endpoints

---

## 🔧 Development

### Backend Commands

```bash
cd apps/api

# Install dependencies
go mod download

# Run with hot reload (using Air)
air

# Run without hot reload
go run main.go

# Build binary
go build -o bin/server main.go

# Run tests (when available)
go test ./... -v
```

### Frontend Commands

```bash
cd apps/web

# Install dependencies
npm install

# Development mode
npm run dev

# Build for production
npm run build

# Start production server
npm start

# Lint code
npm run lint

# Run tests (when available)
npm test
```

### Docker Commands

```bash
# Start all services
docker-compose up -d

# Stop all services
docker-compose down

# View logs
docker-compose logs -f

# View specific service logs
docker-compose logs -f backend-api
docker-compose logs -f postgres
docker-compose logs -f minio

# Restart specific service
docker-compose restart backend-api
docker-compose restart postgres
docker-compose restart minio

# Remove all containers and volumes (⚠️ deletes all data)
docker-compose down -v
```

### MinIO Commands

```bash
# Configure MinIO bucket (first time setup)
make setup-minio

# Check bucket policy
docker exec cakli-storage mc anonymous get local/cakli

# Access MinIO Console
# Open http://localhost:9001
# Login: minioadmin / minioadmin

# List buckets
docker exec cakli-storage mc ls local/

# List files in bucket
docker exec cakli-storage mc ls local/cakli/drivers/
```

### Database Commands

```bash
# Connect to PostgreSQL
docker exec -it cakli-db psql -U cakli -d cakli_db

# Execute SQL file
cat docs/agile-development/v1/db.sql | docker exec -i cakli-db psql -U cakli -d cakli_db

# Backup database
docker exec cakli-db pg_dump -U cakli cakli_db > backup.sql

# Restore database
cat backup.sql | docker exec -i cakli-db psql -U cakli -d cakli_db
```

---

## 🧪 Testing

See [TESTING-GUIDE.md](docs/TESTING-GUIDE.md) for comprehensive testing instructions.

### Quick Test

1. **Login Test**
   ```bash
   curl -X POST http://localhost:8080/api/v1/auth/login \
     -H "Content-Type: application/json" \
     -d '{"email":"master@cakli.id","password":"admin123"}'
   ```

2. **Health Check**
   ```bash
   curl http://localhost:8080/api/v1/health
   ```

3. **Frontend Test**
   - Open http://localhost:3000
   - Login with `master@cakli.id` / `admin123`
   - Navigate to Drivers page
   - Try creating a new driver

---

## 🐛 Troubleshooting

See [TROUBLESHOOTING.md](docs/TROUBLESHOOTING.md) for detailed solutions.

### Common Issues

**1. Port already in use**
```bash
# Check what's using the port
netstat -ano | findstr :8080
netstat -ano | findstr :3000
netstat -ano | findstr :9000

# Kill the process
taskkill /PID <PID> /F
```

**2. Database connection failed**
```bash
# Restart PostgreSQL
docker-compose restart postgres

# Check if running
docker ps | grep cakli-db

# View logs
docker logs cakli-db
```

**3. File upload fails with 403 Forbidden**
```bash
# Configure MinIO bucket policy
make setup-minio

# Verify bucket policy is set to public
docker exec cakli-storage mc anonymous get local/cakli
# Should output: Access permission for `local/cakli` is `public`
```

**4. MinIO bucket not found**
```bash
# Restart API (it will auto-create bucket)
docker-compose restart backend-api

# Or manually create bucket
docker exec cakli-storage mc mb local/cakli
```

**5. Frontend can't connect to backend**
- Check `.env.local` has correct `NEXT_PUBLIC_API_URL`
- Verify backend is running on port 8080
- Check CORS configuration in `apps/api/middleware/cors.go`

**6. Docker containers won't start**
```bash
# Check Docker is running
docker ps

# Remove old containers and start fresh
docker-compose down -v
docker-compose up -d
```

### Additional Resources

- [MinIO Upload Fix Guide](docs/MINIO-403-ERROR-DETAILED.md) - Detailed explanation of upload issues
- [Quick Upload Fix](docs/QUICK-FIX-UPLOAD.md) - Quick reference for upload problems
- [Upload Testing Guide](docs/TEST-UPLOAD.md) - How to test upload functionality

---

## 🗄️ Database

### Tables

1. **banks** - Reference data bank Indonesia (22 banks)
2. **admins** - Admin accounts with soft delete
3. **refresh_tokens** - JWT refresh tokens with rotation
4. **drivers** - Driver data with verification status
5. **users** - User/passenger data

### Seed Data

- **22 banks:** BRI, Mandiri, BNI, BCA, BSI, CIMB Niaga, Danamon, Permata, Maybank, Panin, OCBC NISP, BTN, BII, Mega, Sinarmas, Commonwealth, BTPN, Bukopin, Muamalat, BJB, BPD Jatim, BPD Jabar
- **1 default admin:** `master@cakli.id` / `admin123`

### Schema Highlights

- All main tables have soft delete (`deleted_at`)
- Unique constraints with soft delete support
- Automatic `updated_at` trigger
- Foreign key relationships
- Indexes for performance

---

## 🔐 Security Features

- **JWT Dual Token:** Access token (15 min) + Refresh token (7 days)
- **Token Rotation:** Refresh token rotates on each refresh
- **httpOnly Cookies:** Refresh token stored securely
- **bcrypt Hashing:** Password hashing with cost 12
- **CORS Protection:** Configured allowed origins
- **Soft Delete:** Sensitive data not permanently deleted
- **SQL Injection Protection:** GORM parameterized queries
- **XSS Protection:** React auto-escaping

---

## 📚 Documentation

### Quick References
- [🚀 Quick Start Guide](QUICK-START.md) - Get started in 3 commands
- [📊 Current Status](docs/CURRENT-STATUS.md) - Complete feature list & status
- [📝 Makefile Guide](MAKEFILE-GUIDE.md) - All available make commands

### Detailed Guides
- [🔐 Multi-Role Admin Guide](docs/MULTI-ROLE-ADMIN-GUIDE.md) - Role system & access control
- [👑 Master Admin Access Guide](docs/MASTER-ADMIN-ACCESS-GUIDE.md) - Master admin features
- [🧪 Testing Guide](docs/TESTING-GUIDE.md) - Testing instructions
- [🐛 Troubleshooting](docs/TROUBLESHOOTING.md) - Common issues & solutions
- [🔄 Restart Backend Guide](docs/RESTART-BACKEND-GUIDE.md) - Backend restart instructions

### File Upload & MinIO
- [📤 MinIO Upload Fix](docs/MINIO-403-ERROR-DETAILED.md) - Detailed upload error analysis
- [⚡ Quick Upload Fix](docs/QUICK-FIX-UPLOAD.md) - Quick reference for upload issues
- [🧪 Upload Testing](docs/TEST-UPLOAD.md) - How to test upload functionality

### Technical Documentation
- [API Documentation](docs/agile-development/v1/api.md) - Complete API reference
- [Product Requirements](docs/agile-development/v1/prd.md) - Product specifications
- [Database Schema](docs/agile-development/v1/db.sql) - SQL schema with seed data
- [Implementation Plan](docs/agile-development/v1/implementation-plan.md) - Development roadmap
- [Progress Report](PROGRESS.md) - Current implementation status

---

## 🎯 Implementation Status

### Backend API (100% ✅)

- ✅ Project structure
- ✅ Database connection (GORM)
- ✅ MinIO connection
- ✅ Models (Admin, User, Driver, Bank, RefreshToken)
- ✅ Authentication handlers (Login, Refresh, Logout, Me)
- ✅ Driver CRUD endpoints
- ✅ User CRUD endpoints
- ✅ Upload presigned URL service
- ✅ Banks endpoint
- ✅ Dashboard stats endpoint
- ✅ JWT utilities
- ✅ Password hashing
- ✅ Response utilities
- ✅ Auth middleware
- ✅ CORS middleware
- ✅ Error handling

### Frontend (100% ✅)

- ✅ Next.js 16 setup
- ✅ UI components (shadcn/ui)
- ✅ Login page with auth integration
- ✅ Dashboard with real-time stats
- ✅ Protected routes
- ✅ Sidebar navigation with logout
- ✅ API client (axios with interceptors)
- ✅ Auth context (login, logout, refresh)
- ✅ Driver list page with API integration
- ✅ Driver create/edit modal with file upload
- ✅ User list page with API integration
- ✅ User create/edit modal with file upload
- ✅ Form validation
- ✅ Loading states
- ✅ Error handling
- ✅ Toast notifications

### Infrastructure (100% ✅)

- ✅ Docker setup (PostgreSQL, MinIO)
- ✅ Database schema executed
- ✅ Seed data inserted
- ✅ Environment configuration
- ✅ Hot reload for development
- ✅ Makefile for convenience

### Optional Enhancements (0% - Not Critical)

- ⏳ Delete confirmation dialogs
- ⏳ Detail view dialogs
- ⏳ Automated tests
- ⏳ Rate limiting
- ⏳ Advanced logging
- ⏳ Performance optimization

---

## 🚧 Next Steps (Future Sprints)

1. **Master Admin Module**
   - Area management
   - Partner management
   - Role management
   - Tariff management
   - Audit logs

2. **Reporting Admin Module**
   - Revenue reports
   - Driver performance reports
   - Cancellation reports
   - Analytics dashboard
   - Export to Excel/PDF

3. **Enhancements**
   - Real-time notifications
   - Advanced search & filters
   - Bulk operations
   - Activity logs
   - Email notifications

---

## 📄 License

Private - CAKLI Project

---

## 👥 Team

- **Developer:** [Your Name]
- **Sprint:** v1 (MVP)
- **Status:** Near Complete (~95%)
- **Last Updated:** April 9, 2026

---

## 🎉 Success Metrics

- ✅ Admin can login securely
- ✅ Dashboard shows real-time statistics
- ✅ Driver CRUD complete with file upload
- ✅ User CRUD complete with file upload
- ✅ All 19 API endpoints functional
- ✅ UI responsive & user-friendly
- ✅ No critical bugs
- ✅ Ready for demo

---

**Made with ❤️ for CAKLI**
