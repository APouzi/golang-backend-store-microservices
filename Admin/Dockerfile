# -------- Stage 1: Build --------
FROM golang:1.25 AS builder

WORKDIR /app

# Copy go mod files separately to leverage caching
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Copy rest of source code
COPY . .

# Build binary statically (CGO disabled)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o golang-shop . && chmod +x golang-shop

# -------- Stage 2: Runtime --------
FROM alpine:3.22.1

# Install CA certs (needed for HTTPS, SQL TLS, etc.)
RUN apk --no-cache add ca-certificates

# Create working directory
WORKDIR /app

# Copy built binary from builder
COPY --from=builder /app/golang-shop .

# Copy .env if needed at runtime (optional)
COPY .env .env

# Use a minimal user for safety
USER nobody

# App command
CMD ["./golang-shop", "-initdb", "t", "-initView", "t"]

EXPOSE 8001
