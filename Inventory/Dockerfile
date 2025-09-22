# ---------- Stage 1: Build ----------
FROM golang:1.25 AS builder

WORKDIR /app

# Cache dependencies first
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Copy source code
COPY . .

# Build statically linked binary (safer for minimal base)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o inventory .

# ---------- Stage 2: Runtime ----------
FROM alpine:3.22.1

# Install CA certificates (for HTTPS, DB connections)
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy built binary from the builder stage
COPY --from=builder /app/inventory .

# COPY .env .


# Use a non-root user for safety (Alpine has `nobody`)
USER nobody

# Start your app
CMD ["./inventory", "-initdb", "t", "-initView", "f"]

EXPOSE 8004
