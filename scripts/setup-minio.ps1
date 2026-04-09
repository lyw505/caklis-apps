# Setup MinIO bucket policy for CAKLI
# This script configures the MinIO bucket to allow uploads from the browser

Write-Host "Setting up MinIO bucket policy..." -ForegroundColor Cyan

# Wait for MinIO to be ready
Write-Host "Waiting for MinIO to be ready..." -ForegroundColor Yellow
Start-Sleep -Seconds 5

# Configure mc alias
Write-Host "Configuring MinIO client..." -ForegroundColor Yellow
docker exec cakli-storage mc alias set local http://localhost:9000 minioadmin minioadmin

# Create bucket if not exists
Write-Host "Creating bucket..." -ForegroundColor Yellow
docker exec cakli-storage mc mb local/cakli --ignore-existing

# Set bucket policy to public (allows both GET and PUT)
Write-Host "Setting bucket policy to public..." -ForegroundColor Yellow
docker exec cakli-storage mc anonymous set public local/cakli

# Verify policy
Write-Host "Verifying bucket policy..." -ForegroundColor Yellow
docker exec cakli-storage mc anonymous get local/cakli

Write-Host ""
Write-Host "MinIO bucket policy configured successfully!" -ForegroundColor Green
Write-Host "Bucket: cakli" -ForegroundColor White
Write-Host "Policy: Public (allows GET and PUT operations)" -ForegroundColor White
