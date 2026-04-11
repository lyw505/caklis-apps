FROM minio/minio:latest

# Create entrypoint script directly in the image
RUN echo '#!/bin/sh' > /usr/local/bin/minio-entrypoint.sh && \
    echo 'set -e' >> /usr/local/bin/minio-entrypoint.sh && \
    echo '' >> /usr/local/bin/minio-entrypoint.sh && \
    echo '# Start MinIO server in background' >> /usr/local/bin/minio-entrypoint.sh && \
    echo 'minio server /data --console-address ":9001" &' >> /usr/local/bin/minio-entrypoint.sh && \
    echo 'MINIO_PID=$!' >> /usr/local/bin/minio-entrypoint.sh && \
    echo '' >> /usr/local/bin/minio-entrypoint.sh && \
    echo '# Wait for MinIO to be ready' >> /usr/local/bin/minio-entrypoint.sh && \
    echo 'echo "Waiting for MinIO to start..."' >> /usr/local/bin/minio-entrypoint.sh && \
    echo 'sleep 10' >> /usr/local/bin/minio-entrypoint.sh && \
    echo '' >> /usr/local/bin/minio-entrypoint.sh && \
    echo '# Configure mc client' >> /usr/local/bin/minio-entrypoint.sh && \
    echo 'mc alias set local http://localhost:9000 minioadmin minioadmin' >> /usr/local/bin/minio-entrypoint.sh && \
    echo '' >> /usr/local/bin/minio-entrypoint.sh && \
    echo '# Create bucket if not exists' >> /usr/local/bin/minio-entrypoint.sh && \
    echo 'mc mb local/cakli --ignore-existing' >> /usr/local/bin/minio-entrypoint.sh && \
    echo '' >> /usr/local/bin/minio-entrypoint.sh && \
    echo '# Set bucket policy to public' >> /usr/local/bin/minio-entrypoint.sh && \
    echo 'mc anonymous set public local/cakli' >> /usr/local/bin/minio-entrypoint.sh && \
    echo '' >> /usr/local/bin/minio-entrypoint.sh && \
    echo 'echo "MinIO bucket policy configured successfully!"' >> /usr/local/bin/minio-entrypoint.sh && \
    echo '' >> /usr/local/bin/minio-entrypoint.sh && \
    echo '# Keep container running' >> /usr/local/bin/minio-entrypoint.sh && \
    echo 'wait $MINIO_PID' >> /usr/local/bin/minio-entrypoint.sh && \
    chmod +x /usr/local/bin/minio-entrypoint.sh

# Set entrypoint
ENTRYPOINT ["/usr/local/bin/minio-entrypoint.sh"]

