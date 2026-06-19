package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/labhaus/backend/internal/application/command"
	"github.com/labhaus/backend/internal/application/dto"
	"github.com/labhaus/backend/internal/application/query"
	"github.com/labhaus/backend/internal/infrastructure/auth"
	"github.com/labhaus/backend/internal/infrastructure/http/middleware"
)

// UserHandler handles user-related HTTP requests
type UserHandler struct {
	queryHandler   *query.UserQueryHandler
	commandHandler *command.UserCommandHandler
	jwtService     *auth.JWTService
}

// NewUserHandler creates a new user handler
func NewUserHandler(
	queryHandler *query.UserQueryHandler,
	commandHandler *command.UserCommandHandler,
	jwtService *auth.JWTService,
) *UserHandler {
	return &UserHandler{
		queryHandler:   queryHandler,
		commandHandler: commandHandler,
		jwtService:     jwtService,
	}
}

// Register handles POST /api/users/register
func (h *UserHandler) Register(c *gin.Context) {
	var req dto.RegisterUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.commandHandler.RegisterUser(c.Request.Context(), req)
	if err != nil {
		// Check for specific errors
		if err == command.ErrEmailAlreadyExists {
			c.JSON(http.StatusConflict, gin.H{"error": "email already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, user)
}

// Login handles POST /api/users/login
func (h *UserHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.queryHandler.Authenticate(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
		return
	}

	// Generate JWT token
	token, err := h.jwtService.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user":  user,
	})
}

// GetMe handles GET /api/users/me
func (h *UserHandler) GetMe(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	user, err := h.queryHandler.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// UpdateMe handles PATCH /api/users/me
func (h *UserHandler) UpdateMe(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.commandHandler.UpdateUser(c.Request.Context(), userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}
