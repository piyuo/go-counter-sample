// ===============================================
// Module: rest_handlers.go
// Description: REST API handlers for authentication endpoints
//
// Sections:
//   - Handler Structure
//   - Constructor
//   - Login Endpoint Handler
//   - Register Routes Helper
//   - Response Helper Functions
// ===============================================

package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/piyuo/go-counter-sample/internal/core/domain"
)

// RESTHandler handles REST API requests
type RESTHandler struct {
	authService domain.AuthService
}

// NewRESTHandler creates a new REST handler
func NewRESTHandler(authService domain.AuthService) *RESTHandler {
	return &RESTHandler{
		authService: authService,
	}
}

// LoginHandler handles POST /api/v1/auth/login
func (h *RESTHandler) LoginHandler(c *gin.Context) {
	var request domain.LoginRequest

	// Bind JSON request to struct
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request format",
			"error":   err.Error(),
		})
		return
	}

	// Call business service
	response, err := h.authService.Login(c.Request.Context(), &request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Internal server error",
			"error":   err.Error(),
		})
		return
	}

	// Determine HTTP status code based on success
	statusCode := http.StatusOK
	if !response.Success {
		statusCode = http.StatusUnauthorized
	}

	c.JSON(statusCode, response)
}

// HealthHandler handles GET /api/v1/health
func (h *RESTHandler) HealthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "go-counter-sample",
		"version": "1.0.0",
	})
}

// RegisterRoutes registers all REST routes
func (h *RESTHandler) RegisterRoutes(router *gin.Engine) {
	// API v1 group
	v1 := router.Group("/api/v1")
	{
		// Health check
		v1.GET("/health", h.HealthHandler)

		// Authentication routes
		auth := v1.Group("/auth")
		{
			auth.POST("/login", h.LoginHandler)
		}
	}
}
