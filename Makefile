.PHONY: build run test lint clean docker-up docker-down frontend-dev frontend-build frontend-test

# Build the Go API
build:
	/opt/homebrew/bin/go build -o bin/api ./cmd/api

# Run the API locally
run: build
	./bin/api

# Run Go tests
test:
	/opt/homebrew/bin/go test -v -race -count=1 ./...

# Run Go linter
lint:
	golangci-lint run ./...

# Clean build artifacts
clean:
	rm -rf bin/ tmp/

# Docker operations
docker-up:
	cd docker && docker compose up -d --build

docker-down:
	cd docker && docker compose down

docker-logs:
	cd docker && docker compose logs -f

docker-reset:
	cd docker && docker compose down -v && docker compose up -d --build

# Frontend operations
frontend-dev:
	cd frontend && npm run dev

frontend-build:
	cd frontend && npm run build

frontend-test:
	cd frontend && npm test

frontend-lint:
	cd frontend && npm run lint

# Tidy dependencies
tidy:
	/opt/homebrew/bin/go mod tidy

# Format code
fmt:
	gofmt -s -w .
