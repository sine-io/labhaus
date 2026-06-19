import { describe, it, expect, beforeAll } from 'vitest';
import { StyleRecommendationService } from '../src/services/recommendation.service';
import { Style } from '@labhaus/types';

describe('StyleRecommendationService', () => {
  let service: StyleRecommendationService;
  let mockStyles: Style[];

  beforeAll(() => {
    service = new StyleRecommendationService();

    // Create mock styles
    mockStyles = [
      {
        id: '1',
        case_id: 1,
        title: 'Modern Minimalist UI',
        prompt: 'Clean modern interface with minimalist design and simple navigation',
        prompt_preview: 'Clean modern interface...',
        category: 'UI & Interfaces',
        styles: ['UI', 'Minimalist'],
        scenes: ['Tech'],
        image_url: null,
        source_label: null,
        source_url: null,
        github_url: null,
        featured: false,
        created_at: new Date().toISOString(),
        updated_at: new Date().toISOString(),
      },
      {
        id: '2',
        case_id: 2,
        title: 'Dark Modern Dashboard',
        prompt: 'Professional dark-themed dashboard with modern charts and data visualization',
        prompt_preview: 'Professional dark-themed...',
        category: 'UI & Interfaces',
        styles: ['UI', 'Dashboard'],
        scenes: ['Tech'],
        image_url: null,
        source_label: null,
        source_url: null,
        github_url: null,
        featured: false,
        created_at: new Date().toISOString(),
        updated_at: new Date().toISOString(),
      },
      {
        id: '3',
        case_id: 3,
        title: 'Watercolor Landscape',
        prompt: 'Beautiful watercolor painting of mountains and lakes',
        prompt_preview: 'Beautiful watercolor...',
        category: 'Illustration & Art',
        styles: ['Illustration', 'Watercolor'],
        scenes: ['Travel'],
        image_url: null,
        source_label: null,
        source_url: null,
        github_url: null,
        featured: false,
        created_at: new Date().toISOString(),
        updated_at: new Date().toISOString(),
      },
    ];
  });

  it('should initialize with styles', async () => {
    await service.initialize(mockStyles);
    const results = service.recommend('test', 10);
    expect(results).toBeDefined();
  });

  it('should recommend relevant styles for UI query', async () => {
    await service.initialize(mockStyles);
    const results = service.recommend('modern user interface design', 10);

    expect(results.length).toBeGreaterThan(0);
    expect(results[0].score).toBeGreaterThan(0);
    
    // First result should be UI-related
    expect(['UI & Interfaces']).toContain(results[0].style.category);
  });

  it('should recommend relevant styles for art query', async () => {
    await service.initialize(mockStyles);
    const results = service.recommend('watercolor painting landscape', 10);

    expect(results.length).toBeGreaterThan(0);
    // First result should be art-related
    const firstResult = results[0];
    expect(firstResult.style.category).toContain('Art');
  });

  it('should return empty array for empty query', async () => {
    await service.initialize(mockStyles);
    const results = service.recommend('', 10);
    expect(results).toEqual([]);
  });

  it('should limit results to specified number', async () => {
    await service.initialize(mockStyles);
    const results = service.recommend('design', 2);
    expect(results.length).toBeLessThanOrEqual(2);
  });

  it('should find similar styles', async () => {
    await service.initialize(mockStyles);
    const results = service.recommendSimilar('1', 10);

    // Should return similar styles (excluding the source style itself)
    expect(results.every((r) => r.style.id !== '1')).toBe(true);
    
    if (results.length > 0) {
      expect(results[0].score).toBeGreaterThan(0);
    }
  });

  it('should return empty array for non-existent style ID', async () => {
    await service.initialize(mockStyles);
    const results = service.recommendSimilar('non-existent', 10);
    expect(results).toEqual([]);
  });

  it('should handle styles with no data', async () => {
    await service.initialize([]);
    const results = service.recommend('test', 10);
    expect(results).toEqual([]);
  });
});
