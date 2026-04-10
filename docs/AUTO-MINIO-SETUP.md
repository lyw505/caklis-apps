# Auto MinIO Setup - No Manual Steps!

## ✅ Good News!

`make dev` sekarang **otomatis** configure MinIO bucket policy!

## What Happens Automatically

When you run `make dev`, it will:

1. ✅ Start Docker containers (PostgreSQL, MinIO, Backend API)
2. ✅ Wait 5 seconds for services to be ready
3. ✅ Initialize database schema
4. ✅ Seed database with default data
5. ✅ **Configure MinIO bucket policy** (NEW!)
   - Set mc alias
   - Create bucket if not exists
   - Set bucket policy to public
6. ✅ Sync environment variables
7. ✅ Start frontend development server

## No More Manual Steps!

You don't need to run `make setup-minio` or `make fix-upload` anymore when using `make dev`.

## Workflow

### First Time Setup

```bash
# Clone repository
git clone <repo-url>
cd caklis-apps

# Copy environment file
cp .env.example .env

# Start everything (ONE COMMAND!)
make dev
```

That's it! Everything is configured automatically.

### After Restart

```bash
# Just run make dev again
make dev
```

MinIO bucket policy will be automatically configured every time!

## Manual Fix (If Needed)

If you start containers manually with `docker-compose up -d`, you still need to run:

```bash
make fix-upload
```

But if you use `make dev`, it's automatic! ✅

## Technical Details

### What `make dev` Does for MinIO

```bash
# 1. Configure mc client
docker exec cakli-storage mc alias set local http://localhost:9000 minioadmin minioadmin

# 2. Create bucket (if not exists)
docker exec cakli-storage mc mb local/cakli --ignore-existing

# 3. Set bucket policy to public
docker exec cakli-storage mc anonymous set public local/cakli
```

### Why This Works

- Commands run after containers start
- 5-second wait ensures MinIO is ready
- `--ignore-existing` prevents errors if bucket already exists
- Output redirected to null for clean console

## Troubleshooting

### If Upload Still Fails

1. Check if containers are running:
```bash
docker ps
```

2. Manually run fix:
```bash
make fix-upload
```

3. Check bucket policy:
```bash
docker exec cakli-storage mc anonymous get local/cakli
```

Expected: `Access permission for 'local/cakli' is 'public'`

### If MinIO Container Not Found

Make sure you're using `make dev`, not `docker-compose up -d` directly.

Or start MinIO first:
```bash
docker-compose up -d minio
sleep 5
make fix-upload
```

## Benefits

✅ **Zero manual steps** - Everything automatic
✅ **Works after restart** - No need to remember commands
✅ **Consistent setup** - Same for all developers
✅ **Fast onboarding** - New developers just run `make dev`

## Summary

**Old Way:**
```bash
docker-compose up -d
make setup-minio  # Manual step!
make dev
```

**New Way:**
```bash
make dev  # Everything automatic! 🎉
```

**After Restart:**
```bash
make dev  # Still automatic! 🚀
```

No more "upload not working after restart" issues!
