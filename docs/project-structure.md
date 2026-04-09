# CAKLI — Monorepo Project Structure

> **Platform:** Cakli 1.0 — Transportasi Becak Listrik  
> **Architecture:** Monorepo  
> **Last Updated:** 31 Maret 2026

---

## Tech Stack

| Layer | Technology |
|-------|-----------|
| **API** | Go, GoFiber v3, GORM, JWT |
| **Database** | PostgreSQL 15+ |
| **Object Storage** | MinIO |
| **Web Dashboard** | Next.js, TypeScript, ShadcnUI, TailwindCSS |
| **Mobile (User & Driver)** | Flutter (Dart) |
| **Container** | Docker, Docker Compose |
| **Web Server / Reverse Proxy** | Nginx |

---

## Root Structure

```
cakli/
├── .github/                          # GitHub CI/CD & templates
│   ├── workflows/
│   │   ├── api-ci.yml                # Go lint, test, build
│   │   ├── web-ci.yml                # Next.js lint, test, build
│   │   ├── mobile-ci.yml             # Flutter analyze, test, build
│   │   └── deploy.yml                # Deploy pipeline (staging/prod)
│   ├── PULL_REQUEST_TEMPLATE.md
│   └── ISSUE_TEMPLATE/
│       ├── bug_report.md
│       └── feature_request.md
│
├── apps/                             # ── Semua aplikasi ──
│   ├── api/                          # Backend API (Go)
│   ├── web/                          # Admin Dashboard (Next.js)
│   ├── mobile/                       # Flutter monorepo workspace
│   │   ├── cakli/                    # User App
│   │   └── cakli_driver/             # Driver App
│   └── nginx/                        # Reverse proxy config
│
├── packages/                         # ── Shared packages ──
│   └── proto/                        # Shared API contracts / types
│
├── infra/                            # ── Infrastructure ──
│   ├── docker/                       # Dockerfiles per service
│   ├── scripts/                      # Utility scripts
│   └── k8s/                          # Kubernetes manifests (future)
│
├── docs/                             # ── Dokumentasi ──
│   ├── prd.md
│   ├── feature.md
│   ├── feature-mobile.md
│   ├── api.md
│   ├── api-standarts.md
│   ├── cakli_database.sql
│   ├── cakli_erd.md
│   └── project-structure.md
│
├── docker-compose.yml                # Orchestrator lokal
├── docker-compose.prod.yml           # Override production
├── Makefile                          # Top-level commands
├── .env.example                      # Template environment
├── .gitignore
├── LICENSE
└── README.md
```

---

## 1. Backend API — `apps/api/`

```
apps/api/
├── cmd/
│   └── server/
│       └── main.go                   # Entry point, bootstrap Fiber app
│
├── internal/                         # Private application code
│   ├── config/
│   │   ├── config.go                 # Viper config loader
│   │   ├── database.go               # GORM PostgreSQL connection
│   │   └── minio.go                  # MinIO client init
│   │
│   ├── middleware/
│   │   ├── auth.go                   # JWT Bearer validation
│   │   ├── rbac.go                   # Role-based access (master/operating/reporting)
│   │   ├── cors.go                   # CORS configuration
│   │   ├── logger.go                 # Request logging (zerolog)
│   │   ├── rate_limiter.go           # Rate limiting
│   │   ├── request_id.go            # X-Request-ID injection
│   │   └── error_handler.go          # Global error handler
│   │
│   ├── model/                        # GORM models (1:1 dengan DB tables)
│   │   ├── user.go                   # users, saved_addresses, user_status_history, user_reports, user_audit_logs
│   │   ├── driver.go                 # drivers, driver_documents, driver_status_history, driver_violations, driver_audit_logs
│   │   ├── vehicle.go                # vehicles
│   │   ├── driver_wallet.go          # driver_wallets, driver_wallet_transactions, driver_bank_accounts
│   │   ├── driver_feature.go         # driver_schedules, driver_sanctions, driver_referrals, driver_qr_codes, driver_appointments, partnership_agreements
│   │   ├── driver_earning.go         # driver_earnings, driver_incentives
│   │   ├── zone.go                   # zones, zone_requests
│   │   ├── tariff.go                 # tariff_configs, tariff_version_history, zone_tariff_multipliers
│   │   ├── order.go                  # orders, order_timeline, order_audit_logs
│   │   ├── payment.go                # payments
│   │   ├── rating.go                 # ratings
│   │   ├── chat.go                   # chat_messages
│   │   ├── notification.go           # notifications
│   │   ├── voucher.go                # vouchers, voucher_claims
│   │   ├── capay.go                  # capay_wallets, capay_transactions
│   │   ├── complaint.go              # complaints, complaint_timeline, complaint_escalations, complaint_audit_logs
│   │   ├── driver_activity.go        # driver_activity_alerts, driver_activity_logs
│   │   ├── admin.go                  # admins, admin_activity_logs
│   │   ├── audit_log.go              # master_audit_logs
│   │   ├── revenue.go                # revenue_transactions, settlement_batches
│   │   └── partner_policy.go         # partner_policies, partner_policy_documents
│   │
│   ├── dto/                          # Data Transfer Objects (request/response)
│   │   ├── auth_dto.go
│   │   ├── user_dto.go
│   │   ├── driver_dto.go
│   │   ├── order_dto.go
│   │   ├── zone_dto.go
│   │   ├── tariff_dto.go
│   │   ├── complaint_dto.go
│   │   ├── admin_dto.go
│   │   ├── payment_dto.go
│   │   ├── report_dto.go
│   │   └── common_dto.go            # Pagination, filters, meta
│   │
│   ├── repository/                   # Data access layer (GORM queries)
│   │   ├── user_repository.go
│   │   ├── driver_repository.go
│   │   ├── vehicle_repository.go
│   │   ├── order_repository.go
│   │   ├── zone_repository.go
│   │   ├── tariff_repository.go
│   │   ├── payment_repository.go
│   │   ├── rating_repository.go
│   │   ├── chat_repository.go
│   │   ├── notification_repository.go
│   │   ├── voucher_repository.go
│   │   ├── capay_repository.go
│   │   ├── complaint_repository.go
│   │   ├── driver_wallet_repository.go
│   │   ├── driver_earning_repository.go
│   │   ├── driver_feature_repository.go
│   │   ├── driver_activity_repository.go
│   │   ├── admin_repository.go
│   │   ├── audit_log_repository.go
│   │   ├── revenue_repository.go
│   │   └── partner_policy_repository.go
│   │
│   ├── service/                      # Business logic layer
│   │   ├── auth_service.go
│   │   ├── user_service.go
│   │   ├── driver_service.go
│   │   ├── order_service.go
│   │   ├── zone_service.go
│   │   ├── tariff_service.go
│   │   ├── payment_service.go
│   │   ├── rating_service.go
│   │   ├── chat_service.go
│   │   ├── notification_service.go
│   │   ├── voucher_service.go
│   │   ├── capay_service.go
│   │   ├── complaint_service.go
│   │   ├── driver_wallet_service.go
│   │   ├── driver_earning_service.go
│   │   ├── driver_feature_service.go
│   │   ├── driver_activity_service.go
│   │   ├── admin_service.go
│   │   ├── audit_log_service.go
│   │   ├── revenue_service.go
│   │   ├── partner_policy_service.go
│   │   └── upload_service.go         # MinIO file upload
│   │
│   ├── handler/                      # HTTP handlers (Fiber)
│   │   ├── auth_handler.go
│   │   ├── user_handler.go
│   │   ├── driver_handler.go
│   │   ├── order_handler.go
│   │   ├── zone_handler.go
│   │   ├── tariff_handler.go
│   │   ├── payment_handler.go
│   │   ├── rating_handler.go
│   │   ├── chat_handler.go
│   │   ├── notification_handler.go
│   │   ├── voucher_handler.go
│   │   ├── capay_handler.go
│   │   ├── complaint_handler.go
│   │   ├── driver_wallet_handler.go
│   │   ├── driver_earning_handler.go
│   │   ├── driver_feature_handler.go
│   │   ├── driver_activity_handler.go
│   │   ├── admin_handler.go
│   │   ├── audit_log_handler.go
│   │   ├── revenue_handler.go
│   │   ├── partner_policy_handler.go
│   │   ├── upload_handler.go
│   │   └── realtime_map_handler.go   # WebSocket / SSE map data
│   │
│   └── router/
│       ├── router.go                 # Main router registration
│       ├── auth_routes.go            # /api/v1/auth/*
│       ├── user_routes.go            # /api/v1/users/*
│       ├── driver_routes.go          # /api/v1/drivers/*
│       ├── order_routes.go           # /api/v1/orders/*
│       ├── zone_routes.go            # /api/v1/zones/*
│       ├── tariff_routes.go          # /api/v1/tariffs/*
│       ├── admin_routes.go           # /api/v1/admins/*
│       └── public_routes.go          # Health check, swagger
│
├── pkg/                              # Public shared utilities
│   ├── database/
│   │   └── postgres.go               # GORM connection helper
│   ├── jwt/
│   │   └── jwt.go                    # JWT generate, validate, claims
│   ├── query/
│   │   └── pagination.go             # Paginate, sort, filter helper
│   ├── response/
│   │   └── response.go               # Standard JSON response builder
│   ├── storage/
│   │   └── minio.go                  # MinIO upload/download helper
│   ├── validator/
│   │   └── validator.go              # Request validation helper
│   └── hash/
│       └── hash.go                   # Bcrypt password + PIN hashing
│
├── migrations/
│   ├── 000001_init_schema.up.sql     # Full schema (dari cakli_database.sql)
│   ├── 000001_init_schema.down.sql
│   └── ...                           # Incremental migrations
│
├── seeds/
│   ├── admin_seeder.go               # Default admin accounts
│   ├── zone_seeder.go                # Default zones
│   └── tariff_seeder.go              # Default tariff config
│
├── go.mod
├── go.sum
├── .env.example
├── Makefile                          # build, run, migrate, seed, test
└── Dockerfile
```

### Arsitektur Layer API

```
Request → Router → Middleware → Handler → Service → Repository → Database
                                   ↓
                              DTO (validate)
```

- **Router**: Mendefinisikan routes & middleware per group
- **Middleware**: Auth JWT, RBAC, rate limit, logging, CORS
- **Handler**: Parse request, validasi DTO, panggil service, return response
- **Service**: Business logic, orchestration antar repository
- **Repository**: GORM queries, data access
- **Model**: Struct representasi tabel database
- **DTO**: Request/Response objects, validasi input
- **Pkg**: Reusable utilities (JWT, pagination, response format, MinIO)

---

## 2. Web Admin Dashboard — `apps/web/`

```
apps/web/
├── app/                              # Next.js App Router
│   ├── layout.tsx                    # Root layout (providers, fonts)
│   ├── page.tsx                      # Login page
│   │
│   ├── (auth)/                       # Auth group layout
│   │   ├── login/page.tsx
│   │   └── layout.tsx
│   │
│   ├── (dashboard)/                  # Authenticated dashboard layout
│   │   ├── layout.tsx                # Sidebar + header + auth guard
│   │   │
│   │   ├── master-admin/             # ── Master Admin Pages ──
│   │   │   ├── page.tsx              # Kontrol Sistem Global / Dashboard
│   │   │   ├── areas/page.tsx        # Manajemen Area & Zona
│   │   │   ├── tariffs/page.tsx      # Manajemen Tarif
│   │   │   ├── roles/page.tsx        # Kontrol Akses Admin
│   │   │   ├── audit/page.tsx        # Log Audit
│   │   │   └── partners/page.tsx     # Kebijakan Mitra
│   │   │
│   │   ├── operation-admin/          # ── Operating Admin Pages ──
│   │   │   ├── page.tsx              # Dashboard Operasional
│   │   │   ├── map/page.tsx          # Peta Operasional Real-time
│   │   │   ├── drivers/page.tsx      # Manajemen Driver
│   │   │   ├── orders/page.tsx       # Manajemen Pesanan
│   │   │   ├── complaints/page.tsx   # Keluhan & Sengketa
│   │   │   ├── activity/page.tsx     # Pemantauan Aktivitas Driver
│   │   │   └── users/page.tsx        # Moderasi Pengguna
│   │   │
│   │   └── reporting-admin/          # ── Reporting Admin Pages ──
│   │       ├── page.tsx              # Ikhtisar Bisnis
│   │       ├── drivers/page.tsx      # Wawasan Performa Driver
│   │       ├── reports/
│   │       │   ├── page.tsx          # Pusat Laporan
│   │       │   ├── revenue/page.tsx  # Kokpit Keuangan
│   │       │   ├── driver-performance/page.tsx
│   │       │   └── cancellation/page.tsx
│   │       ├── history/page.tsx      # Riwayat Transaksi
│   │       └── analytics/page.tsx    # Analitik Lintas Area
│   │
│   └── api/                          # Next.js API routes (BFF/proxy, optional)
│       └── ...
│
├── components/                       # Reusable UI components
│   ├── ui/                           # ShadcnUI primitives
│   │   ├── button.tsx
│   │   ├── card.tsx
│   │   ├── dialog.tsx
│   │   ├── table.tsx
│   │   ├── badge.tsx
│   │   ├── input.tsx
│   │   ├── select.tsx
│   │   ├── tabs.tsx
│   │   ├── sheet.tsx                 # Side drawer
│   │   ├── chart.tsx                 # Recharts wrapper
│   │   └── ...
│   │
│   ├── layout/                       # Layout components
│   │   ├── sidebar.tsx               # Navigation sidebar (role-based)
│   │   ├── header.tsx                # Top bar (user info, notifications)
│   │   ├── breadcrumb.tsx
│   │   └── page-header.tsx           # Title + description + actions
│   │
│   ├── shared/                       # Shared business components
│   │   ├── kpi-card.tsx              # KPI stat card
│   │   ├── data-table.tsx            # Generic data table (sorting, filter, pagination)
│   │   ├── status-badge.tsx          # Universal status badge
│   │   ├── confirm-dialog.tsx        # Confirmation dialog
│   │   ├── audit-log-dialog.tsx      # Audit log viewer
│   │   ├── date-range-picker.tsx
│   │   ├── search-filter-bar.tsx
│   │   └── export-button.tsx
│   │
│   └── features/                     # Feature-specific components
│       ├── drivers/
│       │   ├── driver-table.tsx
│       │   ├── add-driver-dialog.tsx  # Multi-step form
│       │   └── driver-action-dropdown.tsx
│       ├── orders/
│       │   ├── order-tabs.tsx
│       │   └── order-intervention-dialog.tsx
│       ├── complaints/
│       │   ├── complaint-table.tsx
│       │   ├── resolution-dialog.tsx
│       │   └── escalation-dialog.tsx
│       ├── zones/
│       │   ├── zone-table.tsx
│       │   ├── zone-map.tsx
│       │   └── add-zone-dialog.tsx
│       ├── tariffs/
│       │   ├── tariff-form.tsx
│       │   ├── impact-simulation.tsx
│       │   └── version-history.tsx
│       └── map/
│           ├── realtime-map.tsx
│           └── map-controls.tsx
│
├── lib/                              # Utilities & configurations
│   ├── api/                          # API client
│   │   ├── client.ts                 # Axios/fetch instance (base URL, interceptors)
│   │   ├── auth.ts                   # Auth endpoints
│   │   ├── users.ts                  # User endpoints
│   │   ├── drivers.ts                # Driver endpoints
│   │   ├── orders.ts                 # Order endpoints
│   │   ├── zones.ts                  # Zone endpoints
│   │   ├── tariffs.ts
│   │   ├── complaints.ts
│   │   ├── reports.ts
│   │   └── admins.ts
│   ├── utils.ts                      # clsx, formatters, helpers
│   ├── constants.ts                  # Role enums, status maps
│   └── validations.ts                # Zod schemas
│
├── hooks/                            # Custom React hooks
│   ├── use-auth.ts                   # Auth state & guard
│   ├── use-pagination.ts
│   ├── use-debounce.ts
│   └── use-realtime.ts               # WebSocket/SSE hook
│
├── stores/                           # Client state (Zustand)
│   ├── auth-store.ts
│   └── ui-store.ts                   # Sidebar, theme, etc.
│
├── types/                            # TypeScript type definitions
│   ├── user.ts
│   ├── driver.ts
│   ├── order.ts
│   ├── zone.ts
│   ├── admin.ts
│   ├── complaint.ts
│   ├── report.ts
│   └── api.ts                        # Generic API response types
│
├── public/
│   ├── logo.svg
│   ├── logo-orange.svg
│   └── favicon.ico
│
├── tailwind.config.ts
├── tsconfig.json
├── next.config.ts
├── postcss.config.mjs
├── components.json                   # ShadcnUI config
├── package.json
├── .env.local.example
└── Dockerfile
```

---

## 3. Mobile Apps — `apps/mobile/`

```
apps/mobile/
├── cakli/                            # ── User App ──
│   ├── lib/
│   │   ├── main.dart                 # App entry point
│   │   ├── app/
│   │   │   ├── routes/
│   │   │   │   ├── app_pages.dart    # Route definitions
│   │   │   │   └── app_routes.dart   # Route names
│   │   │   ├── bindings/
│   │   │   │   └── initial_binding.dart
│   │   │   └── themes/
│   │   │       └── app_theme.dart    # Warna CAKLI (#E04D04)
│   │   │
│   │   ├── core/
│   │   │   ├── constants/
│   │   │   │   ├── api_constants.dart
│   │   │   │   └── app_constants.dart
│   │   │   ├── network/
│   │   │   │   ├── api_client.dart   # Dio HTTP client
│   │   │   │   ├── api_interceptor.dart
│   │   │   │   └── api_exceptions.dart
│   │   │   ├── storage/
│   │   │   │   └── local_storage.dart  # GetStorage / SharedPreferences
│   │   │   └── utils/
│   │   │       ├── formatters.dart
│   │   │       └── validators.dart
│   │   │
│   │   ├── data/
│   │   │   ├── models/               # Data models (JSON serializable)
│   │   │   │   ├── user_model.dart
│   │   │   │   ├── order_model.dart
│   │   │   │   ├── location_model.dart
│   │   │   │   ├── voucher_model.dart
│   │   │   │   ├── chat_model.dart
│   │   │   │   └── rating_model.dart
│   │   │   ├── repositories/         # Repository implementations
│   │   │   │   ├── auth_repository.dart
│   │   │   │   ├── user_repository.dart
│   │   │   │   ├── order_repository.dart
│   │   │   │   ├── location_repository.dart
│   │   │   │   └── chat_repository.dart
│   │   │   └── providers/            # API providers (Dio calls)
│   │   │       ├── auth_provider.dart
│   │   │       ├── user_provider.dart
│   │   │       └── order_provider.dart
│   │   │
│   │   ├── modules/                  # Feature modules (GetX pattern)
│   │   │   ├── splashscreen/
│   │   │   │   ├── bindings/splashscreen_binding.dart
│   │   │   │   ├── controllers/splashscreen_controller.dart
│   │   │   │   └── views/splashscreen_view.dart
│   │   │   ├── home/
│   │   │   ├── setlokasi/
│   │   │   ├── pesanalamat/
│   │   │   ├── pesan/
│   │   │   ├── pesandriver/
│   │   │   ├── chat/
│   │   │   ├── rating/
│   │   │   ├── voucher/
│   │   │   ├── aktivitas/
│   │   │   ├── setting/
│   │   │   ├── editprofile/
│   │   │   ├── editemail/
│   │   │   ├── editpin/
│   │   │   ├── editalamat/
│   │   │   ├── tambahalamat/
│   │   │   ├── carialamat/
│   │   │   └── maptest/
│   │   │
│   │   └── widgets/                  # Shared reusable widgets
│   │       ├── route_location_card.dart
│   │       ├── cakli_button.dart
│   │       ├── loading_overlay.dart
│   │       └── custom_bottom_nav.dart
│   │
│   ├── assets/
│   │   ├── images/
│   │   │   ├── logo.png
│   │   │   └── logo-orange.png
│   │   └── fonts/
│   ├── test/
│   ├── pubspec.yaml
│   ├── analysis_options.yaml
│   └── README.md
│
├── cakli_driver/                     # ── Driver App ──
│   ├── lib/
│   │   ├── main.dart
│   │   ├── app/                      # (same pattern as user)
│   │   ├── core/                     # (same pattern as user)
│   │   ├── data/
│   │   │   ├── models/
│   │   │   │   ├── driver_model.dart
│   │   │   │   ├── order_model.dart
│   │   │   │   ├── earning_model.dart
│   │   │   │   ├── wallet_model.dart
│   │   │   │   └── chat_model.dart
│   │   │   ├── repositories/
│   │   │   └── providers/
│   │   │
│   │   ├── modules/                  # Feature modules
│   │   │   ├── splashscreen/
│   │   │   ├── home/                 # Peta + status + autobid
│   │   │   ├── terimaorder/          # Slide-to-act jemput
│   │   │   ├── antarorder/           # Slide-to-act antar
│   │   │   ├── pendapatan/           # Dashboard keuangan
│   │   │   ├── profil/               # Profil driver + menu
│   │   │   ├── pengaturan/           # Settings
│   │   │   ├── riwayat/              # Trip history + kalender
│   │   │   ├── chat/                 # Chat + foto
│   │   │   └── listchat/             # Daftar percakapan
│   │   │
│   │   └── widgets/
│   │
│   ├── assets/
│   ├── test/
│   ├── pubspec.yaml
│   └── README.md
│
└── packages/                         # Shared Flutter packages
    └── cakli_core/
        ├── lib/
        │   ├── constants/
        │   ├── models/               # Shared models (chat, location)
        │   ├── network/              # Shared Dio client
        │   ├── theme/                # Brand colors & text styles
        │   └── widgets/              # Shared widgets
        └── pubspec.yaml
```

---

## 4. Infrastructure — `infra/`

```
infra/
├── docker/
│   ├── api.Dockerfile                # Go multi-stage build
│   ├── web.Dockerfile                # Next.js standalone build
│   ├── nginx.Dockerfile              # Nginx custom config
│   └── minio/                        # MinIO init scripts
│       └── init-buckets.sh           # Auto-create buckets
│
├── scripts/
│   ├── setup.sh                      # First-time dev setup
│   ├── migrate.sh                    # Run database migrations
│   ├── seed.sh                       # Seed initial data
│   └── backup-db.sh                  # PostgreSQL backup
│
└── k8s/                              # Kubernetes (future scaling)
    ├── namespaces.yaml
    ├── api-deployment.yaml
    ├── web-deployment.yaml
    ├── postgres-statefulset.yaml
    ├── minio-statefulset.yaml
    ├── nginx-ingress.yaml
    └── configmaps/
```

---

## 5. Nginx — `apps/nginx/`

```
apps/nginx/
├── nginx.conf                        # Main config
├── conf.d/
│   ├── api.conf                      # Proxy → api:3000 (/api/v1/*)
│   ├── web.conf                      # Proxy → web:3001 (admin dashboard)
│   ├── minio.conf                    # Proxy → minio:9000 (/storage/*)
│   ├── ssl.conf                      # SSL/TLS settings
│   └── security.conf                 # Headers, rate limit, gzip
├── ssl/
│   ├── cakli.id.crt
│   └── cakli.id.key
└── Dockerfile
```

---

## 6. Docker Compose

```yaml
# docker-compose.yml (development)
services:
  postgres:
    image: postgres:15-alpine
    ports: ["5432:5432"]
    volumes: [postgres_data:/var/lib/postgresql/data]

  minio:
    image: minio/minio:latest
    command: server /data --console-address ":9001"
    ports: ["9000:9000", "9001:9001"]
    volumes: [minio_data:/data]

  api:
    build: { context: ., dockerfile: infra/docker/api.Dockerfile }
    ports: ["3000:3000"]
    depends_on: [postgres, minio]

  web:
    build: { context: ., dockerfile: infra/docker/web.Dockerfile }
    ports: ["3001:3000"]
    depends_on: [api]

  nginx:
    build: { context: ., dockerfile: infra/docker/nginx.Dockerfile }
    ports: ["80:80", "443:443"]
    depends_on: [api, web]

volumes:
  postgres_data:
  minio_data:
```

---

## 7. MinIO Bucket Structure

```
minio/
├── avatars/                          # User & driver profile photos
│   ├── users/{user_id}/avatar.jpg
│   └── drivers/{driver_id}/avatar.jpg
├── documents/                        # Driver documents (KTP, SIM, STNK)
│   └── drivers/{driver_id}/{doc_type}_{timestamp}.jpg
├── chat-media/                       # Chat attachments (photo)
│   └── orders/{order_id}/{message_id}.jpg
├── exports/                          # Generated reports (XLSX, PDF)
│   └── {admin_id}/{report_type}_{date}.xlsx
└── agreements/                       # Partnership agreement PDFs
    └── drivers/{driver_id}/agreement_v{version}.pdf
```

---

## 8. Environment Variables

```env
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=cakli
DB_PASSWORD=secret
DB_NAME=cakli_db
DB_SSL_MODE=disable

# JWT
JWT_SECRET=your-secret-key
JWT_EXPIRY=24h
JWT_REFRESH_EXPIRY=168h

# MinIO
MINIO_ENDPOINT=localhost:9000
MINIO_ACCESS_KEY=minioadmin
MINIO_SECRET_KEY=minioadmin
MINIO_USE_SSL=false
MINIO_BUCKET_AVATARS=avatars
MINIO_BUCKET_DOCUMENTS=documents
MINIO_BUCKET_CHAT=chat-media
MINIO_BUCKET_EXPORTS=exports

# API
API_PORT=3000
API_BASE_URL=http://localhost:3000

# Web
NEXT_PUBLIC_API_URL=http://localhost:3000/api/v1
NEXT_PUBLIC_WS_URL=ws://localhost:3000/ws
```

---

## 9. Makefile (Root)

```makefile
# Development
dev-api:        cd apps/api && go run cmd/server/main.go
dev-web:        cd apps/web && npm run dev
dev-infra:      docker compose up postgres minio -d
dev:            make dev-infra && make dev-api & make dev-web

# Database
db-migrate:     cd apps/api && go run cmd/migrate/main.go up
db-seed:        cd apps/api && go run cmd/seed/main.go
db-reset:       cd apps/api && go run cmd/migrate/main.go down && make db-migrate && make db-seed

# Build
build-api:      cd apps/api && go build -o bin/server cmd/server/main.go
build-web:      cd apps/web && npm run build

# Docker
up:             docker compose up -d --build
down:           docker compose down
logs:           docker compose logs -f

# Test
test-api:       cd apps/api && go test ./...
test-web:       cd apps/web && npm run test
lint-api:       cd apps/api && golangci-lint run
lint-web:       cd apps/web && npm run lint
```

---

## 10. Scalability Notes

| Concern | Strategy |
|---------|----------|
| **API Scaling** | Stateless Go API → horizontal scale via Docker replicas / K8s HPA |
| **Database** | PostgreSQL read replicas, connection pooling (PgBouncer), indexed queries |
| **File Storage** | MinIO distributed mode, CDN caching untuk assets statis |
| **Real-time** | WebSocket/SSE untuk live map, driver tracking, chat — dipisah ke dedicated service jika perlu |
| **Caching** | Redis untuk session, rate limiting, frequently queried data (future) |
| **Multi-Zone** | Zone-based data partitioning, regional filtering pada query level |
| **CI/CD** | GitHub Actions per-app pipeline, Docker image registry, staging → production |
| **Monitoring** | Prometheus + Grafana (metrics), structured logging (zerolog), health endpoints |
| **Mobile** | Shared `cakli_core` package untuk DRY antara user & driver app |
| **Security** | JWT + refresh token rotation, RBAC middleware, audit logging, rate limiting, CORS |
