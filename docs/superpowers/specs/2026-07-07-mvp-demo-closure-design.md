# MVP Demo Closure Design

## Goal

Make the current MVP demonstrable from a clean local checkout:

1. Start infrastructure and the Go API without a real image provider.
2. Seed enough styles for recommendation to return useful results.
3. Verify the end-to-end flow: register/login -> recommend styles -> batch-generate images.

## Scope

This design covers issues #47, #48, and #49.

In scope:

- A local HTTP mock image provider that implements the same contract used by the Go backend's GPT-Image-2 adapter.
- Docker Compose wiring so the API can call that mock provider by default for local demos.
- An idempotent SQL seed file with a compact style dataset.
- A smoke-test script that exercises the MVP API path with curl.
- Documentation updates for the local demo flow.

Out of scope:

- Real provider integration changes.
- Article-to-video, TTS, FFmpeg, visual workflow editing, task monitoring, and template marketplace.
- Browser automation or full production E2E testing.

## Architecture

### Mock image provider

Add a small Go command at `backend/cmd/mock-image-provider`.

It exposes:

- `GET /health` returning healthy JSON.
- `POST /v1/generate` accepting `{ prompt, width, height, quality, style }`.
- `GET /images/:id.png` returning a deterministic PNG.

The generate endpoint returns:

```json
{
  "image_url": "http://mock-image-provider:8089/images/<id>.png",
  "created_at": "2026-07-07T00:00:00Z"
}
```

The service validates `Authorization: Bearer <key>` using `MOCK_IMAGE_PROVIDER_API_KEY`. The Docker Compose default key is only for local demo use.

### Docker Compose

Add a `mock-image-provider` service built from `backend/Dockerfile.mock-image-provider`.

The `api` service uses:

- `LABHAUS_IMAGE_PROVIDER_BASE_URL=${LABHAUS_IMAGE_PROVIDER_BASE_URL:-http://mock-image-provider:8089}`
- `LABHAUS_IMAGE_PROVIDER_API_KEY=${LABHAUS_IMAGE_PROVIDER_API_KEY:-dev-mock-key}`

Backend config still fails if those variables are absent when running the API outside Compose.

### Style seed data

Add `backend/seeds/styles.sql`.

The file uses fixed IDs and `INSERT ... ON CONFLICT (id) DO UPDATE` so it can be run repeatedly without duplicates. The dataset is intentionally small and covers the categories already used in recommendation tests: UI, retro, nature, cyberpunk, art, and luxury.

### Smoke test

Add `scripts/mvp-smoke.sh`.

The script:

1. Checks API health.
2. Registers a deterministic demo user, tolerating "already exists".
3. Logs in and captures a token.
4. Calls `/api/styles/recommend`.
5. Calls `/api/images/generate` with two prompts.
6. Fails if recommendations or generated image results are empty.

## Testing

- Unit-test the mock provider command with `httptest`.
- Syntax-check the smoke script with `bash -n`.
- Keep existing workspace gates:
  - `pnpm format:check`
  - `pnpm lint`
  - `pnpm turbo run typecheck`
  - `pnpm test`
  - `pnpm build`
  - `go test ./...` from `backend/`
  - `go build ./...` from `backend/`
  - `docker compose config`

## Risks

- Docker socket permissions may prevent building or running Compose in this environment. If that happens, verify Compose config and all non-Docker code paths, then report Docker as an environment limitation.
- The API currently initializes the recommender once at startup. Seed styles before starting the API for the demo flow.
