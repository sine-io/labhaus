import { Hono } from 'hono';
import {
  registerSchema,
  loginSchema,
  refreshTokenSchema,
  googleAuthSchema,
} from '@labhaus/types';
import { AuthService } from '../services/auth.service.js';
import { ApiErrors, asyncHandler } from '../middleware/error-handler.js';
import { requireAuth } from '../middleware/auth.js';

const app = new Hono();
const authService = new AuthService();

/**
 * POST /api/auth/register
 * Register a new user
 */
app.post(
  '/register',
  asyncHandler(async (c) => {
    const body = await c.req.json();
    const data = registerSchema.parse(body);

    try {
      const { user, tokens } = await authService.register(data);

      return c.json(
        {
          user: authService.sanitizeUser(user),
          tokens,
        },
        201
      );
    } catch (error) {
      if (error instanceof Error && error.message === 'User already exists') {
        throw ApiErrors.conflict('Email already registered');
      }
      throw error;
    }
  })
);

/**
 * POST /api/auth/login
 * Login with email and password
 */
app.post(
  '/login',
  asyncHandler(async (c) => {
    const body = await c.req.json();
    const data = loginSchema.parse(body);

    try {
      const { user, tokens } = await authService.login(data);

      return c.json({
        user: authService.sanitizeUser(user),
        tokens,
      });
    } catch (error) {
      if (error instanceof Error && error.message === 'Invalid credentials') {
        throw ApiErrors.unauthorized('Invalid email or password');
      }
      throw error;
    }
  })
);

/**
 * POST /api/auth/refresh
 * Refresh access token
 */
app.post(
  '/refresh',
  asyncHandler(async (c) => {
    const body = await c.req.json();
    const data = refreshTokenSchema.parse(body);

    try {
      const tokens = await authService.refreshToken(data.refresh_token);

      return c.json({ tokens });
    } catch (error) {
      throw ApiErrors.unauthorized('Invalid or expired refresh token');
    }
  })
);

/**
 * GET /api/auth/me
 * Get current user
 */
app.get(
  '/me',
  requireAuth,
  asyncHandler(async (c) => {
    const jwtPayload = c.get('user');
    const user = await authService.getUserById(jwtPayload.user_id);

    if (!user) {
      throw ApiErrors.notFound('User');
    }

    return c.json({
      user: authService.sanitizeUser(user),
    });
  })
);

/**
 * POST /api/auth/google
 * Google OAuth login/register
 * 
 * Note: This is a simplified implementation.
 * In production, verify the Google ID token with Google's API.
 */
app.post(
  '/google',
  asyncHandler(async (c) => {
    const body = await c.req.json();
    googleAuthSchema.parse(body);

    // TODO: Verify Google ID token with Google's API
    // For now, we accept the token as-is (NOT PRODUCTION READY)
    
    // Mock: extract user info from token (in production, verify with Google)
    // This is a placeholder - real implementation would call Google's tokeninfo endpoint
    throw ApiErrors.internalError('Google OAuth not fully implemented yet');

    // Example of what the real implementation would look like:
    // const googleUser = await verifyGoogleToken(data.id_token);
    // const { user, tokens } = await authService.googleAuth(
    //   googleUser.sub,
    //   googleUser.email,
    //   googleUser.name,
    //   googleUser.picture
    // );
    // return c.json({ user: authService.sanitizeUser(user), tokens });
  })
);

export default app;
