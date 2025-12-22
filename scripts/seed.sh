#!/bin/sh
set -e

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

ENV_FILE="${PROJECT_ROOT}/.env"

if [ ! -f "$ENV_FILE" ]; then
  echo ".env not found at $ENV_FILE"
  exit 1
fi

set -a
. "$ENV_FILE"
set +a

SEED_FILE="${PROJECT_ROOT}/testdata/couriers.sql"

if [ ! -f "$SEED_FILE" ]; then
  echo "Seed file not found at $SEED_FILE"
  exit 1
fi

echo "Applying couriers seed from ${SEED_FILE}..."

export PGPASSWORD="${POSTGRES_PASSWORD}"
if [ -n "${POSTGRES_SSLMODE}" ]; then
  export PGSSLMODE="${POSTGRES_SSLMODE}"
fi

psql \
  -h "${POSTGRES_HOST}" \
  -p "${POSTGRES_PORT}" \
  -U "${POSTGRES_USER}" \
  -d "${POSTGRES_DB}" \
  -v ON_ERROR_STOP=1 \
  -f "${SEED_FILE}"

echo "✅ Couriers seed applied successfully."

