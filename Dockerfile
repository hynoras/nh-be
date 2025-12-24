# --------------------------------------------------------
# Stage 1: Builder
# --------------------------------------------------------
FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server ./cmd/main.go

# --------------------------------------------------------
# Stage 2: Final image
# --------------------------------------------------------
FROM alpine:3.20

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /app/server .

# Expose your service port
EXPOSE 8080

# Run the server
ENTRYPOINT ["./server"]
