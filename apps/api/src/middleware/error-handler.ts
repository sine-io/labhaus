import { Context } from 'hono';
import { z } from 'zod';

/**
 * Standard API error response
 */
export class ApiError extends Error {
  constructor(
    public statusCode: number,
    message: string,
    public code?: string,
    public details?: unknown
  ) {
    super(message);
    this.name = 'ApiError';
  }
}

/**
 * Common API errors
 */
export const ApiErrors = {
  badRequest: (message: string, details?: unknown) =>
    new ApiError(400, message, 'BAD_REQUEST', details),
  unauthorized: (message = 'Unauthorized') => new ApiError(401, message, 'UNAUTHORIZED'),
  forbidden: (message = 'Forbidden') => new ApiError(403, message, 'FORBIDDEN'),
  notFound: (resource: string) => new ApiError(404, `${resource} not found`, 'NOT_FOUND'),
  conflict: (message: string) => new ApiError(409, message, 'CONFLICT'),
  validationError: (errors: z.ZodError) =>
    new ApiError(400, 'Validation failed', 'VALIDATION_ERROR', errors.errors),
  internalError: (message = 'Internal server error') =>
    new ApiError(500, message, 'INTERNAL_ERROR'),
};

/**
 * Error handler middleware
 */
export const errorHandler = (err: Error, c: Context) => {
  console.error('Error:', err);

  if (err instanceof ApiError) {
    return c.json(
      {
        error: err.code || 'ERROR',
        message: err.message,
        details: err.details,
      },
      err.statusCode as 400 | 401 | 403 | 404 | 409 | 500
    );
  }

  // Zod validation errors
  if (err instanceof z.ZodError) {
    return c.json(
      {
        error: 'VALIDATION_ERROR',
        message: 'Validation failed',
        details: err.errors,
      },
      400
    );
  }

  // Default error
  const isDev = process.env.NODE_ENV !== 'production';
  return c.json(
    {
      error: 'INTERNAL_ERROR',
      message: isDev ? err.message : 'Internal server error',
      stack: isDev ? err.stack : undefined,
    },
    500
  );
};

/**
 * Async handler wrapper with error catching
 */
export const asyncHandler =
  (fn: (c: Context) => Promise<Response>) => async (c: Context) => {
    try {
      return await fn(c);
    } catch (error) {
      throw error;
    }
  };
