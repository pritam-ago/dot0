#!/bin/bash

# Exit on error
set -e

echo "🔍 Running linters..."

# Backend
echo "📝 Linting backend code..."
cd backend
if ! command -v golangci-lint >/dev/null 2>&1; then
    echo "Installing golangci-lint..."
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
fi
golangci-lint run
cd ..

# PC Client
echo "📝 Linting PC client code..."
cd pc-client
npm run lint
cd ..

# Mobile App
echo "📝 Linting mobile app code..."
cd mobile
npm run lint
cd ..

echo "✨ Linting completed successfully!"
