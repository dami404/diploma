FROM golang:1.24-alpine

WORKDIR /app

COPY . .

RUN go mod download
RUN go build -o web ./cmd/web/main.go

CMD ["./web","./config.yaml"]