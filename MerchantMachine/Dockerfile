# ---------- STAGE 1: Build ----------
FROM golang:1.25 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
# Build static binary
RUN CGO_ENABLED=0 GOEXPERIMENT=jsonv2 GOOS=linux GOARCH=amd64 go build -o golang-shop .

# ---------- STAGE 2: Runtime ----------
FROM alpine:3.22.1

# Install optional CA certs or debugging tools if needed
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy the compiled Go binary
COPY --from=builder /app/golang-shop .

# Copy your .env file (optional, if we rely on it in production)
COPY .env .


RUN adduser -D -u 10001 appuser
USER appuser

# App runs with env file assumed in /app
CMD ["./golang-shop", "-initdb", "t", "-initView", "t"]

EXPOSE 8000
