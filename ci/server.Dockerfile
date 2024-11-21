FROM golang:1.22.2 AS builder

WORKDIR /app

COPY . /app

RUN go mod download

RUN GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main ./cmd/server

FROM scratch

WORKDIR /app

COPY --from=builder /app/main /app
COPY --from=builder /app/configs /app/configs

EXPOSE 8888

ENTRYPOINT ["/app/main"]
