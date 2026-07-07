# MVP Demo Closure Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Make the local MVP flow demonstrable without a real image provider.

**Architecture:** Add a standalone Go mock provider service that satisfies the existing GPT-Image-2 HTTP contract, wire Docker Compose to use it by default for local demos, add idempotent style seed SQL, and add a curl-based MVP smoke script.

**Tech Stack:** Go standard library, Docker Compose, PostgreSQL SQL seed, bash/curl/jq, existing pnpm/Turbo workspace.

---

### Task 1: Mock image provider service

**Files:**

- Create: `backend/cmd/mock-image-provider/main_test.go`
- Create: `backend/cmd/mock-image-provider/main.go`
- Create: `backend/Dockerfile.mock-image-provider`

- [ ] Write `main_test.go` first. Cover:
  - `GET /health` returns 200.
  - `POST /v1/generate` without auth returns 401.
  - `POST /v1/generate` with valid auth returns `image_url` and `created_at`.
  - Returned image URL serves PNG bytes.
- [ ] Run `PATH=/tmp/labhaus-go/go/bin:$PATH go test ./cmd/mock-image-provider` from `backend/` and verify it fails because implementation is missing.
- [ ] Implement `main.go` with standard-library HTTP handlers and a deterministic PNG response.
- [ ] Add `backend/Dockerfile.mock-image-provider` that builds `./cmd/mock-image-provider`.
- [ ] Run `PATH=/tmp/labhaus-go/go/bin:$PATH go test ./cmd/mock-image-provider` and verify it passes.

### Task 2: Docker Compose local mock wiring

**Files:**

- Modify: `docker-compose.yml`
- Modify: `backend/.env.example`
- Modify: `docs/guides/quick-start.md`

- [ ] Add `mock-image-provider` service listening on container port `8089`.
- [ ] Change API provider env interpolation to local-demo defaults:
  - `LABHAUS_IMAGE_PROVIDER_BASE_URL=${LABHAUS_IMAGE_PROVIDER_BASE_URL:-http://mock-image-provider:8089}`
  - `LABHAUS_IMAGE_PROVIDER_API_KEY=${LABHAUS_IMAGE_PROVIDER_API_KEY:-dev-mock-key}`
- [ ] Make `api` depend on mock provider health.
- [ ] Document that raw backend startup still requires explicit provider env, while Compose supplies demo defaults.
- [ ] Run `docker compose config` and verify the generated config includes the mock provider service.

### Task 3: Idempotent style seed data

**Files:**

- Create: `backend/seeds/styles.sql`
- Modify: `docs/guides/quick-start.md`

- [ ] Add fixed-ID style rows for UI, retro, nature, cyberpunk, art, and luxury.
- [ ] Use `ON CONFLICT (id) DO UPDATE` to make repeated imports safe.
- [ ] Document:
  - `docker compose exec -T postgres psql -U labhaus -d labhaus < backend/seeds/styles.sql`
  - seed before starting/restarting the API so the startup recommender sees the styles.
- [ ] Run a text check to confirm the SQL contains `ON CONFLICT (id) DO UPDATE`.

### Task 4: MVP smoke script

**Files:**

- Create: `scripts/mvp-smoke.sh`
- Modify: `docs/guides/quick-start.md`

- [ ] Add a bash script that uses `curl` and `jq`.
- [ ] The script should register/login a demo user, call style recommendation, call image generation, and assert non-empty results.
- [ ] Run `bash -n scripts/mvp-smoke.sh`.
- [ ] Document how to run `scripts/mvp-smoke.sh`.

### Task 5: Full verification and issue updates

**Files:**

- Modify only if verification exposes documentation gaps.

- [ ] Run `pnpm format:check`.
- [ ] Run `pnpm lint`.
- [ ] Run `pnpm turbo run typecheck`.
- [ ] Run `pnpm test`.
- [ ] Run `pnpm build`.
- [ ] Run `PATH=/tmp/labhaus-go/go/bin:$PATH go test -count=1 ./...` from `backend/`.
- [ ] Run `PATH=/tmp/labhaus-go/go/bin:$PATH go build ./...` from `backend/`.
- [ ] Run `docker compose config`.
- [ ] If Docker socket is available, run `docker compose build mock-image-provider api`.
- [ ] Update #48 and #49 with implementation notes; update #47 with verification status.
