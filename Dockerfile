# --------------------------------------------------------
# Stage 1: Builder
# --------------------------------------------------------
FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -trimpath -o server ./cmd/main.go
# --------------------------------------------------------
# Stage 2: Final image
# --------------------------------------------------------
FROM alpine:3.20

RUN apk add --no-cache ca-certificates wget \
    && adduser -D -u 1001 appuser

WORKDIR /app

COPY --from=builder /app/server .
COPY --from=builder /app/templates ./templates

USER appuser

# Expose your service port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
    CMD wget --spider -q http://localhost:8080/health/live || exit 1
    
# Run the server
ENTRYPOINT ["./server"]
