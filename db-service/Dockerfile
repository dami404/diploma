# db-service/Dockerfile
FROM golang:1.24-alpine

WORKDIR /app

COPY . .

RUN go mod download
RUN go build -o db-service ./cmd/db/main.go

CMD ["./db-service","./config.yaml"]