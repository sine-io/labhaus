# Authentication API

## Endpoints

### POST /api/auth/register
Register a new user with email and password.

**Request:**
```json
{
  "email": "user@example.com",
  "password": "securepassword123",
  "name": "John Doe" // optional
}
```

**Response (201):**
```json
{
  "user": {
    "id": "uuid",
    "email": "user@example.com",
    "name": "John Doe",
    "avatar_url": null,
    "email_verified": false,
    "created_at": "2026-06-19T...",
    "updated_at": "2026-06-19T..."
  },
  "tokens": {
    "access_token": "eyJhbGc...",
    "refresh_token": "eyJhbGc...",
    "token_type": "Bearer",
    "expires_in": 3600
  }
}
```

### POST /api/auth/login
Login with email and password.

**Request:**
```json
{
  "email": "user@example.com",
  "password": "securepassword123"
}
```

**Response (200):**
```json
{
  "user": { ... },
  "tokens": { ... }
}
```

### POST /api/auth/refresh
Refresh access token using refresh token.

**Request:**
```json
{
  "refresh_token": "eyJhbGc..."
}
```

**Response (200):**
```json
{
  "tokens": {
    "access_token": "eyJhbGc...",
    "refresh_token": "eyJhbGc...",
    "token_type": "Bearer",
    "expires_in": 3600
  }
}
```

### GET /api/auth/me
Get current authenticated user.

**Headers:**
```
Authorization: Bearer <access_token>
```

**Response (200):**
```json
{
  "user": {
    "id": "uuid",
    "email": "user@example.com",
    "name": "John Doe",
    ...
  }
}
```

### POST /api/auth/google
Google OAuth login/register.

**Status:** Not fully implemented (placeholder)

**Request:**
```json
{
  "id_token": "google_id_token"
}
```

## Authentication Flow

1. **Register or Login** → Receive `access_token` and `refresh_token`
2. **Use access_token** in `Authorization: Bearer <token>` header for protected endpoints
3. **When access_token expires** → Use `refresh_token` to get new tokens
4. **Logout** → Client discards tokens (server-side revocation not yet implemented)

## Security Features

- Passwords hashed with bcrypt (10 rounds)
- JWT tokens signed with secret keys
- Access tokens expire in 1 hour (configurable)
- Refresh tokens expire in 7 days (configurable)
- Refresh tokens stored in database for validation

## Environment Variables

```bash
JWT_SECRET=your-secret-key-change-in-production
JWT_REFRESH_SECRET=your-refresh-secret-key-change-in-production
JWT_EXPIRES_IN=1h
JWT_REFRESH_EXPIRES_IN=7d
```

## Protected Endpoints

Use `requireAuth` middleware to protect endpoints:

```typescript
import { requireAuth } from '../middleware/auth.js';

app.get('/protected', requireAuth, async (c) => {
  const user = c.get('user'); // JwtPayload
  // ...
});
```
