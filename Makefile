.PHONY: build run dev test lint generate clean docker-up docker-down migrate

# Build the application
build: generate
	go build -o bin/web ./cmd/web

# Run the application locally
run: build
	./bin/web

# Development mode with live reload (requires air)
dev:
	air

# Generate templ files
generate:
	~/go/bin/templ generate

# Run tests
test:
	go test -v -race -count=1 ./...

# Run linter
lint:
	golangci-lint run ./...

# Clean build artifacts
clean:
	rm -rf bin/ tmp/
	find . -name '*_templ.go' -delete

# Docker operations
docker-up:
	cd docker && docker compose up -d --build

docker-down:
	cd docker && docker compose down

docker-logs:
	cd docker && docker compose logs -f

docker-reset:
	cd docker && docker compose down -v && docker compose up -d --build

# Tidy dependencies
tidy:
	go mod tidy

# Format code
fmt:
	gofmt -s -w .
	~/go/bin/templ fmt .
