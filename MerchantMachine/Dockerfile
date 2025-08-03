# ---------- STAGE 1: Build ----------
FROM golang:1.21 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

# Build static binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o golang-shop .

# ---------- STAGE 2: Runtime ----------
FROM alpine:latest

# Install optional CA certs or debugging tools if needed
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy the compiled Go binary
COPY --from=builder /app/golang-shop .

# Copy your .env file (optional, if we rely on it in production)
COPY .env .

# Use non-root user if needed (Alpine has `nobody`)
USER nobody

# App runs with env file assumed in /app
CMD ["./golang-shop", "-initdb", "t", "-initView", "t"]

EXPOSE 8000
