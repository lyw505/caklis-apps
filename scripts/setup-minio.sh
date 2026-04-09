#!/bin/bash

# Setup MinIO bucket policy for CAKLI
# This script configures the MinIO bucket to allow uploads from the browser

echo "🔧 Setting up MinIO bucket policy..."

# Wait for MinIO to be ready
echo "⏳ Waiting for MinIO to be ready..."
sleep 5

# Install mc (MinIO Client) if not already installed
docker exec cakli-storage mc alias set local http://localhost:9000 minioadmin minioadmin

# Create bucket if not exists
docker exec cakli-storage mc mb local/cakli --ignore-existing

# Set bucket policy to public (allows both GET and PUT)
docker exec cakli-storage mc anonymous set public local/cakli

# Verify policy
docker exec cakli-storage mc anonymous get local/cakli

echo "✅ MinIO bucket policy configured successfully!"
echo "📦 Bucket: cakli"
echo "🔓 Policy: Public (allows GET and PUT operations)"
