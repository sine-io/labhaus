export function buildBackendHeaders(request) {
  const headers = {
    "Content-Type": "application/json",
  };

  const authorization = request.headers.get("authorization");
  if (authorization) {
    headers.Authorization = authorization;
  }

  return headers;
}

export function normalizeGenerateImageResponse(data) {
  const results = Array.isArray(data?.results)
    ? data.results
    : Array.isArray(data?.images)
      ? data.images
      : [];

  return {
    ...data,
    results,
    images: results,
  };
}

export function normalizeRecommendStyleRequest(body) {
  const query = String(body?.query ?? body?.prompt ?? "").trim();
  const limit = numericLimit(body?.limit ?? body?.top_k);

  const normalized = { query };
  if (limit !== undefined) {
    normalized.limit = limit;
  }

  return normalized;
}

function numericLimit(value) {
  if (value === undefined || value === null || value === "") {
    return undefined;
  }

  const parsed = Number(value);
  if (!Number.isFinite(parsed)) {
    return undefined;
  }

  return Math.trunc(parsed);
}
