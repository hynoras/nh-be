# Run your Gin app with hot reload (Air)
dev:
	@echo "🔪 Killing any process on port 8080..."
	-@lsof -ti:8080 | xargs kill -9 2>/dev/null || true
	@echo "🚀 Starting Air..."
	go build -o tmp/noheir ./cmd/main.go
	air

# Run without hot reload
run:
	go run ./cmd/main.go

# Build binary
build:
	go build -o bin/nh-be ./cmd/main.go

# Run tests
test:
	go test ./... -v

# Clean up build files
clean:
	rm -rf tmp bin