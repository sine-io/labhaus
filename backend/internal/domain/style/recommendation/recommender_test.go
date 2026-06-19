package recommendation

import (
	"context"
	"testing"

	"github.com/labhaus/backend/internal/domain/style"
)

func TestRecommender_Recommend(t *testing.T) {
	styles := sampleStyles()
	recommender := NewRecommender(styles)

	recommendations, err := recommender.Recommend(context.Background(), "modern clean user interface design", 3)
	if err != nil {
		t.Fatalf("Recommend() error = %v", err)
	}

	if len(recommendations) != 3 {
		t.Fatalf("Recommend() len = %d, want 3", len(recommendations))
	}
	if recommendations[0].Style.Category != "ui" {
		t.Fatalf("top recommendation category = %q, want ui", recommendations[0].Style.Category)
	}
	if recommendations[0].Score <= 0 {
		t.Fatalf("top recommendation score = %v, want > 0", recommendations[0].Score)
	}
	for i := 1; i < len(recommendations); i++ {
		if recommendations[i-1].Score < recommendations[i].Score {
			t.Fatalf("recommendations not sorted descending: %#v", recommendations)
		}
	}
}

func TestRecommender_AccuracyAbove70Percent(t *testing.T) {
	styles := sampleStyles()
	recommender := NewRecommender(styles)

	queries := []struct {
		query    string
		expected string
	}{
		{query: "modern clean user interface design", expected: "ui"},
		{query: "retro 80s style poster", expected: "retro"},
		{query: "beautiful nature scenery", expected: "nature"},
		{query: "dark cyberpunk neon city", expected: "cyberpunk"},
		{query: "watercolor portrait painting", expected: "art"},
		{query: "luxury premium gold branding", expected: "luxury"},
	}

	correct := 0
	for _, tt := range queries {
		recommendations, err := recommender.Recommend(context.Background(), tt.query, 3)
		if err != nil {
			t.Fatalf("Recommend(%q) error = %v", tt.query, err)
		}

		for _, rec := range recommendations {
			if rec.Style.Category == tt.expected {
				correct++
				break
			}
		}
	}

	accuracy := float64(correct) / float64(len(queries))
	t.Logf("accuracy %.2f%% (%d/%d)", accuracy*100, correct, len(queries))

	if accuracy <= 0.70 {
		t.Fatalf("accuracy = %.2f%%, want > 70%%", accuracy*100)
	}
}

func sampleStyles() []*style.Entity {
	return []*style.Entity{
		{ID: "1", Name: "Modern Minimalist", Category: "ui", Prompt: "clean simple modern interface app dashboard", Description: "minimal user interface design", Tags: []string{"minimal", "clean", "modern", "ui"}},
		{ID: "2", Name: "Material UI", Category: "ui", Prompt: "responsive mobile web components design system", Description: "polished user interface", Tags: []string{"interface", "web", "mobile"}},
		{ID: "3", Name: "Vintage Retro", Category: "retro", Prompt: "1980s retro neon vintage poster cassette", Description: "nostalgic old school graphic style", Tags: []string{"retro", "vintage", "80s"}},
		{ID: "4", Name: "Classic Film", Category: "retro", Prompt: "grainy analog film old photograph", Description: "timeless vintage camera look", Tags: []string{"film", "classic", "vintage"}},
		{ID: "5", Name: "Nature Landscape", Category: "nature", Prompt: "beautiful forest mountain lake sunrise", Description: "outdoor natural scenery", Tags: []string{"nature", "landscape", "outdoor"}},
		{ID: "6", Name: "Botanical Garden", Category: "nature", Prompt: "green plants flowers organic garden", Description: "fresh natural botanical scene", Tags: []string{"plants", "garden", "nature"}},
		{ID: "7", Name: "Cyberpunk City", Category: "cyberpunk", Prompt: "dark futuristic neon city rainy street", Description: "sci fi urban night atmosphere", Tags: []string{"cyberpunk", "neon", "future"}},
		{ID: "8", Name: "Holographic Tech", Category: "cyberpunk", Prompt: "glowing hologram technology matrix", Description: "futuristic digital cyber style", Tags: []string{"tech", "digital", "cyber"}},
		{ID: "9", Name: "Watercolor Portrait", Category: "art", Prompt: "soft watercolor portrait painting brush", Description: "hand painted artistic illustration", Tags: []string{"watercolor", "portrait", "painting"}},
		{ID: "10", Name: "Oil Painting", Category: "art", Prompt: "canvas oil paint classical brush strokes", Description: "traditional fine art look", Tags: []string{"oil", "canvas", "art"}},
		{ID: "11", Name: "Luxury Gold", Category: "luxury", Prompt: "premium gold elegant brand packaging", Description: "high end luxury visual identity", Tags: []string{"luxury", "gold", "premium"}},
		{ID: "12", Name: "Marble Editorial", Category: "luxury", Prompt: "elegant marble fashion editorial", Description: "refined upscale composition", Tags: []string{"elegant", "fashion", "upscale"}},
	}
}
