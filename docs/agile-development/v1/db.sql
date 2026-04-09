-- ============================================================
-- CAKLI Database Schema — v1 (Agile Sprint 1)
-- PostgreSQL — Normalized & Best Practice
-- ============================================================

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- ============================================================
-- 1. BANKS (Reference Table)
-- ============================================================
-- Daftar bank Indonesia untuk dropdown pilihan rekening driver.
-- Data ini di-seed sekali dan jarang berubah.

CREATE TABLE banks (
    id          SERIAL          PRIMARY KEY,
    name        VARCHAR(100)    NOT NULL,               -- nama bank (contoh: "Bank BRI")
    code        VARCHAR(10)     NOT NULL,                -- kode bank (contoh: "002") — unique dihandle GORM
    logo_url    VARCHAR(500),                            -- URL logo bank (opsional)
    is_active   BOOLEAN         NOT NULL DEFAULT TRUE,
    created_at  TIMESTAMPTZ     NOT NULL DEFAULT NOW()
);

-- unique index untuk banks.code dihandle oleh GORM AutoMigrate (uni_banks_code)
CREATE INDEX idx_banks_is_active ON banks(is_active);

-- ============================================================
-- 2. ADMINS
-- ============================================================
-- Akun admin untuk login ke web dashboard.
-- v1 hanya mendukung role 'master_admin'.

CREATE TABLE admins (
    id              UUID            PRIMARY KEY DEFAULT gen_random_uuid(),
    name            VARCHAR(255)    NOT NULL,
    email           VARCHAR(255)    NOT NULL,
    password_hash   VARCHAR(255)    NOT NULL,           -- bcrypt hash
    role            VARCHAR(50)     NOT NULL DEFAULT 'master_admin'
                                    CHECK (role IN ('master_admin', 'operation_admin', 'reporting_admin')),
    is_active       BOOLEAN         NOT NULL DEFAULT TRUE,
    created_at      TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ     -- Langsung tambahkan di sini
);

CREATE UNIQUE INDEX idx_admins_email ON admins(email) WHERE deleted_at IS NULL;

-- Tambahkan soft delete support
-- ALTER TABLE admins ADD COLUMN deleted_at TIMESTAMPTZ; ini garai error bit

-- ============================================================
-- 3. REFRESH TOKENS
-- ============================================================
-- Menyimpan refresh token yang dihash untuk JWT rotation.
-- Satu admin bisa punya beberapa refresh token (multi-device).
-- Token lama direvoke saat rotation.

CREATE TABLE refresh_tokens (
    id          UUID            PRIMARY KEY DEFAULT gen_random_uuid(),
    admin_id    UUID            NOT NULL REFERENCES admins(id) ON DELETE CASCADE,
    token_hash  VARCHAR(255)    NOT NULL,               -- SHA-256 hash dari refresh token
    user_agent  VARCHAR(500),                            -- browser/device info
    ip_address  VARCHAR(45),                             -- IPv4 atau IPv6
    expires_at  TIMESTAMPTZ     NOT NULL,
    revoked_at  TIMESTAMPTZ,                             -- NULL = aktif, ada nilai = sudah direvoke
    created_at  TIMESTAMPTZ     NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_refresh_tokens_admin_id ON refresh_tokens(admin_id);
CREATE INDEX idx_refresh_tokens_token_hash ON refresh_tokens(token_hash);
CREATE INDEX idx_refresh_tokens_expires_at ON refresh_tokens(expires_at);

-- ============================================================
-- 4. DRIVERS
-- ============================================================
-- Data pengemudi becak listrik.
-- File foto disimpan di MinIO, kolom hanya menyimpan object key.

CREATE TABLE drivers (
    id                      UUID            PRIMARY KEY DEFAULT gen_random_uuid(),
    name                    VARCHAR(255)    NOT NULL,
    email                   VARCHAR(255)    NOT NULL,
    phone                   VARCHAR(20)     NOT NULL,
    password_hash           VARCHAR(255)    NOT NULL,           -- bcrypt hash
    nik                     VARCHAR(16)     NOT NULL,           -- Nomor Induk Kependudukan (16 digit)
    birth_place             VARCHAR(100)    NOT NULL,           -- tempat lahir
    birth_date              DATE            NOT NULL,           -- tanggal lahir
    bank_id                 INT             REFERENCES banks(id) ON DELETE SET NULL,
    bank_account_number     VARCHAR(30),                        -- nomor rekening
    photo_profile_key       VARCHAR(500),                       -- MinIO object key
    photo_ktp_key           VARCHAR(500),                       -- MinIO object key
    photo_face_key          VARCHAR(500),                       -- MinIO object key (foto muka)
    verification_status     VARCHAR(20)     NOT NULL DEFAULT 'pending'
                                            CHECK (verification_status IN ('pending', 'accepted', 'rejected')),
    is_active               BOOLEAN         NOT NULL DEFAULT TRUE,
    created_at              TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
    updated_at              TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
    deleted_at              TIMESTAMPTZ                          -- soft delete
);

CREATE UNIQUE INDEX idx_drivers_email ON drivers(email) WHERE deleted_at IS NULL;
CREATE UNIQUE INDEX idx_drivers_phone ON drivers(phone) WHERE deleted_at IS NULL;
CREATE UNIQUE INDEX idx_drivers_nik ON drivers(nik) WHERE deleted_at IS NULL;
CREATE INDEX idx_drivers_verification_status ON drivers(verification_status);
CREATE INDEX idx_drivers_is_active ON drivers(is_active);
CREATE INDEX idx_drivers_deleted_at ON drivers(deleted_at);
CREATE INDEX idx_drivers_bank_id ON drivers(bank_id);

-- ============================================================
-- 5. USERS
-- ============================================================
-- Data pengguna/penumpang aplikasi CAKLI.

CREATE TABLE users (
    id                  UUID            PRIMARY KEY DEFAULT gen_random_uuid(),
    name                VARCHAR(255)    NOT NULL,
    email               VARCHAR(255)    NOT NULL,
    phone               VARCHAR(20)     NOT NULL,
    password_hash       VARCHAR(255)    NOT NULL,               -- bcrypt hash
    photo_profile_key   VARCHAR(500),                            -- MinIO object key
    is_active           BOOLEAN         NOT NULL DEFAULT TRUE,
    created_at          TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
    deleted_at          TIMESTAMPTZ                              -- soft delete
);

CREATE UNIQUE INDEX idx_users_email ON users(email) WHERE deleted_at IS NULL;
CREATE UNIQUE INDEX idx_users_phone ON users(phone) WHERE deleted_at IS NULL;
CREATE INDEX idx_users_is_active ON users(is_active);
CREATE INDEX idx_users_deleted_at ON users(deleted_at);

-- ============================================================
-- 6. SEED DATA — BANKS
-- ============================================================
-- Daftar bank umum di Indonesia

INSERT INTO banks (name, code) VALUES
    ('Bank BRI', '002'),
    ('Bank Mandiri', '008'),
    ('Bank BNI', '009'),
    ('Bank BCA', '014'),
    ('Bank BTN', '200'),
    ('Bank CIMB Niaga', '022'),
    ('Bank Danamon', '011'),
    ('Bank Permata', '013'),
    ('Bank Panin', '019'),
    ('Bank OCBC NISP', '028'),
    ('Bank Mega', '426'),
    ('Bank Bukopin', '441'),
    ('Bank Sinarmas', '153'),
    ('Bank Maybank', '016'),
    ('Bank BSI', '451'),
    ('Bank Jago', '542'),
    ('Bank Neo Commerce', '490'),
    ('Bank Seabank', '535'),
    ('Bank BTPN', '213'),
    ('Bank DKI', '111'),
    ('Bank Jatim', '114'),
    ('Bank Jateng', '113');

-- ============================================================
-- 7. SEED DATA — DEFAULT ADMINS (Multi-Role)
-- ============================================================
-- Password untuk semua admin: admin123 (bcrypt hash)
-- PENTING: Ganti password ini di production!

INSERT INTO admins (name, email, password_hash, role) VALUES
    ('Super Admin', 'master@cakli.id', '$2a$12$V/cPhOf3rnixcFskxCUy/OsFSgmRura4SehRWMVUO6pM3toKaFqbq', 'master_admin'),
    ('Operation Admin', 'operation@cakli.id', '$2a$12$V/cPhOf3rnixcFskxCUy/OsFSgmRura4SehRWMVUO6pM3toKaFqbq', 'operation_admin'),
    ('Reporting Admin', 'reporting@cakli.id', '$2a$12$V/cPhOf3rnixcFskxCUy/OsFSgmRura4SehRWMVUO6pM3toKaFqbq', 'reporting_admin')
ON CONFLICT (email) WHERE deleted_at IS NULL DO NOTHING;

-- ============================================================
-- 8. FUNCTIONS & TRIGGERS
-- ============================================================

-- Auto-update updated_at pada setiap UPDATE
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_admins_updated_at
    BEFORE UPDATE ON admins
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER trigger_drivers_updated_at
    BEFORE UPDATE ON drivers
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER trigger_users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- ============================================================
-- CATATAN DESAIN
-- ============================================================
--
-- 1. SOFT DELETE:
--    Semua tabel utama (admins, drivers, users) menggunakan soft delete
--    via kolom `deleted_at`. Record tidak benar-benar dihapus dari database.
--    Unique index menggunakan partial index (WHERE deleted_at IS NULL)
--    agar email/phone/nik bisa digunakan kembali setelah record dihapus.
--
-- 2. PRESIGNED URL (MinIO):
--    Kolom file (photo_profile_key, photo_ktp_key, photo_face_key)
--    menyimpan OBJECT KEY di MinIO, bukan URL penuh.
--    Contoh: "drivers/photo-profile/550e8400-e29b-41d4-a716-446655440000.jpg"
--    URL akses dihasilkan oleh backend via presigned GET URL saat dibutuhkan.
--
-- 3. PASSWORD:
--    Semua password di-hash menggunakan bcrypt (cost factor 12).
--    Password plain text TIDAK PERNAH disimpan.
--
-- 4. REFRESH TOKEN:
--    Disimpan sebagai SHA-256 hash, bukan plain text.
--    Mendukung token rotation: setiap kali refresh, token lama direvoke
--    dan token baru dikeluarkan.
--
-- 5. UUID:
--    Semua primary key menggunakan UUID v4 via pgcrypto gen_random_uuid()
--    untuk keamanan (tidak bisa ditebak/di-enumerate).
--
-- 6. TIMEZONE:
--    Semua kolom timestamp menggunakan TIMESTAMPTZ (timezone-aware)
--    untuk konsistensi global.

INSERT INTO users (
    id,
    name,
    email,
    phone,
    password_hash,
    photo_profile_key,
    is_active,
    created_at,
    updated_at
) VALUES
(
    gen_random_uuid(),
    'Ahmad Fauzi',
    'ahmad@example.com',
    '081234567801',
    '$2a$10$dummyhashahmad1234567890',
    'users/profile/ahmad.jpg',
    TRUE,
    NOW(),
    NOW()
),
(
    gen_random_uuid(),
    'Rina Kartika',
    'rina@example.com',
    '081234567802',
    '$2a$10$dummyhashrina1234567890',
    'users/profile/rina.jpg',
    TRUE,
    NOW(),
    NOW()
),
(
    gen_random_uuid(),
    'Dimas Saputra',
    'dimas@example.com',
    '081234567803',
    '$2a$10$dummyhashdimas1234567890',
    'users/profile/dimas.jpg',
    TRUE,
    NOW(),
    NOW()
),
(
    gen_random_uuid(),
    'Putri Ayu',
    'putri@example.com',
    '081234567804',
    '$2a$10$dummyhashputri1234567890',
    'users/profile/putri.jpg',
    FALSE,
    NOW(),
    NOW()
),
(
    gen_random_uuid(),
    'Yoga Pratama',
    'yoga@example.com',
    '081234567805',
    '$2a$10$dummyhashyoga1234567890',
    'users/profile/yoga.jpg',
    TRUE,
    NOW(),
    NOW()
);