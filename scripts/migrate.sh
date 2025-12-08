#!/bin/sh
set -e

# Определяем путь к директории, где лежит этот скрипт
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

# Проверяем наличие .env
ENV_FILE="${PROJECT_ROOT}/.env"

if [ ! -f "$ENV_FILE" ]; then
  echo ".env not found at $ENV_FILE"
  exit 1
fi

set -a
. "$ENV_FILE"
set +a

GOOSE_DRIVER="postgres"
GOOSE_DBSTRING="${GOOSE_DRIVER}://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DB}"
MIGRATIONS_DIR="${PROJECT_ROOT}/migrations"

apply_migrations() {
  echo "🚀 Applying migrations from ${MIGRATIONS_DIR}..."
  GOOSE_DRIVER="${GOOSE_DRIVER}" GOOSE_DBSTRING="${GOOSE_DBSTRING}" \
    goose -dir "${MIGRATIONS_DIR}" up
}

rollback_migrations() {
  echo "↩️ Rolling back migrations..."
  GOOSE_DRIVER="${GOOSE_DRIVER}" GOOSE_DBSTRING="${GOOSE_DBSTRING}" \
    goose -dir "${MIGRATIONS_DIR}" down
}

show_status() {
  echo "📋 Migration status in ${MIGRATIONS_DIR}:"
  GOOSE_DRIVER="${GOOSE_DRIVER}" GOOSE_DBSTRING="${GOOSE_DBSTRING}" \
    goose -dir "${MIGRATIONS_DIR}" status
}

apply_test_migrations() {
  echo "Applying test migrations from ${MIGRATIONS_DIR_TEST}..."
  GOOSE_DRIVER="${GOOSE_DRIVER}" GOOSE_DBSTRING="${GOOSE_DBSTRING_TEST}" \
    goose -dir "${MIGRATIONS_DIR}" up
}

case "$1" in
  up)
    apply_migrations
    ;;
  down)
    rollback_migrations
    ;;
  status)
    show_status
    ;;
  *)
    echo "Usage: $0 {up|down|status|test-up}"
    exit 1
    ;;
esac