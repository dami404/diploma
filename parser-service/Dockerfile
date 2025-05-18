FROM golang:1.24 AS builder

WORKDIR /app

COPY . .

RUN go mod download
RUN PWGO_VER=$(grep -oE "playwright-go v\S+" /app/go.mod | sed 's/playwright-go //g') \
    && go install github.com/playwright-community/playwright-go/cmd/playwright@${PWGO_VER}
RUN go build -o parser ./cmd/parser/main.go

FROM ubuntu:noble
COPY --from=builder  /go/bin/playwright /app /
RUN apt-get update && apt-get install -y ca-certificates tzdata \
    # Install dependencies and all browsers (or specify one)
    && /playwright install --with-deps \
    && rm -rf /var/lib/apt/lists/*

CMD ["./parser","./config.yaml"]


