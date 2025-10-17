#!/bin/sh
# Fix ownership of mounted volume files
# Use find to avoid errors on missing/temporary files
find /app -type f -exec chown 1000:1000 {} + 2>/dev/null || true
find /app -type d -exec chown 1000:1000 {} + 2>/dev/null || true
exec su-exec 1000:1000 "$@"