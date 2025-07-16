// ===============================================
// Module: graphql_handlers.go
// Description: Simplified GraphQL-style handler for authentication operations
//
// Sections:
//   - Handler Structure
//   - Constructor
//   - GraphQL-style Request/Response Types
//   - Login Handler
//   - Route Registration
// ===============================================

package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/piyuo/go-counter-sample/internal/core/domain"
)

// GraphQLHandler handles GraphQL-style requests
type GraphQLHandler struct {
	authService domain.AuthService
}

// NewGraphQLHandler creates a new GraphQL-style handler
func NewGraphQLHandler(authService domain.AuthService) *GraphQLHandler {
	return &GraphQLHandler{
		authService: authService,
	}
}

// GraphQL-style request and response types
type GraphQLRequest struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables"`
}

type GraphQLResponse struct {
	Data   interface{} `json:"data,omitempty"`
	Errors []string    `json:"errors,omitempty"`
}

// Handler handles GraphQL-style HTTP requests
func (h *GraphQLHandler) Handler(c *gin.Context) {
	var requestBody GraphQLRequest

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, GraphQLResponse{
			Errors: []string{"Invalid JSON format"},
		})
		return
	}

	// Simple query parsing for demo purposes
	// In production, use a proper GraphQL library like graphql-go/graphql

	response := GraphQLResponse{}

	// Handle login mutation
	if contains(requestBody.Query, "login") {
		// Extract variables for login
		input, ok := requestBody.Variables["input"].(map[string]interface{})
		if !ok {
			response.Errors = []string{"Invalid login input"}
			c.JSON(http.StatusBadRequest, response)
			return
		}

		username, _ := input["username"].(string)
		password, _ := input["password"].(string)

		if username == "" || password == "" {
			response.Errors = []string{"Username and password are required"}
			c.JSON(http.StatusBadRequest, response)
			return
		}

		// Call business service
		domainRequest := &domain.LoginRequest{
			Username: username,
			Password: password,
		}

		loginResult, err := h.authService.Login(c.Request.Context(), domainRequest)
		if err != nil {
			response.Errors = []string{err.Error()}
			c.JSON(http.StatusInternalServerError, response)
			return
		}

		response.Data = map[string]interface{}{
			"login": map[string]interface{}{
				"success": loginResult.Success,
				"message": loginResult.Message,
				"user_id": loginResult.UserID,
				"token":   loginResult.Token,
			},
		}

	} else if contains(requestBody.Query, "health") {
		// Handle health query
		response.Data = map[string]interface{}{
			"health": "GraphQL service is healthy",
		}

	} else {
		response.Errors = []string{"Unknown query"}
	}

	statusCode := http.StatusOK
	if len(response.Errors) > 0 {
		statusCode = http.StatusBadRequest
	}

	c.JSON(statusCode, response)
}

// RegisterRoutes registers GraphQL routes
func (h *GraphQLHandler) RegisterRoutes(router *gin.Engine) {
	router.POST("/graphql", h.Handler)

	// GraphQL playground for development
	router.GET("/graphql", func(c *gin.Context) {
		c.Header("Content-Type", "text/html")
		c.String(http.StatusOK, graphqlPlayground)
	})
}

// Helper function to check if string contains substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr ||
		(len(s) > len(substr) &&
			(s[:len(substr)] == substr ||
				s[len(s)-len(substr):] == substr ||
				indexOfSubstring(s, substr) >= 0)))
}

func indexOfSubstring(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

// Simple GraphQL Playground HTML
const graphqlPlayground = `
<!DOCTYPE html>
<html>
<head>
    <title>GraphQL API</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            margin: 40px;
            background-color: #f5f5f5;
        }
        .container {
            max-width: 800px;
            margin: 0 auto;
            background: white;
            padding: 20px;
            border-radius: 8px;
            box-shadow: 0 2px 10px rgba(0,0,0,0.1);
        }
        pre {
            background: #f8f8f8;
            padding: 15px;
            border-radius: 4px;
            overflow-x: auto;
        }
        .endpoint {
            background: #e3f2fd;
            padding: 10px;
            border-radius: 4px;
            margin: 10px 0;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>GraphQL API</h1>
        <p>Send POST requests to <code>/graphql</code> with the following format:</p>

        <div class="endpoint">
            <strong>Endpoint:</strong> POST /graphql
        </div>

        <h3>Example Login Mutation:</h3>
        <pre>{
  "query": "mutation($input: LoginInput!) { login(input: $input) { success message user_id token } }",
  "variables": {
    "input": {
      "username": "user1",
      "password": "123"
    }
  }
}</pre>

        <h3>Example Health Query:</h3>
        <pre>{
  "query": "{ health }"
}</pre>

        <h3>Test with curl:</h3>
        <pre>curl -X POST http://localhost:8080/graphql \
  -H "Content-Type: application/json" \
  -d '{
    "query": "mutation($input: LoginInput!) { login(input: $input) { success message user_id token } }",
    "variables": {
      "input": {
        "username": "user1",
        "password": "123"
      }
    }
  }'</pre>
    </div>
</body>
</html>`
