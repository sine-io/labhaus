#!/usr/bin/env bash
set -euo pipefail

API_URL="${API_URL:-http://localhost:8080}"
EMAIL="${EMAIL:-demo+mvp@labhaus.io}"
PASSWORD="${PASSWORD:-SecurePassword123!}"
NAME="${NAME:-MVP Demo User}"

HTTP_STATUS=""
BODY_FILE="$(mktemp)"

cleanup() {
  rm -f "$BODY_FILE"
}
trap cleanup EXIT

fail() {
  printf 'ERROR: %s\n' "$1" >&2
  if [[ -s "$BODY_FILE" ]]; then
    printf 'Last response body:\n' >&2
    cat "$BODY_FILE" >&2
    printf '\n' >&2
  fi
  exit 1
}

require_command() {
  if ! command -v "$1" >/dev/null 2>&1; then
    fail "missing required command: $1"
  fi
}

request() {
  local output_file="$1"
  shift
  HTTP_STATUS="$(curl -sS -o "$output_file" -w '%{http_code}' "$@")"
}

expect_status() {
  local expected="$1"
  local label="$2"
  if [[ "$HTTP_STATUS" != "$expected" ]]; then
    fail "$label returned HTTP $HTTP_STATUS, expected $expected"
  fi
}

require_command curl
require_command jq

printf 'Checking API health at %s...\n' "$API_URL"
request "$BODY_FILE" "$API_URL/api/health"
expect_status "200" "health check"

register_payload="$(jq -nc \
  --arg email "$EMAIL" \
  --arg password "$PASSWORD" \
  --arg name "$NAME" \
  '{email: $email, password: $password, name: $name}')"

printf 'Registering demo user %s...\n' "$EMAIL"
request "$BODY_FILE" \
  -X POST "$API_URL/api/users/register" \
  -H "Content-Type: application/json" \
  -d "$register_payload"

if [[ "$HTTP_STATUS" != "201" && "$HTTP_STATUS" != "409" ]]; then
  fail "register returned HTTP $HTTP_STATUS, expected 201 or 409"
fi

login_payload="$(jq -nc \
  --arg email "$EMAIL" \
  --arg password "$PASSWORD" \
  '{email: $email, password: $password}')"

printf 'Logging in...\n'
request "$BODY_FILE" \
  -X POST "$API_URL/api/users/login" \
  -H "Content-Type: application/json" \
  -d "$login_payload"
expect_status "200" "login"

TOKEN="$(jq -r '.token // empty' "$BODY_FILE")"
if [[ -z "$TOKEN" ]]; then
  fail "login response did not include .token"
fi

printf 'Requesting style recommendations...\n'
request "$BODY_FILE" \
  -X POST "$API_URL/api/styles/recommend" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"query":"modern UI dashboard for a SaaS product","limit":5}'
expect_status "200" "style recommendation"

recommendation_count="$(jq '.recommendations | length' "$BODY_FILE")"
if [[ "$recommendation_count" -le 0 ]]; then
  fail "style recommendation returned no results"
fi

style_prompt="$(jq -r '.recommendations[0].prompt // "modern UI dashboard"' "$BODY_FILE")"
image_payload="$(jq -nc \
  --arg style "$style_prompt" \
  '{prompts: ["modern UI dashboard hero", "minimal SaaS product card"], width: 512, height: 512, quality: "standard", style: $style}')"

printf 'Generating images through mock provider...\n'
request "$BODY_FILE" \
  -X POST "$API_URL/api/images/generate" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "$image_payload"
expect_status "200" "image generation"

image_count="$(jq '.results | length' "$BODY_FILE")"
success_count="$(jq '.success // 0' "$BODY_FILE")"
if [[ "$image_count" -le 0 || "$success_count" -le 0 ]]; then
  fail "image generation returned no successful results"
fi

printf 'MVP smoke passed: recommendations=%s generated_images=%s\n' "$recommendation_count" "$image_count"
