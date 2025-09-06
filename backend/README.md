# CloudStore Backend

The backend service for CloudStore, providing REST and WebSocket APIs for file management and synchronization.

## 🛠️ Technologies

- Go 1.21+
- Echo web framework
- PostgreSQL
- Redis
- WebSocket

## 📁 Project Structure

```
backend/
├── cmd/            # Application entrypoints
├── internal/       # Private application code
├── pkg/           # Public library code
├── api/           # API documentation and schemas
├── config/        # Configuration files
└── test/          # Test files
```

## 🚀 Getting Started

1. **Prerequisites**

   - Go 1.21+
   - Docker and Docker Compose
   - Make

2. **Setup**

   ```bash
   # Install dependencies
   make deps

   # Build the application
   make build
   ```

3. **Run**

   ```bash
   # Run locally
   make run

   # Or with Docker
   make docker-build
   make docker-run
   ```

4. **Development**

   ```bash
   # Run tests
   make test

   # Run linter
   make lint
   ```

## 📝 API Documentation

- `GET /health` - Health check endpoint
- `GET /ws` - WebSocket endpoint for real-time updates

## 🔧 Configuration

Environment variables:

- `BACKEND_PORT` - Server port (default: 8080)
- `POSTGRES_*` - PostgreSQL connection settings
- `REDIS_*` - Redis connection settings
- `JWT_*` - JWT authentication settings
