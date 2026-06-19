package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/labhaus/backend/internal/application/dto"
	"github.com/labhaus/backend/internal/domain/style"
	"github.com/stretchr/testify/require"
)

func TestStyleHandler_Recommend(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		requestBody    any
		rawBody        string
		expectedStatus int
		validateResp   func(*testing.T, *dto.RecommendStyleResponse)
	}{
		{
			name: "valid request returns recommendations",
			requestBody: dto.RecommendStyleRequest{
				Query: "modern clean interface",
				Limit: 5,
			},
			expectedStatus: http.StatusOK,
			validateResp: func(t *testing.T, resp *dto.RecommendStyleResponse) {
				require.Equal(t, "modern clean interface", resp.Query)
				require.NotZero(t, resp.Total)
				require.Len(t, resp.Recommendations, resp.Total)
				require.LessOrEqual(t, len(resp.Recommendations), 5)
				for _, rec := range resp.Recommendations {
					require.NotEmpty(t, rec.ID)
					require.NotEmpty(t, rec.Name)
					require.NotEmpty(t, rec.Prompt)
					require.GreaterOrEqual(t, rec.Score, 0.0)
					require.LessOrEqual(t, rec.Score, 1.0)
				}
			},
		},
		{
			name: "empty query returns 400",
			requestBody: dto.RecommendStyleRequest{
				Query: "",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "invalid JSON returns 400",
			rawBody:        `{"query":`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "returns top-k recommendations with scores",
			requestBody: dto.RecommendStyleRequest{
				Query: "warm vintage interior",
				Limit: 2,
			},
			expectedStatus: http.StatusOK,
			validateResp: func(t *testing.T, resp *dto.RecommendStyleResponse) {
				require.Equal(t, 2, resp.Total)
				require.Len(t, resp.Recommendations, 2)
				require.GreaterOrEqual(t, resp.Recommendations[0].Score, resp.Recommendations[1].Score)
			},
		},
		{
			name: "default limit is 10",
			requestBody: dto.RecommendStyleRequest{
				Query: "style",
			},
			expectedStatus: http.StatusOK,
			validateResp: func(t *testing.T, resp *dto.RecommendStyleResponse) {
				require.Equal(t, 10, resp.Total)
				require.Len(t, resp.Recommendations, 10)
			},
		},
		{
			name: "custom limit works",
			requestBody: dto.RecommendStyleRequest{
				Query: "style",
				Limit: 3,
			},
			expectedStatus: http.StatusOK,
			validateResp: func(t *testing.T, resp *dto.RecommendStyleResponse) {
				require.Equal(t, 3, resp.Total)
				require.Len(t, resp.Recommendations, 3)
			},
		},
		{
			name: "limit above 50 returns 400",
			requestBody: dto.RecommendStyleRequest{
				Query: "style",
				Limit: 51,
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "limit below 1 returns 400",
			requestBody: dto.RecommendStyleRequest{
				Query: "style",
				Limit: -1,
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := newTestStyleHandler()
			body := []byte(tt.rawBody)
			if tt.rawBody == "" {
				var err error
				body, err = json.Marshal(tt.requestBody)
				require.NoError(t, err)
			}

			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/api/styles/recommend", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			c, _ := gin.CreateTestContext(rec)
			c.Request = req

			handler.Recommend(c)

			require.Equal(t, tt.expectedStatus, rec.Code)
			if tt.validateResp != nil {
				var resp dto.RecommendStyleResponse
				require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
				tt.validateResp(t, &resp)
			}
		})
	}
}

func TestStyleHandler_Recommend_RecommenderError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := NewStyleHandler(nil, nil, errorStyleRecommender{})
	body, err := json.Marshal(dto.RecommendStyleRequest{Query: "modern"})
	require.NoError(t, err)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/styles/recommend", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	c, _ := gin.CreateTestContext(rec)
	c.Request = req

	handler.Recommend(c)

	require.Equal(t, http.StatusInternalServerError, rec.Code)
}

type errorStyleRecommender struct{}

func (errorStyleRecommender) Recommend(context.Context, string, int) ([]*StyleRecommendation, error) {
	return nil, errors.New("recommendation failed")
}

func newTestStyleHandler() *StyleHandler {
	return NewStyleHandler(nil, nil, NewStyleRecommender(testStyles()))
}

func testStyles() []*style.Entity {
	return []*style.Entity{
		{
			ID:          "style-01",
			Name:        "Modern Clean",
			Prompt:      "modern clean interface with crisp typography and balanced spacing",
			Category:    "ui",
			Description: "Clean modern interface style",
			Tags:        []string{"modern", "clean", "interface"},
		},
		{
			ID:          "style-02",
			Name:        "Warm Vintage",
			Prompt:      "warm vintage interior with aged textures and cozy lighting",
			Category:    "interior",
			Description: "Warm nostalgic interior style",
			Tags:        []string{"warm", "vintage", "interior"},
		},
		{
			ID:          "style-03",
			Name:        "Minimal Product",
			Prompt:      "minimal product photography on neutral background",
			Category:    "product",
			Description: "Minimal product visual style",
			Tags:        []string{"minimal", "product"},
		},
		{
			ID:          "style-04",
			Name:        "Editorial Fashion",
			Prompt:      "editorial fashion shoot with dramatic composition",
			Category:    "fashion",
			Description: "High fashion editorial style",
			Tags:        []string{"editorial", "fashion"},
		},
		{
			ID:          "style-05",
			Name:        "Cinematic Night",
			Prompt:      "cinematic night scene with neon highlights",
			Category:    "cinematic",
			Description: "Cinematic night color grade",
			Tags:        []string{"cinematic", "night"},
		},
		{
			ID:          "style-06",
			Name:        "Soft Pastel",
			Prompt:      "soft pastel palette with gentle gradients",
			Category:    "illustration",
			Description: "Soft pastel illustration style",
			Tags:        []string{"soft", "pastel"},
		},
		{
			ID:          "style-07",
			Name:        "Brutalist Web",
			Prompt:      "brutalist web design with bold grid and stark contrast",
			Category:    "ui",
			Description: "Bold brutalist web style",
			Tags:        []string{"brutalist", "web"},
		},
		{
			ID:          "style-08",
			Name:        "Natural Lifestyle",
			Prompt:      "natural lifestyle photography with candid composition",
			Category:    "photo",
			Description: "Natural lifestyle photo style",
			Tags:        []string{"natural", "lifestyle"},
		},
		{
			ID:          "style-09",
			Name:        "Technical Blueprint",
			Prompt:      "technical blueprint diagram with precise linework",
			Category:    "diagram",
			Description: "Blueprint technical drawing style",
			Tags:        []string{"technical", "blueprint"},
		},
		{
			ID:          "style-10",
			Name:        "Dreamy Fantasy",
			Prompt:      "dreamy fantasy landscape with ethereal atmosphere",
			Category:    "fantasy",
			Description: "Dreamy fantasy landscape style",
			Tags:        []string{"dreamy", "fantasy"},
		},
		{
			ID:          "style-11",
			Name:        "Monochrome Portrait",
			Prompt:      "monochrome portrait with strong shadows",
			Category:    "portrait",
			Description: "High contrast monochrome portrait",
			Tags:        []string{"monochrome", "portrait"},
		},
	}
}
