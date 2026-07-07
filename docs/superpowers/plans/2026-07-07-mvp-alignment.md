# MVP Alignment Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Align Labhaus around a realistic first MVP: Go backend plus `apps/web`, authenticated style recommendation, and batch image generation.

**Architecture:** Keep the Go backend as the only active API runtime and `apps/web` as the active frontend. Keep legacy code as reference, but prevent it from being mistaken for the current product line. Fix contract mismatches through a small frontend API-contract helper and explicit backend provider configuration.

**Tech Stack:** Go/Gin backend, Next App Router frontend, Node built-in test runner for frontend helpers, TypeScript type checking.

---

### Task 1: Mark the active product line

**Files:**

- Modify: `.gitignore`
- Modify: `README.md`
- Modify: `docs/TECH_STACK.md`

- [ ] Add `frontend/` and generated package lock files to ignored legacy/generated paths.
- [ ] State that `backend/` and `apps/web/` are the active runtime.
- [ ] State that `apps/api`, `packages/workflow`, and `frontend` are legacy/reference unless explicitly revived.

### Task 2: Fix frontend proxy and response contracts with TDD

**Files:**

- Create: `apps/web/app/api/_lib/backend-contract.mjs`
- Create: `apps/web/app/api/_lib/backend-contract.test.mjs`
- Modify: `apps/web/app/api/images/generate/route.ts`
- Modify: `apps/web/app/api/styles/recommend/route.ts`
- Modify: `apps/web/app/images/generate/page.tsx`
- Modify: `apps/web/app/styles/recommend/page.tsx`

- [ ] Write failing tests for forwarding Authorization, reading `results`, and converting `top_k` to `limit`.
- [ ] Run `node --test apps/web/app/api/_lib/backend-contract.test.mjs` and verify the expected failures.
- [ ] Implement the helper and wire it into route/page code.
- [ ] Re-run the helper tests and `./node_modules/.bin/tsc --noEmit` in `apps/web`.

### Task 3: Make image provider configuration explicit

**Files:**

- Modify: `backend/cmd/api/main.go`
- Modify: `backend/.env.example`
- Modify: `docker-compose.yml`
- Modify: `backend/README.md`

- [ ] Fail startup when `LABHAUS_IMAGE_PROVIDER_API_KEY` or `LABHAUS_IMAGE_PROVIDER_BASE_URL` is missing.
- [ ] Document both environment variables.
- [ ] Keep provider base URL configurable so a real provider or local mock service can be used.

### Task 4: Narrow MVP documentation

**Files:**

- Modify: `README.md`
- Modify: `docs/planning/mvp-roadmap.md`

- [ ] Reframe the immediate MVP as style recommendation plus batch image generation.
- [ ] Move article-to-video, visual editor, TTS, FFmpeg, and template marketplace to later phases.
- [ ] Keep the broader vision intact but avoid claiming the incomplete workflow is currently available.

### Task 5: Verify and report

**Commands:**

- `node --test apps/web/app/api/_lib/backend-contract.test.mjs`
- `(cd apps/web && ./node_modules/.bin/tsc --noEmit)`
- `(cd packages/workflow && ./node_modules/.bin/vitest run)`
- `git diff --check`

Report any Go or Docker verification that cannot run in the current environment.
