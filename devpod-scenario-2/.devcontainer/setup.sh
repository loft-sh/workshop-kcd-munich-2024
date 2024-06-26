#!/usr/bin/env bash

# Write .env file
OUT=$(pwd)/devpod-scenario-2/.env

touch $OUT

VARS=("DB_HOST" "DB_PORT" "DB_USERNAME" "DB_PASSWORD")

for VAR in "${VARS[@]}"; do
  if ! grep -q "^$VAR=" "$OUT"; then
    echo "$VAR=${!VAR}" >> "$OUT"
  fi
done

