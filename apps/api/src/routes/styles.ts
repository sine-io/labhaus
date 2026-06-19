import { Hono } from 'hono';
import { styleQuerySchema } from '@labhaus/types';
import { StyleRepository } from '../repositories/style.repository.js';

const app = new Hono();
const styleRepo = new StyleRepository();

app.get('/', async (c) => {
  try {
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
  } catch (error) {
    console.error('Error fetching styles:', error);
    if (error instanceof Error && error.name === 'ZodError') {
      return c.json({ error: 'Invalid query parameters', details: error }, 400);
    }
    return c.json({ error: 'Internal server error' }, 500);
  }
});

app.get('/:id', async (c) => {
  try {
    const id = c.req.param('id');
    const style = await styleRepo.findById(id);

    if (!style) {
      return c.json({ error: 'Style not found' }, 404);
    }

    return c.json(style);
  } catch (error) {
    console.error('Error fetching style:', error);
    return c.json({ error: 'Internal server error' }, 500);
  }
});

export default app;
