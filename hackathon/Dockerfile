FROM golang:1.21-bullseye AS builder

# Install build dependencies
RUN apt-get update && apt-get install -y \
    build-essential \
    pkg-config \
    sqlite3 \
    libsqlite3-dev \
    ca-certificates \
    git && rm -rf /var/lib/apt/lists/*

# Set working directory
WORKDIR /build

# Copy go.mod and go.sum files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o /build/app .

# Build final image
FROM debian:bullseye-slim AS runtime
RUN apt-get update && apt-get install -y \
    ca-certificates && rm -rf /var/lib/apt/lists/*

# Copy built binary from builder stage
COPY --from=builder /build/app /app/app

# Set working directory
WORKDIR /app

# Expose port
EXPOSE 8080

# Run the binary
CMD ["./app"]
