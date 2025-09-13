# Build stage
FROM hub.hamdocker.ir/golang:1.24.5-alpine AS builder

# Install git and ca-certificates (needed for go mod download)
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
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o zeus ./cmd/zeus

# Final stage
FROM hub.hamdocker.ir/alpine:3.19

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates tzdata

# Create non-root user
RUN addgroup -g 1001 -S zeus && \
    adduser -u 1001 -S zeus -G zeus

# Set working directory
WORKDIR /app

# Copy the binary from builder stage
COPY --from=builder /app/zeus .

# Copy any additional files if needed (like config files)
# COPY --from=builder /app/sample.env .

# Change ownership to non-root user
RUN chown -R zeus:zeus /app

# Switch to non-root user
USER zeus

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Run the application
CMD ["./zeus"]
