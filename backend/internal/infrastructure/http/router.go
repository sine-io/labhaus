package http

import (
	"github.com/gin-gonic/gin"
	"github.com/labhaus/backend/internal/infrastructure/auth"
	"github.com/labhaus/backend/internal/infrastructure/http/handlers"
	"github.com/labhaus/backend/internal/infrastructure/http/middleware"
	"github.com/labhaus/backend/internal/infrastructure/logger"
)

// Router holds all HTTP handlers and routes
type Router struct {
	engine        *gin.Engine
	healthHandler *handlers.HealthHandler
	styleHandler  *handlers.StyleHandler
	userHandler   *handlers.UserHandler
	jwtService    *auth.JWTService
	logger        *logger.Logger
}

// NewRouter creates a new HTTP router
func NewRouter(
	healthHandler *handlers.HealthHandler,
	styleHandler *handlers.StyleHandler,
	userHandler *handlers.UserHandler,
	jwtService *auth.JWTService,
	logger *logger.Logger,
) *Router {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()

	// Middleware
	engine.Use(ginLogger(logger))
	engine.Use(gin.Recovery())

	return &Router{
		engine:        engine,
		healthHandler: healthHandler,
		styleHandler:  styleHandler,
		userHandler:   userHandler,
		jwtService:    jwtService,
		logger:        logger,
	}
}

// Setup configures all routes
func (r *Router) Setup() {
	api := r.engine.Group("/api")

	// Health check
	api.GET("/health", r.healthHandler.Check)

	// User routes (public)
	users := api.Group("/users")
	{
		users.POST("/register", r.userHandler.Register)
		users.POST("/login", r.userHandler.Login)
	}

	// User routes (authenticated)
	usersAuth := api.Group("/users")
	usersAuth.Use(middleware.AuthMiddleware(r.jwtService))
	{
		usersAuth.GET("/me", r.userHandler.GetMe)
		usersAuth.PATCH("/me", r.userHandler.UpdateMe)
	}

	// Styles (authenticated)
	styles := api.Group("/styles")
	styles.Use(middleware.AuthMiddleware(r.jwtService))
	{
		styles.GET("", r.styleHandler.ListStyles)
		styles.GET("/:id", r.styleHandler.GetStyle)
		styles.POST("", r.styleHandler.CreateStyle)
	}
}

// Engine returns the underlying Gin engine
func (r *Router) Engine() *gin.Engine {
	return r.engine
}

// ginLogger creates a Gin middleware that logs to zerolog
func ginLogger(logger *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Process request
		c.Next()

		// Log after request
		logger.Info("HTTP request",
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"status", c.Writer.Status(),
			"ip", c.ClientIP(),
		)
	}
}
