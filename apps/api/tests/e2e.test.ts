import { describe, it, expect, beforeAll, afterAll } from 'vitest';

const API_URL = process.env.API_URL || 'http://localhost:3001';

describe('E2E: Authentication Flow', () => {
  let accessToken: string;
  let refreshToken: string;
  const testUser = {
    email: `test-${Date.now()}@example.com`,
    password: 'TestPassword123!',
    name: 'Test User',
  };

  it('should register a new user', async () => {
    const response = await fetch(`${API_URL}/api/auth/register`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(testUser),
    });

    expect(response.status).toBe(201);
    const data = await response.json();

    expect(data.user.email).toBe(testUser.email);
    expect(data.user.name).toBe(testUser.name);
    expect(data.tokens.access_token).toBeDefined();
    expect(data.tokens.refresh_token).toBeDefined();

    accessToken = data.tokens.access_token;
    refreshToken = data.tokens.refresh_token;
  });

  it('should not allow duplicate registration', async () => {
    const response = await fetch(`${API_URL}/api/auth/register`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(testUser),
    });

    expect(response.status).toBe(409);
    const data = await response.json();
    expect(data.error).toBe('CONFLICT');
  });

  it('should login with credentials', async () => {
    const response = await fetch(`${API_URL}/api/auth/login`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        email: testUser.email,
        password: testUser.password,
      }),
    });

    expect(response.status).toBe(200);
    const data = await response.json();

    expect(data.user.email).toBe(testUser.email);
    expect(data.tokens.access_token).toBeDefined();
  });

  it('should reject invalid credentials', async () => {
    const response = await fetch(`${API_URL}/api/auth/login`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        email: testUser.email,
        password: 'wrongpassword',
      }),
    });

    expect(response.status).toBe(401);
  });

  it('should get current user with valid token', async () => {
    const response = await fetch(`${API_URL}/api/auth/me`, {
      headers: { Authorization: `Bearer ${accessToken}` },
    });

    expect(response.status).toBe(200);
    const data = await response.json();
    expect(data.user.email).toBe(testUser.email);
  });

  it('should reject request without token', async () => {
    const response = await fetch(`${API_URL}/api/auth/me`);
    expect(response.status).toBe(401);
  });

  it('should refresh access token', async () => {
    const response = await fetch(`${API_URL}/api/auth/refresh`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ refresh_token: refreshToken }),
    });

    expect(response.status).toBe(200);
    const data = await response.json();
    expect(data.tokens.access_token).toBeDefined();
    expect(data.tokens.refresh_token).toBeDefined();
  });
});

describe('E2E: Styles API', () => {
  it('should get health check', async () => {
    const response = await fetch(`${API_URL}/api/health`);
    expect(response.status).toBe(200);
    const data = await response.json();
    expect(data.status).toBe('ok');
  });

  it('should list styles', async () => {
    const response = await fetch(`${API_URL}/api/styles?limit=10`);
    expect(response.status).toBe(200);
    const data = await response.json();

    expect(data.styles).toBeDefined();
    expect(Array.isArray(data.styles)).toBe(true);
    expect(data.pagination).toBeDefined();
    expect(data.pagination.page).toBe(1);
    expect(data.pagination.limit).toBe(10);
  });

  it('should filter styles by category', async () => {
    const response = await fetch(
      `${API_URL}/api/styles?category=${encodeURIComponent('UI & Interfaces')}`
    );
    expect(response.status).toBe(200);
    const data = await response.json();

    if (data.styles.length > 0) {
      expect(data.styles[0].category).toBe('UI & Interfaces');
    }
  });

  it('should search styles', async () => {
    const response = await fetch(`${API_URL}/api/styles?search=portrait`);
    expect(response.status).toBe(200);
    const data = await response.json();
    expect(data.styles).toBeDefined();
  });

  it('should handle pagination', async () => {
    const response = await fetch(`${API_URL}/api/styles?page=2&limit=5`);
    expect(response.status).toBe(200);
    const data = await response.json();

    expect(data.pagination.page).toBe(2);
    expect(data.pagination.limit).toBe(5);
  });

  it('should return 404 for non-existent style', async () => {
    const fakeId = '00000000-0000-0000-0000-000000000000';
    const response = await fetch(`${API_URL}/api/styles/${fakeId}`);
    expect(response.status).toBe(404);
  });
});

describe('E2E: API Gateway', () => {
  it('should redirect root to /api', async () => {
    const response = await fetch(`${API_URL}/`, { redirect: 'manual' });
    expect(response.status).toBe(302);
    expect(response.headers.get('location')).toBe('/api');
  });

  it('should return API info', async () => {
    const response = await fetch(`${API_URL}/api`);
    expect(response.status).toBe(200);
    const data = await response.json();

    expect(data.name).toBe('Labhaus API');
    expect(data.version).toBeDefined();
    expect(data.endpoints).toBeDefined();
  });

  it('should return 404 for unknown route', async () => {
    const response = await fetch(`${API_URL}/api/unknown-route`);
    expect(response.status).toBe(404);
    const data = await response.json();
    expect(data.error).toBe('Not found');
  });

  it('should have security headers', async () => {
    const response = await fetch(`${API_URL}/api/health`);
    
    expect(response.headers.get('x-content-type-options')).toBe('nosniff');
    expect(response.headers.get('x-frame-options')).toBe('DENY');
    expect(response.headers.get('x-xss-protection')).toBe('1; mode=block');
  });
});
