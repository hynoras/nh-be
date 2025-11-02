# --------------------------------------------------------
# Stage 1: Builder
# --------------------------------------------------------
FROM golang:1.25-alpine AS builder

# Install git for go mod download (no need for build-base since no CGO)
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy dependency files first for caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of your source code
COPY . .

# Build the Go binary (no CGO, static build)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server ./cmd/main.go

# --------------------------------------------------------
# Stage 2: Final image
# --------------------------------------------------------
FROM scratch

# Copy certificates for HTTPS
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy binary and any necessary runtime files
COPY --from=builder /app/server .
COPY --from=builder /app/.env ./ 
# COPY --from=builder /app/config.yaml .  # (optional)

# Expose your service port
EXPOSE 8080

# Run the server
ENTRYPOINT ["./server"]
