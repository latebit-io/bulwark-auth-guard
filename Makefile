.PHONY: help test test-up test-down test-run test-logs fmt vet tidy clean

# Default target
help:
	@echo "Available targets:"
	@echo "  make test       - Start services, run tests, stop services"
	@echo "  make test-up    - Start test services (docker-compose up)"
	@echo "  make test-down  - Stop test services (docker-compose down)"
	@echo "  make test-run   - Run tests only (services must be running)"
	@echo "  make test-logs  - View logs from test services"
	@echo "  make fmt        - Format code"
	@echo "  make vet        - Run go vet"
	@echo "  make tidy       - Tidy go modules"
	@echo "  make clean      - Stop services and clean up"

# Full test lifecycle: start services, run tests, stop services
test:
	@echo "Starting test services..."
	@docker-compose up -d
	@echo "Waiting for services to be healthy..."
	@sleep 3
	@echo "Checking bulwark-auth readiness..."
	@for i in 1 2 3 4 5 6 7 8 9 10; do \
		if curl -s http://localhost:8080/health > /dev/null 2>&1 || curl -s http://localhost:8080/ > /dev/null 2>&1; then \
			echo "Services ready!"; \
			break; \
		fi; \
		echo "Waiting for bulwark-auth... ($$i/10)"; \
		sleep 2; \
	done
	@echo "Running tests..."
	@go test -v || (docker-compose down && exit 1)
	@echo "Stopping test services..."
	@docker-compose down

# Start test services
test-up:
	@echo "Starting test services..."
	@docker-compose up -d
	@echo "Waiting for services to be healthy..."
	@sleep 3
	@echo "Checking bulwark-auth readiness..."
	@for i in 1 2 3 4 5 6 7 8 9 10; do \
		if curl -s http://localhost:8080/health > /dev/null 2>&1 || curl -s http://localhost:8080/ > /dev/null 2>&1; then \
			echo "Services ready!"; \
			break; \
		fi; \
		echo "Waiting for bulwark-auth... ($$i/10)"; \
		sleep 2; \
	done
	@echo "  - bulwark-auth: http://localhost:8080"
	@echo "  - MailHog UI:   http://localhost:8025"
	@echo "  - MongoDB:      mongodb://localhost:27017"

# Stop test services
test-down:
	@echo "Stopping test services..."
	@docker-compose down

# Run tests only (assumes services are already running)
test-run:
	@echo "Running tests..."
	@go test -v

# View logs from test services
test-logs:
	@docker-compose logs -f

# Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...

# Run go vet
vet:
	@echo "Running go vet..."
	@go vet ./...

# Tidy go modules
tidy:
	@echo "Tidying go modules..."
	@go mod tidy

# Clean up everything
clean: test-down
	@echo "Cleaned up!"
