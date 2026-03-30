# Makefile for iamzcr blog system

.PHONY: help build build-admin build-frontend run run-admin run-frontend clean all

# Go build variables
GO=go
GO_BUILD=$(GO) build
GO_RUN=$(GO) run

# Output binaries
ADMIN_BINARY=admin
FRONTEND_BINARY=frontend

# Source directories
ADMIN_DIR=./cmd/admin
FRONTEND_DIR=./cmd/frontend

# Help
help:
	@echo "IAMZCR Blog System - Makefile Commands"
	@echo ""
	@echo "Usage:"
	@echo "  make build          - Build both admin and frontend services"
	@echo "  make build-admin    - Build admin service only"
	@echo "  make build-frontend - Build frontend service only"
	@echo "  make run            - Run both services"
	@echo "  make run-admin      - Run admin service"
	@echo "  make run-frontend   - Run frontend service"
	@echo "  make clean          - Remove built binaries"
	@echo "  make all            - Build and run both services"
	@echo ""

# Build both services
build: build-admin build-frontend
	@echo "Build completed!"

# Build admin service
build-admin:
	@echo "Building admin service..."
	$(GO_BUILD) -o $(ADMIN_BINARY) $(ADMIN_DIR)
	@echo "Admin service built: $(ADMIN_BINARY)"

# Build frontend service
build-frontend:
	@echo "Building frontend service..."
	$(GO_BUILD) -o $(FRONTEND_BINARY) $(FRONTEND_DIR)
	@echo "Frontend service built: $(FRONTEND_BINARY)"

# Run both services
run: run-admin run-frontend

# Run admin service
run-admin:
	@echo "Running admin service on port 8081..."
	$(GO_RUN) $(ADMIN_DIR)

# Run frontend service
run-frontend:
	@echo "Running frontend service on port 8082..."
	$(GO_RUN) $(FRONTEND_DIR)

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -f $(ADMIN_BINARY) $(FRONTEND_BINARY)
	@rm -f $(ADMIN_BINARY).exe $(FRONTEND_BINARY).exe
	@echo "Clean completed!"

# Build and run both services
all: build
	@echo "Starting services..."
	@echo "Admin API: http://localhost:8081"
	@echo "Frontend API: http://localhost:8082"
