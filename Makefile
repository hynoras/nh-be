# Run your Gin app with hot reload (Air)
dev:
	mkdir -p bin
	go build -o bin/nh-be ./cmd/main.go
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
	rm -rf tmp bin build main.exe