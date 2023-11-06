#!/bin/sh

set -e

if [ "x$QUAN_BASE_URL" = "x" ]; then
    QUAN_BASE_URL="https://k4qy.com"
fi
if [ "x$QUAN_LENGTH" = "x" ]; then
    QUAN_LENGTH=5
fi
if [ "x$QUAN_DB_FILE" = "x" ]; then
    QUAN_DB_FILE="/data/quan/db"
fi
if [ "x$QUAN_CHAR_RANGE" = "x" ]; then
    QUAN_CHAR_RANGE=62
fi
if [ "x$QUAN_DEFAULT_REDIRECT_URL" = "x" ]; then
    QUAN_DEFAULT_REDIRECT_URL="https://k4qy.com"
fi
if [ "x$QUAN_LIST_SIZE" = "x" ]; then
    QUAN_LIST_SIZE=100
fi
if [ "x$QUAN_PORT" = "x" ]; then
    QUAN_PORT=8080
fi

exec "$@" \
    -base-url $QUAN_BASE_URL \
    -length $QUAN_LENGTH \
    -db-file "$QUAN_DB_FILE" \
    -char-range $QUAN_CHAR_RANGE \
    -default-redirect-url "$QUAN_DEFAULT_REDIRECT_URL" \
    -list-size $QUAN_LIST_SIZE
    -port $QUAN_PORT
