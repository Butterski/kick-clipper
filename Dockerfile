# Multi-stage build for minimal production image
FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o kick-clipper .

# Production stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates tzdata

# Create non-root user
RUN addgroup -g 1001 -S kickbot && \
    adduser -S kickbot -u 1001 -G kickbot

WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /app/kick-clipper .

# Copy proxies file from root directory
COPY proxies.txt .

# Set ownership
RUN chown -R kickbot:kickbot /app

# Switch to non-root user
USER kickbot

# Default command
ENTRYPOINT ["./kick-clipper"]

# Default arguments (can be overridden)
CMD ["-clip", "example", "-workers", "50", "-delay", "10", "-time", "300"]