import natural from 'natural';
import { Style } from '@labhaus/types';

const TfIdf = natural.TfIdf;
const tokenizer = new natural.WordTokenizer();

/**
 * Style recommendation service using TF-IDF and cosine similarity
 */
export class StyleRecommendationService {
  private tfidf: any;
  private styles: Style[] = [];
  private styleTexts: Map<string, string> = new Map();

  /**
   * Initialize the recommendation service with styles
   */
  async initialize(styles: Style[]) {
    this.styles = styles;
    this.tfidf = new TfIdf();

    // Build TF-IDF index
    for (const style of styles) {
      const text = this.extractText(style);
      this.styleTexts.set(style.id, text);
      this.tfidf.addDocument(text);
    }
  }

  /**
   * Extract searchable text from a style
   */
  private extractText(style: Style): string {
    const parts = [
      style.title,
      style.prompt,
      style.category,
      ...style.styles,
      ...style.scenes,
    ];
    return parts.filter(Boolean).join(' ').toLowerCase();
  }

  /**
   * Tokenize and normalize text
   */
  private tokenize(text: string): string[] {
    return tokenizer
      .tokenize(text.toLowerCase())
      .filter((token) => token.length > 2); // Remove very short tokens
  }

  /**
   * Calculate cosine similarity between two vectors
   */
  private cosineSimilarity(vec1: number[], vec2: number[]): number {
    let dotProduct = 0;
    let norm1 = 0;
    let norm2 = 0;

    for (let i = 0; i < vec1.length; i++) {
      dotProduct += vec1[i] * vec2[i];
      norm1 += vec1[i] * vec1[i];
      norm2 += vec2[i] * vec2[i];
    }

    if (norm1 === 0 || norm2 === 0) {
      return 0;
    }

    return dotProduct / (Math.sqrt(norm1) * Math.sqrt(norm2));
  }

  /**
   * Get TF-IDF vector for a text
   */
  private getTfIdfVector(documentIndex: number, terms: string[]): number[] {
    const vector: number[] = [];
    
    for (const term of terms) {
      // Get TF-IDF score for this term in the document
      const measures = this.tfidf.tfidf(term, documentIndex);
      vector.push(measures || 0);
    }

    return vector;
  }

  /**
   * Recommend styles based on query text
   */
  recommend(query: string, limit: number = 10): Array<{ style: Style; score: number }> {
    if (this.styles.length === 0) {
      return [];
    }

    // Tokenize query
    const queryTokens = this.tokenize(query);
    if (queryTokens.length === 0) {
      return [];
    }

    // Add query as a temporary document
    const queryIndex = this.styles.length;
    this.tfidf.addDocument(query);

    // Build vocabulary from all terms in query
    const vocabulary = queryTokens;

    // Get query TF-IDF vector
    const queryVector = this.getTfIdfVector(queryIndex, vocabulary);

    // Calculate similarity with each style
    const scores: Array<{ style: Style; score: number }> = [];

    for (let i = 0; i < this.styles.length; i++) {
      const style = this.styles[i];

      // Get style TF-IDF vector
      const styleVector = this.getTfIdfVector(i, vocabulary);

      // Calculate cosine similarity
      const similarity = this.cosineSimilarity(queryVector, styleVector);

      if (similarity > 0) {
        scores.push({
          style,
          score: Math.round(similarity * 1000) / 1000, // Round to 3 decimal places
        });
      }
    }

    // Remove the temporary query document
    this.tfidf.documents.splice(queryIndex, 1);

    // Sort by score descending and take top N
    scores.sort((a, b) => b.score - a.score);
    return scores.slice(0, limit);
  }

  /**
   * Recommend styles similar to a given style
   */
  recommendSimilar(styleId: string, limit: number = 10): Array<{ style: Style; score: number }> {
    const sourceStyle = this.styles.find((s) => s.id === styleId);
    if (!sourceStyle) {
      return [];
    }

    const sourceText = this.styleTexts.get(styleId) || '';
    return this.recommend(sourceText, limit + 1).filter((r) => r.style.id !== styleId);
  }
}
