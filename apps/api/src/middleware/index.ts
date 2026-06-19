import { Context, Next } from 'hono';

/**
 * Request logging middleware
 */
export const requestLogger = async (c: Context, next: Next) => {
  const start = Date.now();
  const { method, url } = c.req;

  await next();

  const duration = Date.now() - start;
  const status = c.res.status;

  console.log(`[${method}] ${url} - ${status} (${duration}ms)`);
};

/**
 * Security headers middleware
 */
export const securityHeaders = async (c: Context, next: Next) => {
  await next();

  // Add security headers
  c.res.headers.set('X-Content-Type-Options', 'nosniff');
  c.res.headers.set('X-Frame-Options', 'DENY');
  c.res.headers.set('X-XSS-Protection', '1; mode=block');
  c.res.headers.set('Referrer-Policy', 'strict-origin-when-cross-origin');
};

/**
 * Rate limiting (simple in-memory implementation)
 */
const requestCounts = new Map<string, { count: number; resetAt: number }>();

export const rateLimit = (maxRequests = 100, windowMs = 60000) => {
  return async (c: Context, next: Next) => {
    const ip = c.req.header('x-forwarded-for') || c.req.header('x-real-ip') || 'unknown';
    const now = Date.now();

    let record = requestCounts.get(ip);

    // Reset if window expired
    if (!record || now > record.resetAt) {
      record = { count: 0, resetAt: now + windowMs };
      requestCounts.set(ip, record);
    }

    record.count++;

    // Check limit
    if (record.count > maxRequests) {
      return c.json(
        {
          error: 'RATE_LIMIT_EXCEEDED',
          message: 'Too many requests',
          retryAfter: Math.ceil((record.resetAt - now) / 1000),
        },
        429 as const
      );
    }

    // Add rate limit headers
    c.res.headers.set('X-RateLimit-Limit', maxRequests.toString());
    c.res.headers.set('X-RateLimit-Remaining', (maxRequests - record.count).toString());
    c.res.headers.set('X-RateLimit-Reset', record.resetAt.toString());

    return next();
  };
};
