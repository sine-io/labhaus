import { Hono } from 'hono';
import { styleQuerySchema } from '@labhaus/types';
import { StyleRepository } from '../repositories/style.repository.js';
import { ApiErrors, asyncHandler } from '../middleware/error-handler.js';

const app = new Hono();
const styleRepo = new StyleRepository();

app.get(
  '/',
  asyncHandler(async (c) => {
    const query = styleQuerySchema.parse({
      category: c.req.query('category'),
      style: c.req.query('style'),
      scene: c.req.query('scene'),
      featured: c.req.query('featured') ? c.req.query('featured') === 'true' : undefined,
      search: c.req.query('search'),
      page: c.req.query('page') || '1',
      limit: c.req.query('limit') || '20',
    });

    const { styles, total } = await styleRepo.findAll(query);
    const totalPages = Math.ceil(total / query.limit);

    return c.json({
      styles,
      pagination: {
        page: query.page,
        limit: query.limit,
        total,
        totalPages,
      },
    });
  })
);

app.get(
  '/:id',
  asyncHandler(async (c) => {
    const id = c.req.param('id');
    
    if (!id) {
      throw ApiErrors.badRequest('Style ID is required');
    }
    
    const style = await styleRepo.findById(id);

    if (!style) {
      throw ApiErrors.notFound('Style');
    }

    return c.json(style);
  })
);

export default app;
