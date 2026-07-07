# MVP Alignment Design

## Goal

Make the repository move toward a runnable MVP by applying four decisions:

1. Use Go backend plus `apps/web` as the only active product line.
2. Fix the minimum frontend/backend contract for authenticated style recommendation and image generation.
3. Replace the implicit fake image-provider default with explicit provider configuration.
4. Narrow MVP documentation to the first deliverable: authenticated style recommendation plus batch image generation.

## Architecture

The active runtime is:

- Backend: `backend/` (`cmd/api/main.go`) with Gin, PostgreSQL, Redis, and MinIO.
- Frontend: `apps/web/` with Next App Router.
- Legacy/reference code: `apps/api/`, `packages/workflow/`, and `frontend/`.

`apps/web` proxies browser requests to the Go backend through App Router route handlers. The proxy forwards the browser's bearer token when one is present. Client pages read the Go backend response shape directly: image generation returns `results`, and style recommendation accepts `query` plus `limit`.

## Data flow

1. User authenticates against the Go backend user endpoints and receives a bearer token.
2. Frontend requests include `Authorization: Bearer <token>`.
3. Next route handlers forward that header to the Go backend.
4. Go backend calls the configured image provider and stores generated images in MinIO.
5. Frontend renders response data from the Go backend contract.

## Provider configuration

The image provider must not silently default to an example endpoint. Production or local integration must set:

- `LABHAUS_IMAGE_PROVIDER_API_KEY`
- `LABHAUS_IMAGE_PROVIDER_BASE_URL`

If either is missing, backend startup should fail with a clear configuration error before serving traffic.

## MVP scope

Current MVP excludes article-to-video, visual workflow editing, template marketplace, TTS, and FFmpeg rendering. Those stay documented as later phases.

## Testing

- Use small Node tests for frontend proxy/contract helpers.
- Use TypeScript type checking for `apps/web`.
- Use existing TypeScript tests for `packages/workflow` and `apps/api` where possible.
- Go tests/builds require a local Go toolchain or Docker access; if unavailable, report that as an environment limitation.
