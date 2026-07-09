import assert from 'node:assert/strict';
import test from 'node:test';

import {
  buildBackendHeaders,
  normalizeGenerateImageResponse,
  normalizeRecommendStyleRequest,
} from './backend-contract.mjs';

test('buildBackendHeaders forwards bearer authorization from the browser request', () => {
  const request = new Request('http://localhost/api/images/generate', {
    headers: {
      Authorization: 'Bearer token-123',
    },
  });

  assert.deepEqual(buildBackendHeaders(request), {
    'Content-Type': 'application/json',
    Authorization: 'Bearer token-123',
  });
});

test('buildBackendHeaders omits authorization when request has no token', () => {
  const request = new Request('http://localhost/api/styles/recommend');

  assert.deepEqual(buildBackendHeaders(request), {
    'Content-Type': 'application/json',
  });
});

test('normalizeGenerateImageResponse exposes Go backend results as UI images', () => {
  const response = normalizeGenerateImageResponse({
    results: [
      { id: 'img-1.png', url: 'https://cdn.example/img-1.png', prompt: 'first' },
      { id: 'img-2.png', url: 'https://cdn.example/img-2.png', prompt: 'second' },
    ],
    total: 2,
    success: 2,
    failed: 0,
  });

  assert.deepEqual(response.images, [
    { id: 'img-1.png', url: 'https://cdn.example/img-1.png', prompt: 'first' },
    { id: 'img-2.png', url: 'https://cdn.example/img-2.png', prompt: 'second' },
  ]);
  assert.equal(response.total, 2);
  assert.equal(response.success, 2);
  assert.equal(response.failed, 0);
});

test('normalizeRecommendStyleRequest converts legacy top_k into Go backend limit', () => {
  assert.deepEqual(
    normalizeRecommendStyleRequest({
      query: 'cyberpunk city',
      top_k: 7,
    }),
    {
      query: 'cyberpunk city',
      limit: 7,
    }
  );
});

test('normalizeRecommendStyleRequest keeps explicit limit over top_k', () => {
  assert.deepEqual(
    normalizeRecommendStyleRequest({
      query: 'minimal product photo',
      top_k: 20,
      limit: 5,
    }),
    {
      query: 'minimal product photo',
      limit: 5,
    }
  );
});
