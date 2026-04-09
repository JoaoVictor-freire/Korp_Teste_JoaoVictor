#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

# Shared secret so a token issued by stock-service is accepted by billing-service.
export JWT_SECRET="${JWT_SECRET:-dev-secret-change-me}"

export STOCK_SERVICE_HOST="${STOCK_SERVICE_HOST:-0.0.0.0}"
export STOCK_SERVICE_PORT="${STOCK_SERVICE_PORT:-8081}"
export BILLING_SERVICE_HOST="${BILLING_SERVICE_HOST:-0.0.0.0}"
export BILLING_SERVICE_PORT="${BILLING_SERVICE_PORT:-8082}"

echo "JWT_SECRET=$JWT_SECRET"
echo "Starting stock-service on ${STOCK_SERVICE_HOST}:${STOCK_SERVICE_PORT}"
echo "Starting billing-service on ${BILLING_SERVICE_HOST}:${BILLING_SERVICE_PORT}"

stock_pid=""
billing_pid=""

cleanup() {
  echo ""
  echo "Shutting down..."
  if [[ -n "${stock_pid}" ]]; then kill "${stock_pid}" 2>/dev/null || true; fi
  if [[ -n "${billing_pid}" ]]; then kill "${billing_pid}" 2>/dev/null || true; fi
  wait || true
}

trap cleanup INT TERM EXIT

go run ./cmd/stock-service &
stock_pid="$!"

go run ./cmd/billing-service &
billing_pid="$!"

echo ""
echo "Services up:"
echo "- stock-service:   http://localhost:${STOCK_SERVICE_PORT}"
echo "- billing-service: http://localhost:${BILLING_SERVICE_PORT}"
echo ""
echo "Press Ctrl+C to stop."

wait

