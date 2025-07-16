# go-counter-sample

sample web application for counter data collecting

## tech stack

GIN : provides a clean API for routing, middleware, and request/response handling

## 🏗️ Architecture Pattern

```go
Request from REST/GraphQL/gRPC → Protocol Handler → Business Layer → Repository Layer
```

### Protocol Adapters

Each protocol has a thin adapter that:

1. **Parses protocol-specific input** (JSON, GraphQL, protobuf)
2. **Converts to domain model** (`domain.LoginRequest`)
3. **Calls the same business service** (`authService.Login()`)
4. **Returns protocol-specific response**

## 🧪 Verification Tests

### Test Valid Credentials (All Protocols Return Same Result)

- **REST**: `{"success":true,"message":"Login successful","user_id":"1","token":"token_1_user1"}`
- **GraphQL**: `{"data":{"login":{"success":true,"message":"Login successful","user_id":"1","token":"token_1_user1"}}}`
- **gRPC**: `{"success":true,"message":"Login successful","user_id":"1","token":"token_1_user1"}`

### Test Invalid Credentials (Consistent Error Handling)

- **REST**: `{"success":false,"message":"Invalid username or password"}`
- **gRPC**: `{"success":false,"message":"Invalid username or password","user_id":"","token":""}`

## 🚀 Standard Pattern Benefits

✅ **Single Source of Truth**: Business logic implemented once
✅ **Consistent Behavior**: All protocols return identical results
✅ **Easy Testing**: Test business logic independently
✅ **Maintainability**: Changes affect all protocols automatically
✅ **Scalability**: Add new protocols by creating new adapters

## 📁 File Structure

```go
cmd/server/main.go                 # Server setup with all protocols
internal/
├── auth/
│   ├── auth_service.go           # Business logic implementation
│   └── auth_service_test.go      # Business logic tests
├── core/domain/
│   └── user.go                   # Domain models and interfaces
├── handlers/
│   ├── rest_handlers.go          # REST protocol adapter
│   ├── graphql_handlers.go       # GraphQL protocol adapter
│   └── grpc_server.go           # gRPC protocol adapter
└── repository/
    └── user_repository.go        # Data access layer
```

## 🔧 Development Commands

```bash
# Start the server
go run cmd/server/main.go

# Run tests
go test ./...

# Test specific protocol
curl -X POST http://localhost:8080/api/v1/auth/login -H 'Content-Type: application/json' -d '{"username":"user1","password":"123"}'
```

## 🎯 Live Demo

Your server is running on `http://localhost:8080` with the following endpoints:

### 1. REST API

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"user1","password":"123"}'
```

**Response:**

```json
{"success":true,"message":"Login successful","user_id":"1","token":"token_1_user1"}
```

### 2. GraphQL API

```bash
curl -X POST http://localhost:8080/graphql \
  -H "Content-Type: application/json" \
  -d '{"query":"mutation { login(input: {username: \"user1\", password: \"123\"}) { success message user_id token } }","variables":{"input":{"username":"user1","password":"123"}}}'
```

**Response:**

```json
{"data":{"login":{"message":"Login successful","success":true,"token":"token_1_user1","user_id":"1"}}}
```

### 3. gRPC Demo

```bash
curl -X POST http://localhost:8080/grpc-demo \
  -H "Content-Type: application/json" \
  -d '{"username":"user1","password":"123"}'
```

**Response:**

```json
{"success":true,"message":"Login successful","user_id":"1","token":"token_1_user1"}
```

## 🔍 Key Implementation Details
