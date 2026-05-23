#!/bin/sh
set -e

BACKEND_URL="${BACKEND_URL:-http://urlapi:2233}"

export BACKEND_URL
envsubst '${BACKEND_URL}' < /etc/nginx/templates/default.conf.template > /etc/nginx/conf.d/default.conf

exec nginx -g 'daemon off;'
