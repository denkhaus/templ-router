#!/bin/sh
# Fix ownership of mounted volume files
chown -R 1000:1000 /app
exec su-exec 1000:1000 "$@"