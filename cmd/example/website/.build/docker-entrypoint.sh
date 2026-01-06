#!/bin/sh

set -e

export GIN_MODE=release

PIPE_PATH=${NGINX_ACCESS_PIPE:-/tmp/nginx_access.pipe}

if [ -e "$PIPE_PATH" ] && [ ! -p "$PIPE_PATH" ]; then
  mv "$PIPE_PATH" "$PIPE_PATH".bak || true
fi
if [ ! -p "$PIPE_PATH" ]; then
  mkfifo "$PIPE_PATH"
  chmod 0666 "$PIPE_PATH" || true
fi

su-exec appuser /bin/website &
WEBSITE_PID=$!

trap 'kill -TERM "$WEBSITE_PID" 2>/dev/null || true; wait "$WEBSITE_PID"; exit' TERM INT

exec nginx -g 'daemon off;' &
NGINX_PID=$!

wait $WEBSITE_PID $NGINX_PID
kill -TERM $WEBSITE_PID $NGINX_PID 2>/dev/null || true
wait
