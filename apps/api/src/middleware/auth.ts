import { Context, Next } from 'hono';
import { verifyAccessToken } from '../lib/jwt.js';
import { ApiErrors } from './error-handler.js';

/**
 * JWT authentication middleware
 */
export const requireAuth = async (c: Context, next: Next) => {
  const authHeader = c.req.header('Authorization');

  if (!authHeader || !authHeader.startsWith('Bearer ')) {
    throw ApiErrors.unauthorized('Missing or invalid authorization header');
  }

  const token = authHeader.substring(7);

  try {
    const payload = verifyAccessToken(token);
    c.set('user', payload);
    return next();
  } catch (error) {
    throw ApiErrors.unauthorized('Invalid or expired token');
  }
};

/**
 * Optional authentication (doesn't throw if no token)
 */
export const optionalAuth = async (c: Context, next: Next) => {
  const authHeader = c.req.header('Authorization');

  if (authHeader && authHeader.startsWith('Bearer ')) {
    const token = authHeader.substring(7);
    try {
      const payload = verifyAccessToken(token);
      c.set('user', payload);
    } catch (error) {
      // Ignore invalid tokens in optional auth
    }
  }

  return next();
};
