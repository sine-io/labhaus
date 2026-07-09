export const TOKEN_STORAGE_KEY = 'labhaus_bearer_token';

export function buildFrontendHeaders(token) {
  const headers = {
    'Content-Type': 'application/json',
  };

  const normalized = normalizeBearerToken(token);
  if (normalized) {
    headers.Authorization = normalized;
  }

  return headers;
}

export function getStoredBearerToken(storage = globalThis.localStorage) {
  if (!storage) {
    return '';
  }

  return storage.getItem(TOKEN_STORAGE_KEY) ?? '';
}

export function storeBearerToken(token, storage = globalThis.localStorage) {
  const normalized = normalizeBearerToken(token);
  if (!storage || !normalized) {
    return;
  }

  storage.setItem(TOKEN_STORAGE_KEY, normalized);
}

export function clearBearerToken(storage = globalThis.localStorage) {
  if (!storage) {
    return;
  }

  storage.removeItem(TOKEN_STORAGE_KEY);
}

function normalizeBearerToken(token) {
  const trimmed = String(token ?? '').trim();
  if (!trimmed) {
    return '';
  }

  return trimmed.toLowerCase().startsWith('bearer ') ? trimmed : `Bearer ${trimmed}`;
}
