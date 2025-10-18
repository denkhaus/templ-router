#!/bin/sh
# Fix ownership of mounted volume files and Go cache directories
# Use find to avoid errors on missing/temporary files

# Create and fix Go cache directories with proper permissions
mkdir -p /go/pkg/mod/cache 2>/dev/null || true
mkdir -p /home/user/.cache/go-build 2>/dev/null || true

# Fix Go cache permissions recursively
chown -R 1000:1000 /go/pkg/mod 2>/dev/null || true
chown -R 1000:1000 /home/user/.cache 2>/dev/null || true

# Set proper permissions for Go module directories
chmod -R 755 /go/pkg/mod 2>/dev/null || true
chmod -R 755 /home/user/.cache 2>/dev/null || true

exec su-exec 1000:1000 "$@"
