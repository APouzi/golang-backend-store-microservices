FROM golang:1.24

# RUN apk add --no-cache gcc musl-dev

WORKDIR /app

COPY ./go.mod ./go.sum ./
RUN go mod download && go mod verify

# COPY . .
COPY . .

RUN go build -o golang-shop .
# RUN chmod +x app


CMD ["./golang-shop", "-initdb", "t", "-initView", "t"]
EXPOSE 8000