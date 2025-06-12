#!/bin/sh
set -e

DBSTRING="postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=disable"
VERSION="${MIGRATION_VERSION:-latest}"

echo "applying migrations up to version: ${VERSION}"
sleep 2

if [ "$VERSION" = "latest" ]; then
    goose -dir /goose/sql postgres "$DBSTRING" up
else
    goose -dir /goose/sql postgres "$DBSTRING" up-to "$VERSION"
fi
