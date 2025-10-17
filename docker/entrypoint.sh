#!/bin/sh
# Fix ownership of mounted volume files and Go cache directories
# Use find to avoid errors on missing/temporary files
find /app -type f -exec chown 1000:1000 {} + 2>/dev/null || true
find /app -type d -exec chown 1000:1000 {} + 2>/dev/null || true

# Fix Go cache permissions
chown -R 1000:1000 /go/pkg/mod/cache 2>/dev/null || true
chown -R 1000:1000 /home/user/.cache 2>/dev/null || true

exec su-exec 1000:1000 "$@"