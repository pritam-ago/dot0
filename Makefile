.PHONY: help build up down logs clean dev prod restart rebuild

# Default target
help:
	@echo "Available commands:"
	@echo "  make build    - Build all Docker images"
	@echo "  make up       - Start all services (production)"
	@echo "  make dev      - Start all services (development)"
	@echo "  make down     - Stop all services"
	@echo "  make logs     - Show logs for all services"
	@echo "  make clean    - Stop services and remove containers/images"
	@echo "  make restart  - Restart all services"
	@echo "  make rebuild  - Rebuild and start all services"

# Production commands
build:
	docker-compose build

up:
	docker-compose up -d

down:
	docker-compose down

logs:
	docker-compose logs -f

restart:
	docker-compose restart

rebuild:
	docker-compose down
	docker-compose up --build -d

# Development commands
dev:
	docker-compose -f docker-compose.dev.yml up --build

dev-down:
	docker-compose -f docker-compose.dev.yml down

dev-logs:
	docker-compose -f docker-compose.dev.yml logs -f

# Cleanup commands
clean:
	docker-compose down -v --rmi all
	docker system prune -f

# Individual service commands
build-go:
	docker-compose build go-relay

build-web:
	docker-compose build web-client

build-pc:
	docker-compose build pc-app

# Database commands
db-backup:
	docker-compose exec postgres pg_dump -U relay_user relay_db > backup_$(shell date +%Y%m%d_%H%M%S).sql

db-restore:
	@echo "Usage: make db-restore FILE=backup_file.sql"
	@if [ -z "$(FILE)" ]; then echo "Please specify FILE parameter"; exit 1; fi
	docker-compose exec -T postgres psql -U relay_user relay_db < $(FILE)

# Health checks
health:
	@echo "Checking service health..."
	@docker-compose ps
	@echo ""
	@echo "Health check endpoints:"
	@echo "  Go Relay: http://localhost:8080/health"
	@echo "  Web Client: http://localhost:3001/health"
	@echo "  PC App: http://localhost:3002/health"

# Port information
ports:
	@echo "Service ports:"
	@echo "  Go Relay Server: http://localhost:8080"
	@echo "  Web Client: http://localhost:3001"
	@echo "  PC App: http://localhost:3002"
	@echo "  PostgreSQL: localhost:5432" 