package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/labhaus/backend/internal/application/command"
	"github.com/labhaus/backend/internal/application/dto"
	"github.com/labhaus/backend/internal/application/query"
	"github.com/labhaus/backend/internal/infrastructure/http/middleware"
)

// WorkflowHandler handles workflow-related HTTP requests
type WorkflowHandler struct {
	queryHandler   *query.WorkflowQueryHandler
	commandHandler *command.WorkflowCommandHandler
}

// NewWorkflowHandler creates a new workflow handler
func NewWorkflowHandler(
	queryHandler *query.WorkflowQueryHandler,
	commandHandler *command.WorkflowCommandHandler,
) *WorkflowHandler {
	return &WorkflowHandler{
		queryHandler:   queryHandler,
		commandHandler: commandHandler,
	}
}

// CreateWorkflow handles POST /api/workflows
func (h *WorkflowHandler) CreateWorkflow(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req dto.CreateWorkflowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	workflow, err := h.commandHandler.CreateWorkflow(c.Request.Context(), userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, workflow)
}

// ListWorkflows handles GET /api/workflows
func (h *WorkflowHandler) ListWorkflows(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Parse query parameters
	var filter dto.WorkflowFilterDTO
	filter.Limit = 20
	if c.Query("limit") != "" {
		if _, err := c.GetQuery("limit"); err {
			if limit, ok := c.GetQuery("limit"); ok {
				filter.Limit = parseInt(limit, 20)
			}
		}
	}
	if c.Query("offset") != "" {
		if offset, ok := c.GetQuery("offset"); ok {
			filter.Offset = parseInt(offset, 0)
		}
	}
	if state := c.Query("state"); state != "" {
		filter.State = &state
	}

	response, err := h.queryHandler.ListUserWorkflows(c.Request.Context(), userID, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetWorkflow handles GET /api/workflows/:id
func (h *WorkflowHandler) GetWorkflow(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	workflowID := c.Param("id")
	workflow, err := h.queryHandler.GetWorkflowByID(c.Request.Context(), workflowID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "workflow not found"})
		return
	}

	// Check ownership
	if workflow.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	c.JSON(http.StatusOK, workflow)
}

// UpdateWorkflowStatus handles PATCH /api/workflows/:id/status
func (h *WorkflowHandler) UpdateWorkflowStatus(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	workflowID := c.Param("id")
	
	// Check ownership first
	workflow, err := h.queryHandler.GetWorkflowByID(c.Request.Context(), workflowID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "workflow not found"})
		return
	}
	if workflow.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	var req dto.UpdateWorkflowStateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedWorkflow, err := h.commandHandler.UpdateWorkflowState(c.Request.Context(), workflowID, req.State)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedWorkflow)
}

// parseInt parses string to int with default value
func parseInt(s string, defaultVal int) int {
	var val int
	if _, err := fmt.Sscanf(s, "%d", &val); err != nil {
		return defaultVal
	}
	return val
}
