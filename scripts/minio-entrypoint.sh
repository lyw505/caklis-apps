# #!/bin/sh

# # Start MinIO server in background
# minio server /data --console-address ":9001" &

# # Wait for MinIO to be ready
# echo "Waiting for MinIO to start..."
# sleep 5

# # Configure mc client
# mc alias set local http://localhost:9000 minioadmin minioadmin

# # Create bucket if not exists
# mc mb local/cakli --ignore-existing

# # Set bucket policy to public
# mc anonymous set public local/cakli

# echo "MinIO bucket policy configured successfully!"

# # Keep container running
# wait


#!/bin/sh
set -e

# Start MinIO server in background
minio server /data --console-address ":9001" &
MINIO_PID=$!

# Wait for MinIO to be ready
echo "Waiting for MinIO to start..."
sleep 10

# Configure mc client
mc alias set local http://localhost:9000 minioadmin minioadmin

# Create bucket if not exists
mc mb local/cakli --ignore-existing

# Set bucket policy to public
mc anonymous set public local/cakli

echo "MinIO bucket policy configured successfully!"

# Keep container running
wait $MINIO_PID