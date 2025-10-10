# Build stage
FROM golang:1.24-alpine AS builder

# Install git and ca-certificates for fetching dependencies
RUN apk update && apk add --no-cache git ca-certificates tzdata && update-ca-certificates

# Create appuser for security
RUN adduser -D -g '' appuser

# Set working directory
WORKDIR /build

# Copy go mod files and download dependencies
COPY go.mod go.sum ./
RUN go mod download
RUN go mod verify

# Copy source code
COPY . .

# Install swag for generating Swagger documentation
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Generate Swagger docs
RUN swag init

# Build the binary with optimizations
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o edtech-app main.go

# Final stage - minimal runtime image
FROM scratch

# Import from builder
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/passwd /etc/passwd

# Copy the binary
COPY --from=builder /build/edtech-app /app/edtech-app

# Copy static files (docs folder for Swagger)
COPY --from=builder /build/docs /app/docs

# Use non-root user
USER appuser

# Expose port
EXPOSE 8081

# Set working directory
WORKDIR /app

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD ["/app/edtech-app", "--health-check"] || exit 1

# Run the binary
ENTRYPOINT ["/app/edtech-app"]