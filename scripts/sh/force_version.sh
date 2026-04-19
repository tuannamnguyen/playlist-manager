#!/bin/bash

set -euo pipefail

urlencode() {
  local input="$1"
  local output=""
  local i char

  for ((i = 0; i < ${#input}; i++)); do
    char="${input:i:1}"
    case "$char" in
      [a-zA-Z0-9.~_-]) output+="$char" ;;
      *) printf -v output '%s%%%02X' "$output" "'$char" ;;
    esac
  done

  printf '%s' "$output"
}

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
MIGRATIONS_DIR="$SCRIPT_DIR/../migrations"
RDS_CA_BUNDLE="/home/ssm-user/playlist_manager/global-bundle.pem"

read -r -p "Enter PostgreSQL username: " PGUSER
read -r -s -p "Enter PostgreSQL password: " PGPASSWORD
echo
read -r -p "Enter PostgreSQL Host: " PGHOST
read -r -p "Enter PostgreSQL Port [5432]: " PGPORT
PGPORT="${PGPORT:-5432}"
read -r -p "Enter PostgreSQL database [playlist_manager]: " PGDATABASE
PGDATABASE="${PGDATABASE:-playlist_manager}"
read -r -p "Use SSL with RDS global bundle? [y/N]: " USE_SSL
read -r -p "Force version: " FORCE_VERSION

ENCODED_PGUSER="$(urlencode "$PGUSER")"
ENCODED_PGPASSWORD="$(urlencode "$PGPASSWORD")"

if [[ "$USE_SSL" =~ ^[Yy]$ ]]; then
  ENCODED_RDS_CA_BUNDLE="$(urlencode "$RDS_CA_BUNDLE")"
  DATABASE_URL="postgres://${ENCODED_PGUSER}:${ENCODED_PGPASSWORD}@${PGHOST}:${PGPORT}/${PGDATABASE}?sslmode=verify-full&sslrootcert=${ENCODED_RDS_CA_BUNDLE}"
else
  DATABASE_URL="postgres://${ENCODED_PGUSER}:${ENCODED_PGPASSWORD}@${PGHOST}:${PGPORT}/${PGDATABASE}?sslmode=disable"
fi

migrate -path "$MIGRATIONS_DIR" -database "$DATABASE_URL" force "$FORCE_VERSION"
