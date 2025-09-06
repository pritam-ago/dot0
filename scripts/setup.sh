#!/bin/bash

# Exit on error
set -e

echo "🚀 Setting up CloudStore development environment..."

# Check prerequisites
echo "📋 Checking prerequisites..."

command -v docker >/dev/null 2>&1 || { echo "❌ Docker is required but not installed. Aborting." >&2; exit 1; }
command -v docker-compose >/dev/null 2>&1 || { echo "❌ Docker Compose is required but not installed. Aborting." >&2; exit 1; }
command -v go >/dev/null 2>&1 || { echo "❌ Go is required but not installed. Aborting." >&2; exit 1; }
command -v node >/dev/null 2>&1 || { echo "❌ Node.js is required but not installed. Aborting." >&2; exit 1; }
command -v rustc >/dev/null 2>&1 || { echo "❌ Rust is required but not installed. Aborting." >&2; exit 1; }

# Copy environment file if not exists
echo "📝 Setting up environment variables..."
if [ ! -f .env ]; then
    cp .env.example .env
    echo "✅ Created .env file from template"
fi

# Install backend dependencies
echo "📦 Installing backend dependencies..."
cd backend
go mod download
cd ..

# Install PC client dependencies
echo "📦 Installing PC client dependencies..."
cd pc-client
npm install
cd ..

# Install mobile app dependencies
echo "📦 Installing mobile app dependencies..."
cd mobile
npm install
cd ..

# Set up database
echo "🗄️ Setting up database..."
cd db
if ! command -v migrate >/dev/null 2>&1; then
    echo "Installing golang-migrate..."
    go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
fi
cd ..

# Start infrastructure
echo "🚀 Starting infrastructure..."
docker-compose -f infra/docker-compose.yml up -d postgres redis

# Wait for PostgreSQL to be ready
echo "⏳ Waiting for PostgreSQL to be ready..."
sleep 5

# Run database migrations
echo "🔄 Running database migrations..."
migrate -path db/migrations -database "postgresql://cloudstore:development_password@localhost:5432/cloudstore?sslmode=disable" up

echo "✨ Setup completed successfully!"
echo "
Next steps:
1. Start the backend: cd backend && go run main.go
2. Start the PC client: cd pc-client && npm run tauri dev
3. Start the mobile app: cd mobile && npm start
"
