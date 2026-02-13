#!/bin/sh
set -e

# Replace the build-time placeholder with the runtime API_URL env var.
# Default to /v1 (relative) if API_URL is not set.
API_URL="${API_URL:-/v1}"

# Find and replace the placeholder in all JS files served by nginx.
find /usr/share/nginx/html/assets -name '*.js' -exec \
  sed -i "s|__API_URL__|${API_URL}|g" {} +

exec "$@"
