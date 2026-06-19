import { Hono } from 'hono';
import { cors } from 'hono/cors';
import { serve } from '@hono/node-server';
import stylesRouter from './routes/styles.js';
import authRouter from './routes/auth.js';
import { errorHandler } from './middleware/error-handler.js';
import { requestLogger, securityHeaders, rateLimit } from './middleware/index.js';

const app = new Hono();

// Middleware
app.use('*', requestLogger);
app.use('*', securityHeaders);
app.use(
  '*',
  cors({
    origin: process.env.CORS_ORIGIN?.split(',') || ['http://localhost:3000'],
    credentials: true,
    allowMethods: ['GET', 'POST', 'PUT', 'DELETE', 'PATCH', 'OPTIONS'],
    allowHeaders: ['Content-Type', 'Authorization'],
  })
);

// Rate limiting (only in production)
if (process.env.NODE_ENV === 'production') {
  app.use('*', rateLimit(100, 60000)); // 100 requests per minute
}

// API version prefix
const api = new Hono();

// Health check
api.get('/health', (c) => c.json({ status: 'ok', timestamp: new Date().toISOString() }));

// API info
api.get('/', (c) =>
  c.json({
    name: 'Labhaus API',
    version: '0.1.0',
    endpoints: {
      health: '/api/health',
      styles: '/api/styles',
      auth: '/api/auth',
    },
  })
);

// Mount routes
api.route('/styles', stylesRouter);
api.route('/auth', authRouter);

// Mount API under /api prefix
app.route('/api', api);

// Root redirect
app.get('/', (c) => c.redirect('/api'));

// 404 handler
app.notFound((c) => c.json({ error: 'Not found', path: c.req.path }, 404));

// Global error handler
app.onError(errorHandler);

const port = parseInt(process.env.API_PORT || '3001', 10);
const host = process.env.API_HOST || '0.0.0.0';

console.log(`Starting Labhaus API server...`);
console.log(`Environment: ${process.env.NODE_ENV || 'development'}`);

serve({
  fetch: app.fetch,
  port,
  hostname: host,
});

console.log(`✓ Server running at http://${host}:${port}`);
console.log(`✓ API documentation: http://${host}:${port}/api`);
