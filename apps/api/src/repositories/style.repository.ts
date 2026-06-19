import { Style, StyleQuery } from '@labhaus/types';
import { query } from '../db.js';

export class StyleRepository {
  async findAll(filters: StyleQuery): Promise<{ styles: Style[]; total: number }> {
    const { category, style, scene, featured, search, page, limit } = filters;
    const offset = (page - 1) * limit;

    const conditions: string[] = [];
    const params: any[] = [];
    let paramIndex = 1;

    if (category) {
      conditions.push(`category = $${paramIndex++}`);
      params.push(category);
    }

    if (style) {
      conditions.push(`$${paramIndex} = ANY(styles)`);
      params.push(style);
      paramIndex++;
    }

    if (scene) {
      conditions.push(`$${paramIndex} = ANY(scenes)`);
      params.push(scene);
      paramIndex++;
    }

    if (featured !== undefined) {
      conditions.push(`featured = $${paramIndex++}`);
      params.push(featured);
    }

    if (search) {
      conditions.push(
        `(to_tsvector('english', title) @@ plainto_tsquery('english', $${paramIndex}) OR to_tsvector('english', prompt) @@ plainto_tsquery('english', $${paramIndex}))`
      );
      params.push(search);
      paramIndex++;
    }

    const whereClause = conditions.length > 0 ? `WHERE ${conditions.join(' AND ')}` : '';

    // Get total count
    const countQuery = `SELECT COUNT(*) as total FROM styles ${whereClause}`;
    const countResult = await query(countQuery, params);
    const total = parseInt(countResult.rows[0].total, 10);

    // Get paginated results
    params.push(limit, offset);
    const dataQuery = `
      SELECT * FROM styles
      ${whereClause}
      ORDER BY featured DESC, created_at DESC
      LIMIT $${paramIndex++} OFFSET $${paramIndex}
    `;

    const result = await query(dataQuery, params);

    return {
      styles: result.rows as Style[],
      total,
    };
  }

  async findById(id: string): Promise<Style | null> {
    const result = await query('SELECT * FROM styles WHERE id = $1', [id]);
    return result.rows[0] || null;
  }

  async findByCaseId(caseId: number): Promise<Style | null> {
    const result = await query('SELECT * FROM styles WHERE case_id = $1', [caseId]);
    return result.rows[0] || null;
  }

  async create(style: Omit<Style, 'id' | 'created_at' | 'updated_at'>): Promise<Style> {
    const result = await query(
      `INSERT INTO styles (
        case_id, title, prompt, prompt_preview, category,
        styles, scenes, image_url, source_label, source_url,
        github_url, featured
      ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
      RETURNING *`,
      [
        style.case_id,
        style.title,
        style.prompt,
        style.prompt_preview,
        style.category,
        style.styles,
        style.scenes,
        style.image_url,
        style.source_label,
        style.source_url,
        style.github_url,
        style.featured,
      ]
    );
    return result.rows[0] as Style;
  }

  async batchCreate(styles: Omit<Style, 'id' | 'created_at' | 'updated_at'>[]): Promise<number> {
    if (styles.length === 0) return 0;

    const values: string[] = [];
    const params: any[] = [];
    let paramIndex = 1;

    for (const style of styles) {
      values.push(
        `($${paramIndex}, $${paramIndex + 1}, $${paramIndex + 2}, $${paramIndex + 3}, $${paramIndex + 4}, $${paramIndex + 5}, $${paramIndex + 6}, $${paramIndex + 7}, $${paramIndex + 8}, $${paramIndex + 9}, $${paramIndex + 10}, $${paramIndex + 11})`
      );
      params.push(
        style.case_id,
        style.title,
        style.prompt,
        style.prompt_preview,
        style.category,
        style.styles,
        style.scenes,
        style.image_url,
        style.source_label,
        style.source_url,
        style.github_url,
        style.featured
      );
      paramIndex += 12;
    }

    const sql = `
      INSERT INTO styles (
        case_id, title, prompt, prompt_preview, category,
        styles, scenes, image_url, source_label, source_url,
        github_url, featured
      ) VALUES ${values.join(', ')}
      ON CONFLICT (case_id) DO NOTHING
    `;

    const result = await query(sql, params);
    return result.rowCount || 0;
  }
}
