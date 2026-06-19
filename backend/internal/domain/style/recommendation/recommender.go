package recommendation

import (
	"context"
	"sort"
	"strings"

	"github.com/labhaus/backend/internal/domain/style"
)

// Recommender provides style recommendation functionality.
type Recommender struct {
	calculator   *TFIDFCalculator
	styles       []*style.Entity
	styleVectors map[string]map[string]float64
}

// Recommendation represents a recommended style with score.
type Recommendation struct {
	Style *style.Entity
	Score float64
}

// NewRecommender creates a new recommender.
func NewRecommender(styles []*style.Entity) *Recommender {
	documents := make([]*Document, 0, len(styles))
	for _, s := range styles {
		documents = append(documents, styleToDocument(s))
	}

	r := &Recommender{
		calculator:   NewTFIDFCalculator(documents),
		styles:       styles,
		styleVectors: make(map[string]map[string]float64, len(styles)),
	}
	r.buildStyleVectors()

	return r
}

// Recommend returns top-k most similar styles for a query.
func (r *Recommender) Recommend(ctx context.Context, query string, topK int) ([]*Recommendation, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if topK <= 0 || len(r.styles) == 0 {
		return []*Recommendation{}, nil
	}

	queryDoc := &Document{ID: "query", Content: query}
	queryVector := r.calculator.Calculate(queryDoc)

	recommendations := make([]*Recommendation, 0, len(r.styles))
	for _, s := range r.styles {
		if err := ctx.Err(); err != nil {
			return nil, err
		}

		recommendations = append(recommendations, &Recommendation{
			Style: s,
			Score: CosineSimilarity(queryVector, r.styleVectors[s.ID]),
		})
	}

	sort.SliceStable(recommendations, func(i, j int) bool {
		return recommendations[i].Score > recommendations[j].Score
	})

	if topK > len(recommendations) {
		topK = len(recommendations)
	}

	return recommendations[:topK], nil
}

// buildStyleVectors precomputes TF-IDF vectors for all styles.
func (r *Recommender) buildStyleVectors() {
	for _, s := range r.styles {
		r.styleVectors[s.ID] = r.calculator.Calculate(styleToDocument(s))
	}
}

// styleToDocument converts a style entity to a document.
func styleToDocument(s *style.Entity) *Document {
	contentParts := []string{s.Name, s.Prompt, s.Description, strings.Join(s.Tags, " ")}
	content := strings.Join(contentParts, " ")

	return &Document{
		ID:      s.ID,
		Content: content,
		Tokens:  Tokenize(content),
	}
}
