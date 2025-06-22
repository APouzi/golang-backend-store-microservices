FROM golang:1.24-alpine AS builder
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -ldflags="-s -w" -o inventory

FROM scratch
COPY --from=builder /app/inventory /inventory
ENTRYPOINT ["/inventory"]
EXPOSE 8004