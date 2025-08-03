# ---------- Stage 1: Build ----------
FROM golang:latest AS builder

WORKDIR /app

# Cache dependencies first
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Copy source code
COPY . .

# Build statically linked binary (safer for minimal base)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o golang-shop .

# ---------- Stage 2: Runtime ----------
FROM alpine:latest

# Install CA certificates (for HTTPS, DB connections)
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy built binary from the builder stage
COPY --from=builder /app/golang-shop .

COPY .env .

COPY sql ./sql

# Use a non-root user for safety (Alpine has `nobody`)
USER nobody

# Start your app
CMD ["./golang-shop", "-initdb", "t", "-initView", "f"]

EXPOSE 8080
