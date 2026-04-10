-- ============================================================
-- SEED DATA — Users
-- ============================================================
-- Script untuk menambahkan user mobile untuk testing
-- Password: admin123 (bcrypt hash)
-- ============================================================

-- Bersihkan data jika ada duplikasi (opsional)
-- DELETE FROM users WHERE email = 'user@cakli.id';

INSERT INTO users (name, email, phone, password_hash, is_active) VALUES
    ('Test User', 'user@cakli.id', '081234567890', '$2a$12$V/cPhOf3rnixcFskxCUy/OsFSgmRura4SehRWMVUO6pM3toKaFqbq', true)
ON CONFLICT (email) WHERE (deleted_at IS NULL) DO UPDATE SET
    name = EXCLUDED.name,
    phone = EXCLUDED.phone,
    password_hash = EXCLUDED.password_hash,
    is_active = EXCLUDED.is_active,
    updated_at = NOW();

-- Verifikasi
SELECT id, name, email FROM users WHERE email = 'user@cakli.id';
