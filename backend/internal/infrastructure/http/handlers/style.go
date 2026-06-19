package handlers

import (
	"context"
	"errors"
	"net/http"
	"sort"
	"strings"
	"unicode"

	"github.com/gin-gonic/gin"
	"github.com/labhaus/backend/internal/application/command"
	"github.com/labhaus/backend/internal/application/dto"
	"github.com/labhaus/backend/internal/application/query"
	"github.com/labhaus/backend/internal/domain/style"
)

// StyleHandler handles style-related HTTP requests
type StyleHandler struct {
	queryHandler   *query.StyleQueryHandler
	commandHandler *command.StyleCommandHandler
	recommender    StyleRecommender
}

// StyleRecommender recommends styles for a free-text query.
type StyleRecommender interface {
	Recommend(ctx context.Context, query string, topK int) ([]*StyleRecommendation, error)
}

// StyleRecommendation is a single style recommendation with a normalized score.
type StyleRecommendation struct {
	Style *style.Entity
	Score float64
}

// NewStyleHandler creates a new style handler
func NewStyleHandler(
	queryHandler *query.StyleQueryHandler,
	commandHandler *command.StyleCommandHandler,
	recommender StyleRecommender,
) *StyleHandler {
	return &StyleHandler{
		queryHandler:   queryHandler,
		commandHandler: commandHandler,
		recommender:    recommender,
	}
}

// ListStyles handles GET /api/styles
func (h *StyleHandler) ListStyles(c *gin.Context) {
	// Parse query parameters
	var filter dto.StyleFilterDTO
	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.queryHandler.ListStyles(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetStyle handles GET /api/styles/:id
func (h *StyleHandler) GetStyle(c *gin.Context) {
	id := c.Param("id")

	style, err := h.queryHandler.GetStyleByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "style not found"})
		return
	}

	c.JSON(http.StatusOK, style)
}

// CreateStyle handles POST /api/styles
func (h *StyleHandler) CreateStyle(c *gin.Context) {
	var req dto.CreateStyleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	style, err := h.commandHandler.CreateStyle(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, style)
}

// Recommend handles POST /api/styles/recommend
func (h *StyleHandler) Recommend(c *gin.Context) {
	var req dto.RecommendStyleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	limit := req.Limit
	if limit == 0 {
		limit = 10
	}

	if h.recommender == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "style recommender not configured"})
		return
	}

	recommendations, err := h.recommender.Recommend(c.Request.Context(), req.Query, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	results := make([]dto.RecommendedStyle, 0, len(recommendations))
	for _, rec := range recommendations {
		results = append(results, dto.RecommendedStyle{
			ID:          rec.Style.ID,
			Name:        rec.Style.Name,
			Prompt:      rec.Style.Prompt,
			Category:    rec.Style.Category,
			Description: rec.Style.Description,
			Tags:        rec.Style.Tags,
			Score:       rec.Score,
		})
	}

	c.JSON(http.StatusOK, dto.RecommendStyleResponse{
		Query:           req.Query,
		Recommendations: results,
		Total:           len(results),
	})
}

// StaticStyleRecommender ranks an immutable snapshot of styles by query token overlap.
type StaticStyleRecommender struct {
	styles []*style.Entity
}

// NewStyleRecommender creates a recommender from loaded style entities.
func NewStyleRecommender(styles []*style.Entity) *StaticStyleRecommender {
	copied := make([]*style.Entity, len(styles))
	copy(copied, styles)
	return &StaticStyleRecommender{styles: copied}
}

// Recommend returns up to topK styles sorted by normalized token-overlap score.
func (r *StaticStyleRecommender) Recommend(ctx context.Context, query string, topK int) ([]*StyleRecommendation, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if topK <= 0 {
		return nil, errors.New("topK must be greater than 0")
	}

	queryTerms := uniqueTerms(tokenize(query))
	recommendations := make([]*StyleRecommendation, 0, len(r.styles))
	for _, s := range r.styles {
		if s == nil {
			continue
		}
		recommendations = append(recommendations, &StyleRecommendation{
			Style: s,
			Score: scoreStyle(queryTerms, s),
		})
	}

	sort.SliceStable(recommendations, func(i, j int) bool {
		if recommendations[i].Score == recommendations[j].Score {
			return recommendations[i].Style.ID < recommendations[j].Style.ID
		}
		return recommendations[i].Score > recommendations[j].Score
	})

	if topK > len(recommendations) {
		topK = len(recommendations)
	}
	return recommendations[:topK], nil
}

func scoreStyle(queryTerms map[string]struct{}, s *style.Entity) float64 {
	if len(queryTerms) == 0 {
		return 0
	}

	styleTerms := uniqueTerms(tokenize(strings.Join([]string{
		s.Name,
		s.Description,
		s.Prompt,
		s.Category,
		strings.Join(s.Tags, " "),
	}, " ")))

	matches := 0
	for term := range queryTerms {
		if _, ok := styleTerms[term]; ok {
			matches++
		}
	}
	return float64(matches) / float64(len(queryTerms))
}

func tokenize(input string) []string {
	return strings.FieldsFunc(strings.ToLower(input), func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsNumber(r)
	})
}

func uniqueTerms(terms []string) map[string]struct{} {
	unique := make(map[string]struct{}, len(terms))
	for _, term := range terms {
		if term != "" {
			unique[term] = struct{}{}
		}
	}
	return unique
}
