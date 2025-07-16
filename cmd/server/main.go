// ===============================================
// Module: main.go
// Description: Main server application with REST, GraphQL, and gRPC endpoints
//
// Sections:
//   - Imports and Dependencies
//   - Service Initialization
//   - REST API Setup
//   - GraphQL Setup
//   - gRPC Setup
//   - Server Startup
// ===============================================

package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/piyuo/go-counter-sample/internal/auth"
	"github.com/piyuo/go-counter-sample/internal/handlers"
	"github.com/piyuo/go-counter-sample/internal/repository"
)

func main() {
	// ===============================================
	// DEPENDENCY INJECTION - Single Business Layer
	// ===============================================
	// Initialize repository layer
	userRepo := repository.NewInMemoryUserRepository()

	// Initialize business service layer (SINGLE SOURCE OF TRUTH)
	authService := auth.NewService(userRepo)

	// ===============================================
	// WEB SERVER SETUP
	// ===============================================
	// Create Gin router
	router := gin.Default()

	// Add CORS middleware for development
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// ===============================================
	// PROTOCOL HANDLERS - All use same business layer
	// ===============================================

	// 1. REST API Handler
	restHandler := handlers.NewRESTHandler(authService)
	restHandler.RegisterRoutes(router)

	// 2. GraphQL Handler
	graphqlHandler := handlers.NewGraphQLHandler(authService)
	graphqlHandler.RegisterRoutes(router)

	// 3. gRPC Server Handler (for demonstration)
	grpcServer := handlers.NewGRPCServer(authService)

	// Add a REST endpoint that demonstrates gRPC-style calls
	router.POST("/grpc-demo", func(c *gin.Context) {
		var request struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid request format",
			})
			return
		}

		// Use the actual gRPC server implementation
		grpcRequest := &handlers.LoginRequest{
			Username: request.Username,
			Password: request.Password,
		}

		// Call the gRPC handler directly (simulating gRPC call)
		grpcResponse, err := grpcServer.Login(c.Request.Context(), grpcRequest)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, grpcResponse)
	})

	// ===============================================
	// API DOCUMENTATION ENDPOINT
	// ===============================================
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Go Counter Sample API - Multi-Protocol Authentication",
			"version": "1.0.0",
			"architecture": gin.H{
				"pattern": "All protocols route through single business layer",
				"layers": []string{
					"Protocol Layer (REST/GraphQL/gRPC)",
					"Business Layer (AuthService)",
					"Repository Layer (UserRepository)",
				},
			},
			"endpoints": gin.H{
				"rest": gin.H{
					"login":  "POST /api/v1/auth/login",
					"health": "GET /api/v1/health",
				},
				"graphql": gin.H{
					"endpoint":      "POST /graphql",
					"query_example": `{"query": "mutation { login(input: {username: \"user1\", password: \"123\"}) { success message user_id token } }"}`,
				},
				"grpc": gin.H{
					"demo_endpoint": "POST /grpc-demo",
					"note":          "Real gRPC would run on separate port",
				},
			},
			"test_credentials": gin.H{
				"username": "user1",
				"password": "123",
			},
			"sample_requests": gin.H{
				"rest_login": gin.H{
					"curl": `curl -X POST http://localhost:8080/api/v1/auth/login -H 'Content-Type: application/json' -d '{"username":"user1","password":"123"}'`,
				},
				"graphql_login": gin.H{
					"curl": `curl -X POST http://localhost:8080/graphql -H 'Content-Type: application/json' -d '{"query":"mutation { login(input: {username: \"user1\", password: \"123\"}) { success message user_id token } }"}'`,
				},
				"grpc_demo": gin.H{
					"curl": `curl -X POST http://localhost:8080/grpc-demo -H 'Content-Type: application/json' -d '{"username":"user1","password":"123"}'`,
				},
			},
		})
	})

	// ===============================================
	// SERVER STARTUP
	// ===============================================
	// Get port from environment variable, default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start the server
	log.Printf("🚀 Multi-Protocol Authentication Server starting on port %s", port)
	log.Printf("📖 API Documentation: http://localhost:%s/", port)
	log.Printf("🔧 REST API: http://localhost:%s/api/v1/auth/login", port)
	log.Printf("🎯 GraphQL: http://localhost:%s/graphql", port)
	log.Printf("⚡ gRPC Demo: http://localhost:%s/grpc-demo", port)
	log.Printf("")
	log.Printf("💡 Test REST: curl -X POST http://localhost:%s/api/v1/auth/login -H 'Content-Type: application/json' -d '{\"username\":\"user1\",\"password\":\"123\"}'", port)
	log.Printf("💡 Test GraphQL: curl -X POST http://localhost:%s/graphql -H 'Content-Type: application/json' -d '{\"query\":\"mutation { login(input: {username: \\\"user1\\\", password: \\\"123\\\"}) { success message user_id token } }\"}'", port)
	log.Printf("💡 Test gRPC Demo: curl -X POST http://localhost:%s/grpc-demo -H 'Content-Type: application/json' -d '{\"username\":\"user1\",\"password\":\"123\"}'", port)

	log.Fatal(router.Run(":" + port))
}
