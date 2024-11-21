FROM golang:1.22.2 AS builder

WORKDIR /app

COPY . .

RUN go mod download

RUN GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main ./cmd/client

FROM scratch

COPY --from=builder /app/main /app
COPY --from=builder /app/configs/config.yaml /app/configs/config.yaml

ENTRYPOINT ["/app/main"]
