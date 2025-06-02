#!/bin/bash

if [ -z "$POSTGRES_HOST" ] || [ -z "$POSTGRES_PORT" ] || \
   [ -z "$POSTGRES_USER" ] || [ -z "$POSTGRES_PASSWORD" ] || \
   [ -z "$POSTGRES_DB" ]; then
  echo "[$(date)] error: missing required environment variables!" >&2
  exit 1
fi

BACKUP_DIR=/backups
PGPASSWORD="$POSTGRES_PASSWORD"
BACKUP_FILE="${BACKUP_DIR}/backup_$(date +"%Y-%m-%d_%H-%M-%S").sql"

echo "[$(date)] info: creating backup: ${BACKUP_FILE}"
pg_dump -h "$POSTGRES_HOST" -p "$POSTGRES_PORT" -U "$POSTGRES_USER" -d "$POSTGRES_DB" > "$BACKUP_FILE"
echo "[$(date)] info: backup created: ${BACKUP_FILE}"

echo "[$(date)] info: starting cleaning"
ls -t "${BACKUP_DIR}"/*.sql | tail -n +$(("$BACKUP_RETENTION_COUNT" + 1)) | xargs rm -f --
echo "[$(date)] info: keeping last "$BACKUP_RETENTION_COUNT" files"