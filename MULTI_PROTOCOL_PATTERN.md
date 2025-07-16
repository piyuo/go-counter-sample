
# Multi-Protocol  Pattern

This document demonstrates the standard pattern for handling authentication requests from REST, GraphQL, and gRPC protocols through a single business layer.

## 🏗️ Architecture Pattern

```go
┌─────────────────────────────────────────────────────────────┐
│                    Protocol Layer                           │
├─────────────────┬─────────────────┬─────────────────────────┤
│   REST Handler  │ GraphQL Handler │   gRPC Handler          │
│   (Gin/HTTP)    │   (GraphQL)     │   (protobuf)            │
└─────────────────┴─────────────────┴─────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────┐
│                 Business Layer                              │
│              AuthService Interface                          │
│         • Login(request) → response                         │
│         • ValidateUser(username, password)                  │
│         • CreateUser(user)                                  │
└─────────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────┐
│                Repository Layer                             │
│              UserRepository Interface                       │
│         • GetByUsername(username)                           │
│         • Create(user)                                      │
│         • Update(user)                                      │
└─────────────────────────────────────────────────────────────┘
```

## 🔄 Request Flow Example

### User sends login request with: `{"username":"user1","password":"123"}`

**1. REST API Flow:**

```go
POST /api/v1/auth/login
├── RESTHandler.LoginHandler()
├── Converts JSON to domain.LoginRequest
├── Calls authService.Login(ctx, request)
├── Receives domain.LoginResponse
└── Returns HTTP JSON response
```

**2. GraphQL Flow:**

```go
POST /graphql
├── GraphQLHandler.Handler()
├── Parses GraphQL query/variables
├── Converts to domain.LoginRequest
├── Calls authService.Login(ctx, request)
├── Receives domain.LoginResponse
└── Returns GraphQL JSON response
```

**3. gRPC Flow:**

```go
gRPC AuthService.Login()
├── GRPCServer.Login()
├── Converts protobuf to domain.LoginRequest
├── Calls authService.Login(ctx, request)
├── Receives domain.LoginResponse
└── Returns protobuf response
```

## 🧩 Key Components

### 1. Domain Models (Shared)

```go
// Single source of truth for login request/response
type LoginRequest struct {
    Username string `json:"username" validate:"required"`
    Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
    Success bool   `json:"success"`
    Message string `json:"message"`
    UserID  string `json:"user_id,omitempty"`
    Token   string `json:"token,omitempty"`
}
```

### 2. Business Service Interface (Shared)

```go
type AuthService interface {
    Login(ctx context.Context, request *LoginRequest) (*LoginResponse, error)
    ValidateUser(ctx context.Context, username, password string) (*User, error)
    CreateUser(ctx context.Context, user *User) (*User, error)
}
```

### 3. Protocol Adapters

**REST Adapter:**

```go
func (h *RESTHandler) LoginHandler(c *gin.Context) {
    var request domain.LoginRequest

    // 1. Parse protocol-specific input
    if err := c.ShouldBindJSON(&request); err != nil {
        c.JSON(400, gin.H{"error": "Invalid format"})
        return
    }

    // 2. Call business layer (SAME FOR ALL PROTOCOLS)
    response, err := h.authService.Login(c.Request.Context(), &request)
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }

    // 3. Return protocol-specific response
    statusCode := 200
    if !response.Success {
        statusCode = 401
    }
    c.JSON(statusCode, response)
}
```

**GraphQL Adapter:**

```go
func (h *GraphQLHandler) Handler(c *gin.Context) {
    var graphqlRequest GraphQLRequest

    // 1. Parse GraphQL query/variables
    c.ShouldBindJSON(&graphqlRequest)
    input := graphqlRequest.Variables["input"].(map[string]interface{})

    // 2. Convert to domain model
    domainRequest := &domain.LoginRequest{
        Username: input["username"].(string),
        Password: input["password"].(string),
    }

    // 3. Call business layer (SAME FOR ALL PROTOCOLS)
    loginResult, err := h.authService.Login(c.Request.Context(), domainRequest)

    // 4. Return GraphQL response format
    response := GraphQLResponse{
        Data: map[string]interface{}{
            "login": loginResult,
        },
    }
    c.JSON(200, response)
}
```

**gRPC Adapter:**

```go
func (s *GRPCServer) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
    // 1. Convert protobuf to domain model
    domainRequest := &domain.LoginRequest{
        Username: req.Username,
        Password: req.Password,
    }

    // 2. Call business layer (SAME FOR ALL PROTOCOLS)
    response, err := s.authService.Login(ctx, domainRequest)
    if err != nil {
        return nil, err
    }

    // 3. Convert domain response to protobuf
    grpcResponse := &LoginResponse{
        Success: response.Success,
        Message: response.Message,
        UserId:  response.UserID,
        Token:   response.Token,
    }

    return grpcResponse, nil
}
```

## 🎯 Benefits of This Pattern

### 1. **Single Source of Truth**

- Business logic implemented once in AuthService
- All protocols use identical validation and processing
- Consistent behavior across all endpoints

### 2. **Easy Testing**

- Test business logic independently of protocols
- Mock the AuthService interface for protocol tests
- End-to-end tests verify protocol adapters

### 3. **Maintainability**

- Changes to business logic affect all protocols automatically
- Protocol-specific code is minimal and focused on translation
- Clear separation of concerns

### 4. **Scalability**

- Add new protocols by creating new adapters
- Business layer remains unchanged
- Repository layer can be swapped (in-memory → database)

## 🧪 Testing Examples

### Test Business Logic (Protocol-Independent)

```go
func TestAuthService_Login_Success(t *testing.T) {
    // Setup
    userRepo := repository.NewInMemoryUserRepository()
    authService := auth.NewService(userRepo)

    // Test business logic directly
    request := &domain.LoginRequest{
        Username: "user1",
        Password: "123",
    }

    response, err := authService.Login(context.Background(), request)

    // Verify business logic
    assert.NoError(t, err)
    assert.True(t, response.Success)
    assert.Equal(t, "Login successful", response.Message)
}
```

### Test Protocol Adapter

```go
func TestRESTHandler_Login(t *testing.T) {
    // Setup with mock business service
    mockAuthService := &MockAuthService{}
    handler := NewRESTHandler(mockAuthService)

    // Test protocol adapter
    w := httptest.NewRecorder()
    c, _ := gin.CreateTestContext(w)

    // Test REST-specific behavior
    handler.LoginHandler(c)

    // Verify protocol adaptation
    assert.Equal(t, 200, w.Code)
    // Verify business service was called correctly
    mockAuthService.AssertCalled(t, "Login", mock.Anything, mock.Anything)
}
```

## 🚀 Usage Examples

### 1. REST API Request

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"user1","password":"123"}'
```

**Response:**

```json
{
  "success": true,
  "message": "Login successful",
  "user_id": "1",
  "token": "token_1_user1"
}
```

### 2. GraphQL Request

```bash
curl -X POST http://localhost:8080/graphql \
  -H "Content-Type: application/json" \
  -d '{
    "query": "mutation($input: LoginInput!) { login(input: $input) { success message user_id token } }",
    "variables": {
      "input": {
        "username": "user1",
        "password": "123"
      }
    }
  }'
```

**Response:**

```json
{
  "data": {
    "login": {
      "success": true,
      "message": "Login successful",
      "user_id": "1",
      "token": "token_1_user1"
    }
  }
}
```

### 3. gRPC Request (via demo endpoint)

```bash
curl -X POST http://localhost:8080/grpc-demo \
  -H "Content-Type: application/json" \
  -d '{"username":"user1","password":"123"}'
```

**Response:**

```json
{
  "success": true,
  "message": "Login successful",
  "user_id": "1",
  "token": "token_1_user1"
}
```

## 🔍 Key Takeaways

1. **All protocols produce identical results** because they use the same business logic
2. **Protocol adapters are thin** - they only handle format conversion
3. **Business logic is protocol-agnostic** - it works with domain models
4. **Testing is simplified** - test business logic once, test adapters separately
5. **Adding new protocols is easy** - just create a new adapter
