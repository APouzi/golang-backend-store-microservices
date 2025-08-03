# -------- Stage 1: Build --------
FROM golang:1.23 AS builder

WORKDIR /app

# Cache go modules separately
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Copy source code
COPY . .

# Build static binary for portability and minimal runtime base
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o customer-rep-manager .

# -------- Stage 2: Runtime --------
FROM alpine:latest

# Install CA certs if needed for TLS/HTTPS (MySQL, APIs, etc.)
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy built binary from builder
COPY --from=builder /app/customer-rep-manager .

COPY .env .
# Use a minimal user
USER nobody

# Default command
CMD ["./customer-rep-manager", "-initdb", "t", "-initView", "t"]

EXPOSE 8080
