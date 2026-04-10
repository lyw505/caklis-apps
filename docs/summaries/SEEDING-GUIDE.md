# Database Seeding Guide

## Overview

Database seeding adalah proses mengisi database dengan data awal (default data) yang diperlukan untuk development dan testing.

## Commands

### Initialize Database Schema

```bash
make db-init
```

**What it does:**
- Waits for PostgreSQL to be ready
- Creates all tables (admins, users, drivers, banks, refresh_tokens)
- Sets up indexes and constraints
- Only runs if tables don't exist yet

**When to use:**
- First time setup
- After `make down -v` (removes volumes)
- When database is empty

### Seed Database with Default Data

```bash
make db-seed
```

**What it does:**
- Creates 3 admin accounts with different roles
- Inserts 22 Indonesian banks
- Shows credentials after completion

**When to use:**
- After first `make dev`
- When you need fresh admin accounts
- After database reset

### Legacy Command

```bash
make seed  # Same as make db-seed
```

For backward compatibility, `make seed` still works.

## Default Data

### Admin Accounts (3)

| Role | Email | Password | Dashboard |
|------|-------|----------|-----------|
| Master Admin | `master@cakli.id` | `admin123` | /master-admin |
| Operation Admin | `operation@cakli.id` | `admin123` | /operation-admin |
| Reporting Admin | `reporting@cakli.id` | `admin123` | /reporting-admin |

### Banks (22)

Indonesian banks including:
- BRI, Mandiri, BNI, BCA
- BSI, CIMB Niaga, Danamon
- Permata, Maybank, Panin
- And 12 more...

## Workflow

### First Time Setup

```bash
# 1. Start services
make dev

# 2. Seed database (in another terminal)
make db-seed

# 3. Login with admin credentials
# Open http://localhost:3000
# Login: master@cakli.id / admin123
```

### After Restart

```bash
# Just start services
make dev

# No need to seed again (data is persistent)
```

### Reset Database

```bash
# 1. Stop and remove volumes
make down
docker volume rm caklis-apps_postgres_data

# 2. Start fresh
make dev

# 3. Seed again
make db-seed
```

## Why Separate Seeding?

### Before (Old Way)

```bash
make dev  # Always seeds, even if data exists
```

**Problems:**
- ❌ Slow startup (seeding every time)
- ❌ Duplicate data errors if run multiple times
- ❌ Can't control when to seed

### After (New Way)

```bash
make dev      # Fast startup, no seeding
make db-seed  # Seed only when needed
```

**Benefits:**
- ✅ Faster `make dev` (no unnecessary seeding)
- ✅ Seed only when you need it
- ✅ Better control over database state
- ✅ Cleaner workflow

## Troubleshooting

### Error: "relation already exists"

This means tables already exist. You don't need to run `make db-init` again.

**Solution:** Just run `make db-seed` to add data.

### Error: "duplicate key value"

This means data already exists (e.g., admin accounts).

**Solution:** 
- If you want fresh data, reset database first
- Or skip seeding if data is already there

### Can't Login with Default Credentials

**Check if seeding was successful:**

```bash
# Connect to database
docker exec -it cakli-db psql -U cakli -d cakli_db

# Check admins
SELECT email, role FROM admins;

# Should show 3 admins
```

**If no admins found:**

```bash
make db-seed
```

### Database Connection Failed

**Check if PostgreSQL is running:**

```bash
docker ps | grep cakli-db
```

**If not running:**

```bash
make dev
```

## Advanced Usage

### Seed Specific Data

If you want to seed only banks or only admins, you can run SQL directly:

```bash
# Seed only banks
docker exec -i cakli-db psql -U cakli -d cakli_db < docs/agile-development/v1/seed-banks.sql

# Seed only admins
docker exec -i cakli-db psql -U cakli -d cakli_db < docs/agile-development/v1/seed-multi-role-admins.sql
```

### Custom Seed Data

Create your own seed file:

```sql
-- custom-seed.sql
INSERT INTO users (name, email, phone) VALUES
  ('Test User 1', 'test1@example.com', '081234567890'),
  ('Test User 2', 'test2@example.com', '081234567891');
```

Run it:

```bash
docker exec -i cakli-db psql -U cakli -d cakli_db < custom-seed.sql
```

## Summary

**Commands:**
- `make db-init` - Initialize schema (once)
- `make db-seed` - Seed default data (when needed)
- `make seed` - Alias for db-seed

**Workflow:**
1. `make dev` - Start services
2. `make db-seed` - Seed database (first time only)
3. Login and develop

**Benefits:**
- Faster startup
- Better control
- Cleaner workflow

**Default Credentials:**
- Master: `master@cakli.id` / `admin123`
- Operation: `operation@cakli.id` / `admin123`
- Reporting: `reporting@cakli.id` / `admin123`
