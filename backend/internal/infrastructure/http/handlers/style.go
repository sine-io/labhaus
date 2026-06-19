package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/labhaus/backend/internal/application/command"
	"github.com/labhaus/backend/internal/application/dto"
	"github.com/labhaus/backend/internal/application/query"
)

// StyleHandler handles style-related HTTP requests
type StyleHandler struct {
	queryHandler   *query.StyleQueryHandler
	commandHandler *command.StyleCommandHandler
}

// NewStyleHandler creates a new style handler
func NewStyleHandler(
	queryHandler *query.StyleQueryHandler,
	commandHandler *command.StyleCommandHandler,
) *StyleHandler {
	return &StyleHandler{
		queryHandler:   queryHandler,
		commandHandler: commandHandler,
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
