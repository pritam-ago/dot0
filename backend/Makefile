.PHONY: build run test lint clean

# Build the application
build:
	go build -o bin/main .

# Run the application
run:
	go run main.go

# Run tests
test:
	go test -v ./...

# Run linter
lint:
	golangci-lint run

# Clean build artifacts
clean:
	rm -rf bin/

# Download dependencies
deps:
	go mod download

# Update dependencies
deps-update:
	go get -u ./...
	go mod tidy

# Build Docker image
docker-build:
	docker build -t cloudstore-backend .

# Run Docker container
docker-run:
	docker run -p 8080:8080 cloudstore-backend
