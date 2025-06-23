FROM golang:1.24

WORKDIR /app

COPY ./go.mod ./go.sum ./
RUN go mod download && go mod verify

COPY . .

RUN go build -o golang-shop .
# RUN chmod +x app


CMD ["./golang-shop", "-initdb", "t", "-initView", "f"]
EXPOSE 8080