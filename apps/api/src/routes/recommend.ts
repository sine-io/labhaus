import { Hono } from 'hono';
import { z } from 'zod';
import { StyleRepository } from '../repositories/style.repository.js';
import { StyleRecommendationService } from '../services/recommendation.service.js';
import { ApiErrors, asyncHandler } from '../middleware/error-handler.js';

const app = new Hono();
const styleRepo = new StyleRepository();
const recommendationService = new StyleRecommendationService();

// Initialize recommendation service on startup
let initialized = false;

const recommendSchema = z.object({
  query: z.string().min(1, 'Query is required'),
  limit: z.coerce.number().int().positive().max(50).default(10),
});

const similarSchema = z.object({
  style_id: z.string().uuid('Invalid style ID'),
  limit: z.coerce.number().int().positive().max(50).default(10),
});

/**
 * POST /api/styles/recommend
 * Recommend styles based on query text
 */
app.post(
  '/recommend',
  asyncHandler(async (c) => {
    const body = await c.req.json();
    const data = recommendSchema.parse(body);

    // Initialize service if not already done
    if (!initialized) {
      const { styles } = await styleRepo.findAll({ page: 1, limit: 1000 });
      await recommendationService.initialize(styles);
      initialized = true;
    }

    const recommendations = recommendationService.recommend(data.query, data.limit);

    return c.json({
      query: data.query,
      recommendations,
      total: recommendations.length,
    });
  })
);

/**
 * GET /api/styles/:id/similar
 * Find similar styles to a given style
 */
app.get(
  '/:id/similar',
  asyncHandler(async (c) => {
    const styleId = c.req.param('id');
    const limit = parseInt(c.req.query('limit') || '10', 10);

    if (!styleId) {
      throw ApiErrors.badRequest('Style ID is required');
    }

    similarSchema.parse({ style_id: styleId, limit });

    // Check if style exists
    const style = await styleRepo.findById(styleId);
    if (!style) {
      throw ApiErrors.notFound('Style');
    }

    // Initialize service if not already done
    if (!initialized) {
      const { styles } = await styleRepo.findAll({ page: 1, limit: 1000 });
      await recommendationService.initialize(styles);
      initialized = true;
    }

    const recommendations = recommendationService.recommendSimilar(styleId, limit);

    return c.json({
      style_id: styleId,
      recommendations,
      total: recommendations.length,
    });
  })
);

export default app;
