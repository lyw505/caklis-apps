-- ============================================================
-- SEED DATA — Multi-Role Admins
-- ============================================================
-- Script untuk menambahkan admin dengan role berbeda
-- Password untuk semua admin: admin123
-- 
-- Cara menjalankan:
-- cat docs/agile-development/v1/seed-multi-role-admins.sql | docker exec -i cakli-db psql -U cakli -d cakli_db
-- ============================================================

-- Hapus admin lama jika ada (opsional, hati-hati di production!)
-- DELETE FROM admins WHERE email IN ('master@cakli.id', 'operation@cakli.id', 'reporting@cakli.id');

-- Insert 3 admin dengan role berbeda
-- Password: admin123 (bcrypt hash dengan cost 12)
INSERT INTO admins (name, email, password_hash, role, is_active) VALUES
    ('Super Admin', 'master@cakli.id', '$2a$12$V/cPhOf3rnixcFskxCUy/OsFSgmRura4SehRWMVUO6pM3toKaFqbq', 'master_admin', true),
    ('Operation Admin', 'operation@cakli.id', '$2a$12$V/cPhOf3rnixcFskxCUy/OsFSgmRura4SehRWMVUO6pM3toKaFqbq', 'operation_admin', true),
    ('Reporting Admin', 'reporting@cakli.id', '$2a$12$V/cPhOf3rnixcFskxCUy/OsFSgmRura4SehRWMVUO6pM3toKaFqbq', 'reporting_admin', true)
ON CONFLICT (email) WHERE deleted_at IS NULL DO UPDATE SET
    name = EXCLUDED.name,
    password_hash = EXCLUDED.password_hash,
    role = EXCLUDED.role,
    is_active = EXCLUDED.is_active,
    updated_at = NOW();

-- Verifikasi data
SELECT id, name, email, role, is_active, created_at FROM admins ORDER BY role;
