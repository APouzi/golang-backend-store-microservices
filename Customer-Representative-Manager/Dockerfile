FROM golang:1.23

# RUN apk add --no-cache gcc musl-dev

WORKDIR /app

COPY ./go.mod ./go.sum ./
RUN go mod download && go mod verify

# COPY . .
COPY . .

RUN go build -o customer-rep-manager .
# RUN chmod +x app


CMD ["./customer-rep-manager", "-initdb", "t", "-initView", "t"]
