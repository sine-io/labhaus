package handlers

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/labhaus/backend/internal/application/dto"
	imageapp "github.com/labhaus/backend/internal/application/service/image"
	imageprovider "github.com/labhaus/backend/internal/domain/image/provider"
	"github.com/labhaus/backend/internal/infrastructure/storage"
)

const imageURLExpiry = 24 * time.Hour

type ImageHandler struct {
	batchService *imageapp.BatchImageService
	storage      *storage.MinIOImageStorage
}

func NewImageHandler(batchService *imageapp.BatchImageService, storage *storage.MinIOImageStorage) *ImageHandler {
	return &ImageHandler{
		batchService: batchService,
		storage:      storage,
	}
}

func (h *ImageHandler) Generate(c *gin.Context) {
	var req dto.GenerateImageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if len(req.Prompts) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "prompts cannot be empty"})
		return
	}

	results, err := h.batchService.GenerateBatch(c.Request.Context(), req.Prompts, imageprovider.ImageOptions{
		Width:   req.Width,
		Height:  req.Height,
		Quality: imageprovider.Quality(req.Quality),
		Style:   req.Style,
	})
	if err != nil && len(results) == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	responseResults := make([]dto.ImageResult, 0, len(results))
	for _, result := range results {
		if result == nil {
			continue
		}

		id := uuid.New().String() + ".png"
		data, contentType, err := imageData(c.Request.Context(), result)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if err := h.storage.Upload(c.Request.Context(), id, data, contentType); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		url, err := h.storage.GetPresignedURL(c.Request.Context(), id, imageURLExpiry)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		responseResults = append(responseResults, dto.ImageResult{
			ID:        id,
			URL:       url,
			Prompt:    result.Metadata.Prompt,
			CreatedAt: imageCreatedAt(result.Metadata.Timestamp),
		})
	}

	c.JSON(http.StatusOK, dto.GenerateImageResponse{
		Results: responseResults,
		Total:   len(req.Prompts),
		Success: len(responseResults),
		Failed:  len(req.Prompts) - len(responseResults),
	})
}

func (h *ImageHandler) GetImage(c *gin.Context) {
	id := imageIDParam(c)
	if id == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "image not found"})
		return
	}

	exists, err := h.storage.Exists(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "image not found"})
		return
	}

	url, err := h.storage.GetPresignedURL(c.Request.Context(), id, imageURLExpiry)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":  id,
		"url": url,
	})
}

func (h *ImageHandler) GetProgress(c *gin.Context) {
	id := imageIDParam(c)
	c.JSON(http.StatusOK, dto.ImageProgressResponse{
		ID:       id,
		Status:   "completed",
		Progress: 100,
	})
}

func imageIDParam(c *gin.Context) string {
	if id := c.Param("id"); id != "" {
		return id
	}

	path := strings.TrimPrefix(c.Request.URL.Path, "/api/images/")
	path = strings.TrimSuffix(path, "/progress")
	if path == c.Request.URL.Path {
		return ""
	}
	return path
}

func imageData(ctx context.Context, result *imageprovider.ImageResult) ([]byte, string, error) {
	if len(result.Buffer) > 0 {
		return result.Buffer, "image/png", nil
	}
	if strings.TrimSpace(result.URL) == "" {
		return nil, "", storage.ErrEmptyData
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, result.URL, nil)
	if err != nil {
		return nil, "", fmt.Errorf("create image download request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, "", fmt.Errorf("download generated image: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("download generated image: status %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", fmt.Errorf("read generated image: %w", err)
	}
	if len(data) == 0 {
		return nil, "", storage.ErrEmptyData
	}

	contentType := resp.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "image/png"
	}
	return data, contentType, nil
}

func imageCreatedAt(timestamp string) time.Time {
	if timestamp == "" {
		return time.Now().UTC()
	}

	createdAt, err := time.Parse(time.RFC3339Nano, timestamp)
	if err == nil {
		return createdAt
	}

	createdAt, err = time.Parse(time.RFC3339, timestamp)
	if err == nil {
		return createdAt
	}

	return time.Now().UTC()
}
