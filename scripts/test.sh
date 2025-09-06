#!/bin/bash

# Exit on error
set -e

echo "🧪 Running tests..."

# Backend
echo "🔍 Testing backend..."
cd backend
go test -v ./...
cd ..

# PC Client
echo "🔍 Testing PC client..."
cd pc-client
npm test
cd ..

# Mobile App
echo "🔍 Testing mobile app..."
cd mobile
npm test
cd ..

echo "✨ All tests completed successfully!"
