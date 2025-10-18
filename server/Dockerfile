FROM golang:1.24-alpine AS builder

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
WORKDIR /app
COPY . .

RUN go mod download
RUN go build -ldflags="-s -w" -o server

FROM alpine:3.20

WORKDIR /app
COPY --from=builder /app/server ./server

ENTRYPOINT ["./server"]

